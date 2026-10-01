// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package snowflake

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"database/sql"
	"encoding/pem"
	"fmt"
	"net/url"
	"regexp"
	"sync"
	"time"

	"github.com/hashicorp/errwrap"
	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-secure-stdlib/parseutil"
	"github.com/hashicorp/vault/sdk/database/helper/dbutil"

	"github.com/hashicorp/vault/sdk/database/helper/connutil"
	"github.com/mitchellh/mapstructure"
	"github.com/snowflakedb/gosnowflake"
)

var (
	ErrInvalidSnowflakeURL           = fmt.Errorf("invalid connection URL format, expect <account_name>.snowflakecomputing.com/<db_name>")
	ErrInvalidPrivateKey             = fmt.Errorf("failed to read provided private_key")
	accountAndDBNameFromConnURLRegex = regexp.MustCompile(`^(.+)\.snowflakecomputing.com/(.+)$`) // Expected format: <account_name>.snowflakecomputing.com/<db_name>
)

type snowflakeConnectionProducer struct {
	ConnectionURL            string      `json:"connection_url"`
	MaxOpenConnections       int         `json:"max_open_connections"`
	MaxIdleConnections       int         `json:"max_idle_connections"`
	MaxConnectionLifetimeRaw interface{} `json:"max_connection_lifetime"`
	Username                 string      `json:"username"`
	Password                 string      `json:"password"`
	PrivateKey               []byte      `json:"private_key"`
	UsernameTemplate         string      `json:"username_template"`
	DisableEscaping          bool        `json:"disable_escaping"`

	Initialized           bool
	RawConfig             map[string]any
	Type                  string
	maxConnectionLifetime time.Duration
	logger                log.Logger
	snowflakeDB           *sql.DB
	mu                    sync.RWMutex
}

func (c *snowflakeConnectionProducer) secretValues() map[string]string {
	return map[string]string{
		c.Password:           "[password]",
		string(c.PrivateKey): "[private_key]",
	}
}

func (c *snowflakeConnectionProducer) Init(ctx context.Context, initConfig map[string]interface{}, verifyConnection bool) (saveConfig map[string]interface{}, err error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.RawConfig = initConfig

	c.logger = log.New(&log.LoggerOptions{})

	decoderConfig := &mapstructure.DecoderConfig{
		Result:           c,
		WeaklyTypedInput: true,
		TagName:          "json",
	}

	decoder, err := mapstructure.NewDecoder(decoderConfig)
	if err != nil {
		return nil, err
	}

	err = decoder.Decode(initConfig)
	if err != nil {
		return nil, err
	}

	if len(c.ConnectionURL) == 0 {
		return nil, fmt.Errorf("connection_url cannot be empty")
	}

	if len(c.Password) > 0 {
		// Return an error here once Snowflake ends support for password auth.
		c.logger.Warn("[DEPRECATED] Single-factor password authentication is deprecated in Snowflake and will be removed by November 2025. " +
			"Key pair authentication will be required after this date.")

		username := c.Username
		password := c.Password

		if !c.DisableEscaping {
			username = url.PathEscape(c.Username)
			password = url.PathEscape(c.Password)
		}

		// Replace templated username and password in connection URL with actual values
		c.ConnectionURL = dbutil.QueryHelper(c.ConnectionURL, map[string]string{
			"username": username,
			"password": password,
		})
	}

	if c.MaxOpenConnections == 0 {
		c.MaxOpenConnections = 4
	}

	if c.MaxIdleConnections == 0 {
		c.MaxIdleConnections = c.MaxOpenConnections
	}
	if c.MaxIdleConnections > c.MaxOpenConnections {
		c.MaxIdleConnections = c.MaxOpenConnections
	}
	if c.MaxConnectionLifetimeRaw == nil {
		c.MaxConnectionLifetimeRaw = "0s"
	}

	c.maxConnectionLifetime, err = parseutil.ParseDurationSecond(c.MaxConnectionLifetimeRaw)
	if err != nil {
		return nil, errwrap.Wrapf("invalid max_connection_lifetime: {{err}}", err)
	}

	c.Initialized = true

	if verifyConnection {
		if _, err := c.Connection(ctx); err != nil {
			c.close()
			return nil, fmt.Errorf("error verifying connection: %w", err)
		}
	}

	return initConfig, nil
}

func (c *snowflakeConnectionProducer) Initialize(ctx context.Context, config map[string]any, verifyConnection bool) error {
	_, err := c.Init(ctx, config, verifyConnection)
	return err
}

