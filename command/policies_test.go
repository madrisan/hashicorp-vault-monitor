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

import "testing"

func TestContains(t *testing.T) {
	cases := []struct {
		values   []string
		key      string
		shouldbe bool
	}{
		{
			[]string{"first", "second", "third"},
			"first",
			true,
		},
		{
			[]string{"first", "second", "third"},
			"second",
			true,
		},
		{
			[]string{"first", "second", "third"},
			"third",
			true,
		},
		{
			[]string{"first", "second", "third"},
			"fourth",
			false,
		},
		{
			[]string{"first", "second", "third"},
			"",
			false,
		},
	}

	for _, tc := range cases {
		v := contains(tc.values, tc.key)
		if v != tc.shouldbe {
			t.Error("For", tc.values,
				"expected", tc.shouldbe, "got", v,
			)
		}
	}
}
