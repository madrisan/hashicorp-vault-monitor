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
	"os"
	"path/filepath"
	"strings"
)

const defaultVaultAddr = "https://127.0.0.1:8200"
const (
	StateOk byte = iota
	StateWarning
	StateCritical
	StateUnknown
)

var client *api.Client // https://godoc.org/github.com/hashicorp/vault/api
var address string
var status bool
var policies string
var readkey string
var token string

type oracle struct {
	message string
	status  byte
}

func init() {
	flag.StringVar(&address, "address", defaultVaultAddr,
		"The address of the Vault server. "+
			"Overrides the "+api.EnvVaultAddress+" environment variable if set")
	flag.BoolVar(&status, "status", false,
		"Returns the Vault status (sealed/unsealed)")
	flag.StringVar(&policies, "policies", "",
		"Comma-separated list of policies to be checked for existence")
	flag.StringVar(&readkey, "readkey", "",
		"Read a Vault secret")
	flag.StringVar(&token, "token", "",
		"The token to access Vault. "+
			"Overrides the "+api.EnvVaultToken+" environment variable if set")
}

func (o oracle) String() string {
	var status string

	switch o.status {
	case StateOk:
		status = "Ok"
	case StateWarning:
		status = "Warning"
	case StateCritical:
		status = "Critical"
	default:
		status = "Unknown"
	}

	return fmt.Sprintf(status + ": " + o.message)
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

func VaultIsUnsealed(address string) (bool, error) {
	client, err := VaultClientInit(address)
	if err != nil {
		return false, err
	}

	status, err := client.Sys().SealStatus()
	if err != nil {
		return false, err
	}

	return status.Sealed == false, nil
}

func Contains(items []string, item string) bool {
	for _, i := range items {
		if i == item {
			return true
		}
	}
	return false
}

func CheckVaultPolicies(
	address, token string, policies []string) ([]string, error) {

	client, err := VaultClientInit(address)
	if err != nil {
		return nil, err
	}

	if token != "" {
		client.SetToken(token)
	}

	activePolicies, err := client.Sys().ListPolicies()
	if err != nil {
		return nil, err
	}

	for _, policy := range policies {
		if !Contains(activePolicies, policy) {
			return activePolicies,
				errors.New("No such Vault Policy: " + policy)
		}
	}

	return activePolicies, nil
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
		return "", fmt.Errorf("Field '%s' not present in secret", field)
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
		return "", fmt.Errorf("Error reading %s: %s", path, err)
	}
	if secret == nil {
		return "", fmt.Errorf("No value found at %s", path)
	}

	// secret.Data:
	// - data -> map[akey:this-is-a-test]
	// - metadata -> map[created_time:2018-08-31T15:36:31.894655728Z deletion_time: destroyed:false version:3]
	if data, ok := secret.Data["data"]; ok && data != nil {
		value, err := GetRawField(data, key)
		return value, err
	}

	return "", fmt.Errorf("No data found at %s", path)
}

func main() {
	var envAddress string
	var envToken string
	var result oracle

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
		isUnsealed, err := VaultIsUnsealed(address)
		if err != nil {
			result = oracle{
				message: err.Error(),
				status:  StateUnknown,
			}
		}
		if isUnsealed {
			result = oracle{
				message: "Vault is unsealed",
			}
		} else {
			result = oracle{
				message: "Vault is sealed",
				status:  StateCritical,
			}
		}
	} else if policies != "" {
		policies, err := CheckVaultPolicies(
			address, token, strings.Split(policies, ","))
		if err != nil {
			result.message = err.Error()
			if policies != nil {
				result.status = StateCritical
			} else {
				result.status = StateUnknown
			}
		} else {
			result = oracle{
				message: "All the Vault Policies are available",
			}
		}
	} else if readkey != "" {
		secret, err := ReadVaultSecret(readkey, address, token)
		if err != nil {
			result = oracle{
				message: err.Error(),
				status:  StateCritical,
			}
		} else {
			result = oracle{
				message: "Found value: '" + secret + "'",
			}
		}
	} else {
		fmt.Fprintln(os.Stderr,
			"Syntax error: missing -readkey, -status, or -policies flag")
		os.Exit(1)
	}

	fmt.Print(result, "\n")
	os.Exit(int(result.status))
}
