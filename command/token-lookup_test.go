/*
  Copyright 2019 Davide Madrisan <davide.madrisan@gmail.com>

  Licensed under the Mozilla Public License, Version 2.0 (the "License");
  you may not use this file except in compliance with the License.
  You may obtain a copy of the License at

      https://www.mozilla.org/en-US/MPL/2.0/

  Unless required by applicable law or agreed to in writing, software
  distributed under the License is distributed on an "AS IS" BASIS,
  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
  See the License for the specific language governing permissions and
  limitations under the License.

  Credit: This source code is based on the HashiCorp Vault testing code
          (but bugs are mine).
*/

package command

import (
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func testTokenLookupCommand(t *testing.T) (*cli.MockUi, *TokenLookupCommand) {
	ui := cli.NewMockUi()
	return ui, &TokenLookupCommand{
		BaseCommand: &BaseCommand{
			Ui: ui,
		},
	}
}

func TestTokenLookupCommand_Run(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name string
		args []string
		out  string
		code int
	}{
		{
			"help_message",
			[]string{"-help"},
			"Usage: hashicorp-vault-monitor token-lookup [options]",
			StateUndefined,
		},
		{
			"too_many_args",
			[]string{"arg1"},
			"Too many arguments",
			StateUndefined,
		},
		//{
		//	"token_expiration_ok",
		//	[]string{},
		//	// FIXME: we get "Cannot get the expire time of the Vault token"
		//	// because s.Data["expire_time"] == nil
		//	"FIXME",
		//	StateUndefined,
		//},
		{
			"nagios_too_many_args",
			[]string{"-output", "nagios", "arg1"},
			"Too many arguments",
			StateUndefined,
		},
	}

	t.Run("usage", func(t *testing.T) {
		t.Parallel()

		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				client, _, closer := testVaultServerUnseal(t)
				defer closer()

				token := client.Token()
				if token == "" {
					t.Errorf("cannot get the current Vault token")
				} else {
					ui, cmd := testTokenLookupCommand(t)
					cmd.client = client

					code := cmd.Run(tc.args)
					if code != tc.code {
						t.Errorf("expected %d to be %d", code, tc.code)
					}

					combined := ui.OutputWriter.String() + ui.ErrorWriter.String()
					if !strings.Contains(combined, tc.out) {
						t.Errorf("expected %q to contain %q", combined, tc.out)
					}
				}
			})
		}
	})

	t.Run("communication_failure", func(t *testing.T) {
		t.Parallel()

		client, closer := testVaultServerBad(t)
		defer closer()

		ui, cmd := testTokenLookupCommand(t)
		cmd.client = client

		code := cmd.Run([]string{})
		if exp := StateUndefined; code != exp {
			t.Errorf("expected %d to be %d", code, exp)
		}

		expected := "Error making API request."
		combined := ui.OutputWriter.String() + ui.ErrorWriter.String()
		if !strings.Contains(combined, expected) {
			t.Errorf("expected %q to contain %q", combined, expected)
		}
	})
}
