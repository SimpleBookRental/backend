package repositories

import (
	"github.com/SimpleBookRental/backend/internal/models"
)

type UserRepositoryInterface interface {
	GetDB() interface{}
	WithTx(tx interface{}) (UserRepositoryInterface, error)
	Create(user *models.User) error
	FindByID(id string) (*models.User, error)
	FindByEmail(email string) (*models.User, error)
	FindAll() ([]models.User, error)
	Update(user *models.User) error
	Delete(id string) error
}

/*
BookRepositoryInterface defines the interface for book repository.
Add DeleteByUserID to allow deleting all books by a specific user.
*/
type BookRepositoryInterface interface {
	GetDB() interface{}
	WithTx(tx interface{}) (BookRepositoryInterface, error)
	Create(book *models.Book) error
	FindAll() ([]models.Book, error)
	FindByID(id string) (*models.Book, error)
	FindByISBN(isbn string) (*models.Book, error)
	FindByUserID(userID string) ([]models.Book, error)
	Update(book *models.Book) error
	Delete(id string) error
	DeleteByUserID(userID string) error // Delete all books by user ID
}

// TokenRepositoryInterface defines the interface for token repository
type TokenRepositoryInterface interface {
	GetDB() interface{}
	WithTx(tx interface{}) (TokenRepositoryInterface, error)
	CreateToken(token *models.IssuedToken) error
	FindTokenByValue(tokenString string) (*models.IssuedToken, error)
	FindActiveTokensByUserID(userID string) ([]models.IssuedToken, error)
	RevokeToken(token *models.IssuedToken) error
	RevokeAllUserTokens(userID string) error
	DeleteByUserID(userID string) error // Delete all tokens by user ID
	CleanupExpiredTokens() error
}
