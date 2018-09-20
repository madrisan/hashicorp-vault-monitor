EXTERNAL_TOOLS = \
	github.com/mitchellh/gox

CGO_ENABLED = 0
GOFMT_FILES ?= $$(find -name "*.go" -not -path "./vendor/*")
TEST ?= $$(go list ./... | grep -v /vendor/)
TESTARGS = -v

default: dev

# bin generates the releasable binaries for HashiCorp Vault Monitor
bin:
	@CGO_ENABLED=$(CGO_ENABLED) BUILD_TAGS='$(BUILD_TAGS)' sh -c "'$(CURDIR)/scripts/build.sh'"

# dev creates binaries for testing HashiCorp Vault Monitor locally.
# These are put into ./bin/ as well as $GOPATH/bin
dev:
	@CGO_ENABLED=$(CGO_ENABLED) BUILD_TAGS='$(BUILD_TAGS)' VAULT_MONITOR_DEV_BUILD=1 sh -c "'$(CURDIR)/scripts/build.sh'"

# bootstrap the build by downloading additional tools
bootstrap:
	@for tool in  $(EXTERNAL_TOOLS) ; do \
		echo "Installing/Updating $$tool" ; \
		go get -u $$tool; \
	done

fmtcheck:
	./scripts/gofmtcheck.sh

fmt:
	gofmt -w $(GOFMT_FILES)

# test runs the unit tests and vets the code
test:
	@CGO_ENABLED=$(CGO_ENABLED) \
	VAULT_ADDR= \
	VAULT_TOKEN= \
	go test -tags='$(BUILD_TAGS)' $(TEST) $(TESTARGS) -parallel=5