func (c *snowflakeConnectionProducer) Connection(ctx context.Context) (interface{}, error) {
	// This is intentionally not grabbing the lock since the calling functions (e.g. CreateUser)
	// are claiming it.

	if !c.Initialized {
		return nil, connutil.ErrNotInitialized
	}

	// If we already have a DB, test it and return
	if c.snowflakeDB != nil {
		if err := c.snowflakeDB.PingContext(ctx); err == nil {
			return c.snowflakeDB, nil
		}
		// If the ping was unsuccessful, close it and ignore errors as we'll be
		// reestablishing anyways
		c.snowflakeDB.Close()
	}

	var db *sql.DB
	var err error
	if len(c.PrivateKey) > 0 {
		db, err = openSnowflake(c.ConnectionURL, c.Username, c.PrivateKey)
		if err != nil {
			return nil, fmt.Errorf("error opening Snowflake connection using key-pair auth: %w", err)
		}
	} else {
		db, err = sql.Open(snowflakeSQLTypeName, c.ConnectionURL)
		if err != nil {
			return nil, fmt.Errorf("error opening Snowflake connection using user-pass auth: %w", err)
		}
	}

	c.snowflakeDB = db
	c.snowflakeDB.SetMaxOpenConns(c.MaxOpenConnections)
	c.snowflakeDB.SetMaxIdleConns(c.MaxIdleConnections)
	c.snowflakeDB.SetConnMaxLifetime(c.maxConnectionLifetime)

	return c.snowflakeDB, nil
}

// close terminates the database connection without locking
func (c *snowflakeConnectionProducer) close() error {
	if c.snowflakeDB != nil {
		if err := c.snowflakeDB.Close(); err != nil {
			return err
		}
	}

	return nil
}

// Close terminates the database connection with locking
func (c *snowflakeConnectionProducer) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.snowflakeDB = nil

	return c.close()
}

// Open the DB connection to Snowflake or return an error.
func openSnowflake(connectionURL, username string, providedPrivateKey []byte) (*sql.DB, error) {
	cfg, err := getSnowflakeConfig(connectionURL, username, providedPrivateKey)
	if err != nil {
		return nil, fmt.Errorf("error constructing snowflake config: %w", err)
	}
	connector := gosnowflake.NewConnector(gosnowflake.SnowflakeDriver{}, *cfg)

	return sql.OpenDB(connector), nil
}

func getSnowflakeConfig(connectionURL, username string, providedPrivateKey []byte) (*gosnowflake.Config, error) {
	// <account_name>.snowflakecomputing.com/<db_name>?queryParameters...
	u, err := url.Parse(connectionURL)
	if err != nil {
		return nil, fmt.Errorf("error parsing Snowflake connection URL %q: %w", connectionURL, err)
	}

	// add authenticator query param to URL to indicate JWT auth
	// https://pkg.go.dev/github.com/snowflakedb/gosnowflake#hdr-JWT_authentication
	q := u.Query()
	q.Set("authenticator", gosnowflake.AuthTypeJwt.String())
	//q.Set("privateKey", "true") // This is needed to avoid gosnowflake trying to read the private key from a file path
	u.RawQuery = q.Encode()

	// construct dsn for gosnowflake
	// "user:""@<account_name>.snowflakecomputing.com/<db_name>?queryParameters...
	dsn := fmt.Sprintf("%s:%s@%s", username, "", u.String())
	cfg, err := gosnowflake.ParseDSN(dsn)
	if err != nil {
		return nil, fmt.Errorf("error parsing Snowflake DSN %s; err=%w", dsn, err)
	}

	privateKey, err := getPrivateKey(providedPrivateKey)
	if err != nil {
		return nil, err
	}

	cfg.PrivateKey = privateKey

	return cfg, nil
}

// Open and decode the private key file
func getPrivateKey(providedPrivateKey []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(providedPrivateKey)
	if block == nil {
		return nil, ErrInvalidPrivateKey
	}

	// key-type supplied in this part of the workflow has to be private.
	// Public keys are set up directly on the server side in Snowflake.
	// https://docs.snowflake.com/en/user-guide/key-pair-auth#assign-the-public-key-to-a-snowflake-user
	if block.Type != "PRIVATE KEY" {
		return nil, fmt.Errorf("unexpected private key type, expected type 'PRIVATE KEY', got '%s'", block.Type)
	}

	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key to PKCS8: %w", err)
	}

	privateKey, ok := key.(*rsa.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("private key was parsed into an unexpected type")
	}

	return privateKey, nil
}
