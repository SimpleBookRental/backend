# Active Context: Book Rental System

## Current Focus
- Integration testing setup and improvements
- Configuration management for test environments
- Health check endpoints for basic connectivity testing
- Database connection configuration for local development vs Docker environments

## Recent Changes
- Added health check endpoint ("/ping") for basic API connectivity testing
- Created configuration tests to verify environment variable loading
- Updated database host configuration in .env file (localhost to 127.0.0.1)
- Fixed integration test setup for main API server
- Implemented proper test database configuration validation
- Project initialization and structure setup
- Core domain entity definition and interface design
- Repository and service interface implementation
- API endpoint handlers development
- Database schema implementation
- Authentication system implementation with JWT
- API documentation with Swagger

## Next Steps
- Fix configuration loading in test environment
- Address database connection issues in integration tests
- Complete remaining integration tests
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
- Integration tests should run with local database (127.0.0.1) instead of Docker service name
- Health check endpoints should be available without authentication
- Environment-specific configuration loading for tests vs production
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
- Health check endpoints for basic connectivity verification
- Environment-specific configuration loading
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
- Integration tests should be independent of the main application configuration
