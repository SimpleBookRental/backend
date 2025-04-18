package services

import (
	"github.com/SimpleBookRental/backend/internal/models"
	"github.com/SimpleBookRental/backend/internal/repositories"
)

// UserServiceInterface defines the interface for user service
type UserServiceInterface interface {
	Register(userCreate *models.UserCreate) (*models.User, error)
	GetByID(id string) (*models.User, error)
	GetAll() ([]models.User, error)
	Update(id string, userUpdate *models.UserUpdate) (*models.User, error)
	Delete(id string) error
	Login(userLogin *models.UserLogin, tokenRepo repositories.TokenRepositoryInterface) (*models.LoginResponse, error)
}

// BookServiceInterface defines the interface for book service
type BookServiceInterface interface {
	Create(bookCreate *models.BookCreate) (*models.Book, error)
	GetByID(id string) (*models.Book, error)
	GetAll() ([]models.Book, error)
	GetByUserID(userID string) ([]models.Book, error)
	Update(id string, bookUpdate *models.BookUpdate) (*models.Book, error)
	Delete(id string) error
}

// TokenServiceInterface defines the interface for token service
type TokenServiceInterface interface {
	RefreshToken(request *models.RefreshTokenRequest) (*models.RefreshTokenResponse, error)
	Logout(request *models.LogoutRequest) error
}

// BookUserServiceInterface defines the interface for book-user service
type BookUserServiceInterface interface {
	TransferBookOwnership(bookID, fromUserID, toUserID string) error
	CreateBookWithUser(bookCreate *models.BookCreateRequest, userID string) (*models.Book, error)
}
