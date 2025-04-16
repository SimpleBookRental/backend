package services

import (
	"errors"
	"testing"

	"github.com/SimpleBookRental/backend/internal/mocks"
	"github.com/SimpleBookRental/backend/internal/models"
	"github.com/SimpleBookRental/backend/pkg/utils"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestUserService_Create(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepositoryInterface(ctrl)
	service := NewUserService(mockRepo)

	userCreate := &models.UserCreate{
		Name:     "Test User",
		Email:    "test@example.com",
		Password: "password123",
	}

	// Expectations
	mockRepo.EXPECT().FindByEmail(userCreate.Email).Return(nil, nil)
	mockRepo.EXPECT().Create(gomock.AssignableToTypeOf(&models.User{})).Return(nil)

	// Test
	user, err := service.Create(userCreate)
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, userCreate.Name, user.Name)
	assert.Equal(t, userCreate.Email, user.Email)
	assert.Equal(t, models.UserRole, user.Role)            // Role should be USER by default
	assert.NotEqual(t, userCreate.Password, user.Password) // Password should be hashed

	// Verify expectations handled by gomock controller
}

func TestUserService_Create_EmailExists(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepositoryInterface(ctrl)
	service := NewUserService(mockRepo)

	userCreate := &models.UserCreate{
		Name:     "Test User",
		Email:    "test@example.com",
		Password: "password123",
	}

	existingUser := &models.User{
		ID:    "123e4567-e89b-12d3-a456-426614174002",
		Email: userCreate.Email,
	}

	// Expectations
	mockRepo.EXPECT().FindByEmail(userCreate.Email).Return(existingUser, nil)

	// Test
	user, err := service.Create(userCreate)
	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Contains(t, err.Error(), "email already exists")

	// Verify expectations handled by gomock controller
}

func TestUserService_GetByID(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepositoryInterface(ctrl)
	service := NewUserService(mockRepo)

	user := &models.User{
		ID:    "123e4567-e89b-12d3-a456-426614174001",
		Name:  "Test User",
		Email: "test@example.com",
		Role:  models.UserRole,
	}

	// Expectations
	mockRepo.EXPECT().FindByID(user.ID).Return(user, nil)

	// Test
	result, err := service.GetByID(user.ID)
	assert.NoError(t, err)
	assert.Equal(t, user, result)

	// Verify expectations handled by gomock controller
}

func TestUserService_GetByID_InvalidID(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepositoryInterface(ctrl)
	service := NewUserService(mockRepo)

	// Test
	result, err := service.GetByID("invalid-id")
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "invalid user ID")

	// Verify expectations handled by gomock controller
}

func TestUserService_GetByID_NotFound(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepositoryInterface(ctrl)
	service := NewUserService(mockRepo)

	validID := "123e4567-e89b-12d3-a456-426614174001"

	// Expectations
	mockRepo.EXPECT().FindByID(validID).Return(nil, nil)

	// Test
	result, err := service.GetByID(validID)
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "user not found")

	// Verify expectations handled by gomock controller
}

func TestUserService_GetAll(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepositoryInterface(ctrl)
	service := NewUserService(mockRepo)

	users := []models.User{
		{
			ID:    "123e4567-e89b-12d3-a456-426614174001",
			Name:  "User 1",
			Email: "user1@example.com",
			Role:  models.UserRole,
		},
		{
			ID:    "123e4567-e89b-12d3-a456-426614174002",
			Name:  "User 2",
			Email: "user2@example.com",
			Role:  models.UserRole,
		},
	}

	// Expectations
	mockRepo.EXPECT().FindAll().Return(users, nil)

	// Test
	results, err := service.GetAll()
	assert.NoError(t, err)
	assert.Equal(t, users, results)

	// Verify expectations handled by gomock controller
}

