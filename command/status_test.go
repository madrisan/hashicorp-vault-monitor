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

  Credit: This source code is based on the HashiCorp Vault testing code
          (but bugs are mine).
*/

package command

import (
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func testStatusCommand(t *testing.T) (*cli.MockUi, *StatusCommand) {
	ui := cli.NewMockUi()
	return ui, &StatusCommand{
		BaseCommand: &BaseCommand{
			Ui: ui,
		},
	}
}

func TestStatusCommand_Run(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name string
		args []string
		out  string
		code int
	}{
		{
			"too_many_args",
			[]string{"arg1"},
			"Too many arguments",
			StateUndefined,
		},
		{
			"unsealed",
			[]string{},
			"Vault is unsealed",
			StateOk,
		},
		{
			"nagios_too_many_args",
			[]string{"-output", "nagios", "arg1"},
			"Too many arguments",
			StateUndefined,
		},
		{
			"nagios_unsealed",
			[]string{"-output", "nagios"},
			"Vault is unsealed",
			StateOk,
		},
	}

	t.Run("status", func(t *testing.T) {
		t.Parallel()

		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				client, _, closer := testVaultServerUnseal(t)
				defer closer()

				ui, cmd := testStatusCommand(t)
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

		ui, cmd := testStatusCommand(t)
		cmd.client = client

		code := cmd.Run([]string{})
		if exp := StateUndefined; code != exp {
			t.Errorf("expected %d to be %d", code, exp)
		}

		expected := "error checking seal status: "
		combined := ui.OutputWriter.String() + ui.ErrorWriter.String()
		if !strings.Contains(combined, expected) {
			t.Errorf("expected %q to contain %q", combined, expected)
		}
	})
}
