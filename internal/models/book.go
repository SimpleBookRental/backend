package models

import (
	"time"

	"github.com/SimpleBookRental/backend/pkg/utils"
	"gorm.io/gorm"
)

// Book represents a book in the system
type Book struct {
	ID          string    `gorm:"type:uuid;primary_key" json:"id"`
	Title       string    `gorm:"size:200;not null" json:"title"`
	Author      string    `gorm:"size:100;not null" json:"author"`
	ISBN        string    `gorm:"size:20;not null;unique" json:"isbn"`
	Description string    `gorm:"type:text" json:"description"`
	UserID      string    `gorm:"type:uuid;not null" json:"user_id"`
	User        *User     `gorm:"foreignKey:UserID" json:"user,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
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
type BookCreate struct {
	Title       string `json:"title" binding:"required"`
	Author      string `json:"author" binding:"required"`
	ISBN        string `json:"isbn" binding:"required"`
	Description string `json:"description"`
	UserID      string `json:"user_id" binding:"required,uuid"`
}

// BookUpdate is the request body for updating a book
type BookUpdate struct {
	Title       string `json:"title"`
	Author      string `json:"author"`
	ISBN        string `json:"isbn"`
	Description string `json:"description"`
	UserID      string `json:"user_id" binding:"omitempty,uuid"`
}
