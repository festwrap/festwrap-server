ROOT ?= ./
IMAGE_NAME ?= "festwrap-server"
IMAGE_TAG ?= "latest"

.PHONY: pre-commit-install
pre-commit-install:
	pre-commit install
	pre-commit install --hook-type commit-msg

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


.PHONE: build-image
build-image:
	docker build -f Dockerfile -t ${IMAGE_NAME}:${IMAGE_TAG} .
