GO ?= go
GOCACHE ?= /tmp/go-build-cache

GOFMT_FAST_DIRS := protocol tests/harness
GOFMT_ALL_DIRS := client_2finance protocol wallet_manager tests/harness tests/e2e

.PHONY: help fmt fmt-check fmt-check-all vet test-fast test-harness test-unit test-e2e-lifecycle mcp-e2e-check

help:
	@echo "Targets:"
	@echo "  fmt                 Format Go files"
	@echo "  fmt-check           Fail if fast-feedback Go files are not gofmt-formatted"
	@echo "  fmt-check-all       Fail if any Go file in main test/source dirs is not formatted"
	@echo "  vet                 Run go vet on fast packages"
	@echo "  test-fast           Run fast feedback tests"
	@echo "  test-harness        Validate specs and harness checks"
	@echo "  test-unit           Run protocol/wallet tests"
	@echo "  test-e2e-lifecycle  Run live lifecycle e2e tests"
	@echo "  mcp-e2e-check       Run MCP-backed lifecycle e2e"

fmt:
	gofmt -w $(GOFMT_ALL_DIRS)

fmt-check:
	@test -z "$$(gofmt -l $(GOFMT_FAST_DIRS))"

fmt-check-all:
	@test -z "$$(gofmt -l $(GOFMT_ALL_DIRS))"

vet:
	GOCACHE=$(GOCACHE) $(GO) vet ./client_2finance ./protocol ./wallet_manager ./tests/harness

test-fast: fmt-check vet test-harness test-unit

test-harness:
	GOCACHE=$(GOCACHE) $(GO) test ./tests/harness -v

test-unit:
	GOCACHE=$(GOCACHE) $(GO) test ./protocol ./wallet_manager

test-e2e-lifecycle:
	GOCACHE=$(GOCACHE) $(GO) test ./tests/e2e -run 'TestLifecycle|TestSendingLifecycle|TestWalletManager_SignPreparedTransaction' -count=1 -v

mcp-e2e-check:
	MCP_E2E=true MCP_URL=$${MCP_URL:-http://127.0.0.1:8089/mcp} GOCACHE=$(GOCACHE) $(GO) test ./tests/e2e -run TestSendingLifecycle_MCPPrepareSignSubmitGet -count=1 -v
