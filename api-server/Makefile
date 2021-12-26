.DEFAULT_GOAL := help

ROOTDIR := $(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))))

IMAGEER_CONFIG ?= ${ROOTDIR}/configs/dev.yaml

##@ General

.PHONY: help
help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Docs

.PHONY: swag
swag: ## Install swag binary
	@go list -f '{{ join .Imports "\n" }}' tools/tools.go | xargs -t -L1 go install
	@echo Installed swag successfully

.PHONY: api
docs: ## Generate OpenAPI specification
	swag init \
		-d "${ROOTDIR}/cmd,${ROOTDIR}/pkg/adapter/http" \
		-o "${ROOTDIR}/api"
	@echo Generated docs successfully

##@ Development

.PHONY: run-dev
run-dev: ## Run api server in development environment
	IMAGEER_CONFIG=${IMAGEER_CONFIG} go run ${ROOTDIR}/cmd/main.go
