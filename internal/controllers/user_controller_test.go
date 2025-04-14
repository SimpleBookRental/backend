package controllers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/SimpleBookRental/backend/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockTokenRepository is a mock implementation of repositories.TokenRepositoryInterface
type MockTokenRepository struct {
	mock.Mock
}

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

func (m *MockTokenRepository) FindByISBN(isbn string) (*models.Book, error) {
	args := m.Called(isbn)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Book), args.Error(1)
}

func setupUserController() (*gin.Engine, *MockUserService, *MockTokenRepository) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	mockService := new(MockUserService)
	mockTokenRepo := new(MockTokenRepository)
	controller := &UserController{
		userService: mockService,
		tokenRepo:   mockTokenRepo,
	}

	// Setup routes
	v1 := router.Group("/api/v1")
	{
		v1.POST("/users", controller.Create)
		v1.GET("/users", controller.GetAll)
		v1.GET("/users/:id", controller.GetByID)
		v1.PUT("/users/:id", controller.Update)
		v1.DELETE("/users/:id", controller.Delete)
		v1.POST("/login", controller.Login)
	}

	return router, mockService, mockTokenRepo
}

func TestUserController_Create(t *testing.T) {
	// Setup
	router, mockService, _ := setupUserController()

	userCreate := &models.UserCreate{
		Name:     "Test User",
		Email:    "test@example.com",
		Password: "password123",
	}

	user := &models.User{
		ID:       "test-id",
		Name:     userCreate.Name,
		Email:    userCreate.Email,
		Password: "hashed_password",
	}

	// Expectations
	mockService.On("Create", mock.AnythingOfType("*models.UserCreate")).Return(user, nil)

	// Create request
	body, _ := json.Marshal(userCreate)
	req, _ := http.NewRequest("POST", "/api/v1/users", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Perform request
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusCreated, w.Code)
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, true, response["success"])
	assert.Equal(t, "User created successfully", response["message"])

	// Verify expectations
	mockService.AssertExpectations(t)
}

func TestUserController_Create_InvalidRequest(t *testing.T) {
	// Setup
	router, mockService, _ := setupUserController()

	// Invalid request (missing required fields)
	userCreate := map[string]interface{}{
		"name": "Test User",
		// Missing email and password
	}

	// Create request
	body, _ := json.Marshal(userCreate)
	req, _ := http.NewRequest("POST", "/api/v1/users", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Perform request
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, false, response["success"])
	assert.Equal(t, "Invalid request body", response["message"])

	// Verify expectations
	mockService.AssertNotCalled(t, "Create")
}

func TestUserController_Create_ServiceError(t *testing.T) {
	// Setup
	router, mockService, _ := setupUserController()

	userCreate := &models.UserCreate{
		Name:     "Test User",
		Email:    "test@example.com",
		Password: "password123",
	}

	// Expectations
	mockService.On("Create", mock.AnythingOfType("*models.UserCreate")).Return(nil, errors.New("service error"))

	// Create request
	body, _ := json.Marshal(userCreate)
	req, _ := http.NewRequest("POST", "/api/v1/users", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Perform request
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, false, response["success"])
	assert.Equal(t, "Failed to create user", response["message"])

	// Verify expectations
	mockService.AssertExpectations(t)
}

func TestUserController_GetAll(t *testing.T) {
	// Setup
	router, mockService, _ := setupUserController()

	users := []models.User{
		{
			ID:    "user-1",
			Name:  "User 1",
			Email: "user1@example.com",
		},
		{
			ID:    "user-2",
			Name:  "User 2",
			Email: "user2@example.com",
		},
	}

	// Expectations
	mockService.On("GetAll").Return(users, nil)

	// Create request
	req, _ := http.NewRequest("GET", "/api/v1/users", nil)
	w := httptest.NewRecorder()

	// Perform request
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, true, response["success"])
	assert.Equal(t, "Users retrieved successfully", response["message"])

	// Verify expectations
	mockService.AssertExpectations(t)
}

func TestUserController_GetAll_ServiceError(t *testing.T) {
	// Setup
	router, mockService, _ := setupUserController()

	// Expectations
	mockService.On("GetAll").Return([]models.User{}, errors.New("service error"))

	// Create request
	req, _ := http.NewRequest("GET", "/api/v1/users", nil)
	w := httptest.NewRecorder()

	// Perform request
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, false, response["success"])
	assert.Equal(t, "Failed to retrieve users", response["message"])

	// Verify expectations
	mockService.AssertExpectations(t)
}

func TestUserController_GetByID(t *testing.T) {
	// Setup
	router, mockService, _ := setupUserController()

	userID := "123e4567-e89b-12d3-a456-426614174001"
	user := &models.User{
		ID:    userID,
		Name:  "Test User",
		Email: "test@example.com",
	}

	// Expectations
	mockService.On("GetByID", userID).Return(user, nil)

	// Create request
	req, _ := http.NewRequest("GET", "/api/v1/users/"+userID, nil)
	w := httptest.NewRecorder()

	// Perform request
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, true, response["success"])
	assert.Equal(t, "User retrieved successfully", response["message"])

	// Verify expectations
	mockService.AssertExpectations(t)
}

func TestUserController_GetByID_ServiceError(t *testing.T) {
	// Setup
	router, mockService, _ := setupUserController()

	userID := "123e4567-e89b-12d3-a456-426614174001"

	// Expectations
	mockService.On("GetByID", userID).Return(nil, errors.New("service error"))

	// Create request
	req, _ := http.NewRequest("GET", "/api/v1/users/"+userID, nil)
	w := httptest.NewRecorder()

	// Perform request
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusNotFound, w.Code)
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, false, response["success"])
	assert.Equal(t, "service error", response["message"])

	// Verify expectations
	mockService.AssertExpectations(t)
}

