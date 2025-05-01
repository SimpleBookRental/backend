# Rental API Flow Sequence Diagrams

## Get Rental By ID Flow

```mermaid
sequenceDiagram
    participant C as Client
    participant R as Router (Gin)
    participant M as Middleware
    participant H as RentalHandler
    participant S as RentalService
    participant RR as RentalRepository
    participant DB as Database

    C->>R: GET /api/v1/rentals/:id
    R->>M: AuthMiddleware
    M->>M: Validate JWT
    M->>H: GetByID
    H->>H: Parse rental ID
    H->>S: GetByID(id)
    S->>RR: GetByID(id)
    RR->>DB: SELECT FROM rentals WHERE id = ?
    DB-->>RR: Return rental data
    RR-->>S: Return rental
    S-->>H: Return rental
    H->>H: Check if user owns rental or is admin/librarian
    H-->>C: HTTP 200 OK with rental details
```

## List Rentals Flow (Admin/Librarian)

```mermaid
sequenceDiagram
    participant C as Client
    participant R as Router (Gin)
    participant M as Middleware
    participant H as RentalHandler
    participant S as RentalService
    participant RR as RentalRepository
    participant DB as Database

    C->>R: GET /api/v1/rentals?limit=10&offset=0
    R->>M: AuthMiddleware + RoleMiddleware
    M->>M: Validate JWT & check role
    M->>H: List
    H->>H: Parse pagination params
    H->>S: List(limit, offset)
    S->>RR: List(limit, offset)
    RR->>DB: SELECT FROM rentals LIMIT ? OFFSET ?
    DB-->>RR: Return rentals data
    RR-->>S: Return rentals
    S-->>H: Return rentals
    H-->>C: HTTP 200 OK with paginated rentals
```

## List User Rentals Flow

```mermaid
sequenceDiagram
    participant C as Client
    participant R as Router (Gin)
    participant M as Middleware
    participant H as RentalHandler
    participant S as RentalService
    participant RR as RentalRepository
    participant DB as Database

    C->>R: GET /api/v1/rentals/user/:userId
    R->>M: AuthMiddleware
    M->>M: Validate JWT
    M->>H: ListByUser
    H->>H: Parse user ID and pagination params
    H->>H: Check if user requests own rentals or is admin/librarian
    H->>S: ListByUser(userId, limit, offset)
    S->>RR: ListByUser(userId, limit, offset)
    RR->>DB: SELECT FROM rentals WHERE user_id = ? LIMIT ? OFFSET ?
    DB-->>RR: Return rentals data
    RR-->>S: Return rentals
    S-->>H: Return rentals
    H-->>C: HTTP 200 OK with paginated rentals
```

## Create Rental Flow

```mermaid
sequenceDiagram
    participant C as Client
    participant R as Router (Gin)
    participant M as Middleware
    participant H as RentalHandler
    participant S as RentalService
    participant RR as RentalRepository
    participant BR as BookRepository
    participant DB as Database

    C->>R: POST /api/v1/rentals
    R->>M: AuthMiddleware
    M->>M: Validate JWT
    M->>H: Create
    H->>H: Extract userID from JWT context
    H->>H: Validate request body
    H->>S: Create(rental)
    S->>BR: GetByID(rental.BookID)
    BR->>DB: SELECT FROM books WHERE id = ?
    DB-->>BR: Return book data
    BR-->>S: Return book
    S->>S: Check book availability
    S->>S: Update book available copies
    S->>BR: Update(book)
    BR->>DB: UPDATE books SET available_copies = ? WHERE id = ?
    DB-->>BR: Confirm update
    S->>RR: Create(rental)
    RR->>DB: INSERT INTO rentals
    DB-->>RR: Return rental ID
    RR-->>S: Return created rental
    S-->>H: Return created rental
    H-->>C: HTTP 201 Created with rental
```

## Return Rental Flow

```mermaid
sequenceDiagram
    participant C as Client
    participant R as Router (Gin)
    participant M as Middleware
    participant H as RentalHandler
    participant S as RentalService
    participant RR as RentalRepository
    participant BR as BookRepository
    participant DB as Database

    C->>R: PUT /api/v1/rentals/:id/return
    R->>M: AuthMiddleware
    M->>M: Validate JWT
    M->>H: Return
    H->>H: Parse rental ID
    H->>S: GetByID(id)
    S->>RR: GetByID(id)
    RR->>DB: SELECT FROM rentals WHERE id = ?
    DB-->>RR: Return rental data
    RR-->>S: Return rental
    S-->>H: Return rental
    H->>H: Check if user owns rental or is admin/librarian
    H->>S: Return(id)
    S->>RR: GetByID(id)
    RR->>DB: SELECT FROM rentals WHERE id = ?
    DB-->>RR: Return rental data
    RR-->>S: Return rental
    S->>BR: GetByID(rental.BookID)
    BR->>DB: SELECT FROM books WHERE id = ?
    DB-->>BR: Return book data
    BR-->>S: Return book
    S->>S: Update book available copies
    S->>BR: Update(book)
    BR->>DB: UPDATE books SET available_copies = ? WHERE id = ?
    DB-->>BR: Confirm update
    S->>S: Set return date and update status
    S->>RR: Update(rental)
    RR->>DB: UPDATE rentals SET return_date = ?, status = ? WHERE id = ?
    DB-->>RR: Confirm update
    RR-->>S: Return updated rental
    S->>S: CalculateLateFee(rental)
    S-->>H: Return updated rental with late fee
    H-->>C: HTTP 200 OK with rental and late fee
```

## Extend Rental Flow

```mermaid
sequenceDiagram
    participant C as Client
    participant R as Router (Gin)
    participant M as Middleware
    participant H as RentalHandler
    participant S as RentalService
    participant RR as RentalRepository
    participant DB as Database

    C->>R: PUT /api/v1/rentals/:id/extend
    R->>M: AuthMiddleware
    M->>M: Validate JWT
    M->>H: Extend
    H->>H: Parse rental ID
    H->>S: GetByID(id)
    S->>RR: GetByID(id)
    RR->>DB: SELECT FROM rentals WHERE id = ?
    DB-->>RR: Return rental data
    RR-->>S: Return rental
    S-->>H: Return rental
    H->>H: Check if user owns rental or is admin/librarian
    H->>H: Validate request body (days to extend)
    H->>S: Extend(id, days)
    S->>RR: GetByID(id)
    RR->>DB: SELECT FROM rentals WHERE id = ?
    DB-->>RR: Return rental data
    RR-->>S: Return rental
    S->>S: Check if rental can be extended
    S->>S: Calculate new due date
    S->>RR: Update(rental)
    RR->>DB: UPDATE rentals SET due_date = ? WHERE id = ?
    DB-->>RR: Confirm update
    RR-->>S: Return updated rental
    S-->>H: Return updated rental
    H-->>C: HTTP 200 OK with updated rental
