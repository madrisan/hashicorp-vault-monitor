![Release Status](https://img.shields.io/badge/status-beta-yellow.svg)
[![License](https://img.shields.io/badge/License-MPL--2.0-blue.svg)](https://spdx.org/licenses/MPL-2.0.html)
[![Coverage](https://img.shields.io/badge/Go%20Coverage-63.6%25-green.svg?longCache=true&style=flat)](https://github.com/jpoles1/gopherbadger)
[![Go Report Card](https://goreportcard.com/badge/github.com/madrisan/hashicorp-vault-monitor)](https://goreportcard.com/report/github.com/madrisan/hashicorp-vault-monitor)

# Hashicorp Vault Monitor [![GoDoc](https://godoc.org/github.com/madrisan/hashicorp-vault-monitor?status.png)](https://godoc.org/github.com/madrisan/hashicorp-vault-monitor)

![](images/HashiCorp-Vault-logo.png?raw=true "HashiCorp Vault")

## How to build the source code

```
[ "$GOPATH" ] || export GOPATH="$HOME/go"
go get -u github.com/madrisan/hashicorp-vault-monitor

export PATH="$PATH:$GOPATH/bin"
make -C $GOPATH/src/github.com/madrisan/hashicorp-vault-monitor bootstrap dev
$GOPATH/bin/hashicorp-vault-monitor -version
```

#### Run the Test Suite (for developers)

Just run in the top source folder:
```
make test
```

#### Generate Test Coverage Statistics (for developers)

Go to the top source folder and enter the command:
```
make cover
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

#### Monitoring the status (unsealed/sealed)
```
$GOPATH/bin/hashicorp-vault-monitor status \
    -address=http://127.0.0.1:8200
```

Add the output modifier `-output=nagios` if this tool is intented to
be used with the Nagios monitoring.

```
$GOPATH/bin/hashicorp-vault-monitor status \
    -output=nagios -address=http://127.0.0.1:8200
```

###### Example of output

    # default output message
    Vault (vault-cluster-50531563) is unsealed
    
    # with the '-output=nagios' switch
    vault OK - Vault (vault-cluster-50531563) is unsealed

#### Monitoring the HA Cluster Status
```
$GOPATH/bin/hashicorp-vault-monitor hastatus \
    -address=http://127.0.0.1:8200
```

Add -output=nagios as above if you monitor Vault with Nagios.

###### Example of output

    # default output message
    Vault HA (vault-cluster-50531563) is enabled, Standby Node (Active Node Address: https://192.168.1.8:8200)
    
    # with the '-output=nagios' switch
    vault OK - Vault HA (vault-cluster-50531563) is enabled, Standby Node (Active Node Address: https://192.168.1.8:8200)

#### Monitoring the installed Vault policies
```
$GOPATH/bin/hashicorp-vault-monitor policies \
    -address http://127.0.0.1:8200 -token "39d2c714-6dce-6d96-513f-4cb250bf7fe8" \
    root saltstack
```

Add the flag `-output=nagios` if you monitor Vault with Nagios.

#### Monitoring the access to the Vault KV data store

##### Get a secret from Vault KV data store v1
```
$GOPATH/bin/hashicorp-vault-monitor get \
    -address http://127.0.0.1:8200 -token "39d2c714-6dce-6d96-513f-4cb250bf7fe8" \
    -field foo secret/mysecret
```

##### Get a secret from Vault KV data store v2
```
$GOPATH/bin/hashicorp-vault-monitor get \
    -address http://127.0.0.1:8200 -token "39d2c714-6dce-6d96-513f-4cb250bf7fe8" \
    -field foo secret/data/mysecret
```

The `-output=nagios` switch must be added as usual to make the output compliance with Nagios.

###### Example of output

    # default output message
    found value: 'this-is-a-secret-for-monitoring-vault'
    
    # with the '-output=nagios' switch
    vault OK - found value: 'this-is-a-secret-for-monitoring-vault'

Note that you should replace `39d2c7...` with the generated *Root token* from
your output.

You can omit the `-address` and `-token` flags by setting the environment
variables `VAULT_ADDR` and `VAULT_TOKEN`:
```
export VAULT_ADDR="http://127.0.0.1:8200"
export VAULT_TOKEN="39d2c714-6dce-6d96-513f-4cb250bf7fe8"

$GOPATH/bin/hashicorp-vault-monitor status
$GOPATH/bin/hashicorp-vault-monitor policies root saltstack
$GOPATH/bin/hashicorp-vault-monitor get -field foo secret/mysecret
$GOPATH/bin/hashicorp-vault-monitor get -field foo secret/data/mysecret
```

The *Root Token* can also be used to login to the Vault web interface at the
URL
```
http://127.0.0.1:8200/ui
```

![](images/HashiCorp-Vault-web-ui-login.png?raw=true "HashiCorp Vault Web UI Login")
###### Image 1. Screenshot of the web UI login page

![](images/HashiCorp-Vault-web-ui-homepage.png?raw=true "HashiCorp Vault Web UI Homepage")
###### Image 2. Screenshot of the web UI homepage

![](images/HashiCorp-Vault-web-ui-secrets.png?raw=true "HashiCorp Vault Web UI Secrets")
###### Image 3. Screenshot of the web UI secrets management page
