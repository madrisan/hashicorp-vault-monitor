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
	"fmt"

	"github.com/hashicorp/vault/api"
	"github.com/mitchellh/cli"
)

// BaseCommand is a Command that holds the common command options
type BaseCommand struct {
	Address      string
	OutputFormat string
	Token        string
	UI           cli.Ui
	client       *api.Client
}

// Client returs a new HTTP API Vault client for the given configuration
// or the recommended default one, if no custom configuration were provided.
//
// If the environment variable `VAULT_TOKEN` is present, the token will be
// automatically added to the client.
func (c *BaseCommand) Client() (*api.Client, error) {
	// Read the test client if present
	if c.client != nil {
		return c.client, nil
	}

	// Create a default configuration.
	// Note that `api.DefaultConfig` execute `api.ReadEnvironment` and thus
	// loads all the Vault environment variables except `VAULT_TOKEN`
	config := api.DefaultConfig()
	if config == nil {
		return nil, fmt.Errorf("could not create/read default configuration for Vault")
	}
	if config.Error != nil {
		return nil, fmt.Errorf("error encountered setting up default configuration: %s",
			config.Error.Error())
	}

	if c.Address != "" {
		config.Address = c.Address
	}

	client, err := api.NewClient(config)
	if err != nil {
		return nil, err
	}

	if c.Token != "" {
		client.SetToken(c.Token)
	}

	c.client = client

	return client, nil
}
