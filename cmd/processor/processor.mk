PROCESSOR_SSM_ENABLED ?= false
PROCESSOR_SSM_PATH_PREFIX ?= /imageer/processor/local

##@ Processor

.PHONY: processor-run
processor-run: ## Run the processor
	@echo "🔄 Starting processor service..." && \
		PROCESSOR_SSM_ENABLED=$(PROCESSOR_SSM_ENABLED) \
		PROCESSOR_SSM_PATH_PREFIX=$(PROCESSOR_SSM_PATH_PREFIX) \
		go run ./cmd/processor/*.go -configs ./configs/processor && \
		echo "✅ Processor service stopped."

.PHONY: processor-build
processor-build: ## Build the processor
	@echo "📦 Building processor service..." && \
		go build -o ./bin/imageer-processor ./cmd/processor/*.go && \
		echo "✅ Processor service built."