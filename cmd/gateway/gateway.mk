##@ Gateway

.PHONY: gateway-run
gateway-run: ## Run the gateway
	@echo "🔄 Starting gateway service..." && \
		go run ./cmd/gateway/*.go && \
		echo "✅ Gateway service stopped."

.PHONY: gateway-build
gateway-build: ## Build the gateway
	@echo "📦 Building gateway service..." && \
		go build -o ./bin/imageer-gateway ./cmd/gateway/*.go && \
		echo "✅ Gateway service built."
