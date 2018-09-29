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

	"github.com/hashicorp/vault/api"
	"github.com/madrisan/hashicorp-vault-monitor/vault"
	"github.com/mitchellh/cli"
)

// PoliciesCommand is a CLI Command that holds the attributes of the command `policies`.
type PoliciesCommand struct {
	Address  string
	Token    string
	Policies []string
	Ui       cli.Ui
	client   *api.Client
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
Usage: hashicorp-vault-monitor policies [options] POLICIES

  This command check for the active policies of a Vault server.

    $ hashicorp-vault-monitor policies custpolicy1 custpolicy1 ...

  Additional flags and more advanced use cases are detailed below.

    -address=<string>
       Address of the Vault server. The default is https://127.0.0.1:8200. This
       can also be specified via the VAULT_ADDR environment variable.

    -token=<string>
       Specify a token for authentication. This can also be specified via the
       VAULT_TOKEN environment variable.

  The exit code reflects the status of the policies:

      - %d - the secret has been successfully read
      - %d - the secret cannot be found of read
      - %d - error

  For a full list of examples, please see the online documentation.
`
	return fmt.Sprintf(helpText,
		StateOk, StateCritical, StateError)
}

// Run executes the `policies` command with the given CLI instance and command-line arguments.
func (c *PoliciesCommand) Run(args []string) int {
	vaultConfig := api.DefaultConfig()
	if vaultConfig == nil {
		c.Ui.Error("could not create/read default configuration for Vault")
		return StateError
	}
	if vaultConfig.Error != nil {
		c.Ui.Error("error encountered setting up default configuration: " +
			vaultConfig.Error.Error())
		return StateError
	}

	cmdFlags := flag.NewFlagSet("policies", flag.ContinueOnError)
	cmdFlags.Usage = func() { c.Ui.Output(c.Help()) }
	cmdFlags.StringVar(&c.Address, "address", addressDefault, addressDescr)
	cmdFlags.StringVar(&c.Token, "token", tokenDefault, tokenDescr)

	if err := cmdFlags.Parse(args); err != nil {
		c.Ui.Error(err.Error())
		return StateError
	}

	args = cmdFlags.Args()
	if len(args) < 1 {
		c.Ui.Error(fmt.Sprintf(
			"Not enough arguments (expected at list 1)"))
		return StateError
	}

	c.Policies = args[0:]

	// note that `api.DefaultConfig` execute `api.ReadEnvironment` and thus
	// load also the all the Vault environment variables but `VAULT_TOKEN`
	if c.Address != "" {
		vaultConfig.Address = c.Address
	}

	if c.client == nil {
		client, err := vault.NewClient(c.Address)
		if err != nil {
			c.Ui.Error(err.Error())
			return StateError
		}
		c.client = client
	}

	if c.Token != "" {
		c.client.SetToken(c.Token)
	}

	activePolicies, err := c.client.Sys().ListPolicies()
	if err != nil {
		c.Ui.Error(fmt.Sprintf("error checking policies: %s", err))
		return StateError
	}

	for _, policy := range c.Policies {
		if !contains(activePolicies, policy) {
			c.Ui.Error(fmt.Sprintf("no such Vault policy: %s", policy))
			return StateCritical
		}
	}

	c.Ui.Output("all the policies are defined")
	return StateOk
}
