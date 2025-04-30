# Book Rental System API

Book Rental System is a RESTful API that allows users to browse, borrow, and return books from a virtual library. The system tracks book inventory, user accounts, rental history, and fees.

## Technology Stack

- **Backend Language**: Golang
- **Database**: PostgreSQL (running in Docker container)
- **Web Framework**: Gin
- **SQL Interface**: SQLC (for generating type-safe code from SQL)
- **Database Migration**: go-migrate
- **Testing**: GoMock for mocking framework
- **Configuration**: .env file
- **Build System**: Makefile

## Core Features

1. **User Management**
   - User registration and authentication
   - User profile management
   - Role-based access control (Admin, Librarian, Member)

2. **Book Catalog Management**
   - Add, update, and delete books
   - Book categorization and metadata management
   - Search and filter capabilities
   - Book availability status

3. **Rental Operations**
   - Book borrowing process
   - Book return process
   - Rental period extension
   - Overdue book management

4. **Fees and Payments**
   - Rental fee calculation
   - Late return penalty calculation
   - Payment processing
   - Refund handling

5. **Reporting and Analytics**
   - Usage statistics
   - Popular book reports
   - Revenue reports
   - User activity reports

## Installation

### Requirements

- Go 1.21+
- Docker and Docker Compose
- Make

### Development Environment Setup

1. Clone repository:
   ```bash
   git clone https://github.com/yourusername/book-rental-system.git
   cd book-rental-system
   ```

2. Install necessary tools:
   ```bash
   make setup
   ```

3. Create .env file from .env.example:
   ```bash
   cp .env.example .env
   ```

4. Start PostgreSQL using Docker:
   ```bash
   make docker-db
   ```

5. Run database migrations:
   ```bash
   make migrate-up
   ```

6. Generate SQLC code:
   ```bash
   make sqlc
   ```

7. Run development server:
   ```bash
   make run
   ```

### Using Docker

To run the entire application in Docker:

```bash
make docker-run
```

## API Endpoints

See the full API documentation in [docs/prd.md](docs/prd.md).

## Development

### Project Structure

```
.
├── cmd/                  # Application entry points
│   └── api/              # API server
├── docs/                 # Documentation
├── internal/             # Internal code
│   ├── api/              # API handlers
│   ├── domain/           # Domain definitions and interfaces
│   ├── repository/       # Repository implementations
│   └── service/          # Service implementations
├── migrations/           # Database migrations
├── pkg/                  # Reusable packages
│   ├── auth/             # Authentication and authorization
│   ├── config/           # Configuration
│   └── logger/           # Logging
├── sqlc/                 # SQLC configuration and queries
│   └── queries/          # SQL queries
├── .env.example          # Example configuration file
├── docker-compose.yml    # Docker Compose configuration
├── Dockerfile            # Docker configuration
├── go.mod                # Go dependencies
├── go.sum                # Go dependencies checksum
└── Makefile              # Make commands
```

### Useful Commands

- `make run`: Run development server
- `make test`: Run tests
- `make migrate-up`: Run database migrations
- `make migrate-down`: Rollback database migrations
- `make sqlc`: Generate SQLC code
- `make mock`: Generate mocks for testing
- `make build`: Build production binary
- `make docker-build`: Build Docker image
- `make docker-run`: Run Docker containers
- `make docker-stop`: Stop Docker containers
- `make clean`: Clean up

## License

MIT
