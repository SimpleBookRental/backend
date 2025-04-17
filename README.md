# Simple Book Rental API

A RESTful API system for book rental, built with Go, Gin, and GORM, following Clean Architecture. Supports user and book management, role-based access, and book ownership transfer.

## Quick Start

### Requirements

- Go 1.24+
- PostgreSQL
- Docker (optional)

### Run locally

From source:

```bash
cp .env.example .env
make build
make run
```

Or with Docker:

```bash
make start
```

### Run tests

```bash
make test
```

## API Documentation

- See the `swagger.yaml` file or access the `/swagger/index.html` endpoint (if available).

## Project Structure

```
.
├── cmd/api/                 # Application entry point
├── internal/
│   ├── controllers/         # HTTP handlers
│   ├── middleware/          # Auth, role middleware
│   ├── models/              # Data models
│   ├── repositories/        # Data access layer
│   ├── services/            # Business logic
│   ├── mocks/               # Generated mocks for testing
│   └── routes/              # API route definitions
├── pkg/                     # Utilities, DB connection, etc.
├── Dockerfile
├── docker-compose.yml
├── Makefile
├── README.md
├── swagger.yaml
└── .env.example
```
