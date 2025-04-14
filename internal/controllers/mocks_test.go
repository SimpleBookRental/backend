package controllers

import (
	"github.com/SimpleBookRental/backend/internal/models"
	"github.com/SimpleBookRental/backend/internal/repositories"
	"github.com/SimpleBookRental/backend/internal/services"
	"github.com/stretchr/testify/mock"
)

// MockUserService is a mock implementation of services.UserServiceInterface
type MockUserService struct {
	mock.Mock
}

// Ensure MockUserService implements services.UserServiceInterface
var _ services.UserServiceInterface = (*MockUserService)(nil)

func (m *MockUserService) Create(userCreate *models.UserCreate) (*models.User, error) {
	args := m.Called(userCreate)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserService) GetByID(id string) (*models.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserService) GetAll() ([]models.User, error) {
	args := m.Called()
	return args.Get(0).([]models.User), args.Error(1)
}

func (m *MockUserService) Update(id string, userUpdate *models.UserUpdate) (*models.User, error) {
	args := m.Called(id, userUpdate)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserService) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockUserService) Login(userLogin *models.UserLogin, tokenRepo repositories.TokenRepositoryInterface) (*models.LoginResponse, error) {
	args := m.Called(userLogin, tokenRepo)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.LoginResponse), args.Error(1)
}

// MockBookService is a mock implementation of services.BookServiceInterface
type MockBookService struct {
	mock.Mock
}

// Ensure MockBookService implements services.BookServiceInterface
var _ services.BookServiceInterface = (*MockBookService)(nil)

func (m *MockBookService) Create(bookCreate *models.BookCreate) (*models.Book, error) {
	args := m.Called(bookCreate)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Book), args.Error(1)
}

func (m *MockBookService) GetByID(id string) (*models.Book, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Book), args.Error(1)
}

func (m *MockBookService) GetAll() ([]models.Book, error) {
	args := m.Called()
	return args.Get(0).([]models.Book), args.Error(1)
}

func (m *MockBookService) GetByUserID(userID string) ([]models.Book, error) {
	args := m.Called(userID)
	return args.Get(0).([]models.Book), args.Error(1)
}

func (m *MockBookService) Update(id string, bookUpdate *models.BookUpdate) (*models.Book, error) {
	args := m.Called(id, bookUpdate)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Book), args.Error(1)
}

func (m *MockBookService) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

// MockTokenService is a mock implementation of services.TokenServiceInterface
type MockTokenService struct {
	mock.Mock
}

// Ensure MockTokenService implements services.TokenServiceInterface
var _ services.TokenServiceInterface = (*MockTokenService)(nil)

func (m *MockTokenService) RefreshToken(request *models.RefreshTokenRequest) (*models.RefreshTokenResponse, error) {
	args := m.Called(request)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.RefreshTokenResponse), args.Error(1)
}

func (m *MockTokenService) Logout(request *models.LogoutRequest) error {
	args := m.Called(request)
	return args.Error(0)
}
