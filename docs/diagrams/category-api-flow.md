# Category API Flow Sequence Diagrams

## Get Category By ID Flow

```mermaid
sequenceDiagram
    participant C as Client
    participant R as Router (Gin)
    participant H as CategoryHandler
    participant S as CategoryService
    participant CR as CategoryRepository
    participant DB as Database

    C->>R: GET /api/v1/categories/:id
    R->>H: GetByID
    H->>H: Parse category ID
    H->>S: GetByID(id)
    S->>CR: GetByID(id)
    CR->>DB: SELECT FROM categories WHERE id = ?
    DB-->>CR: Return category data
    CR-->>S: Return category
    S-->>H: Return category
    H-->>C: HTTP 200 OK with category details
```

## List Categories Flow

```mermaid
sequenceDiagram
    participant C as Client
    participant R as Router (Gin)
    participant H as CategoryHandler
    participant S as CategoryService
    participant CR as CategoryRepository
    participant DB as Database

    C->>R: GET /api/v1/categories?limit=10&offset=0
    R->>H: List
    H->>H: Parse pagination params
    H->>S: List(limit, offset)
    S->>CR: List(limit, offset)
    CR->>DB: SELECT FROM categories LIMIT ? OFFSET ?
    DB-->>CR: Return categories data
    CR-->>S: Return categories
    S-->>H: Return categories
    H-->>C: HTTP 200 OK with paginated categories
```

## List All Categories Flow

```mermaid
sequenceDiagram
    participant C as Client
    participant R as Router (Gin)
    participant H as CategoryHandler
    participant S as CategoryService
    participant CR as CategoryRepository
    participant DB as Database

    C->>R: GET /api/v1/categories/all
    R->>H: ListAll
    H->>S: ListAll()
    S->>CR: ListAll()
    CR->>DB: SELECT FROM categories
    DB-->>CR: Return all categories data
    CR-->>S: Return all categories
    S-->>H: Return all categories
    H-->>C: HTTP 200 OK with all categories
```

## Create Category Flow

```mermaid
sequenceDiagram
    participant C as Client
    participant R as Router (Gin)
    participant M as Middleware
    participant H as CategoryHandler
    participant S as CategoryService
    participant CR as CategoryRepository
    participant DB as Database

    C->>R: POST /api/v1/categories
    R->>M: AuthMiddleware + RoleMiddleware
    M->>M: Validate JWT & check librarian/admin role
    M->>H: Create
    H->>H: Validate request body
    H->>S: Create(category)
    S->>CR: Create(category)
    CR->>DB: INSERT INTO categories
    DB-->>CR: Return category ID
    CR-->>S: Return created category
    S-->>H: Return created category
    H-->>C: HTTP 201 Created with category
```

## Update Category Flow

```mermaid
sequenceDiagram
    participant C as Client
    participant R as Router (Gin)
    participant M as Middleware
    participant H as CategoryHandler
    participant S as CategoryService
    participant CR as CategoryRepository
    participant DB as Database

    C->>R: PUT /api/v1/categories/:id
    R->>M: AuthMiddleware + RoleMiddleware
    M->>M: Validate JWT & check librarian/admin role
    M->>H: Update
    H->>H: Parse category ID
    H->>H: Validate request body
    H->>S: GetByID(id)
    S->>CR: GetByID(id)
    CR->>DB: SELECT FROM categories WHERE id = ?
    DB-->>CR: Return category data
    CR-->>S: Return category
    S-->>H: Return category
    H->>H: Update category fields
    H->>S: Update(category)
    S->>CR: Update(category)
    CR->>DB: UPDATE categories SET name = ?, description = ? WHERE id = ?
    DB-->>CR: Confirm update
    CR-->>S: Return updated category
    S-->>H: Return updated category
    H-->>C: HTTP 200 OK with updated category
```

## Delete Category Flow

```mermaid
sequenceDiagram
    participant C as Client
    participant R as Router (Gin)
    participant M as Middleware
    participant H as CategoryHandler
    participant S as CategoryService
    participant CR as CategoryRepository
    participant DB as Database

    C->>R: DELETE /api/v1/categories/:id
    R->>M: AuthMiddleware + RoleMiddleware
    M->>M: Validate JWT & check librarian/admin role
    M->>H: Delete
    H->>H: Parse category ID
    H->>S: Delete(id)
    S->>CR: Delete(id)
    CR->>DB: DELETE FROM categories WHERE id = ?
    DB-->>CR: Confirm delete
    CR-->>S: Return success
    S-->>H: Return success
    H-->>C: HTTP 200 OK with success message
