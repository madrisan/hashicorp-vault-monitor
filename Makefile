EXTERNAL_TOOLS = \
	github.com/mitchellh/gox@latest \
	github.com/golangci/golangci-lint/cmd/golangci-lint@latest

CGO_ENABLED = 0
GOFMT_FILES ?= $$(find -name "*.go" -not -path "./vendor/*")
TEST ?= $$(go list ./... | grep -v /vendor/)
TESTARGS = -v

default: dev

# bin generates the releasable binaries for HashiCorp Vault Monitor.
bin: prep
	@CGO_ENABLED=$(CGO_ENABLED) BUILD_TAGS='$(BUILD_TAGS)' sh -c "'$(CURDIR)/scripts/build.sh'"

# dev creates binaries for testing HashiCorp Vault Monitor locally.
# These are put into ./bin/ as well as $GOPATH/bin.
dev: prep
	@CGO_ENABLED=$(CGO_ENABLED) BUILD_TAGS='$(BUILD_TAGS)' VAULT_MONITOR_DEV_BUILD=1 sh -c "'$(CURDIR)/scripts/build.sh'"

# bootstrap the build by downloading additional tools.
bootstrap:
	@for tool in  $(EXTERNAL_TOOLS) ; do \
		echo "Installing/Updating $$tool" ; \
		GO111MODULE=on go get -u $$tool; \
	done

cover:
	./scripts/coverage.sh --html

fmtcheck:
	./scripts/gofmtcheck.sh

fmt:
	gofmt -w $(GOFMT_FILES)

prep: fmtcheck

# test runs the unit tests and vets the code.
test: prep
	@CGO_ENABLED=$(CGO_ENABLED) \
	VAULT_ADDR= \
	VAULT_TOKEN= \
	go test -tags='$(BUILD_TAGS)' $(TEST) $(TESTARGS) -parallel=5

# lint runs vet plus a number of other checkers, it is more comprehensive, but louder.
lint:
	@go list -f '{{.Dir}}' ./... | grep -v /vendor/ \
		| xargs golangci-lint run; if [ $$? -eq 1 ]; then \
			echo ""; \
			echo "Lint found suspicious constructs. Please check the reported constructs"; \
			echo "and fix them if necessary before submitting the code for reviewal."; \
		fi

# vet runs the Go source code static analysis tool `vet` to find any common errors.
vet:
	@go list -f '{{.Dir}}' ./... | grep -v /vendor/ \
		| xargs go vet ; if [ $$? -eq 1 ]; then \
			echo ""; \
			echo "Vet found suspicious constructs. Please check the reported constructs"; \
			echo "and fix them if necessary before submitting the code for reviewal."; \
		fi
