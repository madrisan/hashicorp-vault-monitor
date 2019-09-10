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

// Thresholds in hours for warning and critical status.
const (
	DefaultWarningTokenExpiration  = "168h"
	DefaultCriticalTokenExpiration = "72h"
)

// TokenLookupCommand is a CLI Command that holds the attributes of the command `token-lookup`.
type TokenLookupCommand struct {
	*BaseCommand
	WarningThreshold  string
	CriticalThreshold string
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

    -warning=<string>
       Warning threshold in days (default: %s).

    -critical=<string>
       Critical threshold in days (default: %s).

  The exit code reflects the token expiration time:

      - %d - the token is usable
      - %d - the token will expire in less than the warning threshold
      - %d - the token will expire in less than the critical threshold
      - %d - an error occurred

  For a full list of examples, please see the online documentation.
`
	return fmt.Sprintf(helpText,
		DefaultWarningTokenExpiration,
		DefaultCriticalTokenExpiration,
		StateOk, StateWarning, StateCritical, StateUndefined)
}

// Parse the warning and critical thresholds and return their corresponding Duration
func (c *TokenLookupCommand) GetThresholds() (time.Duration, time.Duration, error) {
	warningThreshold, err := time.ParseDuration(c.WarningThreshold)
	if err != nil {
		return 0, 0, err
	}
	criticalThreshold, err := time.ParseDuration(c.CriticalThreshold)
	if err != nil {
		return 0, 0, err
	}
	return warningThreshold, criticalThreshold, nil
}

// Run executes the `token-lookup` command with the given CLI instance and command-line arguments.
func (c *TokenLookupCommand) Run(args []string) int {
	cmdFlags := flag.NewFlagSet("token-lookup", flag.ContinueOnError)
	cmdFlags.Usage = func() { c.Ui.Output(c.Help()) }
	cmdFlags.StringVar(&c.Address, "address", addressDefault, addressDescr)
	cmdFlags.StringVar(&c.OutputFormat, "output", "default", outputFormatDescr)
	cmdFlags.StringVar(&c.WarningThreshold, "warning",
		DefaultWarningTokenExpiration,
		fmt.Sprintf(warningDescr, DefaultWarningTokenExpiration))
	cmdFlags.StringVar(&c.CriticalThreshold, "critical",
		DefaultCriticalTokenExpiration,
		fmt.Sprintf(criticalDescr, DefaultCriticalTokenExpiration))

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

	warningThreshold, criticalThreshold, err := c.GetThresholds()
	if err != nil {
		out.Undefined(err.Error())
		return StateUndefined
	}

	client, err := c.Client()
	if err != nil {
		out.Undefined(err.Error())
		return StateUndefined
	}

	ta := client.Auth().Token()

	s, err := ta.LookupSelf()
	if err != nil {
		out.Undefined(err.Error())
		return StateUndefined
	}

	if s.Data == nil || s.Data["expire_time"] == nil {
		out.Undefined("Cannot get the expire time of the Vault token")
		return StateUndefined
	}

	expireTimeRaw := s.Data["expire_time"]
	expireTimeStr, ok := expireTimeRaw.(string)
	if !ok {
		out.Undefined("Could not convert expire_time to a string")
		return StateUndefined
	}

	t, _ := time.Parse(time.RFC3339Nano, expireTimeStr)
	delta := time.Until(t)

	pluginMessage := ""
	retCode := StateOk

	if delta > 0 {
		pluginMessage = fmt.Sprintf("The token will expire on %s (%s left)",
			t.Format(time.RFC1123),
			delta.Truncate(time.Second).String())
		if delta < criticalThreshold {
			out.Critical(pluginMessage)
			retCode = StateCritical
		} else if delta < warningThreshold {
			out.Warning(pluginMessage)
			retCode = StateWarning
		} else {
			out.Output(pluginMessage)
		}

	} else {
		out.Critical("The token has expired!")
		retCode = StateCritical
	}

	return retCode
}
