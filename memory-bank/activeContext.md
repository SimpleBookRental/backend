# Active Context: Book Rental System

## Current Focus
- Project initialization and structure setup
- Core domain entity definition and interface design
- Repository and service interface implementation
- API endpoint handlers development
- Database schema implementation

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

## Next Steps
- Implement user authentication flow (registration, login, JWT)
- Complete book catalog management functionality
- Develop rental creation and management operations
- Implement payment processing system
- Add reporting and analytics functionality

## Active Decisions
- Using Clean Architecture with four main layers (domain, repository, service, API)
- Following RESTful API design principles
- Implementing role-based access control (Admin, Librarian, Member)
- Using PostgreSQL for data persistence
- Docker-based development environment
- Interface-first approach for component design

## Technical Patterns
- Domain entities with validation logic
- Interface-driven development for testability
- Repository pattern with complete CRUD operations
- Dependency injection for service composition
- Pagination support for list operations
- Search parameter objects for filtering
- JWT-based authentication
- Middleware chain for request processing

## Implementation Guidelines
- All code comments and documentation in English
- No hardcoded values - use configuration
- Comprehensive error handling with domain-specific errors
- Unit tests for all business logic
- Interface mocking for dependency isolation
- Consistent naming conventions across layers
