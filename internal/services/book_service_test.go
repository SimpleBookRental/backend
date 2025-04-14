package services

import (
	"errors"
	"testing"

	"github.com/SimpleBookRental/backend/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestBookService_Create(t *testing.T) {
	// Setup
	mockBookRepo := new(MockBookRepository)
	mockUserRepo := new(MockUserRepository)
	service := NewBookService(mockBookRepo, mockUserRepo)

	bookCreate := &models.BookCreate{
		Title:       "Test Book",
		Author:      "Test Author",
		ISBN:        "ISBN-123",
		Description: "Test description",
		UserID:      "123e4567-e89b-12d3-a456-426614174001",
	}

	user := &models.User{
		ID: bookCreate.UserID,
	}

	// Expectations
	mockUserRepo.On("FindByID", bookCreate.UserID).Return(user, nil)
	mockBookRepo.On("FindByISBN", bookCreate.ISBN).Return(nil, nil)
	mockBookRepo.On("Create", mock.AnythingOfType("*models.Book")).Return(nil)

	// Test
	book, err := service.Create(bookCreate)
	assert.NoError(t, err)
	assert.NotNil(t, book)
	assert.Equal(t, bookCreate.Title, book.Title)
	assert.Equal(t, bookCreate.Author, book.Author)
	assert.Equal(t, bookCreate.ISBN, book.ISBN)
	assert.Equal(t, bookCreate.Description, book.Description)
	assert.Equal(t, bookCreate.UserID, book.UserID)

	// Verify expectations
	mockBookRepo.AssertExpectations(t)
	mockUserRepo.AssertExpectations(t)
}

func TestBookService_Create_InvalidUserID(t *testing.T) {
	// Setup
	mockBookRepo := new(MockBookRepository)
	mockUserRepo := new(MockUserRepository)
	service := NewBookService(mockBookRepo, mockUserRepo)

	bookCreate := &models.BookCreate{
		Title:       "Test Book",
		Author:      "Test Author",
		ISBN:        "ISBN-123",
		Description: "Test description",
		UserID:      "invalid-id",
	}

	// Expectations
	mockUserRepo.On("FindByID", bookCreate.UserID).Return(nil, errors.New("invalid user ID"))

	// Test
	book, err := service.Create(bookCreate)
	assert.Error(t, err)
	assert.Nil(t, book)
	assert.Contains(t, err.Error(), "invalid user ID")

	// Verify expectations
	mockBookRepo.AssertExpectations(t)
	mockUserRepo.AssertExpectations(t)
}

func TestBookService_Create_UserNotFound(t *testing.T) {
	// Setup
	mockBookRepo := new(MockBookRepository)
	mockUserRepo := new(MockUserRepository)
	service := NewBookService(mockBookRepo, mockUserRepo)

	bookCreate := &models.BookCreate{
		Title:       "Test Book",
		Author:      "Test Author",
		ISBN:        "ISBN-123",
		Description: "Test description",
		UserID:      "123e4567-e89b-12d3-a456-426614174001",
	}

	// Expectations
	mockUserRepo.On("FindByID", bookCreate.UserID).Return(nil, nil)

	// Test
	book, err := service.Create(bookCreate)
	assert.Error(t, err)
	assert.Nil(t, book)
	assert.Contains(t, err.Error(), "user not found")

	// Verify expectations
	mockBookRepo.AssertExpectations(t)
	mockUserRepo.AssertExpectations(t)
}

func TestBookService_GetByID(t *testing.T) {
	// Setup
	mockBookRepo := new(MockBookRepository)
	mockUserRepo := new(MockUserRepository)
	service := NewBookService(mockBookRepo, mockUserRepo)

	book := &models.Book{
		ID:          "123e4567-e89b-12d3-a456-426614174000",
		Title:       "Test Book",
		Author:      "Test Author",
		ISBN:        "ISBN-123",
		Description: "Test description",
		UserID:      "123e4567-e89b-12d3-a456-426614174001",
	}

	// Expectations
	mockBookRepo.On("FindByID", book.ID).Return(book, nil)

	// Test
	result, err := service.GetByID(book.ID)
	assert.NoError(t, err)
	assert.Equal(t, book, result)

	// Verify expectations
	mockBookRepo.AssertExpectations(t)
}

