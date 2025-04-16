package models

import (
	"time"

	"github.com/SimpleBookRental/backend/pkg/utils"
	"gorm.io/gorm"
)

// Role types
const (
	AdminRole string = "ADMIN"
	UserRole  string = "USER"
)

// User represents a user in the system
type User struct {
	ID        string    `gorm:"type:uuid;primary_key" json:"id"`
	Name      string    `gorm:"size:100;not null" json:"name"`
	Email     string    `gorm:"size:100;not null;unique" json:"email"`
	Password  string    `gorm:"size:100;not null" json:"-"`
	Role      string    `gorm:"size:20;not null;default:'USER'" json:"role"`
	Books     []*Book   `gorm:"foreignKey:UserID" json:"books,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TableName overrides the table name
func (User) TableName() string {
	return "br_user"
}

// BeforeCreate is called before creating a user
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == "" {
		u.ID = utils.GenerateUUID()
	}

	// Set default role if not provided
	if u.Role == "" {
		u.Role = UserRole
	}

	return nil
}

// UserCreate is the request body for creating a user
type UserCreate struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// UserUpdate is the request body for updating a user
type UserUpdate struct {
	Name     string `json:"name"`
	Email    string `json:"email" binding:"omitempty,email"`
	Password string `json:"password" binding:"omitempty,min=6"`
}

// UserLogin is the request body for user login
type UserLogin struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse is the response body for user login
type LoginResponse struct {
	User         *User  `json:"user"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresAt    int64  `json:"expires_at"`
}
