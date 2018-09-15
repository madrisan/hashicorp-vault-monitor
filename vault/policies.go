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

package vault

import (
	"errors"
)

func contains(items []string, item string) bool {
	for _, i := range items {
		if i == item {
			return true
		}
	}
	return false
}

func CheckPolicies(
	address, token string, policies []string) ([]string, error) {

	client, err := ClientInit(address)
	if err != nil {
		return nil, err
	}

	if token != "" {
		client.SetToken(token)
	}

	activePolicies, err := client.Sys().ListPolicies()
	if err != nil {
		return nil, err
	}

	for _, policy := range policies {
		if !contains(activePolicies, policy) {
			return activePolicies,
				errors.New("no such Vault policy: " + policy)
		}
	}

	return activePolicies, nil
}
