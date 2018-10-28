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

import "testing"

func TestVaultUtil(t *testing.T) {
	t.Parallel()

	t.Run("sanitize_path", func(t *testing.T) {
		t.Parallel()

		cases := []struct {
			name     string
			value    string
			shouldbe string
		}{
			{
				"empty_string",
				"",
				"",
			},
			{
				"leading_and_trailing_slash",
				"/secret/data/mysecret/",
				"secret/data/mysecret",
			},
			{
				"trailing_slash",
				"secret/data/mysecret/",
				"secret/data/mysecret",
			},
			{
				"leading_slash",
				"/secret/data/mysecret",
				"secret/data/mysecret",
			},
			{
				"no_leading_trailing_slashes",
				"secret/data/mysecret",
				"secret/data/mysecret",
			},
		}

		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				v := SanitizePath(tc.value)
				if v != tc.shouldbe {
					t.Error("For", tc.value,
						"expected", tc.shouldbe, "got", v,
					)
				}
			})
		}
	})

	t.Run("ensure_no_leading_slash", func(t *testing.T) {
		t.Parallel()

		cases := []struct {
			name     string
			value    string
			shouldbe string
		}{
			{
				"empty_string",
				"",
				"",
			},
			{
				"leading_slash",
				"/secret/data/mysecret",
				"secret/data/mysecret",
			},
			{
				"no_leading_slash",
				"secret/data/mysecret",
				"secret/data/mysecret",
			},
		}

		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				v := ensureNoLeadingSlash(tc.value)
				if v != tc.shouldbe {
					t.Error("For", tc.value,
						"expected", tc.shouldbe, "got", v,
					)
				}
			})
		}
	})

	t.Run("ensure_no_trailing_slash", func(t *testing.T) {
		t.Parallel()

		cases := []struct {
			name     string
			value    string
			shouldbe string
		}{
			{
				"empty_string",
				"",
				"",
			},
			{
				"trailing_slash",
				"secret/data/mysecret/",
				"secret/data/mysecret",
			},
			{
				"no_trailing_slash",
				"secret/data/mysecret",
				"secret/data/mysecret",
			},
		}

		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				v := ensureNoTrailingSlash(tc.value)
				if v != tc.shouldbe {
					t.Error("For", tc.value,
						"expected", tc.shouldbe, "got", v,
					)
				}
			})
		}
	})
}