func TestUserService_Update(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepositoryInterface(ctrl)
	service := NewUserService(mockRepo)

	userID := "123e4567-e89b-12d3-a456-426614174001"
	userUpdate := &models.UserUpdate{
		Name:     "Updated User",
		Email:    "updated@example.com",
		Password: "newpassword123",
	}

	existingUser := &models.User{
		ID:       userID,
		Name:     "Test User",
		Email:    "test@example.com",
		Password: "hashedpassword",
		Role:     models.UserRole,
	}

	// Expectations
	mockRepo.EXPECT().FindByID(userID).Return(existingUser, nil)
	mockRepo.EXPECT().FindByEmail(userUpdate.Email).Return(nil, nil)
	mockRepo.EXPECT().Update(gomock.AssignableToTypeOf(&models.User{})).Return(nil)

	// Test
	updatedUser, err := service.Update(userID, userUpdate)
	assert.NoError(t, err)
	assert.Equal(t, userID, updatedUser.ID)
	assert.Equal(t, userUpdate.Name, updatedUser.Name)
	assert.Equal(t, userUpdate.Email, updatedUser.Email)
	assert.NotEqual(t, userUpdate.Password, updatedUser.Password) // Password should be hashed

	// Verify expectations handled by gomock controller
}

func TestUserService_Update_InvalidID(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepositoryInterface(ctrl)
	service := NewUserService(mockRepo)

	userUpdate := &models.UserUpdate{
		Name:  "Updated User",
		Email: "updated@example.com",
	}

	// Test
	updatedUser, err := service.Update("invalid-id", userUpdate)
	assert.Error(t, err)
	assert.Nil(t, updatedUser)
	assert.Contains(t, err.Error(), "invalid user ID")

	// Verify expectations handled by gomock controller
}

func TestUserService_Update_UserNotFound(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepositoryInterface(ctrl)
	service := NewUserService(mockRepo)

	userID := "123e4567-e89b-12d3-a456-426614174001"
	userUpdate := &models.UserUpdate{
		Name:  "Updated User",
		Email: "updated@example.com",
	}

	// Expectations
	mockRepo.EXPECT().FindByID(userID).Return(nil, nil)

	// Test
	updatedUser, err := service.Update(userID, userUpdate)
	assert.Error(t, err)
	assert.Nil(t, updatedUser)
	assert.Contains(t, err.Error(), "user not found")

	// Verify expectations handled by gomock controller
}

func TestUserService_Update_EmailExists(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepositoryInterface(ctrl)
	service := NewUserService(mockRepo)

	userID := "123e4567-e89b-12d3-a456-426614174001"
	userUpdate := &models.UserUpdate{
		Name:  "Updated User",
		Email: "existing@example.com",
	}

	existingUser := &models.User{
		ID:    userID,
		Name:  "Test User",
		Email: "test@example.com",
		Role:  models.UserRole,
	}

	anotherUser := &models.User{
		ID:    "123e4567-e89b-12d3-a456-426614174002",
		Email: userUpdate.Email,
	}

	// Expectations
	mockRepo.EXPECT().FindByID(userID).Return(existingUser, nil)
	mockRepo.EXPECT().FindByEmail(userUpdate.Email).Return(anotherUser, nil)

	// Test
	updatedUser, err := service.Update(userID, userUpdate)
	assert.Error(t, err)
	assert.Nil(t, updatedUser)
	assert.Contains(t, err.Error(), "email already exists")

	// Verify expectations handled by gomock controller
}

func TestUserService_Delete(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepositoryInterface(ctrl)
	service := NewUserService(mockRepo)

	userID := "123e4567-e89b-12d3-a456-426614174001"
	user := &models.User{
		ID:   userID,
		Role: models.UserRole,
	}

	// Expectations
	mockRepo.EXPECT().FindByID(userID).Return(user, nil)
	mockRepo.EXPECT().Delete(userID).Return(nil)

	// Test
	err := service.Delete(userID)
	assert.NoError(t, err)

	// Verify expectations handled by gomock controller
}

func TestUserService_Delete_InvalidID(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepositoryInterface(ctrl)
	service := NewUserService(mockRepo)

	// Test
	err := service.Delete("invalid-id")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid user ID")

	// Verify expectations handled by gomock controller
}

