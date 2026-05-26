SHELL := /bin/sh

GO ?= go
PKG ?= ./...
DRIVER_PKG ?= ./pkg/driver
UNIT_PKGS ?= $(shell $(GO) list ./... | grep -v '/pkg/driver$$')
BIN_DIR ?= bin
CLI_MAIN ?= ./cmd/115driver
MCP_MAIN ?= ./mcp
CLI_BIN ?= $(BIN_DIR)/115driver
MCP_BIN ?= $(BIN_DIR)/115driver-mcp-server
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
CLI_LDFLAGS ?= -s -w -X github.com/SheltonZhu/115driver/cli/cmd.version=$(VERSION)

.DEFAULT_GOAL := help

.PHONY: all
all: build ## Build all binaries

.PHONY: help
help: ## Show available targets
	@awk 'BEGIN {FS = ":.*##"; printf "Usage:\n  make <target>\n\nTargets:\n"} /^[a-zA-Z0-9_.-]+:.*##/ { printf "  %-14s %s\n", $$1, $$2 }' $(MAKEFILE_LIST)

.PHONY: fmt
fmt: ## Format Go packages
	$(GO) fmt $(PKG)

.PHONY: tidy
tidy: ## Tidy Go modules
	$(GO) mod tidy

.PHONY: vet
vet: ## Run go vet
	$(GO) vet $(PKG)

.PHONY: test
test: ## Run local/unit Go tests
	$(GO) test $(UNIT_PKGS)

.PHONY: test-integration
test-integration: ## Run integration tests for pkg/driver (requires COOKIE)
	@test -n "$(COOKIE)" || { echo "COOKIE is required for integration tests"; exit 1; }
	COOKIE="$(COOKIE)" $(GO) test $(DRIVER_PKG)

.PHONY: test-all
test-all: test test-integration ## Run unit and integration tests

.PHONY: build-cli
build-cli: $(BIN_DIR) ## Build the CLI binary
	CGO_ENABLED=0 $(GO) build -trimpath -ldflags '$(CLI_LDFLAGS)' -o $(CLI_BIN) $(CLI_MAIN)

.PHONY: build-mcp
build-mcp: $(BIN_DIR) ## Build the MCP server binary
	CGO_ENABLED=0 $(GO) build -trimpath -o $(MCP_BIN) $(MCP_MAIN)

.PHONY: build
build: build-cli build-mcp ## Build all binaries into bin/

.PHONY: install-cli
install-cli: ## Install the CLI binary with go install
	CGO_ENABLED=0 $(GO) install -trimpath -ldflags '$(CLI_LDFLAGS)' $(CLI_MAIN)

.PHONY: check
check: vet test build ## Run the standard verification suite

.PHONY: pre-commit
pre-commit: ## Run all pre-commit hooks
	pre-commit run --all-files

.PHONY: clean
clean: ## Remove local build artifacts
	rm -rf $(BIN_DIR)

$(BIN_DIR):
	mkdir -p $(BIN_DIR)
