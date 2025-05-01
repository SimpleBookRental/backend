# Book API Flow Sequence Diagrams

## List Books Flow

```mermaid
sequenceDiagram
    participant C as Client
    participant R as Router (Gin)
    participant H as BookHandler
    participant S as BookService
    participant BR as BookRepository
    participant DB as Database

    C->>R: GET /api/v1/books?limit=10&offset=0
    R->>H: List
    H->>H: Parse pagination params
    H->>S: List(limit, offset)
    S->>BR: List(limit, offset)
    BR->>DB: SELECT FROM books LIMIT ? OFFSET ?
    DB-->>BR: Return books data
    BR-->>S: Return books
    S-->>H: Return books
    H-->>C: HTTP 200 OK with paginated books
```

## Get Book By ID Flow

```mermaid
sequenceDiagram
    participant C as Client
    participant R as Router (Gin)
    participant H as BookHandler
    participant S as BookService
    participant BR as BookRepository
    participant DB as Database

    C->>R: GET /api/v1/books/:id
    R->>H: GetByID
    H->>H: Parse book ID
    H->>S: GetByID(id)
    S->>BR: GetByID(id)
    BR->>DB: SELECT FROM books WHERE id = ?
    DB-->>BR: Return book data
    BR-->>S: Return book
    S-->>H: Return book
    H-->>C: HTTP 200 OK with book details
```

## Search Books Flow

```mermaid
sequenceDiagram
    participant C as Client
    participant R as Router (Gin)
    participant H as BookHandler
    participant S as BookService
    participant BR as BookRepository
    participant DB as Database

    C->>R: GET /api/v1/books/search?title=&author=&...
    R->>H: Search
    H->>H: Parse search params
    H->>S: Search(params)
    S->>BR: Search(params)
    BR->>DB: SELECT FROM books WHERE conditions
    DB-->>BR: Return matching books
    BR-->>S: Return books
    S-->>H: Return books
    H-->>C: HTTP 200 OK with search results
```

## Create Book Flow

```mermaid
sequenceDiagram
    participant C as Client
    participant R as Router (Gin)
    participant M as Middleware
    participant H as BookHandler
    participant S as BookService
    participant BR as BookRepository
    participant DB as Database

    C->>R: POST /api/v1/books
    R->>M: AuthMiddleware + RoleMiddleware
    M->>M: Validate JWT & role
    M->>H: Create
    H->>H: Validate request body
    H->>S: Create(book)
    S->>BR: Create(book)
    BR->>DB: INSERT INTO books
    DB-->>BR: Return book ID
    BR-->>S: Return created book
    S-->>H: Return created book
    H-->>C: HTTP 201 Created with book
```

## Update Book Flow

```mermaid
sequenceDiagram
    participant C as Client
    participant R as Router (Gin)
    participant M as Middleware
    participant H as BookHandler
    participant S as BookService
    participant BR as BookRepository
    participant DB as Database

    C->>R: PUT /api/v1/books/:id
    R->>M: AuthMiddleware + RoleMiddleware
    M->>M: Validate JWT & role
    M->>H: Update
    H->>H: Validate request body
    H->>S: GetByID(id)
    S->>BR: GetByID(id)
    BR->>DB: SELECT FROM books WHERE id = ?
    DB-->>BR: Return book data
    BR-->>S: Return book
    S-->>H: Return book
    H->>H: Update book fields
    H->>S: Update(book)
    S->>BR: Update(book)
    BR->>DB: UPDATE books SET ... WHERE id = ?
    DB-->>BR: Confirm update
    BR-->>S: Return updated book
    S-->>H: Return updated book
    H-->>C: HTTP 200 OK with updated book
```

## Update Book Copies Flow

```mermaid
sequenceDiagram
    participant C as Client
    participant R as Router (Gin)
    participant M as Middleware
    participant H as BookHandler
    participant S as BookService
    participant BR as BookRepository
    participant DB as Database

    C->>R: PUT /api/v1/books/:id/copies
    R->>M: AuthMiddleware + RoleMiddleware
    M->>M: Validate JWT & role
    M->>H: UpdateCopies
    H->>H: Validate request body
    H->>S: UpdateCopies(id, totalCopies, availableCopies)
    S->>BR: UpdateCopies(id, totalCopies, availableCopies)
    BR->>DB: UPDATE books SET total_copies = ?, available_copies = ? WHERE id = ?
    DB-->>BR: Confirm update
    BR-->>S: Return updated book
    S-->>H: Return updated book
    H-->>C: HTTP 200 OK with updated book
```

## Delete Book Flow

```mermaid
sequenceDiagram
    participant C as Client
    participant R as Router (Gin)
    participant M as Middleware
    participant H as BookHandler
    participant S as BookService
    participant BR as BookRepository
    participant DB as Database

    C->>R: DELETE /api/v1/books/:id
    R->>M: AuthMiddleware + RoleMiddleware
    M->>M: Validate JWT & role
    M->>H: Delete
    H->>H: Parse book ID
    H->>S: Delete(id)
    S->>BR: Delete(id)
    BR->>DB: DELETE FROM books WHERE id = ?
    DB-->>BR: Confirm delete
    BR-->>S: Return success
    S-->>H: Return success
    H-->>C: HTTP 200 OK with success message
