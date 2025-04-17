.PHONY: build up down logs ps clean help restart test test-coverage start mock swagger

# Default target
.DEFAULT_GOAL := help

# Variables
DOCKER_COMPOSE = docker-compose

# Build the Docker images
build:
	@echo "Building Docker images..."
	$(DOCKER_COMPOSE) build

# Start the services
up:
	@echo "Starting services..."
	$(DOCKER_COMPOSE) up -d

# Stop the services
down:
	@echo "Stopping services..."
	$(DOCKER_COMPOSE) down

# Show logs
logs:
	@echo "Showing logs..."
	$(DOCKER_COMPOSE) logs -f

# Show running containers
ps:
	@echo "Listing containers..."
	$(DOCKER_COMPOSE) ps

# Clean up volumes
clean:
	@echo "Cleaning up volumes..."
	$(DOCKER_COMPOSE) down -v

# Restart services
restart: down up
	@echo "Services restarted"

# Run tests
test: mock
	@echo "Running tests..."
	go test ./... -count=1 -coverprofile=coverage.out

# Run tests with coverage report
test-coverage: test
	@echo "Generating coverage report..."
	go tool cover -html=coverage.out -o coverage.html

# Build and start services
start: build up logs
	@echo "Services started"

# Build the application locally
build-local:
	@echo "Building application locally..."
	go build -o bin/api cmd/api/main.go

# Run the application locally
run:
	@echo "Running application locally..."
	go run cmd/api/main.go

# Generate mocks
mock:
	@echo "Generating mocks..."
	@mockgen -destination=internal/mocks/mock_user_service.go -package=mocks github.com/SimpleBookRental/backend/internal/services UserServiceInterface
	@mockgen -destination=internal/mocks/mock_book_service.go -package=mocks github.com/SimpleBookRental/backend/internal/services BookServiceInterface
	@mockgen -destination=internal/mocks/mock_book_user_service.go -package=mocks github.com/SimpleBookRental/backend/internal/services BookUserServiceInterface
	@mockgen -destination=internal/mocks/mock_token_service.go -package=mocks github.com/SimpleBookRental/backend/internal/services TokenServiceInterface
	@mockgen -destination=internal/mocks/mock_user_repository.go -package=mocks github.com/SimpleBookRental/backend/internal/repositories UserRepositoryInterface
	@mockgen -destination=internal/mocks/mock_book_repository.go -package=mocks github.com/SimpleBookRental/backend/internal/repositories BookRepositoryInterface
	@mockgen -destination=internal/mocks/mock_token_repository.go -package=mocks github.com/SimpleBookRental/backend/internal/repositories TokenRepositoryInterface
	@mockgen -destination=internal/mocks/mock_transaction_manager.go -package=mocks github.com/SimpleBookRental/backend/internal/repositories TransactionManagerInterface
	@echo "Mocks generated successfully"

# Generate swagger.yaml from annotated code using swaggo/swag
swagger:
	@echo "Generating swagger.yaml using swaggo/swag..."
	@swag init -g cmd/api/main.go --outputTypes yaml --output ./ --parseDependency --parseInternal
	@if exist ./docs/swagger.yaml move /Y ./docs/swagger.yaml ./swagger.yaml
	@echo "swagger.yaml generated successfully"

# Help
help:
	@echo "Available commands:"
	@echo "  make build        - Build Docker images"
	@echo "  make build-local  - Build application locally"
	@echo "  make run          - Run application locally"
	@echo "  make up           - Start services"
	@echo "  make down         - Stop services"
	@echo "  make logs         - Show logs"
	@echo "  make ps           - List containers"
	@echo "  make clean        - Clean up volumes"
	@echo "  make restart      - Restart services"
	@echo "  make test         - Run tests"
	@echo "  make test-coverage - Run tests with coverage report"
	@echo "  make start        - Build and start services"
	@echo "  make mock         - Generate mock files for testing"
	@echo "  make swagger      - Generate swagger.yaml from annotated code (requires swag, see README.md)"
	@echo "  make help         - Show this help"
