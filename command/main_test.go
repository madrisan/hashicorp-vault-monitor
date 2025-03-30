/*
  Copyright 2018 Davide Madrisan <davide.madrisan@gmail.com>

  Licensed under the Mozilla Public License, Version 2.0 (the "License");
  you may not use this file except in compliance with the License.
  You may obtain a copy of the License at

      https://www.mozilla.org/en-US/MPL/2.0/

  Unless required by applicable law or agreed to in writing, software
  distributed under the License is distributed on an "AS IS" BASIS,
  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
  See the License for the specific language governing permissions and
  limitations under the License.

  Credit: This source code is based on the HashiCorp Vault testing code
          (but bugs are mine).
*/

package command

import (
	"context"
	"encoding/base64"
	"net"
	"net/http"
	"testing"
	"time"

	log "github.com/hashicorp/go-hclog"
	kv "github.com/hashicorp/vault-plugin-secrets-kv"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/audit"
	"github.com/hashicorp/vault/builtin/logical/pki"
	"github.com/hashicorp/vault/builtin/logical/ssh"
	"github.com/hashicorp/vault/builtin/logical/transit"
	"github.com/hashicorp/vault/helper/builtinplugins"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/vault"

	credUserpass "github.com/hashicorp/vault/builtin/credential/userpass"
	vaulthttp "github.com/hashicorp/vault/http"
)

var (
	defaultVaultLogger             = log.NewNullLogger()
	defaultVaultCredentialBackends = map[string]logical.Factory{
		"userpass": credUserpass.Factory,
	}

	defaultVaultAuditBackends = map[string]audit.Factory{
		"file": audit.NewFileBackend,
	}

	defaultVaultLogicalBackends = map[string]logical.Factory{
		"generic-leased": vault.LeasedPassthroughBackendFactory,
		"pki":            pki.Factory,
		"ssh":            ssh.Factory,
		"transit":        transit.Factory,
		"kv":             kv.Factory,
	}
)

// testVaultServerUnseal creates a test vault cluster and returns a configured
// API client, list of unseal keys (as strings), and a closer function.
func testVaultServerUnseal(tb testing.TB) (*api.Client, []string, func()) {
	return testVaultServerUnsealWithKVVersionWithSeal(tb, "1", nil)
}

// Convert the unseal keys to base64 encoded, since these are how the user
// will get them.
func encodeKeys(rawKeys [][]byte) []string {
	keys := make([]string, len(rawKeys))
	for i := range rawKeys {
		keys[i] = base64.StdEncoding.EncodeToString(rawKeys[i])
	}
	return keys
}

// testVaultServerCoreConfig creates a new vault cluster with the given core
// configuration. This is a lower-level test helper. If the seal config supports recovery keys, then
// recovery keys are returned. Otherwise, unseal keys are returned
func testVaultServerCoreConfigWithOpts(tb testing.TB, coreConfig *vault.CoreConfig, opts *vault.TestClusterOptions) (*api.Client, []string, func()) {
	tb.Helper()

	cluster := vault.NewTestCluster(tb, coreConfig, opts)
	cluster.Start()

	// Make it easy to get access to the active
	core := cluster.Cores[0].Core
	vault.TestWaitActive(tb, core)

	// Get the client already setup for us!
	client := cluster.Cores[0].Client
	client.SetToken(cluster.RootToken)

	var keys [][]byte
	if coreConfig.Seal != nil && coreConfig.Seal.RecoveryKeySupported() {
		keys = cluster.RecoveryKeys
	} else {
		keys = cluster.BarrierKeys
	}

	return client, encodeKeys(keys), cluster.Cleanup
}

// testVaultServerUnseal creates a test vault cluster and returns a configured
// API client, list of unseal keys (as strings), and a closer function.
func testVaultServerUnsealWithKVVersionWithSeal(tb testing.TB, kvVersion string, seal vault.Seal) (*api.Client, []string, func()) {
	tb.Helper()

	return testVaultServerCoreConfigWithOpts(tb, &vault.CoreConfig{
		CredentialBackends: defaultVaultCredentialBackends,
		AuditBackends:      defaultVaultAuditBackends,
		LogicalBackends:    defaultVaultLogicalBackends,
		BuiltinRegistry:    builtinplugins.Registry,
		Seal:               seal,
	}, &vault.TestClusterOptions{
		HandlerFunc: vaulthttp.Handler,
		NumCores:    1,
		KVVersion:   kvVersion,
	})
}

// testVaultServerUnseal creates a test vault cluster and returns a configured
// API client, list of unseal keys (as strings), and a closer function
// configured with the given plugin directory.
func testVaultServerPluginDir(tb testing.TB, pluginDir string) (*api.Client, []string, func()) {
	tb.Helper()

	return testVaultServerCoreConfig(tb, &vault.CoreConfig{
		CredentialBackends: defaultVaultCredentialBackends,
		AuditBackends:      defaultVaultAuditBackends,
		LogicalBackends:    defaultVaultLogicalBackends,
		PluginDirectory:    pluginDir,
		BuiltinRegistry:    builtinplugins.Registry,
	})
}

// testVaultServerCoreConfig creates a new vault cluster with the given core
// configuration. This is a lower-level test helper.
func testVaultServerCoreConfig(tb testing.TB, coreConfig *vault.CoreConfig) (*api.Client, []string, func()) {
	tb.Helper()

	cluster := vault.NewTestCluster(tb, coreConfig, &vault.TestClusterOptions{
		HandlerFunc: vaulthttp.Handler,
		NumCores:    1, // Default is 3, but we don't need that many
	})
	cluster.Start()

	// Make it easy to get access to the active
	core := cluster.Cores[0].Core
	vault.TestWaitActive(tb, core)

	// Get the client already setup for us!
	client := cluster.Cores[0].Client
	client.SetToken(cluster.RootToken)

	// Convert the unseal keys to base64 encoded, since these are how the user
	// will get them.
	unsealKeys := make([]string, len(cluster.BarrierKeys))
	for i := range unsealKeys {
		unsealKeys[i] = base64.StdEncoding.EncodeToString(cluster.BarrierKeys[i])
	}

	return client, unsealKeys, func() { defer cluster.Cleanup() }
}

// testVaultServerBad creates an http server that returns a 500 on each request
// to simulate failures.
func testVaultServerBad(t testing.TB) (*api.Client, func()) {
	t.Helper()

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}

	server := &http.Server{
		Addr: "127.0.0.1:0",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "500 internal server error", http.StatusInternalServerError)
		}),
		ReadTimeout:       1 * time.Second,
		ReadHeaderTimeout: 1 * time.Second,
		WriteTimeout:      1 * time.Second,
		IdleTimeout:       1 * time.Second,
	}

	go func() {
		if err := server.Serve(listener); err != nil && err != http.ErrServerClosed {
			t.Fatal(err)
		}
	}()

	client, err := api.NewClient(&api.Config{
		Address: "http://" + listener.Addr().String(),
	})
	if err != nil {
		t.Fatal(err)
	}

	return client, func() {
		ctx, done := context.WithTimeout(context.Background(), 5*time.Second)
		defer done()

		if err := server.Shutdown(ctx); err != nil {
			t.Fatal(err)
		}
	}
}
