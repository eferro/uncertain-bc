.DEFAULT_GOAL := help

BIN := bin/ubc
PKG ?= ./...

.PHONY: help local-setup test test-pkg format lint validate build

help: ## List available targets
	@grep -hE '^[a-zA-Z_-]+:.*## ' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*## "} {printf "  %-12s %s\n", $$1, $$2}'

local-setup: ## Enable the versioned git hooks
	git config core.hooksPath scripts/hooks

test: ## Run all tests
	go test ./...

test-pkg: ## Run tests of one package: make test-pkg PKG=./internal/domain/...
	go test $(PKG)

format: ## Format all Go code
	gofmt -w .

lint: ## Check formatting and run go vet
	@test -z "$$(gofmt -l .)" || { echo "Not formatted, run 'make format':"; gofmt -l .; exit 1; }
	go vet ./...

validate: lint test ## Everything CI checks; run before every commit

build: ## Build a static binary in bin/
	CGO_ENABLED=0 go build -o $(BIN) ./cmd/ubc
