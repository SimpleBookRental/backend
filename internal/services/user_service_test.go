// Unit tests for UserService using gomock and testify.
package services

import (
	"errors"
	"testing"
	"github.com/SimpleBookRental/backend/internal/models"
	"github.com/SimpleBookRental/backend/internal/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestUserService_Register_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockUserRepo := mocks.NewMockUserRepositoryInterface(ctrl)
	mockBookRepo := mocks.NewMockBookRepositoryInterface(ctrl)
	mockTokenRepo := mocks.NewMockTokenRepositoryInterface(ctrl)
	service := NewUserService(mockUserRepo, mockBookRepo, mockTokenRepo)

	userCreate := &models.UserCreate{
		Name:     "Test User",
		Email:    "test@example.com",
		Password: "password123",
	}
	mockUserRepo.EXPECT().FindByEmail(userCreate.Email).Return(nil, nil)
	mockUserRepo.EXPECT().Create(gomock.Any()).Return(nil)

	user, err := service.Register(userCreate)
	assert.NoError(t, err)
	assert.Equal(t, userCreate.Email, user.Email)
}

func TestUserService_Register_EmailExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockUserRepo := mocks.NewMockUserRepositoryInterface(ctrl)
	mockBookRepo := mocks.NewMockBookRepositoryInterface(ctrl)
	mockTokenRepo := mocks.NewMockTokenRepositoryInterface(ctrl)
	service := NewUserService(mockUserRepo, mockBookRepo, mockTokenRepo)

	userCreate := &models.UserCreate{
		Name:     "Test User",
		Email:    "test@example.com",
		Password: "password123",
	}
	mockUserRepo.EXPECT().FindByEmail(userCreate.Email).Return(&models.User{}, nil)

	user, err := service.Register(userCreate)
	assert.Nil(t, user)
	assert.ErrorContains(t, err, "email already exists")
}

func TestUserService_Register_RepoError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockUserRepo := mocks.NewMockUserRepositoryInterface(ctrl)
	mockBookRepo := mocks.NewMockBookRepositoryInterface(ctrl)
	mockTokenRepo := mocks.NewMockTokenRepositoryInterface(ctrl)
	service := NewUserService(mockUserRepo, mockBookRepo, mockTokenRepo)

	userCreate := &models.UserCreate{
		Name:     "Test User",
		Email:    "test@example.com",
		Password: "password123",
	}
	mockUserRepo.EXPECT().FindByEmail(userCreate.Email).Return(nil, errors.New("db error"))

	user, err := service.Register(userCreate)
	assert.Nil(t, user)
	assert.ErrorContains(t, err, "error checking existing user")
}

// ... (giữ nguyên các test còn lại, không đổi tên hàm test, chỉ đổi service.Create thành service.Register nếu có)
