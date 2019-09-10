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

// StatusCommand is a CLI Command that holds the attributes of the command `status`.
type StatusCommand struct {
	*BaseCommand
}

// Synopsis returns a short synopsis of the `status` command.
func (c *StatusCommand) Synopsis() string {
	return "Returns the Vault status (sealed/unsealed)"
}

// Help returns a long-form help text of the `status` command.
func (c *StatusCommand) Help() string {
	helpText := `
Usage: hashicorp-vault-monitor status [options]

  This command returns the status (sealed/unsealed) of a Vault server.

    $ hashicorp-vault-monitor status

  Additional flags and more advanced use cases are detailed below.

    -address=<string>
       Address of the Vault server. The default is https://127.0.0.1:8200. This
       can also be specified via the VAULT_ADDR environment variable.

    -output=<string>
       Specify an output format. Can be 'default' or 'nagios'.

  The exit code reflects the seal status:

      - %d - the vault node is unsealed
      - %d - the vault node is sealed
      - %d - an error occurred

  For a full list of examples, please see the online documentation.
`
	return fmt.Sprintf(helpText,
		StateOk, StateCritical, StateUndefined)
}

// Run executes the `status` command with the given CLI instance and command-line arguments.
func (c *StatusCommand) Run(args []string) int {
	cmdFlags := flag.NewFlagSet("status", flag.ContinueOnError)
	cmdFlags.Usage = func() { c.UI.Output(c.Help()) }
	cmdFlags.StringVar(&c.Address, "address", addressDefault, addressDescr)
	cmdFlags.StringVar(&c.OutputFormat, "output", "default", outputFormatDescr)

	if err := cmdFlags.Parse(args); err != nil {
		c.UI.Error(err.Error())
		return StateUndefined
	}

	out, err := c.OutputHandle()
	if err != nil {
		c.UI.Error(err.Error())
		return StateUndefined
	}

	args = cmdFlags.Args()
	if len(args) > 0 {
		out.Undefined("Too many arguments (expected 0, got %d)", len(args))
		return StateUndefined
	}

	client, err := c.Client()
	if err != nil {
		out.Undefined(err.Error())
		return StateUndefined
	}

	status, err := client.Sys().SealStatus()
	if err != nil {
		out.Undefined("error checking seal status: %s", err)
		return StateUndefined
	}

	if status.Sealed {
		out.Critical("Vault (%s) is sealed! Unseal Progress: %d/%d",
			status.ClusterName,
			status.Progress,
			status.T)
		return StateCritical
	}

	out.Output("Vault (%s) is unsealed", status.ClusterName)
	return StateOk
}
