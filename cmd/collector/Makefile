.EXPORT_ALL_VARIABLES:
NAME := collector
PKG := github.com/ajbosco/statboard/cmd/collector
BUILD_DIR := $(shell pwd)/build
TARGET := ${BUILD_DIR}/${NAME}

.PHONY: build
build: ## run go build for current OS
	@go build -ldflags "$(LDFLAGS)" -o "${TARGET}" ${PKG}

.PHONY: clean
clean: ## Cleanup any build binaries or packages.
	$(RM) -r $(BUILD_DIR)

.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | sed 's/^[^:]*://g' | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
