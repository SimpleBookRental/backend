package domain

import (
	"time"
)

// RentalStatus defines the status of a rental
type RentalStatus string

const (
	// RentalStatusActive represents an active rental
	RentalStatusActive RentalStatus = "active"
	// RentalStatusReturned represents a returned rental
	RentalStatusReturned RentalStatus = "returned"
	// RentalStatusOverdue represents an overdue rental
	RentalStatusOverdue RentalStatus = "overdue"
)

// Rental represents a book rental in the system
type Rental struct {
	ID           int64        `json:"id"`
	UserID       int64        `json:"user_id"`
	BookID       int64        `json:"book_id"`
	RentalDate   time.Time    `json:"rental_date"`
	DueDate      time.Time    `json:"due_date"`
	ReturnDate   *time.Time   `json:"return_date,omitempty"`
	Status       RentalStatus `json:"status"`
	CreatedAt    time.Time    `json:"created_at"`
	UpdatedAt    time.Time    `json:"updated_at"`
	UserUsername string       `json:"user_username,omitempty"` // For join queries
	BookTitle    string       `json:"book_title,omitempty"`    // For join queries
	BookAuthor   string       `json:"book_author,omitempty"`   // For join queries
}

// RentalRepository defines the interface for rental data access
type RentalRepository interface {
	GetByID(id int64) (*Rental, error)
	List(limit, offset int32) ([]*Rental, error)
	ListByUser(userID int64, limit, offset int32) ([]*Rental, error)
	ListByBook(bookID int64, limit, offset int32) ([]*Rental, error)
	ListActive(limit, offset int32) ([]*Rental, error)
	ListOverdue(limit, offset int32) ([]*Rental, error)
	Create(rental *Rental) (*Rental, error)
	UpdateStatus(id int64, status RentalStatus) (*Rental, error)
	Return(id int64) (*Rental, error)
	Extend(id int64, newDueDate time.Time) (*Rental, error)
	Delete(id int64) error
}

// RentalService defines the interface for rental business logic
type RentalService interface {
	GetByID(id int64) (*Rental, error)
	List(limit, offset int32) ([]*Rental, error)
	ListByUser(userID int64, limit, offset int32) ([]*Rental, error)
	ListByBook(bookID int64, limit, offset int32) ([]*Rental, error)
	ListActive(limit, offset int32) ([]*Rental, error)
	ListOverdue(limit, offset int32) ([]*Rental, error)
	Create(rental *Rental) (*Rental, error)
	Return(id int64) (*Rental, error)
	Extend(id int64, days int) (*Rental, error)
	CalculateLateFee(rental *Rental) (float64, error)
	IsOverdue(rental *Rental) bool
}
