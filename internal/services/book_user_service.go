package services

import (
	"errors"
	"fmt"

	"github.com/SimpleBookRental/backend/internal/models"
	"github.com/SimpleBookRental/backend/internal/repositories"
	"gorm.io/gorm"
)

// BookUserService handles operations that involve both books and users
type BookUserService struct {
	// Ensure BookUserService implements BookUserServiceInterface
	_              BookUserServiceInterface
	txManager      repositories.TransactionManagerInterface
	bookRepository repositories.BookRepositoryInterface
	userRepository repositories.UserRepositoryInterface
}

// NewBookUserService creates a new book-user service
func NewBookUserService(
	txManager repositories.TransactionManagerInterface,
	bookRepository repositories.BookRepositoryInterface,
	userRepository repositories.UserRepositoryInterface,
) *BookUserService {
	return &BookUserService{
		txManager:      txManager,
		bookRepository: bookRepository,
		userRepository: userRepository,
	}
}

// TransferBookOwnership transfers a book from one user to another
func (s *BookUserService) TransferBookOwnership(bookID, fromUserID, toUserID string) error {
	return s.txManager.WithTransaction(func(tx *gorm.DB) error {
		// Create repositories with the shared transaction
		bookRepoTx, err := s.bookRepository.WithTx(tx)
		if err != nil {
			return err
		}
		userRepoTx, err := s.userRepository.WithTx(tx)
		if err != nil {
			return err
		}

		// Get the book
		book, err := bookRepoTx.FindByID(bookID)
		if err != nil {
			return fmt.Errorf("error finding book: %w", err)
		}
		if book == nil {
			return errors.New("book not found")
		}

		// Verify current ownership
		if book.UserID != fromUserID {
			return errors.New("book does not belong to the specified user")
		}

		// Check if the target user exists
		toUser, err := userRepoTx.FindByID(toUserID)
		if err != nil {
			return fmt.Errorf("error finding target user: %w", err)
		}
		if toUser == nil {
			return errors.New("target user not found")
		}

		// Update book ownership
		// Perform direct update of user_id on the book record
		if err := tx.Model(&models.Book{}).
			Where("id = ?", bookID).
			Update("user_id", toUserID).Error; err != nil {
			return fmt.Errorf("error updating book owner: %w", err)
		}

		return nil
	})
}

// CreateBookWithUser creates a book and associates it with a user
func (s *BookUserService) CreateBookWithUser(bookCreate *models.BookCreateRequest, userID string) (*models.Book, error) {
	var book *models.Book

	err := s.txManager.WithTransaction(func(tx *gorm.DB) error {
		// Create repositories with the shared transaction
		bookRepoTx, err := s.bookRepository.WithTx(tx)
		if err != nil {
			return err
		}
		userRepoTx, err := s.userRepository.WithTx(tx)
		if err != nil {
			return err
		}

		// Check if user exists
		user, err := userRepoTx.FindByID(userID)
		if err != nil {
			return fmt.Errorf("error checking user: %w", err)
		}
		if user == nil {
			return errors.New("user not found")
		}

		// Check if ISBN already exists
		existingBook, err := bookRepoTx.FindByISBN(bookCreate.ISBN)
		if err != nil {
			return fmt.Errorf("error checking existing book: %w", err)
		}
		if existingBook != nil {
			return errors.New("ISBN already exists")
		}

		// Create book
		newBook := &models.Book{
			Title:       bookCreate.Title,
			Author:      bookCreate.Author,
			ISBN:        bookCreate.ISBN,
			Description: bookCreate.Description,
			UserID:      userID,
		}

		if err := bookRepoTx.Create(newBook); err != nil {
			return fmt.Errorf("error creating book: %w", err)
		}

		book = newBook
		return nil
	})

	if err != nil {
		return nil, err
	}

	return book, nil
}
