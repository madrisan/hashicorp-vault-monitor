# Hashicorp Vault Seal Monitor

How to build the source code.

```
export GOPATH="$HOME/Development/go"

go get -u github.com/hashicorp/vault/api
go get -u github.com/madrisan/hashicorp-vault-seal-monitor
go install github.com/madrisan/hashicorp-vault-seal-monitor
```

Run the monitoring binary by entering the commands:

```
export VAULT_ADDR="http://127.0.0.1:8200"
$GOPATH/bin/hashicorp-vault-seal-monitor
```

or

```
$GOPATH/bin/hashicorp-vault-monitoring -address=http://127.0.0.1:8200
```
