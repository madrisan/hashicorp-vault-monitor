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

package command

import (
	"flag"
	"fmt"
	"strings"

	"github.com/hashicorp/vault/api"
	"github.com/madrisan/hashicorp-vault-monitor/vault"
	"github.com/mitchellh/cli"
)

const (
	policiesDescr = "Comma-separated list of policies to be checked for existence"
)

type PoliciesCommand struct {
	Address  string
	Token    string
	Policies string
	Ui       cli.Ui
}

func contains(items []string, item string) bool {
	for _, i := range items {
		if i == item {
			return true
		}
	}
	return false
}

// Synopsis returns a short synopsis of the `policies` command.
func (c *PoliciesCommand) Synopsis() string {
	return "Check the active policies of a Vault server"
}

// Help returns a long-form help text of the `policies` command.
func (c *PoliciesCommand) Help() string {
	helpText := `
Usage: hashicorp-vault-monitor policies [options]

  This command check for the active policies of a Vault server.

    $ hashicorp-vault-monitor policies \
        --defined "custpolicy1,custpolicy2" \
        --address https://127.0.0.1:8200 --token "12e2bf2b-3b82-9eff-07e4-8c7ad97715a9"

  The exit code reflects the seal status:

      - 0 - all the comma-separated list of policies are defined
      - 1 - error
      - 2 - at least one of the policies is not defined

  For a full list of examples, please see the documentation.

`
	return strings.TrimSpace(helpText)
}

// Run executes the `policies` command with the given CLI instance and command-line arguments.
func (c *PoliciesCommand) Run(args []string) int {
	vaultConfig := api.DefaultConfig()
	if vaultConfig == nil {
		c.Ui.Error("could not create/read default configuration for Vault")
		return 1
	}
	if vaultConfig.Error != nil {
		c.Ui.Error("error encountered setting up default configuration: " +
			vaultConfig.Error.Error())
		return 1
	}

	cmdFlags := flag.NewFlagSet("policies", flag.ContinueOnError)
	cmdFlags.Usage = func() { c.Ui.Output(c.Help()) }
	cmdFlags.StringVar(&c.Address, "address", addressDefault, addressDescr)
	cmdFlags.StringVar(&c.Token, "token", tokenDefault, tokenDescr)
	cmdFlags.StringVar(&c.Policies, "defined", "", policiesDescr)

	if err := cmdFlags.Parse(args); err != nil {
		c.Ui.Error(err.Error())
		return 1
	}

	// note that `api.DefaultConfig` execute `api.ReadEnvironment` and thus
	// load also the all the Vault environment variables but `VAULT_TOKEN`
	if c.Address != "" {
		vaultConfig.Address = c.Address
	}

	client, err := vault.ClientInit(c.Address)
	if err != nil {
		c.Ui.Error(err.Error())
		return 1
	}

	if c.Token != "" {
		client.SetToken(c.Token)
	}

	activePolicies, err := client.Sys().ListPolicies()
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error checking policies: %s", err))
		return 1
	}

	for _, policy := range strings.Split(c.Policies, ",") {
		if !contains(activePolicies, policy) {
			c.Ui.Error(fmt.Sprintf("No such Vault policy: %s", policy))
			return 2
		}
	}

	c.Ui.Output("All the policies are defined")
	return 0
}
