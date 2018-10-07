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

      - %d - the secret has been successfully read
      - %d - the secret cannot be found of read
      - %d - error

  For a full list of examples, please see the online documentation.
`
	return fmt.Sprintf(helpText,
		StateOk, StateCritical, StateUndefined)
}

// Run executes the `status` command with the given CLI instance and command-line arguments.
func (c *StatusCommand) Run(args []string) int {
	cmdFlags := flag.NewFlagSet("status", flag.ContinueOnError)
	cmdFlags.Usage = func() { c.Ui.Output(c.Help()) }
	cmdFlags.StringVar(&c.Address, "address", addressDefault, addressDescr)
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
	if len(args) > 0 {
		c.Ui.Error(sprintf(
			retCode,
			"Too many arguments (expected 0, got %d)", len(args)))
		return retCode
	}

	client, err := c.Client()
	if err != nil {
		c.Ui.Error(sprintf(retCode, err.Error()))
		return retCode
	}

	status, err := client.Sys().SealStatus()
	if err != nil {
		c.Ui.Error(sprintf(retCode, "error checking seal status: %s", err))
		return retCode
	}

	if status.Sealed {
		retCode = StateCritical
		c.Ui.Output(sprintf(retCode, "Vault is sealed! Unseal Progress: %d/%d",
			status.Progress, status.T))
		return retCode
	}

	retCode = StateOk
	c.Ui.Output(sprintf(retCode, "Vault is unsealed"))
	return retCode
}
