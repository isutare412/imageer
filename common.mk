##@ General

.PHONY: help
help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Development

.PHONY: generate
generate: ## Generate code
	@echo "🔄 Generating code..." && \
		go generate ./... && \
		echo "✅ Code generation complete."

.PHONY: check
check: format lint test ## Run format, lint, and tests

CHECK_IMAGE ?= imageer-check

.PHONY: check-docker
check-docker: ## Run check in a container with libvips (for hosts without libvips)
	@docker build -q -t $(CHECK_IMAGE) - < deployment/check.Dockerfile >/dev/null && \
	docker run --rm --user $$(id -u):$$(id -g) \
		-v $(CURDIR):/src -w /src \
		-v imageer-check-cache:/cache \
		-e HOME=/cache/home -e GOMODCACHE=/cache/mod -e GOCACHE=/cache/build \
		-e GOLANGCI_LINT_CACHE=/cache/lint -e GOFLAGS=-buildvcs=false \
		$(CHECK_IMAGE) make check

.PHONY: test
test: ## Run Go tests
	@echo "🧪 Running Go tests..." && \
		go test ./... && \
		echo "✅ Tests complete."

.PHONY: format
format: ## Check Go code formatting
	@echo "🔍 Checking Go code formatting..." && \
	OUTPUT=$$(gofmt -d $$(find . -name '*.go' -not -path './vendor/*')); \
	if [ "$$OUTPUT" ]; then \
		echo "$$OUTPUT"; \
		echo "❌ Code is not properly formatted. Run 'go fmt ./...' to fix."; \
		exit 1; \
	else \
		echo "✅ All Go files are properly formatted."; \
	fi

.PHONY: lint
lint: golangci-lint ## Run golangci-lint
	@echo "🔍 Running golangci-lint..." && \
		$(GOLANGCI_LINT) run && \
		echo "✅ Linting complete."

##@ Dependencies

# Location to install dependencies to.
LOCALBIN ?= $(shell pwd)/bin
$(LOCALBIN):
	mkdir -p $(LOCALBIN)

# Tool binaries
GOLANGCI_LINT ?= $(LOCALBIN)/golangci-lint

.PHONY: golangci-lint
golangci-lint: ## Check and install golangci-lint if needed
	@if [ ! -f "$(GOLANGCI_LINT)" ]; then \
		echo "⚠️ golangci-lint not found in $$LOCALBIN. Installing..."; \
		curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/HEAD/install.sh | sh -s -- -b $(LOCALBIN) v2.9.0 && \
		echo "✅ golangci-lint installed successfully."; \
	else \
		echo "✅ golangci-lint is already installed."; \
	fi

.PHONY: vips
vips: ## Check and install vips if needed
	@if ! command -v vips >/dev/null 2>&1; then \
		echo "⚠️ vips not found. Please install vips before proceeding."; \
		brew install vips && \
		echo "✅ vips installed successfully."; \
	else \
		echo "✅ vips is already installed."; \
	fi

.PHONY: buf
buf: ## Check and install buf if needed
	@if ! command -v buf >/dev/null 2>&1; then \
		echo "⚠️ buf not found. Please install buf before proceeding."; \
		brew install bufbuild/buf/buf && \
		echo "✅ buf installed successfully."; \
	else \
		echo "✅ buf is already installed."; \
	fi
