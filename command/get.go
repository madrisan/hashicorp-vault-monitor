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
)

const (
	getCommandDescr = "Retrieves data from the KV store"
	getFieldDescr   = "Print only the field with the given name"
)

// GetCommand is a CLI Command that holds the attributes of the command `readsecret`.
type GetCommand struct {
	*BaseCommand
	Field string
	Path  string
}

// Synopsis returns a short synopsis of the `get` command.
func (c *GetCommand) Synopsis() string {
	return getCommandDescr
}

// Help returns a long-form help text of the `get` command.
func (c *GetCommand) Help() string {
	helpText := `
Usage: hashicorp-vault-monitor get [options] -field FIELD KEY

  This command retrieves the value from Vault's key-value store at the given
  key name. If no key exists with that name, an error is returned.

    $ hashicorp-vault-monitor get -field foo secret/test

  Additional flags and more advanced use cases are detailed below.

    -address=<string>
       Address of the Vault server. The default is https://127.0.0.1:8200. This
       can also be specified via the VAULT_ADDR environment variable.

    -token=<string>
       Specify a token for authentication. This can also be specified via the
       VAULT_TOKEN environment variable.

    -output=<string>
       Specify an output format. Can be 'default' or 'nagios'.

  Mandatory Options:

    -field=<string>
       Print only the field with the given name.

  The exit code reflects the result of the read operation:

      - %d - the secret has been successfully read
      - %d - the secret cannot be found of read
      - %d - an error occurred

  For a full list of examples, please see the online documentation.
`
	return fmt.Sprintf(helpText,
		StateOk, StateCritical, StateUndefined)
}

// Run executes the `get` command with the given CLI instance and command-line arguments.
func (c *GetCommand) Run(args []string) int {
	cmdFlags := flag.NewFlagSet("get", flag.ContinueOnError)
	cmdFlags.Usage = func() { c.UI.Output(c.Help()) }
	cmdFlags.StringVar(&c.Address, "address", addressDefault, addressDescr)
	cmdFlags.StringVar(&c.Token, "token", tokenDefault, tokenDescr)
	cmdFlags.StringVar(&c.Field, "field", "", getFieldDescr)
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
	switch {
	case len(args) < 1:
		out.Undefined("Not enough arguments (expected 1, got %d)", len(args))
		return StateUndefined
	case len(args) > 1:
		out.Undefined("Too many arguments (expected 1, got %d)", len(args))
		return StateUndefined
	}

	c.Path = args[0]

	if c.Field == "" {
		out.Undefined("Missing '-field' flag or empty field set")
		return StateUndefined
	}

	client, err := c.Client()
	if err != nil {
		out.Undefined(err.Error())
		return StateUndefined
	}

	secret, err := client.Logical().Read(c.Path)
	if err != nil {
		out.Undefined("error reading %s: %s", c.Path, err)
		return StateUndefined
	}
	if secret == nil {
		out.Undefined("no data found at %s", c.Path)
		return StateUndefined
	}

	// secret.Data in KVv2 is an object of type map[string]interface{} with two entries:
	// - data -> map[foo:bar]
	// - metadata -> map[created_time:2018-08-31T15:36:31.894655728Z deletion_time: destroyed:false version:3]
	// secret.Data in KVv1 is a `map[string]interface{}` object.
	// - map[foo:bar]
	// See: https://godoc.org/github.com/hashicorp/vault/api#Secret
	if data, ok := secret.Data["data"]; ok && data != nil {
		val := data.(map[string]interface{})[c.Field]
		if val == nil {
			out.Undefined("field '%s' not present in secret '%s'", c.Field, c.Path)
			return StateUndefined
		}
		out.Output("found a value for the key %s: '%v'", c.Field, val)
		return StateOk
	} else if val, ok := secret.Data[c.Field]; ok && val != nil {
		out.Output("found value: '%v'", val)
		return StateOk
	}

	out.Critical("field '%s' not present in secret '%s': %s", c.Field, c.Path, strings.Join(secret.Warnings, " "))
	return StateCritical
}
