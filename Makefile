ROOT ?= ./

.PHONY: pre-commit-install
pre-commit-install:
	pre-commit install
	pre-commit install --hook-type commit-msg

.PHONY: run-tests
run-tests:
	go test $(ROOT)/...
