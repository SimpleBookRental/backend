# User API Flow Sequence Diagrams

## Get User By ID Flow

```mermaid
sequenceDiagram
    participant C as Client
    participant R as Router (Gin)
    participant M as Middleware
    participant H as UserHandler
    participant S as UserService
    participant UR as UserRepository
    participant DB as Database

    C->>R: GET /api/v1/users/:id
    R->>M: AuthMiddleware
    M->>M: Validate JWT
    M->>H: GetByID
    H->>H: Parse user ID
    H->>H: Check if user requests own profile or is admin
    H->>S: GetByID(id)
    S->>UR: GetByID(id)
    UR->>DB: SELECT FROM users WHERE id = ?
    DB-->>UR: Return user data
    UR-->>S: Return user
    S-->>H: Return user
    H-->>C: HTTP 200 OK with user details
```

## List Users Flow (Admin Only)

```mermaid
sequenceDiagram
    participant C as Client
    participant R as Router (Gin)
    participant M as Middleware
    participant H as UserHandler
    participant S as UserService
    participant UR as UserRepository
    participant DB as Database

    C->>R: GET /api/v1/users?limit=10&offset=0
    R->>M: AuthMiddleware + RoleMiddleware
    M->>M: Validate JWT & check admin role
    M->>H: List
    H->>H: Parse pagination params
    H->>S: List(limit, offset)
    S->>UR: List(limit, offset)
    UR->>DB: SELECT FROM users LIMIT ? OFFSET ?
    DB-->>UR: Return users data
    UR-->>S: Return users
    S-->>H: Return users
    H-->>C: HTTP 200 OK with paginated users
```

## Update User Flow

```mermaid
sequenceDiagram
    participant C as Client
    participant R as Router (Gin)
    participant M as Middleware
    participant H as UserHandler
    participant S as UserService
    participant UR as UserRepository
    participant DB as Database

    C->>R: PUT /api/v1/users/:id
    R->>M: AuthMiddleware
    M->>M: Validate JWT
    M->>H: Update
    H->>H: Parse user ID
    H->>H: Check if user updates own profile or is admin
    H->>H: Validate request body
    H->>S: GetByID(id)
    S->>UR: GetByID(id)
    UR->>DB: SELECT FROM users WHERE id = ?
    DB-->>UR: Return user data
    UR-->>S: Return user
    S-->>H: Return user
    H->>H: Update user fields
    H->>S: Update(user)
    S->>UR: Update(user)
    UR->>DB: UPDATE users SET ... WHERE id = ?
    DB-->>UR: Confirm update
    UR-->>S: Return updated user
    S-->>H: Return updated user
    H-->>C: HTTP 200 OK with updated user
```

## Delete User Flow (Admin Only)

```mermaid
sequenceDiagram
    participant C as Client
    participant R as Router (Gin)
    participant M as Middleware
    participant H as UserHandler
    participant S as UserService
    participant UR as UserRepository
    participant DB as Database

    C->>R: DELETE /api/v1/users/:id
    R->>M: AuthMiddleware + RoleMiddleware
    M->>M: Validate JWT & check admin role
    M->>H: Delete
    H->>H: Parse user ID
    H->>S: Delete(id)
    S->>UR: Delete(id)
    UR->>DB: DELETE FROM users WHERE id = ?
    DB-->>UR: Confirm delete
    UR-->>S: Return success
    S-->>H: Return success
    H-->>C: HTTP 200 OK with success message
