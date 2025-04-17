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
// @Description User entity
type User struct {
	ID        string    `gorm:"type:uuid;primary_key" json:"id" example:"b1a2c3d4-e5f6-7890-abcd-1234567890ab"` // User ID
	Name      string    `gorm:"size:100;not null" json:"name" example:"John Doe"`                              // User name
	Email     string    `gorm:"size:100;not null;unique" json:"email" example:"john@example.com"`              // User email
	Password  string    `gorm:"size:100;not null" json:"-" swaggerignore:"true"`                               // User password (ignored in docs)
	Role      string    `gorm:"size:20;not null;default:'USER'" json:"role" example:"USER"`                    // User role (ADMIN/USER)
	Books     []*Book   `gorm:"foreignKey:UserID" json:"books,omitempty"`                                      // Books owned by user
	CreatedAt time.Time `json:"created_at" example:"2024-04-17T12:00:00Z"`                                     // Created at
	UpdatedAt time.Time `json:"updated_at" example:"2024-04-17T12:00:00Z"`                                     // Updated at
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
// @Description User create request
type UserCreate struct {
	Name     string `json:"name" binding:"required" example:"John Doe"`           // User name
	Email    string `json:"email" binding:"required,email" example:"john@example.com"` // User email
	Password string `json:"password" binding:"required,min=6" example:"secret123"` // User password
}

// UserUpdate is the request body for updating a user
// @Description User update request
type UserUpdate struct {
	Name     string `json:"name" example:"John Doe"`                              // User name
	Email    string `json:"email" binding:"omitempty,email" example:"john@example.com"` // User email
	Password string `json:"password" binding:"omitempty,min=6" example:"secret123"` // User password
}

// UserLogin is the request body for user login
// @Description User login request
type UserLogin struct {
	Email    string `json:"email" binding:"required,email" example:"john@example.com"` // User email
	Password string `json:"password" binding:"required" example:"secret123"`           // User password
}

// LoginResponse is the response body for user login
// @Description Login response
type LoginResponse struct {
	User         *User  `json:"user"`                        // User info
	AccessToken  string `json:"access_token" example:"jwt-token"` // JWT access token
	RefreshToken string `json:"refresh_token" example:"refresh-token"` // Refresh token
	ExpiresAt    int64  `json:"expires_at" example:"1713345600"` // Expiry timestamp
}
