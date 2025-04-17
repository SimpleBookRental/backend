# Simple Book Rental API

A RESTful API for a book rental system built with Go, Gin, and GORM, following Clean Architecture.

## Main Features

- User & book CRUD
- JWT authentication, role-based access (ADMIN/USER)
- Book ownership & transfer
- PostgreSQL, Docker, CI/CD with GitHub Actions

## Quick Start

### Prerequisites

- Go 1.24+
- PostgreSQL
- Docker (optional)

### Setup

```bash
cp .env.example .env
make build
make run
```

Or with Docker:

```bash
docker-compose up --build
```

### Run Tests

```bash
make test
make test-coverage
make mock
```

## API Documentation

### Generate OpenAPI (swagger.yaml) automatically

This project uses [swaggo/swag](https://github.com/swaggo/swag) to generate the `swagger.yaml` file from annotated Go code.

#### 1. Install swag CLI

```bash
go install github.com/swaggo/swag/cmd/swag@latest
```
Make sure `$GOPATH/bin` (or `$GOBIN`) is in your `PATH`.

#### 2. Annotate your handlers

Add swagger comments to your handler functions. Example for a user controller:

```go
// CreateUser godoc
// @Summary      Create user
// @Description  Create a new user
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        user  body      models.UserCreate  true  "User create payload"
// @Success      201   {object}  models.User
// @Failure      400   {object}  models.ErrorResponse
// @Router       /api/v1/users [post]
func (uc *UserController) CreateUser(c *gin.Context) {
    // ...
}
```
> See [swaggo/swag annotation docs](https://github.com/swaggo/swag#declarative-comments-format) for more details.

#### 3. Generate swagger.yaml

```bash
make swagger
```
This will scan your code and update `swagger.yaml` at the project root.

- User registration, login, JWT refresh/logout
- CRUD for users and books
- Book transfer between users
- Role-based access: ADMIN (all), USER (own resources)

See `swagger.yaml` for full API details.

## Project Structure

```
.
├── cmd/api/                 # Entry point
├── internal/
│   ├── controllers/         # HTTP handlers
│   ├── middleware/          # Auth, role middleware
│   ├── models/              # Data models
│   ├── repositories/        # Data access
│   ├── services/            # Business logic
│   ├── mocks/               # Generated mocks
│   └── routes/              # API routes
├── pkg/                     # Utilities, DB
├── .github/workflows/       # CI/CD
├── Dockerfile, Makefile, etc.
```

## CI/CD

- GitHub Actions: build, test, coverage on push/PR
- Mocks auto-generated before test

## License

MIT
