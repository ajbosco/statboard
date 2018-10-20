.EXPORT_ALL_VARIABLES:
NAME := reads
PKG := github.com/ajbosco/reads/cmd/reads

.PHONY: fmt
fmt: ## Verifies all files have been `gofmt`ed.
	@gofmt -s -l . | grep -v vendor | tee /dev/stderr

.PHONY: lint
lint: ## Verifies `golint` passes.
	@golint ./... | grep -v vendor | tee /dev/stderr

.PHONY: test
test: ## Runs the go tests.
	@go test -cover -race $(shell go list ./... | grep -v vendor)

.PHONY: vet
vet: ## Verifies `go vet` passes.
	@go vet $(shell go list ./... | grep -v vendor) | tee /dev/stderr

.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | sed 's/^[^:]*://g' | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
