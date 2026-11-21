.DEFAULT_GOAL := help

BIN := bin/ubc
PKG ?= ./...
TOOLS := $(CURDIR)/.tools
GOLANGCI_LINT_VERSION := v2.15.0
GOLANGCI_LINT := $(TOOLS)/golangci-lint
COVERAGE_MIN := 90
COVERAGE_PROFILE := coverage.out

.PHONY: help local-setup test test-pkg coverage format lint check-test-style check-architecture validate build

help: ## List available targets
	@grep -hE '^[a-zA-Z_-]+:.*## ' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*## "} {printf "  %-20s %s\n", $$1, $$2}'

local-setup: $(GOLANGCI_LINT) ## Install pinned tools and enable the versioned git hooks
	git config core.hooksPath scripts/hooks

$(GOLANGCI_LINT):
	GOBIN=$(TOOLS) go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION)

test: ## Run all tests
	go test ./...

test-pkg: ## Run tests of one package: make test-pkg PKG=./internal/domain/...
	go test $(PKG)

coverage: ## Run all tests and fail if coverage of internal/ is below the minimum
	go test -coverpkg=./internal/... -coverprofile=$(COVERAGE_PROFILE) ./...
	scripts/checks/coverage.sh $(COVERAGE_PROFILE) $(COVERAGE_MIN)

format: $(GOLANGCI_LINT) ## Format all Go code with gofumpt
	$(GOLANGCI_LINT) fmt ./...

lint: $(GOLANGCI_LINT) ## Check formatting, run go vet and golangci-lint
	$(GOLANGCI_LINT) fmt --diff ./...
	go vet ./...
	$(GOLANGCI_LINT) run ./...

check-test-style: ## Fail on t.Error*/t.Fatal* in tests (ADR 0009)
	scripts/checks/test-style.sh

check-architecture: ## Fail when an import breaks the hexagonal dependency rule (ADR 0010)
	scripts/checks/architecture.sh

validate: lint check-test-style check-architecture coverage ## Everything CI checks; run before every commit

build: ## Build a static binary in bin/
	CGO_ENABLED=0 go build -o $(BIN) ./cmd/ubc
