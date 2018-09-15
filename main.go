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

package main // import "github.com/madrisan/hashicorp-vault-monitor"

import (
	"fmt"
	"os"
	"strings"

	"github.com/madrisan/hashicorp-vault-monitor/command"
	"github.com/madrisan/hashicorp-vault-monitor/vault"
	"github.com/madrisan/hashicorp-vault-monitor/version"
)

const (
	StateOk byte = iota
	StateWarning
	StateCritical
	StateUnknown
)

type oracle struct {
	message string
	status  byte
}

func (o oracle) String() string {
	var status string

	switch o.status {
	case StateOk:
		status = "Ok"
	case StateWarning:
		status = "Warning"
	case StateCritical:
		status = "Critical"
	default:
		status = "Unknown"
	}

	return fmt.Sprintf(status + ": " + o.message)
}

// Version returns the semantic version (see http://semver.org) of the tool.
func Version() string {
	versionInfo := version.GetVersion()
	return versionInfo.FullVersionNumber(true)
}

func main() {
	var result oracle
	opt := command.Run()
	if opt.Status {
		unsealed, err := vault.IsUnsealed(opt.Address)
		if err != nil {
			result = oracle{
				message: err.Error(),
				status:  StateUnknown,
			}
		}
		if unsealed {
			result = oracle{
				message: "Vault is unsealed",
			}
		} else {
			result = oracle{
				message: "Vault is sealed",
				status:  StateCritical,
			}
		}
	} else if opt.Policies != "" {
		policies, err := vault.CheckPolicies(
			opt.Address, opt.Token, strings.Split(opt.Policies, ","))
		if err != nil {
			result.message = err.Error()
			if policies != nil {
				result.status = StateCritical
			} else {
				result.status = StateUnknown
			}
		} else {
			result = oracle{
				message: "all the Vault policies are available",
			}
		}
	} else if opt.ReadKey != "" {
		secret, err := vault.ReadSecret(opt.ReadKey, opt.Address, opt.Token)
		if err != nil {
			result = oracle{
				message: err.Error(),
				status:  StateCritical,
			}
		} else {
			result = oracle{
				message: "found value: '" + secret + "'",
			}
		}
	} else if opt.Infos {
	        fmt.Println(Version())
		os.Exit(0)
	} else {
		fmt.Fprintln(os.Stderr,
			"Syntax error: missing -readkey, -status, or -policies flag")
		os.Exit(1)
	}

	fmt.Print(result, "\n")
	os.Exit(int(result.status))
}
