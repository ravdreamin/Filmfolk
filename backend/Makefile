.PHONY: help run build test clean docker-build docker-up docker-down migrate lint fmt

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

run: ## Run the application locally
	@echo "Starting FilmFolk API..."
	go run cmd/server/main.go

build: ## Build the application
	@echo "Building FilmFolk API..."
	go build -o bin/server cmd/server/main.go

test: ## Run tests
	@echo "Running tests..."
	go test -v ./...

test-coverage: ## Run tests with coverage
	@echo "Running tests with coverage..."
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

clean: ## Clean build artifacts
	@echo "Cleaning..."
	rm -rf bin/
	rm -f coverage.out coverage.html

fmt: ## Format code
	@echo "Formatting code..."
	go fmt ./...

lint: ## Run linter
	@echo "Running linter..."
	golangci-lint run ./...

tidy: ## Tidy go modules
	@echo "Tidying go modules..."
	go mod tidy

# Docker commands
docker-build: ## Build Docker image
	@echo "Building Docker image..."
	docker build -t filmfolk-api:latest .

docker-up: ## Start services with docker-compose
	@echo "Starting services..."
	docker-compose up -d

docker-down: ## Stop services
	@echo "Stopping services..."
	docker-compose down

docker-logs: ## View logs
	docker-compose logs -f api

docker-restart: docker-down docker-up ## Restart services

# Database commands
migrate-up: ## Run database migrations up
	@echo "Running migrations..."
	go run cmd/migrate/main.go

migrate-create: ## Create new migration file (use name=migration_name)
	@echo "Creating migration: $(name)"
	migrate create -ext sql -dir migrations -seq $(name)

db-reset: ## Reset database (WARNING: drops all data)
	@echo "Resetting database..."
	docker-compose down -v
	docker-compose up -d postgres
	sleep 5
	go run cmd/migrate/main.go

# Development commands
dev: ## Run in development mode with hot reload (requires air)
	@echo "Starting development server with hot reload..."
	air

install-tools: ## Install development tools
	@echo "Installing development tools..."
	go install github.com/cosmtrek/air@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
