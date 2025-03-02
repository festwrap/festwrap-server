ROOT ?= ./
IMAGE_NAME ?= "festwrap-server"
IMAGE_TAG ?= "latest"
CONTAINER_NAME ?= "festwrap-server"
PORT ?= 8080

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


.PHONE: run-server
run-server:
	@docker run --name $(CONTAINER_NAME) \
        -d \
		-e FESTWRAP_PORT=$(PORT) \
		-e FESTWRAP_SETLISTFM_APIKEY=$(FESTWRAP_SETLISTFM_APIKEY) \
        -p $(PORT):$(PORT) \
        -t ${IMAGE_NAME}:${IMAGE_TAG}


.PHONE: stop-server
stop-server:
	@docker container stop $(CONTAINER_NAME) && docker container rm $(CONTAINER_NAME)
