.EXPORT_ALL_VARIABLES:
NAME := collector
PKG := github.com/ajbosco/statboard/cmd/collector
REGISTRY ?= ajbosco
COMMIT ?= $(shell git rev-parse HEAD)
BUILD_DIR := $(shell pwd)/build
TARGET := ${BUILD_DIR}/${NAME}
LDFLAGS ?=

.PHONY: build
build: build-linux ## run go build for linux

.PHONY: build-native
build-native: ## run go build for current OS
	@go build --mod=vendor -ldflags "$(LDFLAGS)" -o "${TARGET}" ${PKG}

.PHONY: build-linux
build-linux:
	@GOOS=linux GOARCH=amd64 go build --mod=vendor -ldflags "$(LDFLAGS)" -o "${TARGET}" ${PKG}

.PHONY: docker-build
docker-build: build ## build docker image
	docker build -t $(REGISTRY)/$(NAME):$(COMMIT) .

.PHONY: docker-push
docker-push: ## push docker image to registry
	docker push $(REGISTRY)/$(NAME):$(COMMIT)

.PHONY: docker-run
docker-run: ## run docker image locally
	@docker run --rm \
	-v ${STATBOARD_DIR}:/statboard \
	-e STATBOARD_CONFIGFILEPATH=/statboard/config.yml -e STATBOARD_DBFILEPATH=/statboard/statboard.db \
	$(REGISTRY)/$(NAME):$(COMMIT)