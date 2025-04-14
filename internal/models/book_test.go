package models

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestBookModel(t *testing.T) {
	// Test case 1: Create a new book
	book := &Book{
		ID:          "test-id",
		Title:       "Test Book",
		Author:      "Test Author",
		ISBN:        "ISBN-123",
		Description: "Test description",
		UserID:      "user-id",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Assert
	assert.Equal(t, "test-id", book.ID)
	assert.Equal(t, "Test Book", book.Title)
	assert.Equal(t, "Test Author", book.Author)
	assert.Equal(t, "ISBN-123", book.ISBN)
	assert.Equal(t, "Test description", book.Description)
	assert.Equal(t, "user-id", book.UserID)
	assert.NotZero(t, book.CreatedAt)
	assert.NotZero(t, book.UpdatedAt)
}

func TestBookBeforeCreate(t *testing.T) {
	// Test case 1: BeforeCreate with empty ID
	book := &Book{
		Title:       "Test Book",
		Author:      "Test Author",
		ISBN:        "ISBN-123",
		Description: "Test description",
		UserID:      "user-id",
	}

	err := book.BeforeCreate(&gorm.DB{})
	assert.NoError(t, err)
	assert.NotEmpty(t, book.ID)

	// Test case 2: BeforeCreate with existing ID
	book = &Book{
		ID:          "test-id",
		Title:       "Test Book",
		Author:      "Test Author",
		ISBN:        "ISBN-123",
		Description: "Test description",
		UserID:      "user-id",
	}

	err = book.BeforeCreate(&gorm.DB{})
	assert.NoError(t, err)
	assert.Equal(t, "test-id", book.ID)
}

func TestBookTableName(t *testing.T) {
	// Test case 1: Check table name
	book := Book{}
	assert.Equal(t, "br_book", book.TableName())
}

func TestBookCreate(t *testing.T) {
	// Test case 1: Valid book create
	bookCreate := &BookCreate{
		Title:       "Test Book",
		Author:      "Test Author",
		ISBN:        "ISBN-123",
		Description: "Test description",
		UserID:      "user-id",
	}

	// Assert
	assert.Equal(t, "Test Book", bookCreate.Title)
	assert.Equal(t, "Test Author", bookCreate.Author)
	assert.Equal(t, "ISBN-123", bookCreate.ISBN)
	assert.Equal(t, "Test description", bookCreate.Description)
	assert.Equal(t, "user-id", bookCreate.UserID)
}

func TestBookUpdate(t *testing.T) {
	// Test case 1: Valid book update
	bookUpdate := &BookUpdate{
		Title:       "Updated Book",
		Author:      "Updated Author",
		ISBN:        "ISBN-456",
		Description: "Updated description",
	}

	// Assert
	assert.Equal(t, "Updated Book", bookUpdate.Title)
	assert.Equal(t, "Updated Author", bookUpdate.Author)
	assert.Equal(t, "ISBN-456", bookUpdate.ISBN)
	assert.Equal(t, "Updated description", bookUpdate.Description)
}
