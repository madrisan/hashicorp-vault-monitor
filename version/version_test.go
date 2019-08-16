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

package version

import "testing"

func TestVersion(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name            string
		verInfo         VersionInfo
		verShouldBe     string
		fullVerShouldBe string
	}{
		{
			"all_members",
			VersionInfo{"99ee29e", "1.2.3", "dev"},
			"1.2.3-dev",
			"HashiCorp Vault Monitor v1.2.3-dev (99ee29e)",
		},
		{
			"version_only",
			VersionInfo{"", "2.10", ""},
			"2.10",
			"HashiCorp Vault Monitor v2.10",
		},
		{
			"version_unknown",
			VersionInfo{"", "unknown", "unknown"},
			"(version unknown)",
			"HashiCorp Vault Monitor (version unknown)",
		},
	}

	t.Run("version_number", func(t *testing.T) {
		t.Parallel()

		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				v := tc.verInfo.VersionNumber()
				if v != tc.verShouldBe {
					t.Error("For", tc.verInfo,
						"expected", tc.verShouldBe, "got", v,
					)
				}
			})
		}
	})

	t.Run("full_version_number", func(t *testing.T) {
		t.Parallel()

		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				v := tc.verInfo.FullVersionNumber(true)
				if v != tc.fullVerShouldBe {
					t.Error("For", tc.verInfo,
						"expected", tc.fullVerShouldBe, "got", v,
					)
				}
			})

		}
	})
}
