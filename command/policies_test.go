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
	"strings"
	"testing"

	"github.com/hashicorp/vault/api"
	"github.com/mitchellh/cli"
)

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

func testPoliciesCommand(t *testing.T, token string) (*cli.MockUi, *PoliciesCommand) {
	ui := cli.NewMockUi()
	return ui, &PoliciesCommand{
		BaseCommand: &BaseCommand{
			Token: token,
			Ui:    ui,
		},
	}
}

func TestPoliciesCommand_Run(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name string
		args []string
		out  string
		code int
	}{
		{
			"not_enough_args",
			[]string{},
			"Not enough arguments",
			StateUndefined,
		},
		{
			"default_policy",
			[]string{"default"},
			"all the policies are defined",
			StateOk,
		},
		{
			"default_policies",
			[]string{"default", "root"},
			"all the policies are defined",
			StateOk,
		},
		{
			"non_existent_policy",
			[]string{"nosuchpolicy"},
			"no such Vault policy: nosuchpolicy",
			StateCritical,
		},
		{
			"existing_and_non_existing_policies",
			[]string{"default", "nosuchpolicy"},
			"no such Vault policy: nosuchpolicy",
			StateCritical,
		},
	}

	t.Run("policies", func(t *testing.T) {
		t.Parallel()

		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				client, _, closer := testVaultServerUnseal(t)
				defer closer()

				secret, err := client.Auth().Token().Create(&api.TokenCreateRequest{
					Policies: []string{"policy"},
					TTL:      "30m",
				})
				if err != nil {
					t.Fatal(err)
				}
				if secret == nil || secret.Auth == nil || secret.Auth.ClientToken == "" {
					t.Fatalf("missing auth data: %#v", secret)
				}
				token := secret.Auth.ClientToken

				ui, cmd := testPoliciesCommand(t, token)
				cmd.client = client

				code := cmd.Run(tc.args)
				if code != tc.code {
					t.Errorf("expected %d to be %d", code, tc.code)
				}

				combined := ui.OutputWriter.String() + ui.ErrorWriter.String()
				if !strings.Contains(combined, tc.out) {
					t.Errorf("expected %q to contain %q", combined, tc.out)
				}
			})
		}
	})

	t.Run("communication_failure", func(t *testing.T) {
		t.Parallel()

		client, closer := testVaultServerBad(t)
		defer closer()

		ui, cmd := testPoliciesCommand(t, "")
		cmd.client = client

		code := cmd.Run([]string{"default"})
		if exp := StateUndefined; code != exp {
			t.Errorf("expected %d to be %d", code, exp)
		}

		expected := "error checking policies: "
		combined := ui.OutputWriter.String() + ui.ErrorWriter.String()
		if !strings.Contains(combined, expected) {
			t.Errorf("expected %q to contain %q", combined, expected)
		}
	})
}
