# Payment API Flow Sequence Diagrams

## Get Payment By ID Flow

```mermaid
sequenceDiagram
    participant C as Client
    participant R as Router (Gin)
    participant M as Middleware
    participant H as PaymentHandler
    participant S as PaymentService
    participant PR as PaymentRepository
    participant DB as Database

    C->>R: GET /api/v1/payments/:id
    R->>M: AuthMiddleware
    M->>M: Validate JWT
    M->>H: GetByID
    H->>H: Parse payment ID
    H->>S: GetByID(id)
    S->>PR: GetByID(id)
    PR->>DB: SELECT FROM payments WHERE id = ?
    DB-->>PR: Return payment data
    PR-->>S: Return payment
    S-->>H: Return payment
    H->>H: Check if user owns payment or is admin/librarian
    H-->>C: HTTP 200 OK with payment details
```

## List Payments Flow (Admin/Librarian)

```mermaid
sequenceDiagram
    participant C as Client
    participant R as Router (Gin)
    participant M as Middleware
    participant H as PaymentHandler
    participant S as PaymentService
    participant PR as PaymentRepository
    participant DB as Database

    C->>R: GET /api/v1/payments?limit=10&offset=0
    R->>M: AuthMiddleware + RoleMiddleware
    M->>M: Validate JWT & check admin role
    M->>H: List
    H->>H: Parse pagination params
    H->>S: List(limit, offset)
    S->>PR: List(limit, offset)
    PR->>DB: SELECT FROM payments LIMIT ? OFFSET ?
    DB-->>PR: Return payments data
    PR-->>S: Return payments
    S-->>H: Return payments
    H-->>C: HTTP 200 OK with paginated payments
```

## List User Payments Flow

```mermaid
sequenceDiagram
    participant C as Client
    participant R as Router (Gin)
    participant M as Middleware
    participant H as PaymentHandler
    participant S as PaymentService
    participant PR as PaymentRepository
    participant DB as Database

    C->>R: GET /api/v1/payments/user/:userId
    R->>M: AuthMiddleware
    M->>M: Validate JWT
    M->>H: ListByUser
    H->>H: Parse user ID and pagination params
    H->>H: Check if user requests own payments or is admin/librarian
    H->>S: ListByUser(userId, limit, offset)
    S->>PR: ListByUser(userId, limit, offset)
    PR->>DB: SELECT FROM payments WHERE user_id = ? LIMIT ? OFFSET ?
    DB-->>PR: Return payments data
    PR-->>S: Return payments
    S-->>H: Return payments
    H-->>C: HTTP 200 OK with paginated payments
```

## Create Payment Flow

```mermaid
sequenceDiagram
    participant C as Client
    participant R as Router (Gin)
    participant M as Middleware
    participant H as PaymentHandler
    participant S as PaymentService
    participant PR as PaymentRepository
    participant DB as Database

    C->>R: POST /api/v1/payments
    R->>M: AuthMiddleware
    M->>M: Validate JWT
    M->>H: Create
    H->>H: Extract userID from JWT context
    H->>H: Validate request body
    H->>S: Create(payment)
    S->>PR: Create(payment)
    PR->>DB: INSERT INTO payments
    DB-->>PR: Return payment ID
    PR-->>S: Return created payment
    S-->>H: Return created payment
    H-->>C: HTTP 201 Created with payment
```

## Process Payment Flow

```mermaid
sequenceDiagram
    participant C as Client
    participant R as Router (Gin)
    participant M as Middleware
    participant H as PaymentHandler
    participant S as PaymentService
    participant PR as PaymentRepository
    participant RR as RentalRepository
    participant DB as Database

    C->>R: POST /api/v1/payments/process
    R->>M: AuthMiddleware
    M->>M: Validate JWT
    M->>H: Process
    H->>H: Extract userID from JWT context
    H->>H: Validate request body
    H->>S: ProcessPayment(payment)
    
    alt Payment for Rental
        S->>RR: GetByID(payment.RentalID)
        RR->>DB: SELECT FROM rentals WHERE id = ?
        DB-->>RR: Return rental data
        RR-->>S: Return rental
        S->>S: Validate rental belongs to user
    end
    
    S->>S: Process payment (external service integration if any)
    S->>S: Update payment status to COMPLETED
    S->>PR: Create(payment)
    PR->>DB: INSERT INTO payments
    DB-->>PR: Return payment ID
    PR-->>S: Return processed payment
    S-->>H: Return processed payment
    H-->>C: HTTP 200 OK with processed payment
```

## Refund Payment Flow

```mermaid
sequenceDiagram
    participant C as Client
    participant R as Router (Gin)
    participant M as Middleware
    participant H as PaymentHandler
    participant S as PaymentService
    participant PR as PaymentRepository
    participant DB as Database

    C->>R: PUT /api/v1/payments/:id/refund
    R->>M: AuthMiddleware + RoleMiddleware
    M->>M: Validate JWT & check librarian/admin role
    M->>H: Refund
    H->>H: Parse payment ID
    H->>S: RefundPayment(id)
    S->>PR: GetByID(id)
    PR->>DB: SELECT FROM payments WHERE id = ?
    DB-->>PR: Return payment data
    PR-->>S: Return payment
    S->>S: Validate payment can be refunded
    S->>S: Process refund (external service integration if any)
    S->>S: Update payment status to REFUNDED
    S->>PR: Update(payment)
    PR->>DB: UPDATE payments SET status = 'refunded' WHERE id = ?
    DB-->>PR: Confirm update
    PR-->>S: Return refunded payment
    S-->>H: Return refunded payment
    H-->>C: HTTP 200 OK with refunded payment
