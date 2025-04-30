package service

import (
	"errors"

	"github.com/SimpleBookRental/backend/internal/domain"
	"github.com/SimpleBookRental/backend/pkg/logger"
	"go.uber.org/zap"
)

// BookServiceImpl implements domain.BookService
type BookServiceImpl struct {
	repo         domain.BookRepository
	categoryRepo domain.CategoryRepository
	logger       *logger.Logger
}

// NewBookService creates a new BookService
func NewBookService(repo domain.BookRepository, categoryRepo domain.CategoryRepository, logger *logger.Logger) domain.BookService {
	return &BookServiceImpl{
		repo:         repo,
		categoryRepo: categoryRepo,
		logger:       logger,
	}
}

// GetByID retrieves a book by ID
func (s *BookServiceImpl) GetByID(id int64) (*domain.Book, error) {
	book, err := s.repo.GetByID(id)
	if err != nil {
		s.logger.Error("Failed to get book by ID", zap.Int64("id", id), zap.Error(err))
		return nil, err
	}
	return book, nil
}

// GetByISBN retrieves a book by ISBN
func (s *BookServiceImpl) GetByISBN(isbn string) (*domain.Book, error) {
	book, err := s.repo.GetByISBN(isbn)
	if err != nil {
		s.logger.Error("Failed to get book by ISBN", zap.String("isbn", isbn), zap.Error(err))
		return nil, err
	}
	return book, nil
}

// List retrieves a list of books with pagination
func (s *BookServiceImpl) List(limit, offset int32) ([]*domain.Book, error) {
	books, err := s.repo.List(limit, offset)
	if err != nil {
		s.logger.Error("Failed to list books", zap.Error(err))
		return nil, err
	}
	return books, nil
}

// ListByCategory retrieves a list of books by category with pagination
func (s *BookServiceImpl) ListByCategory(categoryID int64, limit, offset int32) ([]*domain.Book, error) {
	// Check if category exists
	_, err := s.categoryRepo.GetByID(categoryID)
	if err != nil {
		s.logger.Error("Failed to get category by ID", zap.Int64("categoryID", categoryID), zap.Error(err))
		return nil, err
	}

	books, err := s.repo.ListByCategory(categoryID, limit, offset)
	if err != nil {
		s.logger.Error("Failed to list books by category", zap.Int64("categoryID", categoryID), zap.Error(err))
		return nil, err
	}
	return books, nil
}

// Search searches for books based on search parameters
func (s *BookServiceImpl) Search(params domain.BookSearchParams) ([]*domain.Book, error) {
	// Validate category ID if provided
	if params.CategoryID != 0 {
		_, err := s.categoryRepo.GetByID(params.CategoryID)
		if err != nil {
			s.logger.Error("Invalid category ID in search params", zap.Int64("categoryID", params.CategoryID), zap.Error(err))
			return nil, err
		}
	}

	books, err := s.repo.Search(params)
	if err != nil {
		s.logger.Error("Failed to search books", zap.Error(err))
		return nil, err
	}
	return books, nil
}

// Create creates a new book
func (s *BookServiceImpl) Create(book *domain.Book) (*domain.Book, error) {
	// Check if ISBN already exists
	existingBook, err := s.repo.GetByISBN(book.ISBN)
	if err == nil && existingBook != nil {
		return nil, domain.ErrBookAlreadyExists
	}
	if err != nil && !errors.Is(err, domain.ErrBookNotFound) {
		s.logger.Error("Error checking ISBN existence", zap.String("isbn", book.ISBN), zap.Error(err))
		return nil, err
	}

	// Validate category if provided
	if book.CategoryID != 0 {
		_, err := s.categoryRepo.GetByID(book.CategoryID)
		if err != nil {
			s.logger.Error("Invalid category ID", zap.Int64("categoryID", book.CategoryID), zap.Error(err))
			return nil, err
		}
	}

	// Ensure available copies doesn't exceed total copies
	if book.AvailableCopies > book.TotalCopies {
		book.AvailableCopies = book.TotalCopies
	}

	// Create book
	createdBook, err := s.repo.Create(book)
	if err != nil {
		s.logger.Error("Failed to create book", zap.Error(err))
		return nil, err
	}

	return createdBook, nil
}

