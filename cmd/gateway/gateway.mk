GATEWAY_SSM_ENABLED ?= false
GATEWAY_SSM_PATH_PREFIX ?= /imageer/gateway/local

##@ Gateway

.PHONY: gateway-run
gateway-run: ## Run the gateway
	@echo "🔄 Starting gateway service..." && \
		GATEWAY_SSM_ENABLED=$(GATEWAY_SSM_ENABLED) \
		GATEWAY_SSM_PATH_PREFIX=$(GATEWAY_SSM_PATH_PREFIX) \
		go run ./cmd/gateway/*.go -configs ./configs/gateway && \
		echo "✅ Gateway service stopped."

.PHONY: gateway-build
gateway-build: ## Build the gateway
	@echo "📦 Building gateway service..." && \
		go build -o ./bin/imageer-gateway ./cmd/gateway/*.go && \
		echo "✅ Gateway service built."
