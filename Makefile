.PHONY: setup docker-db migrate-up migrate-down mock test run build docker-build docker-run

# Default target
all: setup

# Setup development environment
setup:
	@echo "Setting up development environment..."
	@go mod tidy
	@go install github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	@go install github.com/golang/mock/mockgen@latest

# Start PostgreSQL container
docker-db:
	@echo "Starting PostgreSQL container..."
	@docker-compose up -d postgres

# Run database migrations
migrate-up:
	@echo "Running database migrations..."
	@migrate -path migrations -database "postgres://postgres:postgres@localhost:5432/book_rental?sslmode=disable" up

# Rollback database migrations
migrate-down:
	@echo "Rolling back database migrations..."
	@migrate -path migrations -database "postgres://postgres:postgres@localhost:5432/book_rental?sslmode=disable" down

# Generate mocks for testing
mock:
	@echo "Generating mocks..."
	@mockgen -source=internal/domain/user.go -destination=internal/mocks/user_mock.go -package=mocks
	@mockgen -source=internal/domain/category.go -destination=internal/mocks/category_mock.go -package=mocks
	@mockgen -source=internal/domain/book.go -destination=internal/mocks/book_mock.go -package=mocks
	@mockgen -source=internal/domain/rental.go -destination=internal/mocks/rental_mock.go -package=mocks
	@mockgen -source=internal/domain/payment.go -destination=internal/mocks/payment_mock.go -package=mocks

# Run tests
test:
	@echo "Running tests..."
	@go test -v ./...

# Run development server
run:
	@echo "Running development server..."
	@go run cmd/api/main.go

# Build production binary
build:
	@echo "Building production binary..."
	@go build -o bin/book-rental-api cmd/api/main.go

# Build Docker image
docker-build:
	@echo "Building Docker image..."
	@docker-compose build

# Run Docker containers
docker-run:
	@echo "Running Docker containers..."
	@docker-compose up -d

# Stop Docker containers
docker-stop:
	@echo "Stopping Docker containers..."
	@docker-compose down

# Clean up
clean:
	@echo "Cleaning up..."
	@rm -rf bin
	@docker-compose down -v
