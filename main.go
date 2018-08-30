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
	"strings"
)

const defaultVaultAddr = "https://127.0.0.1:8200"

var client *api.Client // https://godoc.org/github.com/hashicorp/vault/api
var address string
var status bool
var policies string
var token string

func init() {
	flag.StringVar(&address, "address", defaultVaultAddr,
		"The address of the Vault server. "+
			"Overrides the "+api.EnvVaultAddress+" environment variable if set")
	flag.BoolVar(&status, "status", false,
		"Returns the Vault status (sealed/unsealed)")
	flag.StringVar(&policies, "policies", "",
		"Comma-separated list of policies to be checked for existance")
	flag.StringVar(&token, "token", "",
		"The token to access Vault. "+
			"Overrides the "+api.EnvVaultToken+" environment variable if set")
}

func initClient(address string) (*api.Client, error) {
	client, err := api.NewClient(&api.Config{
		Address: address,
	})
	if err != nil {
		return nil, err
	}

	return client, nil
}

func checkSealStatus(address string) (bool, error) {
	client, err := initClient(address)
	if err != nil {
		return true, err
	}

	status, err := client.Sys().SealStatus()
	if err != nil {
		return true, err
	}

	return status.Sealed, nil
}

func contains(items []string, item string) bool {
	for _, i := range items {
		if i == item {
			return true
		}
	}
	return false
}

func checkForPolicies(address, token string, policies []string) error {
	client, err := initClient(address)
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
		if !contains(activePolicies, policy) {
			return errors.New("No such Vault Policy: " + policy)
		}
	}

	return nil
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
		sealStatus, err := checkSealStatus(address)
		if err != nil {
			log.Fatal(err)
		}
		if sealStatus {
			fmt.Println("Error: Vault sealed")
		} else {
			fmt.Println("Vault unsealed")
		}
	} else if policies != "" {
		err := checkForPolicies(
			address, token, strings.Split(policies, ","))
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		} else {
			fmt.Println("All the Vault Policies are available")
		}
	} else {
		fmt.Fprintln(os.Stderr, "Syntax error: missing -status or -policies flag")
		os.Exit(1)
	}
}
