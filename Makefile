# See all available commands
.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'

.PHONY: install-air
install-air: ## Installs 'air' for live reloading
	@echo "Installing air..."
	go install github.com/air-verse/air@latest

.PHONY: dev
dev: ## Starts the dev server with 'air' for live reload
	@echo "Starting dev server with air..."
	@air

.PHONY: build
build: ## Builds the production binary
	@echo "Building binary..."
	go build -o ./bin/orchestrator ./cmd/orchestrator/main.go

.PHONY: tidy
tidy: ## Tidy go modules
	go mod tidy