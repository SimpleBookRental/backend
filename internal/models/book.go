package models

import (
	"time"

	"github.com/SimpleBookRental/backend/pkg/utils"
	"gorm.io/gorm"
)

// Book represents a book in the system
// @Description Book entity
type Book struct {
	ID          string    `gorm:"type:uuid;primary_key" json:"id" example:"b1a2c3d4-e5f6-7890-abcd-1234567890ab"` // Book ID
	Title       string    `gorm:"size:200;not null" json:"title" example:"The Great Gatsby"`                      // Book title
	Author      string    `gorm:"size:100;not null" json:"author" example:"F. Scott Fitzgerald"`                  // Book author
	ISBN        string    `gorm:"size:20;not null;unique" json:"isbn" example:"978-1234567890"`                   // Book ISBN
	Description string    `gorm:"type:text" json:"description" example:"A classic novel"`                         // Book description
	UserID      string    `gorm:"type:uuid;not null" json:"user_id" example:"b1a2c3d4-e5f6-7890-abcd-1234567890ab"` // Owner user ID
	User        *User     `gorm:"foreignKey:UserID" json:"user,omitempty"`                                        // Owner user
	CreatedAt   time.Time `json:"created_at" example:"2024-04-17T12:00:00Z"`                                      // Created at
	UpdatedAt   time.Time `json:"updated_at" example:"2024-04-17T12:00:00Z"`                                      // Updated at
}

// TableName overrides the table name
func (Book) TableName() string {
	return "br_book"
}

// BeforeCreate is called before creating a book
func (b *Book) BeforeCreate(tx *gorm.DB) error {
	if b.ID == "" {
		b.ID = utils.GenerateUUID()
	}
	return nil
}

// BookCreate is the request body for creating a book
// @Description Book create request
type BookCreate struct {
	Title       string `json:"title" binding:"required" example:"The Great Gatsby"`         // Book title
	Author      string `json:"author" binding:"required" example:"F. Scott Fitzgerald"`     // Book author
	ISBN        string `json:"isbn" binding:"required" example:"978-1234567890"`            // Book ISBN
	Description string `json:"description" example:"A classic novel"`                       // Book description
	UserID      string `json:"user_id" binding:"omitempty,uuid" example:"b1a2c3d4-e5f6-7890-abcd-1234567890ab"` // Owner user ID
}

// BookUpdate is the request body for updating a book
// @Description Book update request
type BookUpdate struct {
	Title       string `json:"title" example:"The Great Gatsby"`         // Book title
	Author      string `json:"author" example:"F. Scott Fitzgerald"`     // Book author
	ISBN        string `json:"isbn" example:"978-1234567890"`            // Book ISBN
	Description string `json:"description" example:"A classic novel"`    // Book description
	UserID      string `json:"user_id" binding:"omitempty,uuid" example:"b1a2c3d4-e5f6-7890-abcd-1234567890ab"` // Owner user ID
}
