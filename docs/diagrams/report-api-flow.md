# Report API Flow Sequence Diagrams

## Get Popular Books Report Flow

```mermaid
sequenceDiagram
    participant C as Client
    participant R as Router (Gin)
    participant M as Middleware
    participant H as ReportHandler
    participant S as ReportService
    participant BR as BookRepository
    participant RR as RentalRepository
    participant DB as Database

    C->>R: GET /api/v1/reports/books/popular?limit=10&offset=0
    R->>M: AuthMiddleware + RoleMiddleware
    M->>M: Validate JWT & check librarian/admin role
    M->>H: GetPopularBooks
    H->>H: Parse pagination params
    H->>S: GetPopularBooks(limit, offset)
    S->>RR: GetMostRentedBooks(limit, offset)
    RR->>DB: SELECT b.*, COUNT(r.book_id) as rent_count FROM books b JOIN rentals r ON b.id = r.book_id GROUP BY b.id ORDER BY rent_count DESC LIMIT ? OFFSET ?
    DB-->>RR: Return books data with rental counts
    RR-->>S: Return popular books
    S-->>H: Return popular books
    H-->>C: HTTP 200 OK with popular books
```

## Get Revenue Report Flow

```mermaid
sequenceDiagram
    participant C as Client
    participant R as Router (Gin)
    participant M as Middleware
    participant H as ReportHandler
    participant S as ReportService
    participant PR as PaymentRepository
    participant DB as Database

    C->>R: GET /api/v1/reports/revenue?start_date=2023-01-01&end_date=2023-12-31
    R->>M: AuthMiddleware + RoleMiddleware
    M->>M: Validate JWT & check admin role
    M->>H: GetRevenueReport
    H->>H: Validate date parameters
    H->>S: GetRevenueReport(startDate, endDate)
    S->>PR: GetRevenueReport(startDate, endDate)
    PR->>DB: SELECT DATE_TRUNC('month', payment_date) as month, SUM(amount) as revenue FROM payments WHERE status = 'completed' AND payment_date >= ? AND payment_date <= ? GROUP BY month ORDER BY month
    DB-->>PR: Return revenue data by month
    PR-->>S: Return revenue report
    S-->>H: Return revenue report
    H-->>C: HTTP 200 OK with revenue report
```

## Get Overdue Books Report Flow

```mermaid
sequenceDiagram
    participant C as Client
    participant R as Router (Gin)
    participant M as Middleware
    participant H as ReportHandler
    participant S as ReportService
    participant RR as RentalRepository
    participant BR as BookRepository
    participant UR as UserRepository
    participant DB as Database

    C->>R: GET /api/v1/reports/overdue?limit=10&offset=0
    R->>M: AuthMiddleware + RoleMiddleware
    M->>M: Validate JWT & check librarian/admin role
    M->>H: GetOverdueBooks
    H->>H: Parse pagination params
    H->>S: GetOverdueBooks(limit, offset)
    S->>RR: GetOverdueRentals(limit, offset)
    RR->>DB: SELECT r.*, b.title, u.username FROM rentals r JOIN books b ON r.book_id = b.id JOIN users u ON r.user_id = u.id WHERE r.status = 'active' AND r.due_date < CURRENT_TIMESTAMP LIMIT ? OFFSET ?
    DB-->>RR: Return overdue rentals data
    RR-->>S: Return overdue rentals
    S-->>H: Return overdue rentals
    H-->>C: HTTP 200 OK with overdue books
