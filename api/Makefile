.DEFAULT_GOAL := help

ROOTDIR:=$(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))))

.PHONY: swag
swag: ## Install swag binary
	@go list -f '{{ join .Imports "\n" }}' tools/tools.go | xargs -t -L1 go install
	@echo Installed swag successfully

.PHONY: api
docs: ## Generate OpenAPI specification
	swag init \
		-d "${ROOTDIR}/cmd,${ROOTDIR}/internal/adapter/http" \
		-o "${ROOTDIR}/api"
	@echo Generated docs successfully

.PHONY: help
help: ## Print this help message
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'