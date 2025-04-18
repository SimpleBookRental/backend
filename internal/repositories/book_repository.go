package repositories

import (
	"errors"
	"fmt"

	"github.com/SimpleBookRental/backend/internal/models"
	"gorm.io/gorm"
)

// BookRepository handles database operations for books
type BookRepository struct {
	db *gorm.DB
}

// GetDB returns the database connection
func (r *BookRepository) GetDB() interface{} {
	return r.db
}


// WithTx returns a new BookRepository with the given transaction
func (r *BookRepository) WithTx(tx interface{}) (BookRepositoryInterface, error) {
	db, ok := tx.(*gorm.DB)
	if !ok {
		return nil, fmt.Errorf("invalid transaction type")
	}
	return &BookRepository{db: db}, nil
}

// Ensure BookRepository implements BookRepositoryInterface
var _ BookRepositoryInterface = (*BookRepository)(nil)

// NewBookRepository creates a new book repository
func NewBookRepository(db *gorm.DB) *BookRepository {
	return &BookRepository{db: db}
}

// Create creates a new book
func (r *BookRepository) Create(book *models.Book) error {
	return r.db.Create(book).Error
}

// FindByID finds a book by ID
func (r *BookRepository) FindByID(id string) (*models.Book, error) {
	var book models.Book
	err := r.db.Preload("User").First(&book, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &book, nil
}

// FindByISBN finds a book by ISBN
func (r *BookRepository) FindByISBN(isbn string) (*models.Book, error) {
	var book models.Book
	err := r.db.First(&book, "isbn = ?", isbn).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &book, nil
}

// FindAll finds all books
func (r *BookRepository) FindAll() ([]models.Book, error) {
	var books []models.Book
	err := r.db.Preload("User").Find(&books).Error
	return books, err
}

// FindByUserID finds all books by user ID
func (r *BookRepository) FindByUserID(userID string) ([]models.Book, error) {
	var books []models.Book
	err := r.db.Preload("User").Where("user_id = ?", userID).Find(&books).Error
	return books, err
}

// Update updates a book
func (r *BookRepository) Update(book *models.Book) error {
	return r.db.Save(book).Error
}

/*
Delete deletes a book by its ID.
*/
func (r *BookRepository) Delete(id string) error {
	return r.db.Delete(&models.Book{}, "id = ?", id).Error
}

/*
DeleteByUserID deletes all books belonging to a specific user.
*/
func (r *BookRepository) DeleteByUserID(userID string) error {
	return r.db.Delete(&models.Book{}, "user_id = ?", userID).Error
}
