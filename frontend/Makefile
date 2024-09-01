APP_PORT ?= 3000
ENV_FILE ?= .env

include ${ENV_FILE}

define extract_endpoint
$(shell echo ${1} | sed -E 's|^https?://[^/]*||' | tr -d ' ')
endef


.PHONY: pre-commit-install
pre-commit-install:
	pre-commit install
	pre-commit install --hook-type commit-msg

.PHONY: run-app
run-app:
	@echo "Starting the frontend app..."
	@set -a; source ${ENV_FILE}; set +a; \
	npm run dev

.PHONY: run-tests
run-tests:
	@echo "Running tests..."
	@npm run test
