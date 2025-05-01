# Book Rental System - System Overview

This diagram provides a comprehensive view of the Book Rental System architecture, showing how different components interact with each other.

```mermaid
graph TD
    subgraph "Client Layer"
        Client[Client Applications]
    end

    subgraph "API Layer"
        API[Gin Router/API Endpoints]
        Middleware[Authentication & Authorization]
    end

    subgraph "Handler Layer"
        AuthHandler[Auth Handler]
        UserHandler[User Handler]
        CategoryHandler[Category Handler]
        BookHandler[Book Handler]
        RentalHandler[Rental Handler]
        PaymentHandler[Payment Handler]
        ReportHandler[Report Handler]
    end

    subgraph "Service Layer"
        AuthService[Auth Service]
        UserService[User Service]
        CategoryService[Category Service]
        BookService[Book Service]
        RentalService[Rental Service]
        PaymentService[Payment Service]
        ReportService[Report Service]
        JWTService[JWT Service]
    end

    subgraph "Repository Layer"
        UserRepo[User Repository]
        CategoryRepo[Category Repository]
        BookRepo[Book Repository]
        RentalRepo[Rental Repository]
        PaymentRepo[Payment Repository]
    end

    subgraph "Database Layer"
        DB[PostgreSQL Database]
    end

    %% Client to API connections
    Client -->|HTTP Requests| API
    API -->|Authentication Check| Middleware

    %% API to Handler connections
    Middleware -->|Auth Requests| AuthHandler
    Middleware -->|User Requests| UserHandler
    Middleware -->|Category Requests| CategoryHandler
    Middleware -->|Book Requests| BookHandler
    Middleware -->|Rental Requests| RentalHandler
    Middleware -->|Payment Requests| PaymentHandler
    Middleware -->|Report Requests| ReportHandler

    %% Handler to Service connections
    AuthHandler -->|Method Calls| AuthService
    UserHandler -->|Method Calls| UserService
    CategoryHandler -->|Method Calls| CategoryService
    BookHandler -->|Method Calls| BookService
    RentalHandler -->|Method Calls| RentalService
    PaymentHandler -->|Method Calls| PaymentService
    ReportHandler -->|Method Calls| ReportService

    %% Service dependencies
    AuthService -->|Uses| JWTService
    AuthService -->|Uses| UserService
    BookService -->|Uses| CategoryService
    RentalService -->|Uses| BookService
    PaymentService -->|Uses| RentalService
    ReportService -->|Uses| BookService
    ReportService -->|Uses| RentalService
    ReportService -->|Uses| PaymentService

    %% Service to Repository connections
    UserService -->|Data Access| UserRepo
    CategoryService -->|Data Access| CategoryRepo
    BookService -->|Data Access| BookRepo
    RentalService -->|Data Access| RentalRepo
    PaymentService -->|Data Access| PaymentRepo

    %% Repository to Database connections
    UserRepo -->|SQL Queries| DB
    CategoryRepo -->|SQL Queries| DB
    BookRepo -->|SQL Queries| DB
    RentalRepo -->|SQL Queries| DB
    PaymentRepo -->|SQL Queries| DB

    %% Database connections
    DB -->|Returns Data| UserRepo
    DB -->|Returns Data| CategoryRepo
    DB -->|Returns Data| BookRepo
    DB -->|Returns Data| RentalRepo
    DB -->|Returns Data| PaymentRepo
```

## Architectural Patterns

The Book Rental System implements:

1. **Clean Architecture** - Separating concerns with distinct layers:
   - External interfaces (API)
   - Business rules (Services)
   - Data access (Repositories)

2. **Dependency Injection** - Each layer receives its dependencies through constructors

3. **Repository Pattern** - Abstracting data access logic

4. **Service Layer Pattern** - Encapsulating business logic

5. **Middleware Pattern** - For cross-cutting concerns like authentication and logging

## Data Flow

1. Client makes HTTP request to API endpoints
2. Middleware processes authentication and authorization
3. Handler receives validated request
4. Handler calls appropriate Service methods
5. Service implements business logic
6. Service calls Repository methods for data access
7. Repository executes database operations
8. Data flows back up through the layers
9. Handler formats response
10. Client receives HTTP response

This architecture ensures separation of concerns, maintainability, and testability of the system.
