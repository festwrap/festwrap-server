ROOT ?= ./
IMAGE_NAME ?= "festwrap-server"
IMAGE_TAG ?= "latest"
CONTAINER_NAME ?= "festwrap-server"
PORT ?= 8080

# Load environment variables from .env file if it exists
ifneq (,$(wildcard ./.env))
    include .env
    export
endif

.PHONY: pre-commit-install
pre-commit-install:
	pre-commit install
	pre-commit install --hook-type commit-msg


.PHONY: run-local-server
run-local-server:
	@if [ ! -f .env ]; then echo "Error: .env file not found. Please create one with your environment variables."; exit 1; fi
	@export $(shell cat .env | xargs) && go run ./cmd


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
run-server:
	@if [ ! -f .env ]; then echo "Error: .env file not found. Please create one with your environment variables."; exit 1; fi
	@if [ -z "$(FESTWRAP_SETLISTFM_APIKEY)" ]; then echo "Error: FESTWRAP_SETLISTFM_APIKEY not set in .env file"; exit 1; fi
	@if [ -z "$(SPOTIFY_CLIENT_ID)" ]; then echo "Error: SPOTIFY_CLIENT_ID not set in .env file"; exit 1; fi
	@if [ -z "$(SPOTIFY_CLIENT_SECRET)" ]; then echo "Error: SPOTIFY_CLIENT_SECRET not set in .env file"; exit 1; fi
	@if [ -z "$(SPOTIFY_REFRESH_TOKEN)" ]; then echo "Error: SPOTIFY_REFRESH_TOKEN not set in .env file"; exit 1; fi
	@docker run --name $(CONTAINER_NAME) \
        -d \
		-e FESTWRAP_PORT=$(PORT) \
		-e FESTWRAP_SETLISTFM_APIKEY=$(FESTWRAP_SETLISTFM_APIKEY) \
		-e SPOTIFY_CLIENT_ID=$(SPOTIFY_CLIENT_ID) \
		-e SPOTIFY_CLIENT_SECRET=$(SPOTIFY_CLIENT_SECRET) \
		-e SPOTIFY_REFRESH_TOKEN=$(SPOTIFY_REFRESH_TOKEN) \
		-e FESTWRAP_MAX_CONNS_PER_HOST=$(FESTWRAP_MAX_CONNS_PER_HOST) \
		-e FESTWRAP_SETLISTFM_NUM_SEARCH_PAGES=$(FESTWRAP_SETLISTFM_NUM_SEARCH_PAGES) \
		-e FESTWRAP_MAX_UPDATE_ARTISTS=$(FESTWRAP_MAX_UPDATE_ARTISTS) \
		-e FESTWRAP_ADD_SETLIST_SLEEP_MS=$(FESTWRAP_ADD_SETLIST_SLEEP_MS) \
		-e FESTWRAP_HTTP_CLIENT_TIMEOUT_S=$(FESTWRAP_HTTP_CLIENT_TIMEOUT_S) \
        -p $(PORT):$(PORT) \
        -t ${IMAGE_NAME}:${IMAGE_TAG}


.PHONY: stop-server
stop-server:
	@docker container stop $(CONTAINER_NAME) && docker container rm $(CONTAINER_NAME)
