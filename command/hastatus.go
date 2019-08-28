/*
  Copyright 2019 Davide Madrisan <davide.madrisan@gmail.com>

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

// HAStatusCommand is a CLI Command that holds the attributes of the command `hastatus`.
type HAStatusCommand struct {
	*BaseCommand
}

// Synopsis returns a short synopsis of the `hastatus` command.
func (c *HAStatusCommand) Synopsis() string {
	return "Returns the Vault HA Cluster status"
}

// Help returns a long-form help text of the `status` command.
func (c *HAStatusCommand) Help() string {
	helpText := `
Usage: hashicorp-vault-monitor hastatus [options]

  This command returns the HA cluster status.

    $ hashicorp-vault-monitor hastatus

  Additional flags and more advanced use cases are detailed below.

    -address=<string>
       Address of the Vault server. The default is https://127.0.0.1:8200. This
       can also be specified via the VAULT_ADDR environment variable.

    -output=<string>
       Specify an output format. Can be 'default' or 'nagios'.

  The exit code reflects the seal status:

      - %d - the HA cluster is enabled and the node is active or in standby mode
      - %d - the HA cluster is enabled, the node is in standby mode but the active node is unknown
      - %d - the HA cluster node is not enabled
      - %d - an error occurred

  For a full list of examples, please see the online documentation.
`
	return fmt.Sprintf(helpText,
		StateOk, StateWarning, StateCritical, StateUndefined)
}

// Run executes the `hastatus` command with the given CLI instance and command-line arguments.
func (c *HAStatusCommand) Run(args []string) int {
	cmdFlags := flag.NewFlagSet("hastatus", flag.ContinueOnError)
	cmdFlags.Usage = func() { c.Ui.Output(c.Help()) }
	cmdFlags.StringVar(&c.Address, "address", addressDefault, addressDescr)
	cmdFlags.StringVar(&c.OutputFormat, "output", "default", outputFormatDescr)

	if err := cmdFlags.Parse(args); err != nil {
		c.Ui.Error(err.Error())
		return StateUndefined
	}

	out, err := c.OutputHandle()
	if err != nil {
		c.Ui.Error(err.Error())
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

	leaderStatus, err := client.Sys().Leader()
	if err != nil {
		out.Critical("Error checking leader status: %s", err)
		return StateCritical
	}

	if !leaderStatus.HAEnabled {
		out.Critical("Vault HA (%s) is not enabled", status.ClusterName)
		return StateCritical
	}

	modeInfo := "Active Node"
	retCode := StateOk

	if !leaderStatus.IsSelf {
		if leaderStatus.LeaderAddress == "" {
			leaderStatus.LeaderAddress = "<none>"
			retCode = StateWarning
		}
		modeInfo = fmt.Sprintf("Standby Node (Active Node Address: %s)",
			leaderStatus.LeaderAddress)
	}

	out.Output("Vault HA (%s) is enabled, %s",
		status.ClusterName,
		modeInfo)

	return retCode
}
