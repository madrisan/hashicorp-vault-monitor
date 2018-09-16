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

// StatusCommand is a CLI Command that holds the attributes of the command `status`.
type StatusCommand struct {
	Address string
	Ui      cli.Ui
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

    $ hashicorp-vault-monitor status --address https://127.0.0.1:8200

  The exit code reflects the seal status:

      - 0 - unsealed
      - 1 - error
      - 2 - sealed

  For a full list of examples, please see the documentation.

`
	return strings.TrimSpace(helpText)
}

// Run executes the `status` command with the given CLI instance and command-line arguments.
func (c *StatusCommand) Run(args []string) int {
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

	cmdFlags := flag.NewFlagSet("status", flag.ContinueOnError)
	cmdFlags.Usage = func() { c.Ui.Output(c.Help()) }
	cmdFlags.StringVar(&c.Address, "address", addressDefault, addressDescr)
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

	status, err := client.Sys().SealStatus()
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error checking seal status: %s", err))
		return 1
	}

	if status.Sealed {
		c.Ui.Output(fmt.Sprintf("Vault is sealed! Unseal Progress: %d/%d",
			status.Progress, status.T))
		return 2
	}

	c.Ui.Output("Vault is unsealed")
	return 0
}
