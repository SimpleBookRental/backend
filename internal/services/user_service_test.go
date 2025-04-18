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

func TestUserService_Register_InvalidInput(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockUserRepo := mocks.NewMockUserRepositoryInterface(ctrl)
	mockBookRepo := mocks.NewMockBookRepositoryInterface(ctrl)
	mockTokenRepo := mocks.NewMockTokenRepositoryInterface(ctrl)
	service := NewUserService(mockUserRepo, mockBookRepo, mockTokenRepo)

	// Empty email
	userCreate := &models.UserCreate{
		Name:     "Test User",
		Email:    "",
		Password: "password123",
	}
	user, err := service.Register(userCreate)
	assert.Nil(t, user)
	assert.Error(t, err)

	// Empty name
	userCreate = &models.UserCreate{
		Name:     "",
		Email:    "test@example.com",
		Password: "password123",
	}
	user, err = service.Register(userCreate)
	assert.Nil(t, user)
	assert.Error(t, err)

	// Empty password
	userCreate = &models.UserCreate{
		Name:     "Test User",
		Email:    "test@example.com",
		Password: "",
	}
	user, err = service.Register(userCreate)
	assert.Nil(t, user)
	assert.Error(t, err)
}

func TestUserService_Update_InvalidID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockUserRepo := mocks.NewMockUserRepositoryInterface(ctrl)
	mockBookRepo := mocks.NewMockBookRepositoryInterface(ctrl)
	mockTokenRepo := mocks.NewMockTokenRepositoryInterface(ctrl)
	service := NewUserService(mockUserRepo, mockBookRepo, mockTokenRepo)

	userUpdate := &models.UserUpdate{Name: "New Name"}
	user, err := service.Update("invalid-uuid", userUpdate)
	assert.Nil(t, user)
	assert.ErrorContains(t, err, "invalid user ID")
}

func TestUserService_Update_UserNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockUserRepo := mocks.NewMockUserRepositoryInterface(ctrl)
	mockBookRepo := mocks.NewMockBookRepositoryInterface(ctrl)
	mockTokenRepo := mocks.NewMockTokenRepositoryInterface(ctrl)
	service := NewUserService(mockUserRepo, mockBookRepo, mockTokenRepo)

	mockUserRepo.EXPECT().FindByID("11111111-1111-1111-1111-111111111111").Return(nil, nil)
	userUpdate := &models.UserUpdate{Name: "New Name"}
	user, err := service.Update("11111111-1111-1111-1111-111111111111", userUpdate)
	assert.Nil(t, user)
	assert.ErrorContains(t, err, "user not found")
}

func TestUserService_Update_EmailExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockUserRepo := mocks.NewMockUserRepositoryInterface(ctrl)
	mockBookRepo := mocks.NewMockBookRepositoryInterface(ctrl)
	mockTokenRepo := mocks.NewMockTokenRepositoryInterface(ctrl)
	service := NewUserService(mockUserRepo, mockBookRepo, mockTokenRepo)

	existing := &models.User{ID: "11111111-1111-1111-1111-111111111111", Email: "old@example.com"}
	mockUserRepo.EXPECT().FindByID(existing.ID).Return(existing, nil)
	mockUserRepo.EXPECT().FindByEmail("new@example.com").Return(&models.User{}, nil)
	userUpdate := &models.UserUpdate{Email: "new@example.com"}
	user, err := service.Update(existing.ID, userUpdate)
	assert.Nil(t, user)
	assert.ErrorContains(t, err, "email already exists")
}

func TestUserService_Update_RepoError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockUserRepo := mocks.NewMockUserRepositoryInterface(ctrl)
	mockBookRepo := mocks.NewMockBookRepositoryInterface(ctrl)
	mockTokenRepo := mocks.NewMockTokenRepositoryInterface(ctrl)
	service := NewUserService(mockUserRepo, mockBookRepo, mockTokenRepo)

	existing := &models.User{ID: "11111111-1111-1111-1111-111111111111", Email: "test@example.com"}
	mockUserRepo.EXPECT().FindByID(existing.ID).Return(existing, nil)
	mockUserRepo.EXPECT().Update(existing).Return(errors.New("update error"))
	userUpdate := &models.UserUpdate{Name: "New Name"}
	user, err := service.Update(existing.ID, userUpdate)
	assert.Nil(t, user)
	assert.ErrorContains(t, err, "error updating user")
}

func TestUserService_Delete_InvalidID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockUserRepo := mocks.NewMockUserRepositoryInterface(ctrl)
	mockBookRepo := mocks.NewMockBookRepositoryInterface(ctrl)
	mockTokenRepo := mocks.NewMockTokenRepositoryInterface(ctrl)
	service := NewUserService(mockUserRepo, mockBookRepo, mockTokenRepo)

	err := service.Delete("invalid-uuid")
	assert.ErrorContains(t, err, "invalid user ID")
}

func TestUserService_Delete_UserNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockUserRepo := mocks.NewMockUserRepositoryInterface(ctrl)
	mockBookRepo := mocks.NewMockBookRepositoryInterface(ctrl)
	mockTokenRepo := mocks.NewMockTokenRepositoryInterface(ctrl)
	service := NewUserService(mockUserRepo, mockBookRepo, mockTokenRepo)

	mockUserRepo.EXPECT().FindByID("11111111-1111-1111-1111-111111111111").Return(nil, nil)
	err := service.Delete("11111111-1111-1111-1111-111111111111")
	assert.ErrorContains(t, err, "user not found")
}

func TestUserService_Delete_AdminForbidden(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockUserRepo := mocks.NewMockUserRepositoryInterface(ctrl)
	mockBookRepo := mocks.NewMockBookRepositoryInterface(ctrl)
	mockTokenRepo := mocks.NewMockTokenRepositoryInterface(ctrl)
	service := NewUserService(mockUserRepo, mockBookRepo, mockTokenRepo)

	admin := &models.User{ID: "11111111-1111-1111-1111-111111111111", Role: "ADMIN"}
	mockUserRepo.EXPECT().FindByID(admin.ID).Return(admin, nil)
	err := service.Delete(admin.ID)
	assert.ErrorContains(t, err, "cannot delete user with ADMIN role")
}

func TestUserService_Delete_RepoError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockUserRepo := mocks.NewMockUserRepositoryInterface(ctrl)
	mockBookRepo := mocks.NewMockBookRepositoryInterface(ctrl)
	mockTokenRepo := mocks.NewMockTokenRepositoryInterface(ctrl)
	service := NewUserService(mockUserRepo, mockBookRepo, mockTokenRepo)

	user := &models.User{ID: "11111111-1111-1111-1111-111111111112", Role: "USER"}
	mockUserRepo.EXPECT().FindByID(user.ID).Return(user, nil)
	mockBookRepo.EXPECT().DeleteByUserID(user.ID).Return(nil)
	mockTokenRepo.EXPECT().DeleteByUserID(user.ID).Return(nil)
	mockUserRepo.EXPECT().Delete(user.ID).Return(errors.New("delete error"))
	err := service.Delete(user.ID)
	assert.ErrorContains(t, err, "delete error")
}
