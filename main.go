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
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/hashicorp/vault/api"
	"github.com/madrisan/hashicorp-vault-monitor/command"
	"github.com/madrisan/hashicorp-vault-monitor/version"
)

const (
	StateOk byte = iota
	StateWarning
	StateCritical
	StateUnknown
)

var client *api.Client // https://godoc.org/github.com/hashicorp/vault/api

type oracle struct {
	message string
	status  byte
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

func vaultClientInit(address string) (*api.Client, error) {
	client, err := api.NewClient(&api.Config{
		Address: address,
	})
	if err != nil {
		return nil, err
	}

	return client, nil
}

func VaultIsUnsealed(address string) (bool, error) {
	client, err := vaultClientInit(address)
	if err != nil {
		return false, err
	}

	status, err := client.Sys().SealStatus()
	if err != nil {
		return false, err
	}

	return status.Sealed == false, nil
}

func contains(items []string, item string) bool {
	for _, i := range items {
		if i == item {
			return true
		}
	}
	return false
}

func CheckVaultPolicies(
	address, token string, policies []string) ([]string, error) {

	client, err := vaultClientInit(address)
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
		if !contains(activePolicies, policy) {
			return activePolicies,
				errors.New("no such Vault policy: " + policy)
		}
	}

	return activePolicies, nil
}

func getRawField(data interface{}, field string) (string, error) {
	var val interface{}
	switch data.(type) {
	case *api.Secret:
		val = data.(*api.Secret).Data[field]
	case map[string]interface{}:
		val = data.(map[string]interface{})[field]
	}

	if val == nil {
		return "", fmt.Errorf("field '%s' not present in secret", field)
	}

	return val.(string), nil
}

func ReadVaultSecret(keypath, address, token string) (string, error) {
	client, err := vaultClientInit(address)
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
		return "", fmt.Errorf("error reading %s: %s", path, err)
	}
	if secret == nil {
		return "", fmt.Errorf("no value found at %s", path)
	}

	// secret.Data:
	// - data -> map[akey:this-is-a-test]
	// - metadata -> map[created_time:2018-08-31T15:36:31.894655728Z deletion_time: destroyed:false version:3]
	if data, ok := secret.Data["data"]; ok && data != nil {
		value, err := getRawField(data, key)
		return value, err
	}

	return "", fmt.Errorf("no data found at %s", path)
}

// Version returns the semantic version (see http://semver.org) of the tool.
func Version() string {
	versionInfo := version.GetVersion()
	return versionInfo.FullVersionNumber(true)
}

func main() {
	var result oracle
	opt := command.Run()
	if opt.Status {
		isUnsealed, err := VaultIsUnsealed(opt.Address)
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
	} else if opt.Policies != "" {
		policies, err := CheckVaultPolicies(
			opt.Address, opt.Token, strings.Split(opt.Policies, ","))
		if err != nil {
			result.message = err.Error()
			if policies != nil {
				result.status = StateCritical
			} else {
				result.status = StateUnknown
			}
		} else {
			result = oracle{
				message: "all the Vault policies are available",
			}
		}
	} else if opt.ReadKey != "" {
		secret, err := ReadVaultSecret(opt.ReadKey, opt.Address, opt.Token)
		if err != nil {
			result = oracle{
				message: err.Error(),
				status:  StateCritical,
			}
		} else {
			result = oracle{
				message: "found value: '" + secret + "'",
			}
		}
	} else if opt.Infos {
	        fmt.Println(Version())
		os.Exit(0)
	} else {
		fmt.Fprintln(os.Stderr,
			"Syntax error: missing -readkey, -status, or -policies flag")
		os.Exit(1)
	}

	fmt.Print(result, "\n")
	os.Exit(int(result.status))
}
