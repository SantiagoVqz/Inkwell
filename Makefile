BINARY  := inkwell
CMD     := ./cmd/inkwell
BIN_DIR := bin

# VERSION is computed from `git describe` if available, else "dev".
# Override at invocation: make build VERSION=v0.1.0
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)

# LDFLAGS rewrites the value of the package-level `version` const in
# cmd/inkwell/main.go at link time. `-s -w` strip debug info to shrink the
# binary; remove them if you want delve to work cleanly.
LDFLAGS := -s -w -X main.version=$(VERSION)

GO ?= go

.PHONY: build run test vet fmt tidy generate clean help

build: ## Build the inkwell binary into ./bin/
	@mkdir -p $(BIN_DIR)
	$(GO) build -trimpath -ldflags="$(LDFLAGS)" -o $(BIN_DIR)/$(BINARY) $(CMD)

run: ## Run from source (pass args via ARGS="...")
	$(GO) run $(CMD) $(ARGS)

test: ## Run all tests
	$(GO) test ./...

vet: ## Static analysis
	$(GO) vet ./...

fmt: ## Format all Go files in place
	$(GO) fmt ./...

tidy: ## Add/remove module dependencies as imports change
	$(GO) mod tidy

generate: ## Regenerate sqlc code from migrations + query.sql
	$(GO) tool sqlc generate

clean: ## Remove build artifacts
	rm -rf $(BIN_DIR)

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-10s\033[0m %s\n", $$1, $$2}'
