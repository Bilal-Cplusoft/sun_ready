.PHONY: help build run test clean docker-build docker-up docker-down docker-logs

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Build the application binary
	@echo "Building sunready..."
	@go build -o sunready ./cmd/sunready

run: ## Run the application locally
	@echo "Running sunready..."
	@go run ./cmd/sunready/main.go

test: ## Run tests
	@echo "Running tests..."
	@go test -v ./...

clean: ## Clean build artifacts
	@echo "Cleaning..."
	@rm -f sunready
	@go clean

fmt: ## Format code
	@echo "Formatting code..."
	@go fmt ./...

lint: ## Run linter
	@echo "Running linter..."
	@golangci-lint run

docker-build: ## Build Docker image
	@echo "Building Docker image..."
	@docker-compose build

docker-up: ## Start Docker containers
	@echo "Starting Docker containers..."
	@docker-compose up -d

docker-down: ## Stop Docker containers
	@echo "Stopping Docker containers..."
	@docker-compose down

docker-logs: ## View Docker logs
	@docker-compose logs -f

docker-clean: ## Remove Docker containers and volumes
	@echo "Cleaning Docker resources..."
	@docker-compose down -v

deps: ## Download dependencies
	@echo "Downloading dependencies..."
	@go mod download

tidy: ## Tidy go.mod
	@echo "Tidying go.mod..."
	@go mod tidy

dev: ## Run in development mode with hot reload (requires air)
	@air

install-tools: ## Install development tools
	@echo "Installing development tools..."
	@go install github.com/cosmtrek/air@latest
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
