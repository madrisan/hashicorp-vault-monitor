[![License](https://img.shields.io/badge/License-MPL--2.0-blue.svg)](https://spdx.org/licenses/MPL-2.0.html)

# Hashicorp Vault Monitor

![](images/HashiCorp-Vault-logo.png?raw=true "HashiCorp Vault")

## How to build the source code

```
export GOPATH="$HOME/Development/go"

go get -u github.com/hashicorp/vault/api
go get -u github.com/madrisan/hashicorp-vault-monitor
go install github.com/madrisan/hashicorp-vault-monitor
```
## How to test the monitoring binary

If you do not have a running Vault server, you can run a dockerized version of
the latest version.
```
docker run -it -p 8200:8200 --cap-add=IPC_LOCK vault:latest

    Unable to find image 'vault:latest' locally
    latest: Pulling from library/vault
    8e3ba11ec2a2: Pull complete
    9d3c08966c5f: Pull complete
    8c2e0f2bce8e: Pull complete
    5752743f26bd: Pull complete
    fd7271f646fb: Pull complete
    Digest: sha256:85d4e6f0a52ba10d5f1d07c3f06aa64469c209237c04ed3fe1d5728f7c11fba6
    Status: Downloaded newer image for vault:latest
    ==> Vault server configuration:

                Api Address: http://0.0.0.0:8200
                        Cgo: disabled
            Cluster Address: https://0.0.0.0:8201
                Listener 1: tcp (addr: "0.0.0.0:8200", cluster address: "0.0.0.0:8201", max_request_duration: "999999h0m0s", max_request_size: "33554432", tls: "disabled")
                Log Level: info
                    Mlock: supported: true, enabled: false
                    Storage: inmem
                    Version: Vault v0.10.4
                Version Sha: e21712a687889de1125e0a12a980420b1a4f72d3

    WARNING! dev mode is enabled! In this mode, Vault runs entirely in-memory
    and starts unsealed with a single unseal key. The root token is already
    authenticated to the CLI, so you can immediately begin using Vault.

    You may need to set the following environment variable:

        $ export VAULT_ADDR='http://0.0.0.0:8200'

    The unseal key and root token are displayed below in case you want to
    seal/unseal the Vault or re-authenticate.

    Unseal Key: VJCtYcgcsmAUFaT70tZoS4uYliEz6XVbxRvcNvg/hqQ=
    Root Token: 39d2c714-6dce-6d96-513f-4cb250bf7fe8

    Development mode should NOT be used in production installations!
```

You can now run the monitoring binary by entering the commands:

```
$GOPATH/bin/hashicorp-vault-monitor \
    -address=http://127.0.0.1:8200 \
    -status
$GOPATH/bin/hashicorp-vault-monitor \
    -address=http://127.0.0.1:8200 \
    -token="39d2c714-6dce-6d96-513f-4cb250bf7fe8" \
    -policies="root,saltstack"
```

Note that you should replace `39d2c7...` with the generated *Root token* from
your output.

You can omit the `-address` and `-policies` flags by setting the environment
variables `VAULT_ADDR` and `VAULT_TOKEN`:
```
export VAULT_ADDR="http://127.0.0.1:8200"
export VAULT_TOKEN="39d2c714-6dce-6d96-513f-4cb250bf7fe8"

$GOPATH/bin/hashicorp-vault-monitor -status
$GOPATH/bin/hashicorp-vault-monitor -policies="root,saltstack"
```

The *Root Token* can also be used to login to the Vault web interface at the
URL
```
http://127.0.0.1:8200/ui
```
![](images/HashiCorp-Vault-web-ui.png?raw=true "HashiCorp Vault Web UI")
