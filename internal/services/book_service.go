package services

import (
	"errors"
	"fmt"

	"github.com/SimpleBookRental/backend/internal/models"
	"github.com/SimpleBookRental/backend/internal/repositories"
	"github.com/SimpleBookRental/backend/pkg/utils"
)

// BookService handles business logic for books
type BookService struct {
	bookRepo *repositories.BookRepository
	userRepo *repositories.UserRepository
}

// NewBookService creates a new book service
func NewBookService(bookRepo *repositories.BookRepository, userRepo *repositories.UserRepository) *BookService {
	return &BookService{
		bookRepo: bookRepo,
		userRepo: userRepo,
	}
}

// Create creates a new book
func (s *BookService) Create(bookCreate *models.BookCreate) (*models.Book, error) {
	// Check if user exists
	user, err := s.userRepo.FindByID(bookCreate.UserID)
	if err != nil {
		return nil, fmt.Errorf("error checking user: %w", err)
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	// Check if ISBN already exists
	existingBook, err := s.bookRepo.FindByISBN(bookCreate.ISBN)
	if err != nil {
		return nil, fmt.Errorf("error checking existing book: %w", err)
	}
	if existingBook != nil {
		return nil, errors.New("ISBN already exists")
	}

	// Create book
	book := &models.Book{
		Title:       bookCreate.Title,
		Author:      bookCreate.Author,
		ISBN:        bookCreate.ISBN,
		Description: bookCreate.Description,
		UserID:      bookCreate.UserID,
	}

	err = s.bookRepo.Create(book)
	if err != nil {
		return nil, fmt.Errorf("error creating book: %w", err)
	}

	return book, nil
}

// GetByID gets a book by ID
func (s *BookService) GetByID(id string) (*models.Book, error) {
	if !utils.IsValidUUID(id) {
		return nil, errors.New("invalid book ID")
	}

	book, err := s.bookRepo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("error finding book: %w", err)
	}
	if book == nil {
		return nil, errors.New("book not found")
	}

	return book, nil
}

// GetAll gets all books
func (s *BookService) GetAll() ([]models.Book, error) {
	return s.bookRepo.FindAll()
}

// GetByUserID gets all books by user ID
func (s *BookService) GetByUserID(userID string) ([]models.Book, error) {
	if !utils.IsValidUUID(userID) {
		return nil, errors.New("invalid user ID")
	}

	// Check if user exists
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, fmt.Errorf("error checking user: %w", err)
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	return s.bookRepo.FindByUserID(userID)
}

// Update updates a book
func (s *BookService) Update(id string, bookUpdate *models.BookUpdate) (*models.Book, error) {
	if !utils.IsValidUUID(id) {
		return nil, errors.New("invalid book ID")
	}

	book, err := s.bookRepo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("error finding book: %w", err)
	}
	if book == nil {
		return nil, errors.New("book not found")
	}

	// Update fields if provided
	if bookUpdate.Title != "" {
		book.Title = bookUpdate.Title
	}
	if bookUpdate.Author != "" {
		book.Author = bookUpdate.Author
	}
	if bookUpdate.ISBN != "" {
		// Check if ISBN already exists
		if bookUpdate.ISBN != book.ISBN {
			existingBook, err := s.bookRepo.FindByISBN(bookUpdate.ISBN)
			if err != nil {
				return nil, fmt.Errorf("error checking existing book: %w", err)
			}
			if existingBook != nil {
				return nil, errors.New("ISBN already exists")
			}
			book.ISBN = bookUpdate.ISBN
		}
	}
	if bookUpdate.Description != "" {
		book.Description = bookUpdate.Description
	}
	if bookUpdate.UserID != "" {
		// Check if user exists
		user, err := s.userRepo.FindByID(bookUpdate.UserID)
		if err != nil {
			return nil, fmt.Errorf("error checking user: %w", err)
		}
		if user == nil {
			return nil, errors.New("user not found")
		}
		book.UserID = bookUpdate.UserID
	}

	err = s.bookRepo.Update(book)
	if err != nil {
		return nil, fmt.Errorf("error updating book: %w", err)
	}

	return book, nil
}

// Delete deletes a book
func (s *BookService) Delete(id string) error {
	if !utils.IsValidUUID(id) {
		return errors.New("invalid book ID")
	}

	book, err := s.bookRepo.FindByID(id)
	if err != nil {
		return fmt.Errorf("error finding book: %w", err)
	}
	if book == nil {
		return errors.New("book not found")
	}

	return s.bookRepo.Delete(id)
}
