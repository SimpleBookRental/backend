package domain

import (
	"time"
)

// PaymentStatus defines the status of a payment
type PaymentStatus string

const (
	// PaymentStatusPending represents a pending payment
	PaymentStatusPending PaymentStatus = "pending"
	// PaymentStatusCompleted represents a completed payment
	PaymentStatusCompleted PaymentStatus = "completed"
	// PaymentStatusFailed represents a failed payment
	PaymentStatusFailed PaymentStatus = "failed"
	// PaymentStatusRefunded represents a refunded payment
	PaymentStatusRefunded PaymentStatus = "refunded"
)

// Payment represents a payment in the system
type Payment struct {
	ID            int64         `json:"id"`
	UserID        int64         `json:"user_id"`
	RentalID      *int64        `json:"rental_id,omitempty"`
	Amount        float64       `json:"amount"`
	PaymentDate   time.Time     `json:"payment_date"`
	PaymentMethod string        `json:"payment_method,omitempty"`
	Status        PaymentStatus `json:"status"`
	TransactionID string        `json:"transaction_id,omitempty"`
	CreatedAt     time.Time     `json:"created_at"`
	UpdatedAt     time.Time     `json:"updated_at"`
	UserUsername  string        `json:"user_username,omitempty"` // For join queries
	BookTitle     string        `json:"book_title,omitempty"`    // For join queries
}

// RevenueReport represents a revenue report entry
type RevenueReport struct {
	Month        time.Time `json:"month"`
	TotalRevenue float64   `json:"total_revenue"`
	PaymentCount int64     `json:"payment_count"`
}

// PaymentRepository defines the interface for payment data access
type PaymentRepository interface {
	GetByID(id int64) (*Payment, error)
	List(limit, offset int32) ([]*Payment, error)
	ListByUser(userID int64, limit, offset int32) ([]*Payment, error)
	ListByRental(rentalID int64) ([]*Payment, error)
	Create(payment *Payment) (*Payment, error)
	UpdateStatus(id int64, status PaymentStatus) (*Payment, error)
	Delete(id int64) error
	GetRevenueReport(startDate, endDate time.Time) ([]*RevenueReport, error)
}

// PaymentService defines the interface for payment business logic
type PaymentService interface {
	GetByID(id int64) (*Payment, error)
	List(limit, offset int32) ([]*Payment, error)
	ListByUser(userID int64, limit, offset int32) ([]*Payment, error)
	ListByRental(rentalID int64) ([]*Payment, error)
	Create(payment *Payment) (*Payment, error)
	ProcessPayment(payment *Payment) (*Payment, error)
	RefundPayment(id int64) (*Payment, error)
	GetRevenueReport(startDate, endDate time.Time) ([]*RevenueReport, error)
}
