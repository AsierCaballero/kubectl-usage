SHELL := /bin/bash
.DEFAULT_GOAL := help

BIN_DIR := bin
BIN_NAME := kubectl-usage
GO_MODULE := github.com/AsierCaballero/kubectl-usage
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")

.PHONY: help
help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

.PHONY: build
build: ## Build the binary
	go build -ldflags="-X main.Version=$(VERSION)" -o $(BIN_DIR)/$(BIN_NAME) .

.PHONY: install
install: build ## Install as kubectl plugin
	cp $(BIN_DIR)/$(BIN_NAME) /usr/local/bin/kubectl-usage

.PHONY: run
run: build ## Run locally
	./$(BIN_DIR)/$(BIN_NAME)

.PHONY: test
test: fmt vet ## Run tests
	go test ./... -v -coverprofile cover.out

.PHONY: lint
lint: ## Run linters
	golangci-lint run ./...

.PHONY: fmt
fmt: ## Format Go code
	go fmt ./...

.PHONY: vet
vet: ## Run go vet
	go vet ./...

.PHONY: coverage
coverage: test ## Show test coverage
	go tool cover -html=cover.out -o cover.html

.PHONY: clean
clean: ## Clean build artifacts
	rm -rf $(BIN_DIR) cover.out cover.html

.PHONY: docker-build
docker-build: ## Build Docker image
	docker build -t ghcr.io/asiercaballero/kubectl-usage:latest .

.PHONY: docker-push
docker-push: ## Push Docker image
	docker push ghcr.io/asiercaballero/kubectl-usage:latest
