.DEFAULT_GOAL := help

GO ?= go
HUGO ?= hugo
NPM ?= npm
BINARY ?= terraview
CMD ?= ./cmd/$(BINARY)
OUT ?= bin/$(BINARY)
DOCS_DIR ?= docs/hugo-pages

.PHONY: help fmt fmt-check lint test build check docs-build docs-serve

help: ## Show available targets.
	@awk 'BEGIN {FS = ":.*## "} /^[a-zA-Z_-]+:.*## / {printf "%-12s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

fmt: ## Format Go source files.
	$(GO) fmt ./...

fmt-check: ## Verify Go source formatting.
	@test -z "$$(gofmt -l .)" || (gofmt -l . && exit 1)

lint: ## Run Go static analysis.
	$(GO) vet ./...

test: ## Run all Go tests.
	$(GO) test ./...

build: ## Build the terraview binary.
	@mkdir -p $(dir $(OUT))
	$(GO) build -o $(OUT) $(CMD)

check: fmt lint test ## Run the standard verification suite.

docs-build: ## Build the Hugo docs locally.
	$(NPM) --prefix $(DOCS_DIR) install
	$(GO) -C $(DOCS_DIR) mod download
	$(HUGO) --source $(DOCS_DIR)

docs-serve: ## Serve the Hugo docs locally.
	$(NPM) --prefix $(DOCS_DIR) install
	$(GO) -C $(DOCS_DIR) mod download
	$(HUGO) server --source $(DOCS_DIR) --disableFastRender
