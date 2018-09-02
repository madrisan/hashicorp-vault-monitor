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
*/

package main // import "github.com/madrisan/hashicorp-vault-monitor"

import (
	"errors"
	"flag"
	"fmt"
	"github.com/hashicorp/vault/api"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const defaultVaultAddr = "https://127.0.0.1:8200"
const (
	StateOk byte = iota
	StateWarning
	StateCritical
)

var client *api.Client // https://godoc.org/github.com/hashicorp/vault/api
var address string
var status bool
var policies string
var readkey string
var token string

type oracle struct {
	message string
	status byte
}

func init() {
	flag.StringVar(&address, "address", defaultVaultAddr,
		"The address of the Vault server. "+
			"Overrides the "+api.EnvVaultAddress+" environment variable if set")
	flag.BoolVar(&status, "status", false,
		"Returns the Vault status (sealed/unsealed)")
	flag.StringVar(&policies, "policies", "",
		"Comma-separated list of policies to be checked for existance")
	flag.StringVar(&readkey, "readkey", "",
		"Read a Vault secret")
	flag.StringVar(&token, "token", "",
		"The token to access Vault. "+
			"Overrides the "+api.EnvVaultToken+" environment variable if set")
}

func (o oracle) String() string {
	if o.status == StateCritical {
		return fmt.Sprintf("Critical: " + o.message)
	} else {
		return fmt.Sprintf(o.message)
	}
}

func VaultClientInit(address string) (*api.Client, error) {
	client, err := api.NewClient(&api.Config{
		Address: address,
	})
	if err != nil {
		return nil, err
	}

	return client, nil
}

func CheckVaultSealStatus(address string) (bool, error) {
	client, err := VaultClientInit(address)
	if err != nil {
		return true, err
	}

	status, err := client.Sys().SealStatus()
	if err != nil {
		return true, err
	}

	return status.Sealed, nil
}

func Contains(items []string, item string) bool {
	for _, i := range items {
		if i == item {
			return true
		}
	}
	return false
}

func CheckVaultPolicies(address, token string, policies []string) error {
	client, err := VaultClientInit(address)
	if err != nil {
		return err
	}

	if token != "" {
		client.SetToken(token)
	}

	activePolicies, err := client.Sys().ListPolicies()
	if err != nil {
		return err
	}

	for _, policy := range policies {
		if !Contains(activePolicies, policy) {
			return errors.New("No such Vault Policy: " + policy)
		}
	}

	return nil
}

func GetRawField(data interface{}, field string) (string, error) {
	var val interface{}
	switch data.(type) {
	case *api.Secret:
		val = data.(*api.Secret).Data[field]
	case map[string]interface{}:
		val = data.(map[string]interface{})[field]
	}

	if val == nil {
		return "", errors.New(
			fmt.Sprintf("Field '%s' not present in secret", field))
	}

	return val.(string), nil
}

func ReadVaultSecret(keypath, address, token string) (string, error) {
	client, err := VaultClientInit(address)
	if err != nil {
		return "", err
	}

	if token != "" {
		client.SetToken(token)
	}

	// see: https://godoc.org/github.com/hashicorp/vault/api#Secret
	path, key := filepath.Split(keypath)
	secret, err := client.Logical().Read(path)
	if err != nil {
		return "", errors.New(
			fmt.Sprintf("Error reading %s: %s", path, err))
	}
	if secret == nil {
		return "", errors.New(
			fmt.Sprintf("No value found at %s", path))
	}

	// secret.Data:
	// - data -> map[akey:this-is-a-test]
	// - metadata -> map[created_time:2018-08-31T15:36:31.894655728Z deletion_time: destroyed:false version:3]
	if data, ok := secret.Data["data"]; ok && data != nil {
		value, err := GetRawField(data, key)
		return value, err
	} else {
		return "", errors.New(fmt.Sprintf("No data found at %s", path))
	}
}

func main() {
	var envAddress string
	var envToken string

	// Parse the environment variables
	if v := os.Getenv(api.EnvVaultAddress); v != "" {
		envAddress = v
	}
	if v := os.Getenv(api.EnvVaultToken); v != "" {
		envToken = v
	}

	if envAddress != "" && address != "" {
		address = envAddress
	}
	if envToken != "" && token != "" {
		token = envToken
	}

	flag.Parse()

	if status {
		sealStatus, err := CheckVaultSealStatus(address)
		if err != nil {
			log.Fatal(err)
		}
		if sealStatus {
			fmt.Print(oracle{
				message: "Vault sealed",
				status: StateCritical,
			}, "\n")
		} else {
			fmt.Print(oracle{
				message: "Vault unsealed",
				status: StateOk,
			}, "\n")
		}
	} else if policies != "" {
		err := CheckVaultPolicies(
			address, token, strings.Split(policies, ","))
		if err != nil {
			fmt.Print(oracle{
				message: err.Error(),
				status: StateCritical,
			}, "\n")
		} else {
			fmt.Print(oracle{
				message: "All the Vault Policies are available",
				status: StateOk,
			}, "\n")
		}
	} else if readkey != "" {
		secret, err := ReadVaultSecret(readkey, address, token)
		if err != nil {
			fmt.Print(oracle{
				message: err.Error(),
				status: StateCritical,
			}, "\n")
		} else {
			// export VAULT_ADDR="http://127.0.0.1:8200"
			// export VAULT_TOKEN="..."
			// $GOPATH/bin/hashicorp-vault-monitor -readkey secret/data/test/testkey
			//   -> /v1/map[testkey:this-is-a-secret] -> "this-is-a-secret"
			fmt.Print(oracle{
				message: "Found value: '" + secret + "'",
				status: StateOk,
			}, "\n")
		}
	} else {
		fmt.Fprintln(os.Stderr,
			"Syntax error: missing -readkey, -status, or -policies flag")
		os.Exit(1)
	}
}
