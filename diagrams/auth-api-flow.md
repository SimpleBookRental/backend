# Auth API Flow Sequence Diagrams

## Register Flow

```mermaid
sequenceDiagram
    participant C as Client
    participant R as Router (Gin)
    participant H as AuthHandler
    participant S as AuthService
    participant US as UserService
    participant UR as UserRepository
    participant DB as Database
    participant JWT as JWTService

    C->>R: POST /api/v1/auth/register
    R->>H: Register
    H->>H: Validate request data
    H->>S: Register(user, password)
    S->>S: Hash password
    S->>UR: Create(user)
    UR->>DB: INSERT INTO users
    DB-->>UR: Return user ID
    UR-->>S: Return created user
    S-->>H: Return user data
    H-->>C: HTTP 201 Created
```

## Login Flow

```mermaid
sequenceDiagram
    participant C as Client
    participant R as Router (Gin)
    participant H as AuthHandler
    participant S as AuthService
    participant UR as UserRepository
    participant DB as Database
    participant JWT as JWTService

    C->>R: POST /api/v1/auth/login
    R->>H: Login
    H->>H: Validate credentials
    H->>S: Login(username, password)
    S->>UR: GetByUsername(username)
    UR->>DB: SELECT FROM users WHERE username = ?
    DB-->>UR: Return user data
    UR-->>S: Return user
    S->>S: Verify password hash
    S->>JWT: GenerateTokens(userId, role)
    JWT-->>S: Return access & refresh tokens
    S-->>H: Return tokens
    H-->>C: HTTP 200 OK with tokens
```

## Refresh Token Flow

```mermaid
sequenceDiagram
    participant C as Client
    participant R as Router (Gin)
    participant H as AuthHandler
    participant S as AuthService
    participant JWT as JWTService

    C->>R: POST /api/v1/auth/refresh
    R->>H: RefreshToken
    H->>H: Validate refresh token
    H->>S: RefreshToken(token)
    S->>JWT: ValidateRefreshToken(token)
    JWT-->>S: Return userId, role
    S->>JWT: GenerateTokens(userId, role)
    JWT-->>S: Return new access & refresh tokens
    S-->>H: Return new tokens
    H-->>C: HTTP 200 OK with new tokens
```

## Logout Flow

```mermaid
sequenceDiagram
    participant C as Client
    participant R as Router (Gin)
    participant M as Middleware
    participant H as AuthHandler
    participant S as AuthService
    participant JWT as JWTService

    C->>R: POST /api/v1/auth/logout
    R->>M: AuthMiddleware
    M->>M: Validate JWT
    M->>H: Logout
    H->>S: Logout(token)
    S->>JWT: InvalidateToken(token)
    JWT-->>S: Confirm invalidation
    S-->>H: Return success
    H-->>C: HTTP 200 OK
