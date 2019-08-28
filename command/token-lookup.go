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
	"time"
)

// TokenLookupCommand is a CLI Command that holds the attributes of the command `token-lookup`.
type TokenLookupCommand struct {
	*BaseCommand
}

// Synopsis returns a short synopsis of the `token-lookup` command.
func (c *TokenLookupCommand) Synopsis() string {
	return "Returns information about the current client token"
}

// Help returns a long-form help text of the `token-lookup` command.
func (c *TokenLookupCommand) Help() string {
	helpText := `
Usage: hashicorp-vault-monitor token-lookup [options]

  This command returns information about the current client token.

    $ hashicorp-vault-monitor token-lookup

  Additional flags and more advanced use cases are detailed below.

    -address=<string>
       Address of the Vault server. The default is https://127.0.0.1:8200. This
       can also be specified via the VAULT_ADDR environment variable.

    -output=<string>
       Specify an output format. Can be 'default' or 'nagios'.

  The exit code reflects the token expiration time:

      - %d - the token is usable
      - %d - the token will expire in less than a week
      - %d - the token will expire in less than 3 days
      - %d - an error occurred

  For a full list of examples, please see the online documentation.
`
	return fmt.Sprintf(helpText,
		StateOk, StateWarning, StateCritical, StateUndefined)
}

// Thresholds in hours for warning and critical status.
const (
	WarningExpirationHours  = 7 * 24
	CriticalExpirationHours = 3 * 24
)

// Run executes the `token-lookup` command with the given CLI instance and command-line arguments.
func (c *TokenLookupCommand) Run(args []string) int {
	cmdFlags := flag.NewFlagSet("token-lookup", flag.ContinueOnError)
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
		out.Error("Too many arguments (expected 0, got %d)", len(args))
		return StateUndefined
	}

	client, err := c.Client()
	if err != nil {
		out.Error(err.Error())
		return StateUndefined
	}

	ta := client.Auth().Token()

	s, err := ta.LookupSelf()
	if err != nil {
		out.Error(err.Error())
		return StateUndefined
	}

	if s.Data == nil || s.Data["expire_time"] == nil {
		out.Error("Cannot get the expire time of the Vault token")
		return StateUndefined
	}

	expireTimeRaw := s.Data["expire_time"]
	expireTimeStr, ok := expireTimeRaw.(string)
	if !ok {
		out.Error("Could not convert expire_time to a string")
		return StateUndefined
	}

	t, err := time.Parse(time.RFC3339Nano, expireTimeStr)
	delta := time.Until(t)
	deltaHours := delta.Hours()

	pluginMessage := ""
	retCode := StateOk
	outStream := out.Output

	if deltaHours > 0 {
		pluginMessage = fmt.Sprintf("The token will expire the %s, in %s",
			expireTimeStr,
			delta.String())
		if deltaHours < CriticalExpirationHours {
			retCode = StateCritical
			outStream = out.Error
		} else if deltaHours < WarningExpirationHours {
			retCode = StateWarning
		}
	} else {
		pluginMessage = fmt.Sprintf("The token has expired!")
		retCode = StateCritical
		outStream = out.Error
	}

	outStream(pluginMessage)
	return retCode
}
