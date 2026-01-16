# Load environment variables from .env.local if it exists
ifneq (,$(wildcard ./.env.local))
    include .env.local
    export $(shell sed 's/=.*//' ./.env.local)
endif

# Image tags
TAG_GATEWAY ?= latest
TAG_PROCESSOR ?= latest
TAG_UI ?= latest

# Image names
IMAGE_GATEWAY = redshoore/imageer-gateway:$(TAG_GATEWAY)
IMAGE_PROCESSOR = redshoore/imageer-processor:$(TAG_PROCESSOR)
IMAGE_UI = redshoore/imageer-ui:$(TAG_UI)

# Docker Hub credentials
DOCKER_USER ?= <docker_hub_username>
DOCKER_PASSWORD ?= <docker_hub_password>

##@ Deployments

.PHONY: gateway-docker-build
gateway-docker-build: ## Build docker image of gateway module.
	@echo "🔨 Building gateway image: $(IMAGE_GATEWAY)..." && \
		docker build --platform linux/amd64 -f deployment/gateway.Dockerfile -t $(IMAGE_GATEWAY) . && \
		echo "✅ Gateway image built successfully."

.PHONY: gateway-docker-push
gateway-docker-push: ## Push docker image of gateway module.
	@echo "🚀 Pushing gateway image: $(IMAGE_GATEWAY)..." && \
		echo $(DOCKER_PASSWORD) | docker login -u $(DOCKER_USER) --password-stdin && \
		docker push $(IMAGE_GATEWAY) && \
		echo "✅ Gateway image pushed successfully."

.PHONY: processor-docker-build
processor-docker-build: ## Build docker image of processor module.
	@echo "🔨 Building processor image: $(IMAGE_PROCESSOR)..." && \
		docker build --platform linux/amd64 -f deployment/processor.Dockerfile -t $(IMAGE_PROCESSOR) . && \
		echo "✅ Processor image built successfully."

.PHONY: processor-docker-push
processor-docker-push: ## Push docker image of processor module.
	@echo "🚀 Pushing processor image: $(IMAGE_PROCESSOR)..." && \
		echo $(DOCKER_PASSWORD) | docker login -u $(DOCKER_USER) --password-stdin && \
		docker push $(IMAGE_PROCESSOR) && \
		echo "✅ Processor image pushed successfully."

.PHONY: ui-docker-build
ui-docker-build: ## Build docker image of UI module.
	@echo "🔨 Building UI image: $(IMAGE_UI)..." && \
		docker build --platform linux/amd64 -f deployment/ui.Dockerfile -t $(IMAGE_UI) ui && \
		echo "✅ UI image built successfully."

.PHONY: ui-docker-push
ui-docker-push: ## Push docker image of UI module.
	@echo "🚀 Pushing UI image: $(IMAGE_UI)..." && \
		echo $(DOCKER_PASSWORD) | docker login -u $(DOCKER_USER) --password-stdin && \
		docker push $(IMAGE_UI) && \
		echo "✅ UI image pushed successfully."
