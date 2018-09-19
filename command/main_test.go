package command

import (
	"encoding/base64"
	"testing"

	log "github.com/hashicorp/go-hclog"
	kv "github.com/hashicorp/vault-plugin-secrets-kv"
	"github.com/hashicorp/vault/api"
        "github.com/hashicorp/vault/audit"
        "github.com/hashicorp/vault/builtin/logical/pki"
        "github.com/hashicorp/vault/builtin/logical/ssh"
        "github.com/hashicorp/vault/builtin/logical/transit"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/vault"

	auditFile "github.com/hashicorp/vault/builtin/audit/file"
	credUserpass "github.com/hashicorp/vault/builtin/credential/userpass"
	vaulthttp "github.com/hashicorp/vault/http"
)

var (
	defaultVaultLogger             = log.NewNullLogger()
	defaultVaultCredentialBackends = map[string]logical.Factory{
		"userpass": credUserpass.Factory,
	}

	defaultVaultAuditBackends = map[string]audit.Factory{
		"file": auditFile.Factory,
	}

	defaultVaultLogicalBackends = map[string]logical.Factory{
		"generic-leased": vault.LeasedPassthroughBackendFactory,
		"pki":            pki.Factory,
		"ssh":            ssh.Factory,
		"transit":        transit.Factory,
		"kv":             kv.Factory,
	}
)

// testVaultServerUnseal creates a test vault cluster and returns a configured
// API client, list of unseal keys (as strings), and a closer function.
func testVaultServerUnseal(tb testing.TB) (*api.Client, []string, func()) {
	tb.Helper()

	return testVaultServerCoreConfig(tb, &vault.CoreConfig{
		DisableMlock:       true,
		DisableCache:       true,
		Logger:             defaultVaultLogger,
		CredentialBackends: defaultVaultCredentialBackends,
		AuditBackends:      defaultVaultAuditBackends,
		LogicalBackends:    defaultVaultLogicalBackends,
	})
}

// testVaultServerCoreConfig creates a new vault cluster with the given core
// configuration. This is a lower-level test helper.
func testVaultServerCoreConfig(tb testing.TB, coreConfig *vault.CoreConfig) (*api.Client, []string, func()) {
	tb.Helper()

	cluster := vault.NewTestCluster(tb, coreConfig, &vault.TestClusterOptions{
		HandlerFunc: vaulthttp.Handler,
		NumCores:    1, // Default is 3, but we don't need that many
	})
	cluster.Start()

	// Make it easy to get access to the active
	core := cluster.Cores[0].Core
	vault.TestWaitActive(tb, core)

	// Get the client already setup for us!
	client := cluster.Cores[0].Client
	client.SetToken(cluster.RootToken)

	// Convert the unseal keys to base64 encoded, since these are how the user
	// will get them.
	unsealKeys := make([]string, len(cluster.BarrierKeys))
	for i := range unsealKeys {
		unsealKeys[i] = base64.StdEncoding.EncodeToString(cluster.BarrierKeys[i])
	}

	return client, unsealKeys, func() { defer cluster.Cleanup() }
}