// Update updates an existing book
func (s *BookServiceImpl) Update(book *domain.Book) (*domain.Book, error) {
	// Check if book exists
	existingBook, err := s.repo.GetByID(book.ID)
	if err != nil {
		s.logger.Error("Failed to get book by ID", zap.Int64("id", book.ID), zap.Error(err))
		return nil, err
	}

	// Check if ISBN is being changed and if it already exists
	if book.ISBN != existingBook.ISBN {
		bookByISBN, err := s.repo.GetByISBN(book.ISBN)
		if err == nil && bookByISBN != nil {
			return nil, domain.ErrBookAlreadyExists
		}
		if err != nil && !errors.Is(err, domain.ErrBookNotFound) {
			s.logger.Error("Error checking ISBN existence", zap.String("isbn", book.ISBN), zap.Error(err))
			return nil, err
		}
	}

	// Validate category if provided
	if book.CategoryID != 0 {
		_, err := s.categoryRepo.GetByID(book.CategoryID)
		if err != nil {
			s.logger.Error("Invalid category ID", zap.Int64("categoryID", book.CategoryID), zap.Error(err))
			return nil, err
		}
	}

	// Preserve copies information
	book.TotalCopies = existingBook.TotalCopies
	book.AvailableCopies = existingBook.AvailableCopies

	// Update book
	updatedBook, err := s.repo.Update(book)
	if err != nil {
		s.logger.Error("Failed to update book", zap.Int64("id", book.ID), zap.Error(err))
		return nil, err
	}

	return updatedBook, nil
}

// UpdateCopies updates the total and available copies of a book
func (s *BookServiceImpl) UpdateCopies(id int64, totalCopies, availableCopies int32) (*domain.Book, error) {
	// Check if book exists
	existingBook, err := s.repo.GetByID(id)
	if err != nil {
		s.logger.Error("Failed to get book by ID", zap.Int64("id", id), zap.Error(err))
		return nil, err
	}

	// Validate copies
	if totalCopies < 0 {
		return nil, domain.NewInvalidInputError("total copies cannot be negative")
	}

	if availableCopies < 0 {
		return nil, domain.NewInvalidInputError("available copies cannot be negative")
	}

	if availableCopies > totalCopies {
		return nil, domain.NewInvalidInputError("available copies cannot exceed total copies")
	}

	// Calculate the difference in available copies
	availableDiff := availableCopies - existingBook.AvailableCopies
	totalDiff := totalCopies - existingBook.TotalCopies

	// If reducing total copies, ensure we're not removing copies that are currently borrowed
	if totalDiff < 0 && availableDiff > 0 {
		return nil, domain.NewInvalidInputError("cannot increase available copies while decreasing total copies")
	}

	// If reducing total copies, ensure we're not removing more copies than are available
	if totalDiff < 0 && -totalDiff > existingBook.AvailableCopies {
		return nil, domain.NewInvalidInputError("cannot remove more copies than are currently available")
	}

	// Update copies
	updatedBook, err := s.repo.UpdateCopies(id, totalCopies, availableCopies)
	if err != nil {
		s.logger.Error("Failed to update book copies", zap.Int64("id", id), zap.Error(err))
		return nil, err
	}

	return updatedBook, nil
}

// Delete deletes a book
func (s *BookServiceImpl) Delete(id int64) error {
	// Check if book exists
	book, err := s.repo.GetByID(id)
	if err != nil {
		s.logger.Error("Failed to get book by ID", zap.Int64("id", id), zap.Error(err))
		return err
	}

	// Check if all copies are available (no active rentals)
	if book.AvailableCopies < book.TotalCopies {
		return domain.NewInvalidInputError("cannot delete book with active rentals")
	}

	err = s.repo.Delete(id)
	if err != nil {
		s.logger.Error("Failed to delete book", zap.Int64("id", id), zap.Error(err))
		return err
	}
	return nil
}

// IsAvailable checks if a book is available for rental
func (s *BookServiceImpl) IsAvailable(id int64) (bool, error) {
	book, err := s.repo.GetByID(id)
	if err != nil {
		s.logger.Error("Failed to get book by ID", zap.Int64("id", id), zap.Error(err))
		return false, err
	}
	return book.AvailableCopies > 0, nil
}
