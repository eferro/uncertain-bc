.DEFAULT_GOAL := help

BIN := bin/ubc
PKG ?= ./...
TOOLS := $(CURDIR)/.tools
GOLANGCI_LINT_VERSION := v2.15.0
GOLANGCI_LINT := $(TOOLS)/golangci-lint
GITLEAKS_VERSION := v8.30.1
GITLEAKS := $(TOOLS)/gitleaks
LOG_OPTS ?=
COVERAGE_MIN := 90
VERSION ?= $(shell git describe --tags --always --dirty)
GORELEASER_VERSION := v2.20.3
GORELEASER := $(TOOLS)/goreleaser
COVERAGE_PROFILE := coverage.out
GREMLINS_VERSION := v0.6.0
GREMLINS := $(TOOLS)/gremlins
MUTATION_MIN := 85
MUTATION_MEMORY_KB := 4000000
DOCS_VENV := $(TOOLS)/docs-venv
MKDOCS := $(DOCS_VENV)/bin/mkdocs

.PHONY: help local-setup test test-pkg coverage format lint check-test-style check-architecture check-binary check-scripts secrets-scan test-mutation validate build release-check release-snapshot release-notes docs

help: ## List available targets
	@grep -hE '^[a-zA-Z_-]+:.*## ' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*## "} {printf "  %-20s %s\n", $$1, $$2}'

local-setup: $(GOLANGCI_LINT) $(GITLEAKS) ## Install pinned tools and enable the versioned git hooks
	git config core.hooksPath scripts/hooks

$(GOLANGCI_LINT):
	GOTOOLCHAIN=auto GOBIN=$(TOOLS) go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION)

$(GITLEAKS):
	GOTOOLCHAIN=auto GOBIN=$(TOOLS) go install github.com/zricethezav/gitleaks/v8@$(GITLEAKS_VERSION)

test: ## Run all tests
	go test ./...

test-pkg: ## Run tests of one package: make test-pkg PKG=./internal/domain/...
	go test $(PKG)

coverage: ## Run all tests and fail if coverage of internal/ is below the minimum
	go test -count=1 -coverpkg=./internal/... -coverprofile=$(COVERAGE_PROFILE) ./...
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

secrets-scan: $(GITLEAKS) ## Scan git history for secrets; narrow it with LOG_OPTS="A..B" (ADR 0019)
	$(GITLEAKS) git --no-banner --redact --log-opts="$(LOG_OPTS)" .

check-binary: ## Build the static binary with a known version and smoke-test it
	$(MAKE) --no-print-directory build VERSION=v0.0.0-smoke
	scripts/checks/binary-smoke.sh $(BIN) v0.0.0-smoke

check-scripts: ## Lint every shell script with shellcheck (runs in CI, where shellcheck is installed)
	shellcheck install.sh scripts/hooks/* scripts/checks/*.sh

validate: lint check-test-style check-architecture coverage check-binary ## Everything CI checks; run before every commit

build: ## Build a static binary in bin/
	CGO_ENABLED=0 go build -trimpath -ldflags "-s -w -X main.version=$(VERSION)" -o $(BIN) ./cmd/ubc

$(GREMLINS):
	GOTOOLCHAIN=auto GOBIN=$(TOOLS) go install github.com/go-gremlins/gremlins/cmd/gremlins@$(GREMLINS_VERSION)

test-mutation: $(GREMLINS) ## Mutation-test internal/domain and fail below the efficacy minimum (ADR 0012, 0020)
	ulimit -v $(MUTATION_MEMORY_KB) && scripts/checks/mutation.sh $(GREMLINS) $(MUTATION_MIN) ./internal/domain

$(GORELEASER):
	GOTOOLCHAIN=auto GOBIN=$(TOOLS) go install github.com/goreleaser/goreleaser/v2@$(GORELEASER_VERSION)

release-check: $(GORELEASER) ## Validate the goreleaser configuration
	$(GORELEASER) check

release-snapshot: $(GORELEASER) ## Build all release archives locally into dist/ without publishing
	$(GORELEASER) release --snapshot --clean

release-notes: ## Print the CHANGELOG.md section of a version, used as its release notes: make -s release-notes VERSION=v1.2.0 (ADR 0060)
	@go run ./tools/releasenotes CHANGELOG.md $(VERSION)

$(MKDOCS): tools/docsite/requirements.txt
	python3 -m venv $(DOCS_VENV)
	$(DOCS_VENV)/bin/pip install --quiet --requirement tools/docsite/requirements.txt
	touch $@

docs: $(MKDOCS) ## Build the documentation website into build/site; fails on broken links (ADR 0061)
	rm -rf build/docs
	go run ./tools/docsite mkdocs.yml
	$(MKDOCS) build --strict
