.DEFAULT_GOAL := help

ROOTDIR:=$(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))))

.PHONY: swag
swag: ## Install swag binary
	@cat ${ROOTDIR}/tools/tools.go | grep _ | awk -F'"' '{print $$2}' | xargs -t go install
	@echo Installed swag successfully

.PHONY: api
docs: swag ## Generate OpenAPI specification
	swag init \
		-d "${ROOTDIR}/cmd,${ROOTDIR}/internal/adapter/http" \
		-o "${ROOTDIR}/api"
	@echo Generated docs successfully

.PHONY: help
help: ## Print this help message
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'