func TestUserService_Delete_UserNotFound(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepositoryInterface(ctrl)
	service := NewUserService(mockRepo)

	userID := "123e4567-e89b-12d3-a456-426614174001"

	// Expectations
	mockRepo.EXPECT().FindByID(userID).Return(nil, nil)

	// Test
	err := service.Delete(userID)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "user not found")

	// Verify expectations handled by gomock controller
}

func TestUserService_Login(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepositoryInterface(ctrl)
	mockTokenRepo := mocks.NewMockTokenRepositoryInterface(ctrl)
	service := NewUserService(mockRepo)

	userLogin := &models.UserLogin{
		Email:    "test@example.com",
		Password: "password123",
	}

	// Create a user with a hashed password that will match
	hashedPassword, _ := utils.HashPassword(userLogin.Password)
	user := &models.User{
		ID:       "123e4567-e89b-12d3-a456-426614174001",
		Email:    userLogin.Email,
		Password: hashedPassword,
		Role:     models.UserRole,
	}

	// Expectations
	mockRepo.EXPECT().FindByEmail(userLogin.Email).Return(user, nil)
	mockTokenRepo.EXPECT().CreateToken(gomock.AssignableToTypeOf(&models.IssuedToken{})).Return(nil).Times(2)

	// Test
	response, err := service.Login(userLogin, mockTokenRepo)
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, user, response.User)
	assert.NotEmpty(t, response.AccessToken)
	assert.NotEmpty(t, response.RefreshToken)
	assert.NotZero(t, response.ExpiresAt)

	// Verify expectations handled by gomock controller
}

func TestUserService_Login_UserNotFound(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepositoryInterface(ctrl)
	mockTokenRepo := mocks.NewMockTokenRepositoryInterface(ctrl)
	service := NewUserService(mockRepo)

	userLogin := &models.UserLogin{
		Email:    "test@example.com",
		Password: "password123",
	}

	// Expectations
	mockRepo.EXPECT().FindByEmail(userLogin.Email).Return(nil, nil)

	// Test
	response, err := service.Login(userLogin, mockTokenRepo)
	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Contains(t, err.Error(), "invalid email or password")

	// Verify expectations handled by gomock controller
}

func TestUserService_Login_InvalidPassword(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepositoryInterface(ctrl)
	mockTokenRepo := mocks.NewMockTokenRepositoryInterface(ctrl)
	service := NewUserService(mockRepo)

	userLogin := &models.UserLogin{
		Email:    "test@example.com",
		Password: "password123",
	}

	// Create a user with a different password
	hashedPassword, _ := utils.HashPassword("differentpassword")
	user := &models.User{
		ID:       "123e4567-e89b-12d3-a456-426614174001",
		Email:    userLogin.Email,
		Password: hashedPassword,
	}

	// Expectations
	mockRepo.EXPECT().FindByEmail(userLogin.Email).Return(user, nil)

	// Test
	response, err := service.Login(userLogin, mockTokenRepo)
	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Contains(t, err.Error(), "invalid email or password")

	// Verify expectations handled by gomock controller
}

func TestUserService_Login_TokenCreationError(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepositoryInterface(ctrl)
	mockTokenRepo := mocks.NewMockTokenRepositoryInterface(ctrl)
	service := NewUserService(mockRepo)

	userLogin := &models.UserLogin{
		Email:    "test@example.com",
		Password: "password123",
	}

	// Create a user with a hashed password that will match
	hashedPassword, _ := utils.HashPassword(userLogin.Password)
	user := &models.User{
		ID:       "123e4567-e89b-12d3-a456-426614174001",
		Email:    userLogin.Email,
		Password: hashedPassword,
		Role:     models.UserRole,
	}

	// Expectations
	mockRepo.EXPECT().FindByEmail(userLogin.Email).Return(user, nil)
	mockTokenRepo.EXPECT().CreateToken(gomock.AssignableToTypeOf(&models.IssuedToken{})).Return(errors.New("token creation error"))

	// Test
	response, err := service.Login(userLogin, mockTokenRepo)
	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Contains(t, err.Error(), "error saving access token")

	// Verify expectations handled by gomock controller
}
