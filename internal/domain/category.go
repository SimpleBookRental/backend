package domain

import (
	"time"
)

// Category represents a book category in the system
type Category struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// CategoryRepository defines the interface for category data access
type CategoryRepository interface {
	GetByID(id int64) (*Category, error)
	GetByName(name string) (*Category, error)
	List(limit, offset int32) ([]*Category, error)
	ListAll() ([]*Category, error)
	Create(category *Category) (*Category, error)
	Update(category *Category) (*Category, error)
	Delete(id int64) error
}

// CategoryService defines the interface for category business logic
type CategoryService interface {
	GetByID(id int64) (*Category, error)
	GetByName(name string) (*Category, error)
	List(limit, offset int32) ([]*Category, error)
	ListAll() ([]*Category, error)
	Create(category *Category) (*Category, error)
	Update(category *Category) (*Category, error)
	Delete(id int64) error
}
