.DEFAULT_GOAL := help

GO        := go
GOFLAGS   := -trimpath
GOBIN     := $(CURDIR)/bin

NODE_BIN  := node_modules/.bin
NODE_DEPS := node_modules/.package-lock.json

# Tool versions — bump these on upgrade
GOLINT_VERSION     := v2.13.2
GOVULN_VERSION     := v1.1.4
TFPLUGINDOCS_VERSION := v0.25.0

GOLINT        := $(GOBIN)/golangci-lint
GOVULN        := $(GOBIN)/govulncheck
TFPLUGINDOCS  := $(GOBIN)/tfplugindocs

# Provider
PLUGIN_DIR  := $(HOME)/.terraform.d/plugins/registry.terraform.io/sojournerdev/platctx/0.0.0
OS          := $(shell uname -s | tr A-Z a-z)
ARCH        := $(shell uname -m | sed 's/x86_64/amd64/' | sed 's/aarch64/arm64/')
PLUGIN_PATH := $(PLUGIN_DIR)/$(OS)_$(ARCH)

# Overridable: make build VERSION=1.2.3
VERSION   ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
LDFLAGS   := -s -w -X main.version=$(VERSION)

.PHONY: help tools build install test test-acc lint fmt vet vuln clean \
        check check-go check-docs \
        docs docs-generate docs-validate docs-lint docs-spell docs-links docs-fmt \
        dev

help:
	@printf '\n  build              Compile all packages\n'
	@printf '  install            Install provider to local plugins directory\n'
	@printf '  dev                Build and install for local development\n'
	@printf '  test               Run unit tests with race detector\n'
	@printf '  test-acc           Run acceptance tests\n'
	@printf '  lint               Run golangci-lint\n'
	@printf '  fmt                Check formatting (exits non-zero on diff)\n'
	@printf '  vet                Run go vet\n'
	@printf '  vuln               Scan for known vulnerabilities\n'
	@printf '  tools              Install all development dependencies\n'
	@printf '  clean              Remove build artifacts and caches\n'
	@printf '  check              Full CI suite (Go + docs)\n'
	@printf '  check-go           All Go quality gates\n'
	@printf '  check-docs         All documentation quality gates\n'
	@printf '  docs               Alias for check-docs\n'
	@printf '  docs-generate      Generate Registry documentation\n'
	@printf '  docs-validate      Validate generated documentation\n'
	@printf '  docs-lint          Lint Markdown structure and style\n'
	@printf '  docs-spell         Check documentation spelling\n'
	@printf '  docs-links         Check documentation links\n'
	@printf '  docs-fmt           Apply safe Markdown lint fixes\n'

tools: $(NODE_DEPS) $(GOLINT) $(GOVULN) $(TFPLUGINDOCS)

$(NODE_DEPS): package.json package-lock.json .node-version .npmrc
	npm ci --ignore-scripts --no-audit --no-fund

$(GOBIN):
	@mkdir -p $@

$(GOLINT): | $(GOBIN)
	GOBIN=$(GOBIN) $(GO) install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@$(GOLINT_VERSION)

$(GOVULN): | $(GOBIN)
	GOBIN=$(GOBIN) $(GO) install golang.org/x/vuln/cmd/govulncheck@$(GOVULN_VERSION)

$(TFPLUGINDOCS): | $(GOBIN)
	GOBIN=$(GOBIN) $(GO) install github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs@$(TFPLUGINDOCS_VERSION)

build:
	$(GO) build $(GOFLAGS) -ldflags '$(LDFLAGS)' ./...

install: build
	@rm -f $(PLUGIN_PATH)/terraform-provider-platctx
	@mkdir -p $(PLUGIN_PATH)
	$(GO) build $(GOFLAGS) -ldflags '$(LDFLAGS)' -o $(PLUGIN_PATH)/terraform-provider-platctx .

dev: install
	@mkdir -p $(CURDIR)/.terraform-dev
	@cp $(PLUGIN_PATH)/terraform-provider-platctx $(CURDIR)/.terraform-dev/
	@printf 'Provider installed to %s\n' $(PLUGIN_PATH)

test:
	$(GO) test $(GOFLAGS) -count=1 -race ./internal/...

test-acc:
	TF_ACC=1 $(GO) test $(GOFLAGS) -count=1 -v ./tests/...

lint: $(GOLINT)
	$(GOLINT) run ./...

fmt:
	@DIFF=$$($(GO) fmt ./... 2>&1); \
	if [ -n "$$DIFF" ]; then \
		printf 'Formatted:\n%s\n' "$$DIFF"; exit 1; \
	fi

vet:
	$(GO) vet ./...

vuln: $(GOVULN)
	$(GOVULN) ./...

clean:
	$(GO) clean -cache -testcache
	rm -rf $(GOBIN)
	rm -rf $(PLUGIN_PATH)

check: check-go check-docs
	@echo "\n✓ All checks passed"

check-go:
	@echo "\n=== Building ==="
	@$(MAKE) build
	@echo "\n=== Testing ==="
	@$(MAKE) test
	@echo "\n=== Vet ==="
	@$(MAKE) vet
	@echo "\n=== Format ==="
	@$(MAKE) fmt
	@echo "\n=== Lint ==="
	@$(MAKE) lint
	@echo "\n=== Vulnerabilities ==="
	@$(MAKE) vuln

check-docs:
	@echo "\n=== Docs Lint ==="
	@$(MAKE) docs-lint
	@echo "\n=== Docs Spell ==="
	@$(MAKE) docs-spell
	@echo "\n=== Docs Links ==="
	@$(MAKE) docs-links
	@echo "\n=== Docs Validate ==="
	@$(MAKE) docs-validate

docs: check-docs

docs-lint: $(NODE_DEPS)
	$(NODE_BIN)/markdownlint-cli2

docs-spell: $(NODE_DEPS)
	$(NODE_BIN)/cspell lint "**/*.md"

docs-links: $(NODE_DEPS)
	$(NODE_BIN)/markdown-link-check --quiet --config .markdown-link-check.json docs

docs-fmt: $(NODE_DEPS)
	$(NODE_BIN)/markdownlint-cli2 --fix

docs-generate: $(TFPLUGINDOCS)
	$(TFPLUGINDOCS) generate --provider-name platctx

docs-validate: $(TFPLUGINDOCS)
	$(TFPLUGINDOCS) validate --provider-name platctx
