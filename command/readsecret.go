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
	"errors"
	"flag"
	"fmt"
	"strings"

	"github.com/hashicorp/vault/api"
	"github.com/madrisan/hashicorp-vault-monitor/vault"
	"github.com/mitchellh/cli"
)

const (
	readSecretCommandDescr = "Read a Vault secret"
)

// ReadSecretCommand is a CLI Command that holds the attributes of the command `readsecret`.
type ReadSecretCommand struct {
	Address string
	Token   string
	KeyPath string
	Ui      cli.Ui
	client  *api.Client
}

// Synopsis returns a short synopsis of the `readsecret` command.
func (c *ReadSecretCommand) Synopsis() string {
	return "Try to read a secret stored in Vault"
}

// Help returns a long-form help text of the `readkey` command.
func (c *ReadSecretCommand) Help() string {
	helpText := `
Usage: hashicorp-vault-monitor readsecret [options]

  This command try to get a secret stored in a Vault server storage.

    $ hashicorp-vault-monitor readkey \
        --secret foo@secret/test \
        --address https://127.0.0.1:8200 --token "12e2bf2b-3b82-9eff-07e4-8c7ad97715a9"

  The exit code reflects the seal status:

      - 0 - the secret has been successfully read
      - 2 - the secret cannot be found of read
      - 3 - error

  For a full list of examples, please see the documentation.

`
	return strings.TrimSpace(helpText)
}

// Run executes the `readsecret` command with the given CLI instance and command-line arguments.
func (c *ReadSecretCommand) Run(args []string) int {
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

	cmdFlags := flag.NewFlagSet("readkey", flag.ContinueOnError)
	cmdFlags.Usage = func() { c.Ui.Output(c.Help()) }
	cmdFlags.StringVar(&c.Address, "address", addressDefault, addressDescr)
	cmdFlags.StringVar(&c.Token, "token", tokenDefault, tokenDescr)
	cmdFlags.StringVar(&c.KeyPath, "secret", "", policiesDescr)

	if err := cmdFlags.Parse(args); err != nil {
		c.Ui.Error(err.Error())
		return StateError
	}

	field, path, err := keyPathSplit(c.KeyPath)
	if err != nil {
		c.Ui.Error(err.Error())
		return StateError
	}

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

	secret, err := c.client.Logical().Read(path)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("error reading %s: %s", path, err))
		return StateError
	}
	if secret == nil {
		c.Ui.Error(fmt.Sprintf("no data found at %s", path))
		return StateCritical
	}

	// secret.Data in KVv2 is an object of type map[string]interface{} with two entries:
	// - data -> map[foo:bar]
	// - metadata -> map[created_time:2018-08-31T15:36:31.894655728Z deletion_time: destroyed:false version:3]
	// secret.Data in KVv1 is a `map[string]interface{}` object.
	// - map[foo:bar]
	// See: https://godoc.org/github.com/hashicorp/vault/api#Secret
	if data, ok := secret.Data["data"]; ok && data != nil {
		val := data.(map[string]interface{})[field]
		if val == nil {
			c.Ui.Error(fmt.Sprintf(
                               "field '%s' not present in secret", field))
			return StateCritical
		}
		c.Ui.Output(fmt.Sprintf("found value: '%v'", val))
		return StateOk
	} else if val, ok := secret.Data[field]; ok && val != nil {
		c.Ui.Output(fmt.Sprintf("found value: '%v'", val))
		return StateOk
	}
	c.Ui.Error(fmt.Sprintf("field '%s' not present in secret", field))
	return StateCritical
}

// keyPathSplit split a KeyPath string into the key and path parts
func keyPathSplit(keypath string) (string, string, error) {
	s := strings.Split(keypath, "@")
	if len(s) < 2 {
		return "" , "", errors.New("wrong format for a KeyPath argument")
	}
	return s[0], vault.SanitizePath(s[1]), nil
}
