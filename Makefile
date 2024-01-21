.DEFAULT_GOAL := help

DOCKER_IMAGE_NAME := go-lambda-docker

.PHONY: deploy
deploy: ## deploy to aws lambda
	sls deploy --verbose

.PHONY: build
build: ## Build docker image
	docker image build -t $(DOCKER_IMAGE_NAME) . --platform linux/arm64 --no-cache

.PHONY: build-amd64
build-amd64: ## Build docker image
	docker image build -t $(DOCKER_IMAGE_NAME) . --platform linux/amd64 --no-cache

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
