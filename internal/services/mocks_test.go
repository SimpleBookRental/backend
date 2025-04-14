package services

import (
	"github.com/SimpleBookRental/backend/internal/models"
	"github.com/SimpleBookRental/backend/internal/repositories"
	"github.com/stretchr/testify/mock"
)

// MockUserRepository is a mock implementation of repositories.UserRepositoryInterface
type MockUserRepository struct {
	mock.Mock
}

// Ensure MockUserRepository implements repositories.UserRepositoryInterface
var _ repositories.UserRepositoryInterface = (*MockUserRepository)(nil)

func (m *MockUserRepository) Create(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) FindByID(id string) (*models.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) FindByEmail(email string) (*models.User, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) FindAll() ([]models.User, error) {
	args := m.Called()
	return args.Get(0).([]models.User), args.Error(1)
}

func (m *MockUserRepository) Update(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

// MockBookRepository is a mock implementation of repositories.BookRepositoryInterface
type MockBookRepository struct {
	mock.Mock
}

// Ensure MockBookRepository implements repositories.BookRepositoryInterface
var _ repositories.BookRepositoryInterface = (*MockBookRepository)(nil)

func (m *MockBookRepository) Create(book *models.Book) error {
	args := m.Called(book)
	return args.Error(0)
}

func (m *MockBookRepository) FindByID(id string) (*models.Book, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Book), args.Error(1)
}

func (m *MockBookRepository) FindByISBN(isbn string) (*models.Book, error) {
	args := m.Called(isbn)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Book), args.Error(1)
}

func (m *MockBookRepository) FindAll() ([]models.Book, error) {
	args := m.Called()
	return args.Get(0).([]models.Book), args.Error(1)
}

func (m *MockBookRepository) FindByUserID(userID string) ([]models.Book, error) {
	args := m.Called(userID)
	return args.Get(0).([]models.Book), args.Error(1)
}

func (m *MockBookRepository) Update(book *models.Book) error {
	args := m.Called(book)
	return args.Error(0)
}

func (m *MockBookRepository) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

// MockTokenRepository is a mock implementation of repositories.TokenRepositoryInterface
type MockTokenRepository struct {
	mock.Mock
}

// Ensure MockTokenRepository implements repositories.TokenRepositoryInterface
var _ repositories.TokenRepositoryInterface = (*MockTokenRepository)(nil)

func (m *MockTokenRepository) CreateToken(token *models.IssuedToken) error {
	args := m.Called(token)
	return args.Error(0)
}

func (m *MockTokenRepository) FindTokenByValue(tokenString string) (*models.IssuedToken, error) {
	args := m.Called(tokenString)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.IssuedToken), args.Error(1)
}

func (m *MockTokenRepository) FindActiveTokensByUserID(userID string) ([]models.IssuedToken, error) {
	args := m.Called(userID)
	return args.Get(0).([]models.IssuedToken), args.Error(1)
}

func (m *MockTokenRepository) RevokeToken(token *models.IssuedToken) error {
	args := m.Called(token)
	return args.Error(0)
}

func (m *MockTokenRepository) RevokeAllUserTokens(userID string) error {
	args := m.Called(userID)
	return args.Error(0)
}

func (m *MockTokenRepository) CleanupExpiredTokens() error {
	args := m.Called()
	return args.Error(0)
}
