APP_PORT ?= 3000
ENV_FILE ?= .env

include ${ENV_FILE}

define extract_endpoint
$(shell echo ${1} | sed -E 's|^https?://[^/]*||' | tr -d ' ')
endef

.PHONY: run-app
run-app:
	@set -a; source ${ENV_FILE}; set +a; \
	REDIRECT_ENDPOINT="$(call extract_endpoint,${REDIRECT_URI})" node index.js
