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
	"flag"
	"os"

	"github.com/hashicorp/vault/api"
)

const defaultVaultAddr = "https://127.0.0.1:8200"

// Command line switches
type RunOptions struct {
	Address  string
	Infos    bool
	Policies string
	ReadKey  string
	Status   bool
	Token    string
}

func Run() *RunOptions {
	return runCustom(nil)
}

func runCustom(runOpts *RunOptions) *RunOptions {
	if runOpts == nil {
		runOpts = &RunOptions{}
	}

	var envAddress string
	var envToken string

	// Parses the environment variables
	if v := os.Getenv(api.EnvVaultAddress); v != "" {
		envAddress = v
	}
	if v := os.Getenv(api.EnvVaultToken); v != "" {
		envToken = v
	}

	if envAddress != "" && runOpts.Address != "" {
		runOpts.Address = envAddress
	}
	if envToken != "" && runOpts.Token != "" {
		runOpts.Token = envToken
	}

	flag.StringVar(&runOpts.Address, "address", defaultVaultAddr,
		"The address of the Vault server. "+
			"Overrides the "+api.EnvVaultAddress+" environment variable if set")
	flag.BoolVar(&runOpts.Infos, "version", false,
		"Print the tool version number and exit")
	flag.BoolVar(&runOpts.Status, "status", false,
		"Returns the Vault status (sealed/unsealed)")
	flag.StringVar(&runOpts.Policies, "policies", "",
		"Comma-separated list of policies to be checked for existence")
	flag.StringVar(&runOpts.ReadKey, "readkey", "",
		"Read a Vault secret")
	flag.StringVar(&runOpts.Token, "token", "",
		"The token to access Vault. "+
			"Overrides the "+api.EnvVaultToken+" environment variable if set")

	flag.Parse()

	return runOpts
}
