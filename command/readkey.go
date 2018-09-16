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
	"path/filepath"
	"strings"

	"github.com/hashicorp/vault/api"
	"github.com/madrisan/hashicorp-vault-monitor/vault"
	"github.com/mitchellh/cli"
)

const (
	readKeyCommandDescr = "Read a Vault secret"
)

// ReadKeyCommand is a CLI Command that holds the attributes of the command `readkey`.
type ReadKeyCommand struct {
	Address string
	Token   string
	KeyPath string
	Ui      cli.Ui
}

// Synopsis returns a short synopsis of the `readkey` command.
func (c *ReadKeyCommand) Synopsis() string {
	return "Try to read a secret stored in Vault"
}

// Help returns a long-form help text of the `readkey` command.
func (c *ReadKeyCommand) Help() string {
	helpText := `
Usage: hashicorp-vault-monitor readkey [options]

  This command check if it is possible to read a secret stored in Vault.

    $ hashicorp-vault-monitor readkey \
        --path secret/data/test/testkey \
        --address https://127.0.0.1:8200 --token "12e2bf2b-3b82-9eff-07e4-8c7ad97715a9"

  The exit code reflects the seal status:

      - 0 - the secret has been successfully read
      - 1 - error
      - 2 - the secret cannot be found of read

  For a full list of examples, please see the documentation.

`
	return strings.TrimSpace(helpText)
}

// Run executes the `readkey` command with the given CLI instance and command-line arguments.
func (c *ReadKeyCommand) Run(args []string) int {
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

	cmdFlags := flag.NewFlagSet("readkey", flag.ContinueOnError)
	cmdFlags.Usage = func() { c.Ui.Output(c.Help()) }
	cmdFlags.StringVar(&c.Address, "address", addressDefault, addressDescr)
	cmdFlags.StringVar(&c.Token, "token", tokenDefault, tokenDescr)
	cmdFlags.StringVar(&c.KeyPath, "path", "", policiesDescr)

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

	if c.Token != "" {
		client.SetToken(c.Token)
	}

	// see: https://godoc.org/github.com/hashicorp/vault/api#Secret
	path, key := filepath.Split(c.KeyPath)
	secret, err := client.Logical().Read(path)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("error reading %s: %s", c.KeyPath, err))
		return 2
	}
	if secret == nil {
		c.Ui.Error(fmt.Sprintf("no value found at %s", c.KeyPath))
		return 2
	}

	// secret.Data:
	// - data -> map[akey:this-is-a-test]
	// - metadata -> map[created_time:2018-08-31T15:36:31.894655728Z deletion_time: destroyed:false version:3]
	if data, ok := secret.Data["data"]; ok && data != nil {
		value, err := getRawField(data, key)
		if err != nil {
			c.Ui.Error(err.Error())
			return 2
		}

		c.Ui.Output(fmt.Sprintf("found value: '%s'", value))
		return 0
	}

	c.Ui.Error(fmt.Sprintf("no value found at %s", c.KeyPath))
	return 1
}

func getRawField(data interface{}, field string) (string, error) {
	var val interface{}
	switch data.(type) {
	case *api.Secret:
		val = data.(*api.Secret).Data[field]
	case map[string]interface{}:
		val = data.(map[string]interface{})[field]
	}

	if val == nil {
		return "", fmt.Errorf("field '%s' not present in secret", field)
	}

	return val.(string), nil
}
