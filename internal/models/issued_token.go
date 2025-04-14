package models

import (
	"time"

	"github.com/SimpleBookRental/backend/pkg/utils"
	"gorm.io/gorm"
)

// IssuedToken represents an issued token in the system
type IssuedToken struct {
	ID        string     `gorm:"type:uuid;primary_key" json:"id"`
	UserID    string     `gorm:"type:uuid;not null;index" json:"user_id"`
	User      *User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Token     string     `gorm:"size:500;not null;uniqueIndex" json:"token"`
	TokenType string     `gorm:"size:20;not null" json:"token_type"` // "access" or "refresh"
	ExpiresAt time.Time  `gorm:"not null" json:"expires_at"`
	IsRevoked bool       `gorm:"default:false" json:"is_revoked"`
	RevokedAt *time.Time `json:"revoked_at,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

// TableName overrides the table name
func (IssuedToken) TableName() string {
	return "br_issued_token"
}

// BeforeCreate is called before creating an issued token
func (t *IssuedToken) BeforeCreate(tx *gorm.DB) error {
	if t.ID == "" {
		t.ID = utils.GenerateUUID()
	}
	return nil
}

// LogoutRequest is the request body for user logout
type LogoutRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// TokenType defines the type of token
type TokenType string

// Token types
const (
	AccessToken  TokenType = "access"
	RefreshToken TokenType = "refresh"
)