func TestBookService_GetByID_InvalidID(t *testing.T) {
	// Setup
	mockBookRepo := new(MockBookRepository)
	mockUserRepo := new(MockUserRepository)
	service := NewBookService(mockBookRepo, mockUserRepo)

	// Test
	result, err := service.GetByID("invalid-id")
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "invalid book ID")

	// Verify expectations
	mockBookRepo.AssertExpectations(t)
}

func TestBookService_GetByID_BookNotFound(t *testing.T) {
	// Setup
	mockBookRepo := new(MockBookRepository)
	mockUserRepo := new(MockUserRepository)
	service := NewBookService(mockBookRepo, mockUserRepo)

	bookID := "123e4567-e89b-12d3-a456-426614174000"

	// Expectations
	mockBookRepo.On("FindByID", bookID).Return(nil, nil)

	// Test
	result, err := service.GetByID(bookID)
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "book not found")

	// Verify expectations
	mockBookRepo.AssertExpectations(t)
}

func TestBookService_GetAll(t *testing.T) {
	// Setup
	mockBookRepo := new(MockBookRepository)
	mockUserRepo := new(MockUserRepository)
	service := NewBookService(mockBookRepo, mockUserRepo)

	books := []models.Book{
		{
			ID:     "123e4567-e89b-12d3-a456-426614174002",
			Title:  "Book 1",
			Author: "Author 1",
		},
		{
			ID:     "123e4567-e89b-12d3-a456-426614174003",
			Title:  "Book 2",
			Author: "Author 2",
		},
	}

	// Expectations
	mockBookRepo.On("FindAll").Return(books, nil)

	// Test
	results, err := service.GetAll()
	assert.NoError(t, err)
	assert.Equal(t, books, results)

	// Verify expectations
	mockBookRepo.AssertExpectations(t)
}

func TestBookService_GetByUserID(t *testing.T) {
	// Setup
	mockBookRepo := new(MockBookRepository)
	mockUserRepo := new(MockUserRepository)
	service := NewBookService(mockBookRepo, mockUserRepo)

	userID := "123e4567-e89b-12d3-a456-426614174001"
	user := &models.User{
		ID: userID,
	}
	books := []models.Book{
		{
			ID:     "123e4567-e89b-12d3-a456-426614174002",
			Title:  "Book 1",
			UserID: userID,
		},
		{
			ID:     "123e4567-e89b-12d3-a456-426614174003",
			Title:  "Book 2",
			UserID: userID,
		},
	}

	// Expectations
	mockUserRepo.On("FindByID", userID).Return(user, nil)
	mockBookRepo.On("FindByUserID", userID).Return(books, nil)

	// Test
	results, err := service.GetByUserID(userID)
	assert.NoError(t, err)
	assert.Equal(t, books, results)

	// Verify expectations
	mockBookRepo.AssertExpectations(t)
	mockUserRepo.AssertExpectations(t)
}

func TestBookService_GetByUserID_InvalidUserID(t *testing.T) {
	// Setup
	mockBookRepo := new(MockBookRepository)
	mockUserRepo := new(MockUserRepository)
	service := NewBookService(mockBookRepo, mockUserRepo)

	// Test
	results, err := service.GetByUserID("invalid-id")
	assert.Error(t, err)
	assert.Nil(t, results)
	assert.Contains(t, err.Error(), "invalid user ID")

	// Verify expectations
	mockBookRepo.AssertExpectations(t)
	mockUserRepo.AssertExpectations(t)
}

func TestBookService_GetByUserID_UserNotFound(t *testing.T) {
	// Setup
	mockBookRepo := new(MockBookRepository)
	mockUserRepo := new(MockUserRepository)
	service := NewBookService(mockBookRepo, mockUserRepo)

	userID := "123e4567-e89b-12d3-a456-426614174001"

	// Expectations
	mockUserRepo.On("FindByID", userID).Return(nil, nil)

	// Test
	results, err := service.GetByUserID(userID)
	assert.Error(t, err)
	assert.Nil(t, results)
	assert.Contains(t, err.Error(), "user not found")

	// Verify expectations
	mockBookRepo.AssertExpectations(t)
	mockUserRepo.AssertExpectations(t)
}

