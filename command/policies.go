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
)

// PoliciesCommand is a CLI Command that holds the attributes of the command `policies`.
type PoliciesCommand struct {
	*BaseCommand
	Policies []string
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

    -output=<string>
       Specify an output format. Can be 'default' or 'nagios'.

  The exit code reflects the status of the policies:

      - %d - the secret has been successfully read
      - %d - the secret cannot be found of read
      - %d - error

  For a full list of examples, please see the online documentation.
`
	return fmt.Sprintf(helpText,
		StateOk, StateCritical, StateUndefined)
}

// Run executes the `policies` command with the given CLI instance and command-line arguments.
func (c *PoliciesCommand) Run(args []string) int {
	cmdFlags := flag.NewFlagSet("policies", flag.ContinueOnError)
	cmdFlags.Usage = func() { c.Ui.Output(c.Help()) }
	cmdFlags.StringVar(&c.Address, "address", addressDefault, addressDescr)
	cmdFlags.StringVar(&c.Token, "token", tokenDefault, tokenDescr)
	cmdFlags.StringVar(&c.OutputFormat, "output", "default", outputFormatDescr)

	retCode := StateUndefined

	if err := cmdFlags.Parse(args); err != nil {
		c.Ui.Error(err.Error())
		return retCode
	}

	sprintf, err := c.Outputter()
	if err != nil {
		c.Ui.Error(err.Error())
		return retCode
	}

	args = cmdFlags.Args()
	if len(args) < 1 {
		c.Ui.Error(sprintf(
			retCode,
			"Not enough arguments (expected at list 1)"))
		return retCode
	}

	c.Policies = args[0:]

	client, err := c.Client()
	if err != nil {
		c.Ui.Error(sprintf(retCode, err.Error()))
		return retCode
	}

	activePolicies, err := client.Sys().ListPolicies()
	if err != nil {
		c.Ui.Error(sprintf(retCode, "error checking policies: %s", err))
		return retCode
	}

	for _, policy := range c.Policies {
		if !contains(activePolicies, policy) {
			retCode = StateCritical
			c.Ui.Error(sprintf(retCode, "no such Vault policy: %s", policy))
			return retCode
		}
	}

	retCode = StateOk
	c.Ui.Output(sprintf(retCode, "all the policies are defined"))
	return retCode
}
