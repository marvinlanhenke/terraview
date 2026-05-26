.DEFAULT_GOAL := help

GO ?= go
BINARY ?= terraview
CMD ?= ./cmd/$(BINARY)
OUT ?= bin/$(BINARY)

.PHONY: help fmt lint test build check

help: ## Show available targets.
	@awk 'BEGIN {FS = ":.*## "} /^[a-zA-Z_-]+:.*## / {printf "%-12s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

fmt: ## Format Go source files.
	$(GO) fmt ./...

lint: ## Run Go static analysis.
	$(GO) vet ./...

test: ## Run all Go tests.
	$(GO) test ./...

build: ## Build the terraview binary.
	@mkdir -p $(dir $(OUT))
	$(GO) build -o $(OUT) $(CMD)

check: fmt lint test ## Run the standard verification suite.
