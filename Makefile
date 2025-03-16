.PHONY: build run test clean migrate swagger docker-build docker-run

# Application name
APP_NAME=todo-api

# Build the application
build:
	@echo "Building $(APP_NAME)..."
	@go build -o bin/$(APP_NAME) ./cmd/api

# Run the application
run: build
	@echo "Running $(APP_NAME)..."
	@./bin/$(APP_NAME)

# Run the application in development mode with hot reload (requires air: https://github.com/cosmtrek/air)
dev:
	@echo "Running $(APP_NAME) in development mode..."
	@if command -v air > /dev/null; then \
		air; \
	else \
		echo "Error: air is not installed. Install it with 'go install github.com/cosmtrek/air@latest'"; \
		exit 1; \
	fi

# Run database migrations
migrate:
	@echo "Running database migrations..."
	@go run ./scripts/migrate.go

# Run tests
test:
	@echo "Running tests..."
	@go test -v ./...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	@go test -v -cover ./...
	@go tool cover -html=coverage.out -o coverage.html

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf bin/
	@rm -f coverage.out coverage.html

# Generate Swagger documentation (requires swag: https://github.com/swaggo/swag)
swagger:
	@echo "Generating Swagger documentation..."
	@if command -v swag > /dev/null; then \
		swag init -g cmd/api/main.go; \
	else \
		echo "Error: swag is not installed. Install it with 'go install github.com/swaggo/swag/cmd/swag@latest'"; \
		exit 1; \
	fi

# Build Docker image
docker-build:
	@echo "Building Docker image..."
	@docker build -t $(APP_NAME) .

# Run Docker container
docker-run: docker-build
	@echo "Running Docker container..."
	@docker run -p 3000:3000 $(APP_NAME)

# Help command
help:
	@echo "Available commands:"
	@echo "  make build          - Build the application"
	@echo "  make run            - Build and run the application"
	@echo "  make dev            - Run with hot reload (requires air)"
	@echo "  make migrate        - Run database migrations"
	@echo "  make test           - Run tests"
	@echo "  make test-coverage  - Run tests with coverage"
	@echo "  make clean          - Clean build artifacts"
	@echo "  make swagger        - Generate Swagger documentation"
	@echo "  make docker-build   - Build Docker image"
	@echo "  make docker-run     - Run Docker container"

# Default target
default: help
