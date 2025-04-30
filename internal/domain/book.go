package domain

import (
	"time"
)

// Book represents a book in the system
type Book struct {
	ID              int64     `json:"id"`
	Title           string    `json:"title"`
	Author          string    `json:"author"`
	ISBN            string    `json:"isbn"`
	Description     string    `json:"description,omitempty"`
	PublishedYear   int32     `json:"published_year,omitempty"`
	Publisher       string    `json:"publisher,omitempty"`
	TotalCopies     int32     `json:"total_copies"`
	AvailableCopies int32     `json:"available_copies"`
	CategoryID      int64     `json:"category_id,omitempty"`
	CategoryName    string    `json:"category_name,omitempty"` // For join queries
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// BookSearchParams represents parameters for searching books
type BookSearchParams struct {
	Title         string `json:"title,omitempty"`
	Author        string `json:"author,omitempty"`
	ISBN          string `json:"isbn,omitempty"`
	PublishedYear int32  `json:"published_year,omitempty"`
	CategoryID    int64  `json:"category_id,omitempty"`
	Available     bool   `json:"available,omitempty"`
	Limit         int32  `json:"limit,omitempty"`
	Offset        int32  `json:"offset,omitempty"`
}

// BookRepository defines the interface for book data access
type BookRepository interface {
	GetByID(id int64) (*Book, error)
	GetByISBN(isbn string) (*Book, error)
	List(limit, offset int32) ([]*Book, error)
	ListByCategory(categoryID int64, limit, offset int32) ([]*Book, error)
	Search(params BookSearchParams) ([]*Book, error)
	Create(book *Book) (*Book, error)
	Update(book *Book) (*Book, error)
	UpdateCopies(id int64, totalCopies, availableCopies int32) (*Book, error)
	DecrementAvailableCopies(id int64) (*Book, error)
	IncrementAvailableCopies(id int64) (*Book, error)
	Delete(id int64) error
}

// BookService defines the interface for book business logic
type BookService interface {
	GetByID(id int64) (*Book, error)
	GetByISBN(isbn string) (*Book, error)
	List(limit, offset int32) ([]*Book, error)
	ListByCategory(categoryID int64, limit, offset int32) ([]*Book, error)
	Search(params BookSearchParams) ([]*Book, error)
	Create(book *Book) (*Book, error)
	Update(book *Book) (*Book, error)
	UpdateCopies(id int64, totalCopies, availableCopies int32) (*Book, error)
	Delete(id int64) error
	IsAvailable(id int64) (bool, error)
}
