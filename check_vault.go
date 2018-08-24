package main

import (
	"flag"
	"fmt"
	"github.com/hashicorp/vault/api"
	"log"
	"os"
	"sort"
	"strings"
)

const (
	defaultVaultAddr string = "https://127.0.0.1:8200"
	defaultToken     string = ""
)

var client *api.Client // https://godoc.org/github.com/hashicorp/vault/api
var status bool
var policies string

func init() {
	const (
		statusMsg   = "Returns the Vault status (sealed/unsealed)"
		policiesMgs = "Comma-separated list of policies " +
			"to be checked for existance"
	)

	flag.BoolVar(&status, "status", false, statusMsg)
	flag.StringVar(&policies, "policies", "", policiesMgs)
}

func initClient(address string) (*api.Client, error) {
	client, err := api.NewClient(nil)
	if err != nil {
		return nil, err
	}

	err = client.SetAddress(address)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func checkSealStatus(address string) (bool, error) {
	client, err := initClient(address)
	if err != nil {
		return true, err
	}

	status, err := client.Sys().SealStatus()
	if err != nil {
		return true, err
	}

	return status.Sealed, nil
}

func stringInSlice(x string, data []string) bool {
	i := sort.Search(
		len(data), func(i int) bool { return data[i] >= x })
	return i < len(data) && data[i] == x
}

func checkForPolicies(address, token string, policies []string) (bool, error) {
	client, err := initClient(address)
	if err != nil {
		return false, err
	}

	if token != "" {
		client.SetToken(token)
	}

	installedPolicies, err := client.Sys().ListPolicies()
	if err != nil {
		return false, err
	}

	for _, policy := range policies {
		if stringInSlice(policy, installedPolicies) == false {
			return false, nil
		}
	}

	return true, nil
}

func main() {
	vaultAddress := defaultVaultAddr
	if envAddress := os.Getenv("VAULT_ADDR"); envAddress != "" {
		vaultAddress = envAddress
	}
	address := flag.String("address", vaultAddress,
		"The address of the Vault server. "+
			"Overrides the VAULT_ADDR environment variable if set.")

	vaultToken := defaultToken
	if envToken := os.Getenv("VAULT_TOKEN"); envToken != "" {
		vaultToken = envToken
	}
	token := flag.String("token", vaultToken,
		"The token to access Vault. "+
			"Overrides the VAULT_TOKEN environment variable if set.")

	flag.Parse()

	if status {
		sealStatus, err := checkSealStatus(*address)
		if err != nil {
			log.Fatal(err)
		}
		if sealStatus {
			fmt.Println("Error: Vault sealed")
		} else {
			fmt.Println("Vault unsealed")
		}
	} else if policies != "" {
		ok, err := checkForPolicies(
			*address, *token, strings.Split(policies, ","))
		if err != nil {
			log.Fatal(err)
		}
		if ok {
			fmt.Println("All the policies are available")
		} else {
			fmt.Fprintln(os.Stderr,
				"At least one policy is not installed")
		}
	} else {
		fmt.Fprintln(os.Stderr, "Syntax error")
		os.Exit(1)
	}
}
