# API Flow Sequence Diagram

```mermaid
sequenceDiagram
    participant C as Client
    participant R as Router (Gin)
    participant M as Middleware
    participant H as Handler
    participant S as Service
    participant Repo as Repository
    participant DB as Database

    C->>R: HTTP Request
    R->>M: Pass request
    M->>M: Apply middleware (CORS, logging, etc.)
    
    alt Authentication Required
        M->>M: Validate JWT
        
        alt Role Authorization Required
            M->>M: Check user role
        end
    end
    
    M->>H: Forward to appropriate handler
    H->>H: Validate request data
    H->>S: Call service method
    S->>S: Apply business logic
    S->>Repo: Call repository method
    Repo->>DB: Execute database query
    DB-->>Repo: Return data
    Repo-->>S: Return data
    S-->>H: Return result
    H-->>C: HTTP Response with data

    alt Error occurs at any point
        Note over C,DB: Error handling flow
        H-->>C: Return appropriate error response
    end
