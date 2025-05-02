# Active Context: Book Rental System

## Current Focus
- Project initialization and structure setup
- Core domain entity definition and interface design
- Repository and service interface implementation
- API endpoint handlers development
- Database schema implementation
- Authentication system implementation with JWT
- API documentation with Swagger

## Recent Changes
- Initial project setup with Clean Architecture structure
- Core domains defined with detailed interfaces:
  - User: Authentication and profile management
  - Book: Catalog management with availability tracking
  - Category: Book categorization system
  - Rental: Borrowing records and status tracking
  - Payment: Financial transaction recording
- Database schema designed with migration files ready
- API handlers structured for main endpoints
- README and documentation created
- Server configuration and startup sequence implemented
- Logger integration for system-wide logging
- Custom IP-based rate limiting implemented for API protection
- Swagger API documentation implemented:
  - Added annotations to all handler functions across all domains
  - Added examples to request/response models
  - Created ErrorResponse struct for consistent API documentation
  - Set up Swagger UI endpoint at /swagger/*
  - Added Makefile commands for Swagger generation and formatting

## Next Steps
- Implement user authentication flow (registration, login, JWT)
- Complete book catalog management functionality
- Develop rental creation and management operations
- Implement payment processing system
- Add reporting and analytics functionality
- Implement remaining repository layer implementations
- Expand service layer with business logic
- Add comprehensive unit and integration tests
- Keep Swagger documentation in sync with API changes
- Add more detailed examples to API documentation

## Active Decisions
- Using Clean Architecture with four main layers (domain, repository, service, API)
- Following RESTful API design principles
- Implementing role-based access control (Admin, Librarian, Member)
- Using PostgreSQL for data persistence
- Docker-based development environment
- Interface-first approach for component design
- Graceful server shutdown with context timeout
- Structured logging throughout application
- Environment-based configuration management
- OpenAPI/Swagger for API documentation and exploration

## Technical Patterns
- Domain entities with validation logic
- Interface-driven development for testability
- Repository pattern with complete CRUD operations
- Dependency injection for service composition
- Pagination support for list operations
- Search parameter objects for filtering
- JWT-based authentication
- Middleware chain for request processing
- Structured error handling with domain-specific types
- Context-based timeout management
- In-memory request rate limiting with IP tracking
- Swagger annotation pattern for API documentation
- Example-rich model documentation

## Implementation Guidelines
- All code comments and documentation in English
- No hardcoded values - use configuration
- Comprehensive error handling with domain-specific errors
- Unit tests for all business logic
- Interface mocking for dependency isolation
- Consistent naming conventions across layers
- Graceful shutdown handling for production stability
- Keep Swagger annotations up-to-date with code changes
- Provide examples for all request/response models
- Document security requirements for all endpoints
