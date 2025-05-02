# System Patterns: Book Rental System

## Architecture
- Clean Architecture implementation with clear dependency layers
- Core domain entities independent of external frameworks
- Dependency injection for service and repository connections
- RESTful API patterns using Gin framework
- Repository pattern for data access abstraction

## Component Structure
- **Domain Layer**: Core business entities and interfaces
  - Well-defined entities (Book, User, Category, Rental, Payment)
  - Clear interface segregation for repositories and services
  - Consistent patterns across domain models
- **Repository Layer**: Data persistence implementations
  - Complete CRUD operations
  - Query methods with filtering capabilities
  - Transaction management
- **Service Layer**: Business logic and use cases
  - Domain-specific operations
  - Business rule enforcement
  - Cross-entity coordination
- **API Layer**: HTTP handlers and middleware
  - Route definitions
  - Request validation
  - Response formatting
  - Authentication verification
- **Main Package**: Application configuration and wiring

## Key Design Patterns
- **Repository Pattern**: Data access abstraction with specialized query methods
- **Dependency Injection**: Service instantiation and composition
- **Middleware Chain**: Request processing pipeline for auth and logging
- **DTO Pattern**: Data transfer between layers with clear separation
- **Factory Pattern**: Object creation with validation
- **Strategy Pattern**: Fee calculation strategies for different scenarios
- **Status Enum Pattern**: String-based enums for entity states (RentalStatus, etc.)

## Critical Paths
- **Authentication Flow**: Registration → Login → JWT Issuance → Protected Routes
- **Rental Flow**: Book Selection → Availability Check → Rental Creation → Due Date Assignment
- **Return Flow**: Book Return → Fee Calculation → Payment Processing → Inventory Update
- **Payment Flow**: Fee Determination → Payment Method Selection → Transaction Processing → Receipt Generation
- **Book Management Flow**: Creation → Categorization → Inventory Tracking → Availability Updates

## Error Handling Strategy
- Domain-specific error types
- Consistent error responses
- Proper HTTP status code mapping
- Transactional boundaries for data operations
- Validation at entity creation/update points

## Testing Approach
- Mocked interfaces for testability
- Separation of unit and integration tests
- Repository abstraction for test data management
- Service-level business logic testing
- API endpoint integration testing
