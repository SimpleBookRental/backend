# Simple Book Rental API

A RESTful API for a simple book rental system built with Go, Gin, and GORM.

## Features

- CRUD operations for users and books
- RESTful API design
- PostgreSQL database with GORM ORM
- Clean architecture
- Docker support

## Requirements

- Go 1.24 or higher
- PostgreSQL
- Docker (optional)

## Getting Started

### Environment Setup

Copy the example environment file and modify it as needed:

```bash
cp .env.example .env
```

### Running Locally

```bash
# Build the application
make build

# Run the application
make run
```

### Running with Docker

```bash
# Build and run with Docker Compose
make docker-run
```

## API Endpoints

### Users

- `POST /api/v1/users` - Create a new user
- `GET /api/v1/users` - Get all users
- `GET /api/v1/users/:id` - Get a user by ID
- `PUT /api/v1/users/:id` - Update a user
- `DELETE /api/v1/users/:id` - Delete a user

### Books

- `POST /api/v1/books` - Create a new book
- `GET /api/v1/books` - Get all books
- `GET /api/v1/books/:id` - Get a book by ID
- `PUT /api/v1/books/:id` - Update a book
- `DELETE /api/v1/books/:id` - Delete a book

### User's Books

- `GET /api/v1/users/:user_id/books` - Get all books by user ID

## Testing

```bash
# Run tests
make test

# Run tests with coverage
make test-coverage
```

## Project Structure

```
.
├── cmd
│   └── api
│       └── main.go           # Application entry point
├── internal
│   ├── config                # Configuration
│   ├── controllers           # HTTP request handlers
│   ├── models                # Data models
│   ├── repositories          # Data access layer
│   ├── routes                # API routes
│   └── services              # Business logic
├── pkg
│   ├── database              # Database connection
│   └── utils                 # Utility functions
├── .env                      # Environment variables
├── docker-compose.yml        # Docker Compose configuration
├── Dockerfile                # Docker configuration
└── Makefile                  # Build commands
```

## License

MIT