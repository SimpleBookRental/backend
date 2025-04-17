// Unit tests for UserService using gomock and testify.
package services

import (
	"errors"
	"testing"
	"github.com/SimpleBookRental/backend/internal/models"
	"github.com/SimpleBookRental/backend/internal/mocks"
	"github.com/SimpleBookRental/backend/pkg/utils"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestUserService_Create_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockUserRepo := mocks.NewMockUserRepositoryInterface(ctrl)
	service := NewUserService(mockUserRepo)

	userCreate := &models.UserCreate{
		Name:     "Test User",
		Email:    "test@example.com",
		Password: "password123",
	}
	mockUserRepo.EXPECT().FindByEmail(userCreate.Email).Return(nil, nil)
	mockUserRepo.EXPECT().Create(gomock.Any()).Return(nil)

	user, err := service.Create(userCreate)
	assert.NoError(t, err)
	assert.Equal(t, userCreate.Email, user.Email)
}

func TestUserService_Create_EmailExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockUserRepo := mocks.NewMockUserRepositoryInterface(ctrl)
	service := NewUserService(mockUserRepo)

	userCreate := &models.UserCreate{
		Name:     "Test User",
		Email:    "test@example.com",
		Password: "password123",
	}
	mockUserRepo.EXPECT().FindByEmail(userCreate.Email).Return(&models.User{}, nil)

	user, err := service.Create(userCreate)
	assert.Nil(t, user)
	assert.ErrorContains(t, err, "email already exists")
}

func TestUserService_Create_RepoError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockUserRepo := mocks.NewMockUserRepositoryInterface(ctrl)
	service := NewUserService(mockUserRepo)

	userCreate := &models.UserCreate{
		Name:     "Test User",
		Email:    "test@example.com",
		Password: "password123",
	}
	mockUserRepo.EXPECT().FindByEmail(userCreate.Email).Return(nil, errors.New("db error"))

	user, err := service.Create(userCreate)
	assert.Nil(t, user)
	assert.ErrorContains(t, err, "error checking existing user")
}

func TestUserService_GetByID_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockUserRepo := mocks.NewMockUserRepositoryInterface(ctrl)
	service := NewUserService(mockUserRepo)

	user := &models.User{ID: "11111111-1111-1111-1111-111111111111"}
	mockUserRepo.EXPECT().FindByID(user.ID).Return(user, nil)

	got, err := service.GetByID(user.ID)
	assert.NoError(t, err)
	assert.Equal(t, user.ID, got.ID)
}

func TestUserService_GetByID_InvalidID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockUserRepo := mocks.NewMockUserRepositoryInterface(ctrl)
	service := NewUserService(mockUserRepo)

	got, err := service.GetByID("invalid")
	assert.Nil(t, got)
	assert.ErrorContains(t, err, "invalid user ID")
}

func TestUserService_GetByID_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockUserRepo := mocks.NewMockUserRepositoryInterface(ctrl)
	service := NewUserService(mockUserRepo)

	mockUserRepo.EXPECT().FindByID("11111111-1111-1111-1111-111111111111").Return(nil, nil)

	got, err := service.GetByID("11111111-1111-1111-1111-111111111111")
	assert.Nil(t, got)
	assert.ErrorContains(t, err, "user not found")
}

func TestUserService_GetAll_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockUserRepo := mocks.NewMockUserRepositoryInterface(ctrl)
	service := NewUserService(mockUserRepo)

	users := []models.User{{ID: "1"}, {ID: "2"}}
	mockUserRepo.EXPECT().FindAll().Return(users, nil)

	got, err := service.GetAll()
	assert.NoError(t, err)
	assert.Equal(t, users, got)
}

func TestUserService_Update_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockUserRepo := mocks.NewMockUserRepositoryInterface(ctrl)
	service := NewUserService(mockUserRepo)

	user := &models.User{ID: validUUID(), Name: "Old"}
	userUpdate := &models.UserUpdate{Name: "New"}
	mockUserRepo.EXPECT().FindByID(user.ID).Return(user, nil)
	mockUserRepo.EXPECT().Update(user).Return(nil)

	got, err := service.Update(user.ID, userUpdate)
	assert.NoError(t, err)
	assert.Equal(t, "New", got.Name)
}

