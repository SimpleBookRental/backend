# Include environment variables from .env file
-include .env

# Default target
.PHONY: all
all: setup

# Setup development environment
.PHONY: setup
setup:
	@echo "Setting up development environment..."
	@go mod tidy
	@go install github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	@go install github.com/golang/mock/mockgen@latest

# Start PostgreSQL container
.PHONY: db
db:
	@echo "Starting PostgreSQL container..."
	@docker-compose up -d

# Start PostgreSQL container
.PHONY: db-down
db-down:
	@echo "Starting PostgreSQL container..."
	@docker-compose down postgres

# Run database migrations
.PHONY: migrate-up
migrate-up:
	@echo "Running database migrations..."
	@migrate -path migrations -database "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSL_MODE)" up

# Rollback database migrations
.PHONY: migrate-down
migrate-down:
	@echo "Rolling back database migrations..."
	@migrate -path migrations -database "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSL_MODE)" down

# Generate mocks for testing
.PHONY: mock
mock:
	@echo "Generating mocks..."
	@mockgen -source=internal/domain/user.go -destination=internal/mocks/user_mock.go -package=mocks
	@mockgen -source=internal/domain/category.go -destination=internal/mocks/category_mock.go -package=mocks
	@mockgen -source=internal/domain/book.go -destination=internal/mocks/book_mock.go -package=mocks
	@mockgen -source=internal/domain/rental.go -destination=internal/mocks/rental_mock.go -package=mocks
	@mockgen -source=internal/domain/payment.go -destination=internal/mocks/payment_mock.go -package=mocks

# Run tests
.PHONY: test
test:
	@echo "Running tests..."
	@go test -v ./...

# Build production binary
.PHONY: build
build:
	@echo "Building production binary..."
	@go build -o bin/book-rental-api cmd/api/main.go

# Run development server
.PHONY: run
run:
	@echo "Running development server..."
	@go run cmd/api/main.go

# Install Swagger tooling
.PHONY: swagger-install
swagger-install:
	@echo "Installing Swagger tooling..."
	@go install github.com/swaggo/swag/cmd/swag@latest

# Generate Swagger documentation
.PHONY: swagger
swagger:
	@echo "Generating Swagger documentation..."
	@swag init -g cmd/api/main.go -o docs

# Format Swagger annotations
.PHONY: swagger-fmt
swagger-fmt:
	@echo "Formatting Swagger annotations..."
	@swag fmt
