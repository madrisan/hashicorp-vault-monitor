![Release Status](https://img.shields.io/badge/status-stable-brightgreen.svg)
[![License](https://img.shields.io/badge/License-MPL--2.0-blue.svg)](https://spdx.org/licenses/MPL-2.0.html)
[![Coverage](https://img.shields.io/badge/Go%20Coverage-67.1%25-green.svg?longCache=true&style=flat)](https://github.com/jpoles1/gopherbadger)
[![Go Report Card](https://goreportcard.com/badge/github.com/madrisan/hashicorp-vault-monitor)](https://goreportcard.com/report/github.com/madrisan/hashicorp-vault-monitor)
[![Build Status](https://travis-ci.org/madrisan/hashicorp-vault-monitor.svg?branch=master)](https://travis-ci.org/madrisan/hashicorp-vault-monitor)
[![GolangCI](https://img.shields.io/badge/GolangCI-A+-success)](https://github.com/madrisan/hashicorp-vault-monitor/)

# Hashicorp Vault Monitor [![GoDoc](https://godoc.org/github.com/madrisan/hashicorp-vault-monitor?status.png)](https://godoc.org/github.com/madrisan/hashicorp-vault-monitor)

![](images/HashiCorp-Vault-logo.png?raw=true "HashiCorp Vault")

## How to build the source code
```
[ "$GOPATH" ] || export GOPATH="$HOME/go"
go get -u github.com/madrisan/hashicorp-vault-monitor

export PATH="$PATH:$GOPATH/bin"
$GOPATH/bin/hashicorp-vault-monitor -version
```
Optionally if you want to compile this tool for all the supported operating systems:
```
make -C $GOPATH/src/github.com/madrisan/hashicorp-vault-monitor bootstrap dev
```
You'll find the compiled binaries in the folder `$GOPATH//src/github.com/madrisan/hashicorp-vault-monitor/pkg/`.

## How to use the hashicorp-vault-monitor tool

Export the Hashicorp Vault server url in the variable `VAULT_ADDR`
```
export VAULT_ADDR='https://myvaultserver.mydomain.com:8200'
```

If you do not have a running Vault server and you want to test this monitoring tool,
you can run a dockerized version of the latest version
(requires [Docker](https://www.docker.com/) or [Podman](https://podman.io/)).
Then run the export `VAULT_ADDR ...` command from the terminal output
(replace `docker` by `podman` if you use the latter and do not have the
podman Docker CLI emulation configured).
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


export VAULT_ADDR='http://0.0.0.0:8200'
```

We can now create (in a different terminal, if you run the dockerized version of the Vault server) a Vault *policy* that we'll use later in the examples (you can find the binary *vault* [here](https://www.vaultproject.io/downloads)):
```
cat > accessor_lookup_policy.hcl <<__END
path "auth/token/lookup-accessor" {
    capabilities = ["update", "sudo", "read", "list"]
}
__END

vault login
   # enter the root token (or an admin token with sufficient permissions)

vault policy write accessor-policy accessor_lookup_policy.hcl
```
and two extra (non-root) tokens:
```
vault token create -policy=accessor-policy -renewable -period=768h
       Key                  Value
        ---                  -----
        token                s.iJPhLRp25r9FRwg4vrxfd0I7
        token_accessor       NzHyqTGPITcSYMiA31goyEXh
        token_duration       768h
        token_renewable      true
        token_policies       ["default" "accessor-policy"]
        identity_policies    []
        policies             ["default" "accessor-policy"]

vault token create -policy=default -renewable -period=768h
       Key                  Value
        ---                  -----
        token                s.EFI8PMCZF1KInfCj1yyI7Rpy
        token_accessor       ljXiSqQDdSZBYthO7IsrFMD2
        token_duration       768h
        token_renewable      true
        token_policies       ["default"]
        identity_policies    []
        policies             ["default"]
```
Note that the policies applied to the tokens are different.
 
### Monitoring the status (unsealed/sealed)
```
$GOPATH/bin/hashicorp-vault-monitor status \
    -address=$VAULT_ADDR
```

Add the output modifier `-output=nagios` if this tool is intented to
be used with the Nagios monitoring.

```
$GOPATH/bin/hashicorp-vault-monitor status \
    -output=nagios -address=$VAULT_ADDR
```

##### Example of output

    # default output message
    Vault (vault-cluster-50531563) is unsealed
    
    # with the '-output=nagios' switch
    vault OK - Vault (vault-cluster-50531563) is unsealed

### Monitoring the HA Cluster Status
```
$GOPATH/bin/hashicorp-vault-monitor hastatus \
    -address=$VAULT_ADDR
```

Add `-output=nagios` as above if you monitor Vault with Nagios.

##### Example of output

    # default output message
    Vault HA (vault-cluster-50531563) is enabled, Standby Node (Active Node Address: https://192.168.1.8:8200)

    # error message displayed when the HA mode is not enabled
    Vault HA (vault-cluster-50531563) is not enabled
    
    # with the '-output=nagios' switch
    vault OK - Vault HA (vault-cluster-50531563) is enabled, Standby Node (Active Node Address: https://192.168.1.8:8200)
    vault CRITICAL - Vault HA (vault-cluster-50531563) is not enabled

### Monitoring the installed Vault policies
```
$GOPATH/bin/hashicorp-vault-monitor policies \
    -address $VAULT_ADDR -token "39d2c714-6dce-6d96-513f-4cb250bf7fe8" \
    root saltstack
```

Add the flag `-output=nagios` if you monitor Vault with Nagios.

### Monitoring the access to the Vault KV data store

#### Get a secret from Vault KV data store v1
```
$GOPATH/bin/hashicorp-vault-monitor get \
    -address $VAULT_ADDR -token "39d2c714-6dce-6d96-513f-4cb250bf7fe8" \
    -field foo secret/mysecret
```

#### Get a secret from Vault KV data store v2
```
$GOPATH/bin/hashicorp-vault-monitor get \
    -address $VAULT_ADDR -token "39d2c714-6dce-6d96-513f-4cb250bf7fe8" \
    -field foo secret/data/mysecret
```

The `-output=nagios` switch must be added as usual to make the output compliance with Nagios.

##### Example of output

    # default output message
    found a value for the key foo: 'this-is-a-secret'
    
    # with the '-output=nagios' switch
    vault OK - found a value for the key foo: 'this-is-a-secret-for-checking-vault'

### Monitoring the expiration date of a Vault token
```
$GOPATH/bin/hashicorp-vault-monitor token-lookup \
    -address=$VAULT_ADDR -token="s.EFI8PMCZF1KInfCj1yyI7Rpy" \
    -warning=120h -critical=72h
```
The `-warning` and `-critical` switches are optional and default to *168h* (7 days)
and *72h* (3 days) respectively.

As usual, add `-output=nagios` to get an output compliant with the Nagios specifications.

##### Example of output

    # default output message
    This (renewable) token will expire on Mon, 07 Oct 2019 14:25:06 UTC (4 weeks 3 days 23 hours 55 minutes 35 seconds left)
    
    # with the '-output=nagios' switch
    vault OK - This (renewable) token will expire on Mon, 07 Oct 2019 14:25:06 UTC (4 weeks 3 days 23 hours 55 minutes 35 seconds left)

### Monitoring the expiration date of a Vault token via its associated token accessor

To avoid exposing the tokens in your monitoring setup, you can make use of their associated *Token Accessors*.
```
$GOPATH/bin/hashicorp-vault-monitor token-lookup \
    -address=$VAULT_ADDR -token="s.iJPhLRp25r9FRwg4vrxfd0I7" \
    -token-accessor="ljXiSqQDdSZBYthO7IsrFMD2" \
    -warning=120h -critical=72h
```
The `-warning` and `-critical` switches are optional and default to *168h* (7 days)
and *72h* (3 days) respectively, as described above.

Add `-output=nagios` to get an output compliant with the Nagios specifications.

---

Note that you should replace `39d2c7...` with the generated *Root token* from
your output. The same for the values of the two other tokens used in the examples.

You can omit the `-address` and `-token` flags by setting the environment
variables `VAULT_ADDR` and `VAULT_TOKEN` as shown in the following example:
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
https://myvaultserver.mydomain.com:8200/ui
```
or, if you're using the dockerized Vault server:
```
http://127.0.0.1:8200/ui
```

![](images/HashiCorp-Vault-web-ui-login.png?raw=true "HashiCorp Vault Web UI Login")
###### Image 1. Screenshot of the web UI login page

![](images/HashiCorp-Vault-web-ui-homepage.png?raw=true "HashiCorp Vault Web UI Homepage")
###### Image 2. Screenshot of the web UI homepage

![](images/HashiCorp-Vault-web-ui-secrets.png?raw=true "HashiCorp Vault Web UI Secrets")
###### Image 3. Screenshot of the web UI secrets management page

## Developers' corner

Some extra actions that may be usefull to project developers.

#### Run the Test Suite

Just run in the top source folder (`$GOPATH/src/github.com/madrisan/hashicorp-vault-monitor`):
```
make test
```

#### Generate Test Coverage Statistics

Go to the top source folder and enter the command:
```
make cover
```

#### Style and Static Code Analyzers

##### Golint

Run the `golint`, the official linter for Go source code:
```
go get -u golang.org/x/lint/golint
  # the golint binary is now available in:
  #   go list -f {{.Target}} golang.org/x/lint/golint
cd $GOPATH/src/github.com/madrisan/hashicorp-vault-monitor

export PATH="$PATH:$GOPATH/bin"
golint -set_exit_status ./... | grep -v ^vendor
```

##### GolangCI-Lint

Run the `GolangCI-Lint` linters aggregator:
```
go get -u github.com/golangci/golangci-lint/cmd/golangci-lint
cd $GOPATH/src/github.com/madrisan/hashicorp-vault-monitor

export PATH="$PATH:$GOPATH/bin"
golangci-lint run ./...
```

or just execute (requires `hashicorp-vault-monitor` version > 0.8.2):
```
make -C $GOPATH/src/github.com/madrisan/hashicorp-vault-monitor lint
```

##### Go Vet

Run the Go source code static analysis tool `vet` to find any common errors.
```
make -C $GOPATH/src/github.com/madrisan/hashicorp-vault-monitor vet
```

This command is available with `hashicorp-vault-monitor` version > 0.8.2.
