.DEFAULT_GOAL := help

DOCKER_IMAGE_NAME := go-lambda-docker

.PHONY: build
build: ## Build docker image
	docker image build -t $(DOCKER_IMAGE_NAME) .

.PHONY: build-no-cache
build-no-cache: ## Build docker image
	docker image build -t $(DOCKER_IMAGE_NAME) . --no-cache

.PHONY: run
run: ## Do docker compose up with hot reload
	docker container run --rm $(DOCKER_IMAGE_NAME)

.PHONY: help
help: ## Show options
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'
