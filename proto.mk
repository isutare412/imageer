##@ Protocol Buffers

.PHONY: proto-format
proto-format: buf ## Format protobuf files
	@echo "🔧 Formatting protobuf files..." && \
		buf format -d --exit-code && \
		echo "✅ Protobuf formatting complete."

.PHONY: proto-lint
proto-lint: buf ## Lint protobuf files
	@echo "🔍 Linting protobuf files..." && \
		buf lint && \
		echo "✅ Protobuf linting complete."

.PHONY: proto-check
proto-check: proto-format proto-lint ## Run proto format, lint

.PHONY: proto-generate
proto-generate: buf ## Generate protobuf code
	@echo "🔄 Generating protobuf code..." && \
		buf generate && \
		echo "✅ Protobuf code generation complete."