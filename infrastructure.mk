##@ Infrastructure

.PHONY: infra-up
infra-up: ## Start infrastructure services
	@echo "🔄 Starting infrastructure services..." && \
		docker compose -f ./compose.yaml up -d && \
		echo "✅ Infrastructure services are up."

.PHONY: infra-down
infra-down: ## Stop infrastructure services
	@echo "🔄 Stopping infrastructure services..." && \
		docker compose -f ./compose.yaml down && \
		echo "✅ Infrastructure services are stopped."
