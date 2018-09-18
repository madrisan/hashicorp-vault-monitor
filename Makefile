EXTERNAL_TOOLS=\
	github.com/mitchellh/gox

CGO_ENABLED=0

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

