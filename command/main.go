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

// Package command provides the logic for parsing command-line arguments
// and run the available monitoring commands
package command

import (
	"fmt"
	"os"

	"github.com/hashicorp/vault/api"
	"github.com/madrisan/hashicorp-vault-monitor/version"
	"github.com/mitchellh/cli"
)

const (
	addressDefault = "https://127.0.0.1:8200"
	addressDescr   = "The address of the Vault server. " +
		"Overrides the " + api.EnvVaultAddress + " environment variable if set"

	tokenDefault = ""
	tokenDescr   = "The token to access Vault. " +
		"Overrides the " + api.EnvVaultToken + " environment variable if set"

	tokenAccessorDefault = ""
	tokenAccessorDescr   = "The token accessor to lookup"

	warningDescr  = "Warning threshold (default: %s)"
	criticalDescr = "Critical threshold (default: %s)"

	outputFormatDescr = "Select an output format ('default' or 'nagios')"
)

// Run initializes a CLI instance and its command state engine.
func Run(args []string) int {
	ui := &cli.BasicUi{
		Reader:      os.Stdin,
		Writer:      os.Stdout,
		ErrorWriter: os.Stderr,
	}

	verInfo := version.GetVersion()
	version := verInfo.FullVersionNumber(true)

	c := cli.NewCLI("hashicorp-vault-monitor", version)
	c.Args = args

	c.Commands = map[string]cli.CommandFactory{
		"get": func() (cli.Command, error) {
			return &GetCommand{
				BaseCommand: &BaseCommand{
					UI: &cli.ColoredUi{
						Ui:          ui,
						ErrorColor:  cli.UiColorRed,
						OutputColor: cli.UiColorGreen,
						WarnColor:   cli.UiColorYellow,
					},
					OutputFormat: "default",
				},
			}, nil
		},
		"hastatus": func() (cli.Command, error) {
			return &HAStatusCommand{
				BaseCommand: &BaseCommand{
					UI: &cli.ColoredUi{
						Ui:          ui,
						ErrorColor:  cli.UiColorRed,
						OutputColor: cli.UiColorGreen,
						WarnColor:   cli.UiColorYellow,
					},
					OutputFormat: "default",
				},
			}, nil
		},
		"policies": func() (cli.Command, error) {
			return &PoliciesCommand{
				BaseCommand: &BaseCommand{
					UI: &cli.ColoredUi{
						Ui:          ui,
						ErrorColor:  cli.UiColorRed,
						OutputColor: cli.UiColorGreen,
						WarnColor:   cli.UiColorYellow,
					},
					OutputFormat: "default",
				},
			}, nil
		},
		"status": func() (cli.Command, error) {
			return &StatusCommand{
				BaseCommand: &BaseCommand{
					UI: &cli.ColoredUi{
						Ui:          ui,
						ErrorColor:  cli.UiColorRed,
						OutputColor: cli.UiColorGreen,
						WarnColor:   cli.UiColorYellow,
					},
					OutputFormat: "default",
				},
			}, nil
		},
		"token-lookup": func() (cli.Command, error) {
			return &TokenLookupCommand{
				BaseCommand: &BaseCommand{
					UI: &cli.ColoredUi{
						Ui:          ui,
						ErrorColor:  cli.UiColorRed,
						OutputColor: cli.UiColorGreen,
						WarnColor:   cli.UiColorYellow,
					},
					OutputFormat: "default",
				},
			}, nil
		},
	}

	exitStatus, err := c.Run()

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error executing CLI: %s\n", err.Error())
		return StateUndefined
	}

	return exitStatus
}