func TestUserService_Update_EmailExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockUserRepo := mocks.NewMockUserRepositoryInterface(ctrl)
	service := NewUserService(mockUserRepo)

	user := &models.User{ID: validUUID(), Email: "old@example.com"}
	userUpdate := &models.UserUpdate{Email: "new@example.com"}
	mockUserRepo.EXPECT().FindByID(user.ID).Return(user, nil)
	mockUserRepo.EXPECT().FindByEmail("new@example.com").Return(&models.User{}, nil)

	got, err := service.Update(user.ID, userUpdate)
	assert.Nil(t, got)
	assert.ErrorContains(t, err, "email already exists")
}

func TestUserService_Update_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockUserRepo := mocks.NewMockUserRepositoryInterface(ctrl)
	service := NewUserService(mockUserRepo)

	mockUserRepo.EXPECT().FindByID("11111111-1111-1111-1111-111111111111").Return(nil, nil)

	got, err := service.Update("11111111-1111-1111-1111-111111111111", &models.UserUpdate{Name: "New"})
	assert.Nil(t, got)
	assert.ErrorContains(t, err, "user not found")
}

func TestUserService_Delete_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockUserRepo := mocks.NewMockUserRepositoryInterface(ctrl)
	service := NewUserService(mockUserRepo)

	user := &models.User{ID: "11111111-1111-1111-1111-111111111111"}
	mockUserRepo.EXPECT().FindByID(user.ID).Return(user, nil)
	mockUserRepo.EXPECT().Delete(user.ID).Return(nil)

	err := service.Delete(user.ID)
	assert.NoError(t, err)
}

func TestUserService_Delete_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockUserRepo := mocks.NewMockUserRepositoryInterface(ctrl)
	service := NewUserService(mockUserRepo)

	mockUserRepo.EXPECT().FindByID("11111111-1111-1111-1111-111111111111").Return(nil, nil)

	err := service.Delete("11111111-1111-1111-1111-111111111111")
	assert.ErrorContains(t, err, "user not found")
}

func TestUserService_Login_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockUserRepo := mocks.NewMockUserRepositoryInterface(ctrl)
	mockTokenRepo := mocks.NewMockTokenRepositoryInterface(ctrl)
	service := NewUserService(mockUserRepo)

	hashed, _ := utils.HashPassword("password123")
	user := &models.User{
		ID:       "11111111-1111-1111-1111-111111111111",
		Email:    "test@example.com",
		Password: hashed,
	}
	userLogin := &models.UserLogin{
		Email:    "test@example.com",
		Password: "password123",
	}
	mockUserRepo.EXPECT().FindByEmail(userLogin.Email).Return(user, nil)
	mockTokenRepo.EXPECT().CreateToken(gomock.Any()).Return(nil).AnyTimes()

	resp, err := service.Login(userLogin, mockTokenRepo)
	assert.NoError(t, err)
	assert.NotEmpty(t, resp.AccessToken)
	assert.NotEmpty(t, resp.RefreshToken)
}

func TestUserService_Login_InvalidPassword(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockUserRepo := mocks.NewMockUserRepositoryInterface(ctrl)
	mockTokenRepo := mocks.NewMockTokenRepositoryInterface(ctrl)
	service := NewUserService(mockUserRepo)

	user := &models.User{
		ID:       validUUID(),
		Email:    "test@example.com",
		Password: "$2a$10$7Qw1Qw1Qw1Qw1Qw1Qw1QwOeQw1Qw1Qw1Qw1Qw1Qw1Qw1Qw1Qw1Qw1Q", // bcrypt hash for "password123"
	}
	userLogin := &models.UserLogin{
		Email:    "test@example.com",
		Password: "wrongpassword",
	}
	mockUserRepo.EXPECT().FindByEmail(userLogin.Email).Return(user, nil)

	resp, err := service.Login(userLogin, mockTokenRepo)
	assert.Nil(t, resp)
	assert.ErrorContains(t, err, "invalid email or password")
}

func TestUserService_Login_UserNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockUserRepo := mocks.NewMockUserRepositoryInterface(ctrl)
	mockTokenRepo := mocks.NewMockTokenRepositoryInterface(ctrl)
	service := NewUserService(mockUserRepo)

	userLogin := &models.UserLogin{
		Email:    "notfound@example.com",
		Password: "password123",
	}
	mockUserRepo.EXPECT().FindByEmail(userLogin.Email).Return(nil, nil)

	resp, err := service.Login(userLogin, mockTokenRepo)
	assert.Nil(t, resp)
	assert.ErrorContains(t, err, "invalid email or password")
}
