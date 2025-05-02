# Technical Context: Book Rental System

## Technology Stack
- **Backend Language**: Golang
- **Database**: PostgreSQL (Docker containerized)
- **Web Framework**: Gin
- **Database Migration**: go-migrate
- **Testing Framework**: GoMock
- **Configuration**: Environment variables via .env file
- **Build System**: Makefile
- **Authentication**: JWT (JSON Web Tokens)
- **Documentation**: Markdown files
- **Version Control**: Git
- **Logging**: Structured logging system

## Project Structure
- **cmd/**: Application entry points
  - **api/**: Main API server (main.go)
- **docs/**: Documentation and diagrams
  - API flow diagrams for each domain area
  - System overview documentation
  - Product requirements documentation
- **internal/**: Core application code
  - **api/**: HTTP handlers and middleware
  - **domain/**: Domain entities and interfaces
  - **repository/**: Database implementations
  - **service/**: Business logic implementations
  - **mocks/**: Test mock objects
- **migrations/**: SQL database migrations
  - Sequential migration files for schema changes
- **scripts/**: Utility scripts
- **tests/**: Integration tests
- **pkg/**: Shared packages and utilities
  - **auth/**: JWT authentication service
  - **config/**: Configuration loading
  - **logger/**: Structured logging system

## Domain Model Design
- **User**: Authentication and profile management with role-based access
- **Book**: Catalog management with inventory tracking
- **Category**: Taxonomic organization of book collection
- **Rental**: Borrowing record with status tracking
- **Payment**: Financial transaction recording and reporting

## Interface Design Patterns
- Repository interfaces for data access with specialized query methods
- Service interfaces for business logic implementation
- Clear separation between data access and business rules
- Consistent method signatures across similar entities
- Status enums for state management (RentalStatus, PaymentStatus)

## Development Setup
- Local Go environment (1.21+)
- Docker and Docker Compose for containerization
- PostgreSQL database running in Docker
- Environment variables configured in .env file

## Key Dependencies
- Go standard library
- Gin web framework for routing and middleware
- JWT-Go for authentication
- Testify for assertions in tests
- GoMock for mocking interfaces
- PostgreSQL driver for database connectivity
- go-migrate for database schema migrations

## Development Workflow
- Database setup via Docker
- Migrations applied using make commands
- Code changes in relevant layers
- Unit tests with mocked dependencies
- Full build with Makefile commands
- Containerized deployment

## Build Processes
- `make run`: Start development server
- `make test`: Run unit tests
- `make migrate-up/down`: Manage database migrations
- `make mock`: Generate mock interfaces
- `make build`: Create production binary
- `make docker-build`: Create Docker image
- `make docker-run`: Run application in Docker

## Server Configuration
- Mode-based configuration (development/release)
- Configurable server port
- Read and write timeouts
- Graceful shutdown with context timeout
- Signal handling for proper shutdown
- Maximum request header size limitation
