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
	"fmt"
	"path/filepath"

	"github.com/hashicorp/vault/api"
)

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

func ReadSecret(keypath, address, token string) (string, error) {
	client, err := ClientInit(address)
	if err != nil {
		return "", err
	}

	if token != "" {
		client.SetToken(token)
	}

	// see: https://godoc.org/github.com/hashicorp/vault/api#Secret
	path, key := filepath.Split(keypath)
	secret, err := client.Logical().Read(path)
	if err != nil {
		return "", fmt.Errorf("error reading %s: %s", path, err)
	}
	if secret == nil {
		return "", fmt.Errorf("no value found at %s", path)
	}

	// secret.Data:
	// - data -> map[akey:this-is-a-test]
	// - metadata -> map[created_time:2018-08-31T15:36:31.894655728Z deletion_time: destroyed:false version:3]
	if data, ok := secret.Data["data"]; ok && data != nil {
		value, err := getRawField(data, key)
		return value, err
	}

	return "", fmt.Errorf("no data found at %s", path)
}
