package main

import (
	"flag"
	"fmt"
	"github.com/hashicorp/vault/api"
	"log"
	"os"
)

const defaultVaultAddr string = "https://127.0.0.1:8200"

// see: https://godoc.org/github.com/hashicorp/vault/api
var client *api.Client

func checkSealStatus(addr string) (bool, error) {
	client, err := api.NewClient(nil)
	if err != nil {
		return true, err
	}

	err = client.SetAddress(addr)
	if err != nil {
		return true, err
	}

	status, err := client.Sys().SealStatus()
	if err != nil {
		return true, err
	}

	return status.Sealed, nil
}

func main() {
	vaultAddress := defaultVaultAddr

	if envAddress := os.Getenv("VAULT_ADDR"); envAddress != "" {
		vaultAddress = envAddress
	}
	addr := flag.String("address", vaultAddress,
		"The address of the Vault server. "+
			"Overrides the VAULT_ADDR environment variable if set.")

	flag.Parse()

	sealStatus, err := checkSealStatus(*addr)
	if err != nil {
		log.Fatal(err)
	}
	if sealStatus {
		fmt.Println("Vault sealed")
	} else {
		fmt.Println("Vault unsealed")
	}
}
