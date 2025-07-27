ROOT ?= ./
IMAGE_NAME ?= "festwrap-server"
IMAGE_TAG ?= "latest"
CONTAINER_NAME ?= "festwrap-server"
PORT ?= 8080
ENV_FILE ?= ".env"
ENV_VARS := $(shell cat ${ENV_FILE} | xargs)

PUBSUB_PORT?=8085
PUBSUB_TEST_SUBSCRIPTION=test-consumer
DOCKER_COMPOSE_NETWORK?=integration_backing_services


.PHONY: pre-commit-install
pre-commit-install:
	pre-commit install
	pre-commit install --hook-type commit-msg


.PHONY: run-backing-services
run-backing-services:
	@export $(ENV_VARS) && \
	PUBSUB_PORT=$(PUBSUB_PORT) \
	FESTWRAP_PUBSUB_TEST_SUBSCRIPTION=$(PUBSUB_TEST_SUBSCRIPTION) \
	docker compose -f integration/docker-compose.yml up -d


.PHONY: clean-backing-services
clean-backing-services:
	@export $(ENV_VARS) && \
	PUBSUB_PORT=$(PUBSUB_PORT) \
	FESTWRAP_PUBSUB_TEST_SUBSCRIPTION=$(PUBSUB_TEST_SUBSCRIPTION) \
	docker compose -f integration/docker-compose.yml down --volumes --remove-orphans


.PHONY: run-local-server
run-local-server: run-backing-services
	@(trap "$(MAKE) clean-backing-services" EXIT; \
	export $(ENV_VARS) && PUBSUB_EMULATOR_HOST=localhost:$(PUBSUB_PORT) go run ./cmd)


.PHONY: run-unit-tests
run-unit-tests:
	@echo "Running unit tests..."
	go test $(ROOT)/... -v -short


.PHONY: run-integration-tests
run-integration-tests:
	@echo "Running integration tests..."
	go test $(ROOT)/... -v -run Integration


.PHONY: run-tests
run-tests: run-unit-tests run-integration-tests


.PHONY: build-image
build-image:
	docker build -f Dockerfile -t ${IMAGE_NAME}:${IMAGE_TAG} .


.PHONY: run-server
run-server: build-image run-backing-services
	@(trap "$(MAKE) clean-backing-services" EXIT; \
	docker run --name $(CONTAINER_NAME) \
		--env-file ${ENV_FILE} \
		-e PUBSUB_EMULATOR_HOST=pubsub:$(PUBSUB_PORT) \
		-p $(PORT):$(PORT) \
		--network $(DOCKER_COMPOSE_NETWORK) \
		-t ${IMAGE_NAME}:${IMAGE_TAG})


.PHONY: stop-server
stop-server:
	@docker container stop $(CONTAINER_NAME) && docker container rm $(CONTAINER_NAME)
