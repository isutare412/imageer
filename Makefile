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

.PHONY: format
format: ## Format Go code
	@echo "🔄 Formatting Go code..." && \
		gofmt -w . && \
		echo "✅ Code formatting complete."

.PHONY: lint
lint: check-golangci-lint ## Run golangci-lint
	@echo "🔍 Running golangci-lint..." && \
		golangci-lint run && \
		echo "✅ Linting complete."

.PHONY: lint-fix
lint-fix: check-golangci-lint ## Run golangci-lint with auto-fix
	@echo "🔧 Running golangci-lint with auto-fix..." && \
		golangci-lint run --fix && \
		echo "✅ Linting with auto-fix complete."

.PHONY: test
test: ## Run Go tests
	@echo "🧪 Running Go tests..." && \
		go test ./... && \
		echo "✅ Tests complete."

##@ Dependencies

.PHONY: check-golangci-lint
check-golangci-lint: ## Check and install golangci-lint if needed
	@if ! command -v golangci-lint >/dev/null 2>&1; then \
		echo "⚠️ golangci-lint not found. Installing via Homebrew..." && \
		brew install golangci-lint && \
		echo "✅ golangci-lint installed successfully."; \
	else \
		echo "✅ golangci-lint is already installed."; \
	fi
