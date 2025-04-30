# Book Rental System - Product Requirements Document (PRD)

## Overview
The Book Rental System is a RESTful API service that allows users to browse, borrow, and return books from a virtual library. The system will track book inventory, user accounts, rental history, and fees.

## Technical Stack
- **Backend Language**: Golang
- **Database**: PostgreSQL (running in Docker container)
- **Web Framework**: Gin
- **SQL Interface**: SQLC (for generating type-safe code from SQL)
- **Database Migration**: go-migrate
- **Testing**: GoMock for mocking framework
- **Configuration**: .env file
- **Build System**: Makefile

## Core Features

### 1. User Management
- User registration and authentication
- User profile management
- Role-based access control (Admin, Librarian, Member)

### 2. Book Catalog Management
- Add, update, and delete books
- Book categorization and metadata management
- Search and filter capabilities
- Book availability status

### 3. Rental Operations
- Book borrowing process
- Book return process
- Rental period extension
- Overdue book management

### 4. Fees and Payments
- Rental fee calculation
- Late return penalty calculation
- Payment processing
- Refund handling

### 5. Reporting and Analytics
- Usage statistics
- Popular book reports
- Revenue reports
- User activity reports

## API Endpoints

### Authentication
- `POST /api/v1/auth/register` - Register a new user
- `POST /api/v1/auth/login` - User login
- `POST /api/v1/auth/refresh` - Refresh authentication token
- `POST /api/v1/auth/logout` - User logout

### Users
- `GET /api/v1/users` - Get all users (admin only)
- `GET /api/v1/users/:id` - Get user by ID
- `PUT /api/v1/users/:id` - Update user profile
- `DELETE /api/v1/users/:id` - Delete user (admin only)

### Books
- `GET /api/v1/books` - Get all books with filters and pagination
- `GET /api/v1/books/:id` - Get book details
- `POST /api/v1/books` - Add a new book (admin/librarian only)
- `PUT /api/v1/books/:id` - Update book details (admin/librarian only)
- `DELETE /api/v1/books/:id` - Remove a book (admin/librarian only)
- `GET /api/v1/books/categories` - Get book categories
- `GET /api/v1/books/search` - Search books by various parameters

### Rentals
- `GET /api/v1/rentals` - Get all rentals (admin/librarian only)
- `GET /api/v1/rentals/user/:userId` - Get rentals for specific user
- `POST /api/v1/rentals` - Create a new rental
- `PUT /api/v1/rentals/:id/return` - Process book return
- `PUT /api/v1/rentals/:id/extend` - Extend rental period
- `GET /api/v1/rentals/:id` - Get rental details

### Payments
- `GET /api/v1/payments` - Get all payments (admin only)
- `GET /api/v1/payments/user/:userId` - Get user payments
- `POST /api/v1/payments` - Process a new payment
- `GET /api/v1/payments/:id` - Get payment details

### Reports
- `GET /api/v1/reports/books/popular` - Get popular books report
- `GET /api/v1/reports/revenue` - Get revenue report (admin only)
- `GET /api/v1/reports/overdue` - Get overdue books report

## Database Schema

### Users Table
```sql
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    first_name VARCHAR(100),
    last_name VARCHAR(100),
    role VARCHAR(20) NOT NULL DEFAULT 'member',
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);
```

### Books Table
```sql
CREATE TABLE books (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    author VARCHAR(255) NOT NULL,
    isbn VARCHAR(20) UNIQUE NOT NULL,
    description TEXT,
    published_year INT,
    publisher VARCHAR(255),
    total_copies INT NOT NULL DEFAULT 1,
    available_copies INT NOT NULL DEFAULT 1,
    category_id INT REFERENCES categories(id),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);
```

### Categories Table
```sql
CREATE TABLE categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) UNIQUE NOT NULL,
    description TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);
```

### Rentals Table
```sql
CREATE TABLE rentals (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id),
    book_id INT NOT NULL REFERENCES books(id),
    rental_date TIMESTAMP NOT NULL DEFAULT NOW(),
    due_date TIMESTAMP NOT NULL,
    return_date TIMESTAMP,
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);
```

### Payments Table
```sql
CREATE TABLE payments (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id),
    rental_id INT REFERENCES rentals(id),
    amount DECIMAL(10, 2) NOT NULL,
    payment_date TIMESTAMP NOT NULL DEFAULT NOW(),
    payment_method VARCHAR(50),
    status VARCHAR(20) NOT NULL,
    transaction_id VARCHAR(255),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);
```

## Authentication and Authorization

The API will use JWT (JSON Web Token) for authentication. Each endpoint will have appropriate role-based access controls:
- Public endpoints: registration, login, book browsing
- Member endpoints: rental operations, profile management
- Librarian endpoints: book management, rental approvals
- Admin endpoints: user management, system reports

## Development Workflow

### Environment Setup
1. Clone repository
2. Set up Docker and run PostgreSQL container
3. Configure .env file
4. Run database migrations
5. Start development server

### Build System
The project will use a Makefile with the following commands:
- `make setup` - Set up development environment
- `make docker-db` - Start PostgreSQL container
- `make migrate-up` - Run database migrations
- `make migrate-down` - Rollback database migrations
- `make sqlc` - Generate SQLC code
- `make mock` - Generate mocks for testing
- `make test` - Run unit tests
- `make run` - Start development server
- `make build` - Build production binary

## Testing Strategy

- Unit tests for all core business logic
- Integration tests for API endpoints
- Mock external dependencies using GoMock
- Database tests using a test database instance

## Deployment Considerations

- Docker container for the application
- Database connection pooling
- Environment variable configuration
- Database backup strategy
- API rate limiting
- Logging and monitoring

## Future Enhancements

- Push notifications for due dates
- Book recommendations based on rental history
- Integration with external book APIs
- Mobile application support
- Social features (reviews, ratings)
- E-book rental support

## Project Timeline

1. **Phase 1: Core Setup** (2 weeks)
   - Project structure setup
   - Database schema design and migration
   - Basic API structure and authentication

2. **Phase 2: Book Management** (2 weeks)
   - Book and category endpoints
   - Search functionality
   - Admin book management features

3. **Phase 3: Rental System** (3 weeks)
   - Rental workflow implementation
   - Due date management
   - Fee calculation system

4. **Phase 4: Payment Integration** (2 weeks)
   - Payment processing
   - Receipt generation
   - Financial reporting

5. **Phase 5: Testing and Optimization** (1 week)
   - Comprehensive testing
   - Performance optimization
   - Documentation finalization