func TestUserController_Update(t *testing.T) {
	// Setup
	router, mockService, _ := setupUserController()

	userID := "123e4567-e89b-12d3-a456-426614174001"
	userUpdate := &models.UserUpdate{
		Name:     "Updated User",
		Email:    "updated@example.com",
		Password: "newpassword123",
	}

	updatedUser := &models.User{
		ID:       userID,
		Name:     userUpdate.Name,
		Email:    userUpdate.Email,
		Password: "hashed_password",
	}

	// Expectations
	mockService.On("Update", userID, mock.AnythingOfType("*models.UserUpdate")).Return(updatedUser, nil)

	// Create request
	body, _ := json.Marshal(userUpdate)
	req, _ := http.NewRequest("PUT", "/api/v1/users/"+userID, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Perform request
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, true, response["success"])
	assert.Equal(t, "User updated successfully", response["message"])

	// Verify expectations
	mockService.AssertExpectations(t)
}

func TestUserController_Update_InvalidRequest(t *testing.T) {
	// Setup
	router, mockService, _ := setupUserController()

	userID := "test-id"

	// Invalid request (invalid JSON)
	req, _ := http.NewRequest("PUT", "/api/v1/users/"+userID, bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Perform request
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, false, response["success"])
	assert.Equal(t, "Invalid request body", response["message"])

	// Verify expectations
	mockService.AssertNotCalled(t, "Update")
}

func TestUserController_Update_ServiceError(t *testing.T) {
	// Setup
	router, mockService, _ := setupUserController()

	userID := "123e4567-e89b-12d3-a456-426614174001"
	userUpdate := &models.UserUpdate{
		Name:     "Updated User",
		Email:    "updated@example.com",
		Password: "newpassword123",
	}

	// Expectations
	mockService.On("Update", userID, mock.AnythingOfType("*models.UserUpdate")).Return(nil, errors.New("service error"))

	// Create request
	body, _ := json.Marshal(userUpdate)
	req, _ := http.NewRequest("PUT", "/api/v1/users/"+userID, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Perform request
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, false, response["success"])
	assert.Equal(t, "Failed to update user", response["message"])

	// Verify expectations
	mockService.AssertExpectations(t)
}

func TestUserController_Delete(t *testing.T) {
	// Setup
	router, mockService, _ := setupUserController()

	userID := "123e4567-e89b-12d3-a456-426614174001"

	// Expectations
	mockService.On("Delete", userID).Return(nil)

	// Create request
	req, _ := http.NewRequest("DELETE", "/api/v1/users/"+userID, nil)
	w := httptest.NewRecorder()

	// Perform request
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, true, response["success"])
	assert.Equal(t, "User deleted successfully", response["message"])

	// Verify expectations
	mockService.AssertExpectations(t)
}

func TestUserController_Delete_ServiceError(t *testing.T) {
	// Setup
	router, mockService, _ := setupUserController()

	userID := "123e4567-e89b-12d3-a456-426614174001"

	// Expectations
	mockService.On("Delete", userID).Return(errors.New("service error"))

	// Create request
	req, _ := http.NewRequest("DELETE", "/api/v1/users/"+userID, nil)
	w := httptest.NewRecorder()

	// Perform request
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, false, response["success"])
	assert.Equal(t, "Failed to delete user", response["message"])

	// Verify expectations
	mockService.AssertExpectations(t)
}

func TestUserController_Login(t *testing.T) {
	// Setup
	router, mockService, mockTokenRepo := setupUserController()

	userLogin := &models.UserLogin{
		Email:    "test@example.com",
		Password: "password123",
	}

	loginResponse := &models.LoginResponse{
		User: &models.User{
			ID:    "test-id",
			Name:  "Test User",
			Email: userLogin.Email,
		},
		AccessToken:  "access-token",
		RefreshToken: "refresh-token",
		ExpiresAt:    time.Now().Unix(),
	}

	// Expectations
	mockService.On("Login", mock.AnythingOfType("*models.UserLogin"), mockTokenRepo).Return(loginResponse, nil)

	// Create request
	body, _ := json.Marshal(userLogin)
	req, _ := http.NewRequest("POST", "/api/v1/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Perform request
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, true, response["success"])
	assert.Equal(t, "Login successful", response["message"])

	// Verify expectations
	mockService.AssertExpectations(t)
}

func TestUserController_Login_InvalidRequest(t *testing.T) {
	// Setup
	router, mockService, _ := setupUserController()

	// Invalid request (missing required fields)
	userLogin := map[string]interface{}{
		"email": "test@example.com",
		// Missing password
	}

	// Create request
	body, _ := json.Marshal(userLogin)
	req, _ := http.NewRequest("POST", "/api/v1/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Perform request
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, false, response["success"])
	assert.Equal(t, "Invalid request body", response["message"])

	// Verify expectations
	mockService.AssertNotCalled(t, "Login")
}

func TestUserController_Login_ServiceError(t *testing.T) {
	// Setup
	router, mockService, mockTokenRepo := setupUserController()

	userLogin := &models.UserLogin{
		Email:    "test@example.com",
		Password: "password123",
	}

	// Expectations
	mockService.On("Login", mock.AnythingOfType("*models.UserLogin"), mockTokenRepo).Return(nil, errors.New("service error"))

	// Create request
	body, _ := json.Marshal(userLogin)
	req, _ := http.NewRequest("POST", "/api/v1/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Perform request
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, false, response["success"])
	assert.Equal(t, "Login failed", response["message"])

	// Verify expectations
	mockService.AssertExpectations(t)
}
