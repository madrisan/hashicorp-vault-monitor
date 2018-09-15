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

// See: https://godoc.org/github.com/hashicorp/vault/api

package vault

import (
	"github.com/hashicorp/vault/api"
)

// ClientInit returs a new Vault client for the given configuration.
//
// If the configuration is nil, Vault will use configuration from
// api.DefaultConfig(), which is the recommended starting configuration.
//
// If the environment variable `VAULT_TOKEN` is present, the token will be
// automatically added to the client. Otherwise, you must manually call
// `api.SetToken()`.
func ClientInit(address string) (*api.Client, error) {
	client, err := api.NewClient(&api.Config{
		Address: address,
	})
	if err != nil {
		return nil, err
	}

	return client, nil
}
