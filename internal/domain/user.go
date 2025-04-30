package domain

import (
	"time"
)

// UserRole defines the role of a user in the system
type UserRole string

const (
	// RoleAdmin represents an administrator user
	RoleAdmin UserRole = "admin"
	// RoleLibrarian represents a librarian user
	RoleLibrarian UserRole = "librarian"
	// RoleMember represents a regular member user
	RoleMember UserRole = "member"
)

// User represents a user in the system
type User struct {
	ID           int64     `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"` // Never expose password hash in JSON responses
	FirstName    string    `json:"first_name,omitempty"`
	LastName     string    `json:"last_name,omitempty"`
	Role         UserRole  `json:"role"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// UserRepository defines the interface for user data access
type UserRepository interface {
	GetByID(id int64) (*User, error)
	GetByUsername(username string) (*User, error)
	GetByEmail(email string) (*User, error)
	List(limit, offset int32) ([]*User, error)
	Create(user *User) (*User, error)
	Update(user *User) (*User, error)
	UpdatePassword(id int64, passwordHash string) error
	Delete(id int64) error
}

// UserService defines the interface for user business logic
type UserService interface {
	GetByID(id int64) (*User, error)
	GetByUsername(username string) (*User, error)
	GetByEmail(email string) (*User, error)
	List(limit, offset int32) ([]*User, error)
	Create(user *User, password string) (*User, error)
	Update(user *User) (*User, error)
	ChangePassword(id int64, currentPassword, newPassword string) error
	Delete(id int64) error
	ValidateCredentials(username, password string) (*User, error)
}