func TestBookService_Update(t *testing.T) {
	// Setup
	mockBookRepo := new(MockBookRepository)
	mockUserRepo := new(MockUserRepository)
	service := NewBookService(mockBookRepo, mockUserRepo)

	bookID := "123e4567-e89b-12d3-a456-426614174000"
	bookUpdate := &models.BookUpdate{
		Title:       "Updated Book",
		Author:      "Updated Author",
		ISBN:        "ISBN-456",
		Description: "Updated description",
	}

	existingBook := &models.Book{
		ID:          bookID,
		Title:       "Test Book",
		Author:      "Test Author",
		ISBN:        "ISBN-123",
		Description: "Test description",
		UserID:      "123e4567-e89b-12d3-a456-426614174001",
	}

	// Expectations
	mockBookRepo.On("FindByID", bookID).Return(existingBook, nil)
	mockBookRepo.On("FindByISBN", bookUpdate.ISBN).Return(nil, nil)
	mockBookRepo.On("Update", mock.AnythingOfType("*models.Book")).Return(nil)

	// Test
	updatedBook, err := service.Update(bookID, bookUpdate)
	assert.NoError(t, err)
	assert.Equal(t, bookID, updatedBook.ID)
	assert.Equal(t, bookUpdate.Title, updatedBook.Title)
	assert.Equal(t, bookUpdate.Author, updatedBook.Author)
	assert.Equal(t, bookUpdate.ISBN, updatedBook.ISBN)
	assert.Equal(t, bookUpdate.Description, updatedBook.Description)
	assert.Equal(t, existingBook.UserID, updatedBook.UserID) // UserID should not change

	// Verify expectations
	mockBookRepo.AssertExpectations(t)
}

func TestBookService_Update_InvalidID(t *testing.T) {
	// Setup
	mockBookRepo := new(MockBookRepository)
	mockUserRepo := new(MockUserRepository)
	service := NewBookService(mockBookRepo, mockUserRepo)

	bookUpdate := &models.BookUpdate{
		Title: "Updated Book",
	}

	// Test
	updatedBook, err := service.Update("invalid-id", bookUpdate)
	assert.Error(t, err)
	assert.Nil(t, updatedBook)
	assert.Contains(t, err.Error(), "invalid book ID")

	// Verify expectations
	mockBookRepo.AssertExpectations(t)
}

func TestBookService_Update_BookNotFound(t *testing.T) {
	// Setup
	mockBookRepo := new(MockBookRepository)
	mockUserRepo := new(MockUserRepository)
	service := NewBookService(mockBookRepo, mockUserRepo)

	bookID := "123e4567-e89b-12d3-a456-426614174000"
	bookUpdate := &models.BookUpdate{
		Title: "Updated Book",
	}

	// Expectations
	mockBookRepo.On("FindByID", bookID).Return(nil, nil)

	// Test
	updatedBook, err := service.Update(bookID, bookUpdate)
	assert.Error(t, err)
	assert.Nil(t, updatedBook)
	assert.Contains(t, err.Error(), "book not found")

	// Verify expectations
	mockBookRepo.AssertExpectations(t)
}

func TestBookService_Delete(t *testing.T) {
	// Setup
	mockBookRepo := new(MockBookRepository)
	mockUserRepo := new(MockUserRepository)
	service := NewBookService(mockBookRepo, mockUserRepo)

	bookID := "123e4567-e89b-12d3-a456-426614174000"
	book := &models.Book{
		ID: bookID,
	}

	// Expectations
	mockBookRepo.On("FindByID", bookID).Return(book, nil)
	mockBookRepo.On("Delete", bookID).Return(nil)

	// Test
	err := service.Delete(bookID)
	assert.NoError(t, err)

	// Verify expectations
	mockBookRepo.AssertExpectations(t)
}

func TestBookService_Delete_InvalidID(t *testing.T) {
	// Setup
	mockBookRepo := new(MockBookRepository)
	mockUserRepo := new(MockUserRepository)
	service := NewBookService(mockBookRepo, mockUserRepo)

	// Test
	err := service.Delete("invalid-id")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid book ID")

	// Verify expectations
	mockBookRepo.AssertExpectations(t)
}

func TestBookService_Delete_BookNotFound(t *testing.T) {
	// Setup
	mockBookRepo := new(MockBookRepository)
	mockUserRepo := new(MockUserRepository)
	service := NewBookService(mockBookRepo, mockUserRepo)

	bookID := "123e4567-e89b-12d3-a456-426614174000"

	// Expectations
	mockBookRepo.On("FindByID", bookID).Return(nil, nil)

	// Test
	err := service.Delete(bookID)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "book not found")

	// Verify expectations
	mockBookRepo.AssertExpectations(t)
}
