.DEFAULT_GOAL := help

ROOTDIR:=$(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))))

SWAG_IMAGE_NAME?=tmp/swag

.PHONY: swag
swag: ## Build swag docker image
	docker build \
		-t ${SWAG_IMAGE_NAME} \
		- < ${ROOTDIR}/build/docker/swag.Dockerfile

.PHONY: help
help: ## Print this help message
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'