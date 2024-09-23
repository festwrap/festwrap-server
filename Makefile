ROOT ?= ./
IMAGE_NAME ?= "festwrap-server"
IMAGE_TAG ?= "latest"

.PHONY: pre-commit-install
pre-commit-install:
	pre-commit install
	pre-commit install --hook-type commit-msg

.PHONY: run-tests
run-tests:
	go test $(ROOT)/...

.PHONE: build-image
build-image:
	docker build -f Dockerfile -t ${IMAGE_NAME}:${IMAGE_TAG} .
