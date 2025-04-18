// Unit tests for UserController using gomock and testify.
package controllers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/SimpleBookRental/backend/internal/middleware"
	"github.com/SimpleBookRental/backend/internal/mocks"
	"github.com/SimpleBookRental/backend/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func setupGinUser() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.New()
}

func TestUserController_Register_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := mocks.NewMockUserServiceInterface(ctrl)
	mockTokenRepo := mocks.NewMockTokenRepositoryInterface(ctrl)
	controller := NewUserController(mockService, mockTokenRepo)
	router := setupGinUser()
	router.POST("/register", controller.Register)

	reqBody := models.UserCreate{
		Name:     "Test User",
		Email:    "test@example.com",
		Password: "password123",
	}
	expectedUser := &models.User{
		ID:    "11111111-1111-1111-1111-111111111111",
		Name:  "Test User",
		Email: "test@example.com",
	}
	mockService.EXPECT().Register(&reqBody).Return(expectedUser, nil)

	body, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestUserController_Register_BadRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := mocks.NewMockUserServiceInterface(ctrl)
	mockTokenRepo := mocks.NewMockTokenRepositoryInterface(ctrl)
	controller := NewUserController(mockService, mockTokenRepo)
	router := setupGinUser()
	router.POST("/register", controller.Register)

	// Invalid JSON
	req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer([]byte("{invalid")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUserController_Register_ServiceError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := mocks.NewMockUserServiceInterface(ctrl)
	mockTokenRepo := mocks.NewMockTokenRepositoryInterface(ctrl)
	controller := NewUserController(mockService, mockTokenRepo)
	router := setupGinUser()
	router.POST("/register", controller.Register)

	reqBody := models.UserCreate{
		Name:     "Test User",
		Email:    "test@example.com",
		Password: "password123",
	}
	mockService.EXPECT().Register(&reqBody).Return(nil, errors.New("service error"))

	body, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUserController_GetByID_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := mocks.NewMockUserServiceInterface(ctrl)
	mockTokenRepo := mocks.NewMockTokenRepositoryInterface(ctrl)
	controller := NewUserController(mockService, mockTokenRepo)
	router := setupGinUser()
	router.GET("/users/:id", controller.GetByID)

	expectedUser := &models.User{
		ID:    "11111111-1111-1111-1111-111111111111",
		Name:  "Test User",
		Email: "test@example.com",
	}
	mockService.EXPECT().GetByID("11111111-1111-1111-1111-111111111111").Return(expectedUser, nil)

	req, _ := http.NewRequest("GET", "/users/11111111-1111-1111-1111-111111111111", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestUserController_GetByID_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := mocks.NewMockUserServiceInterface(ctrl)
	mockTokenRepo := mocks.NewMockTokenRepositoryInterface(ctrl)
	controller := NewUserController(mockService, mockTokenRepo)
	router := setupGinUser()
	router.GET("/users/:id", controller.GetByID)

	mockService.EXPECT().GetByID("notfound").Return(nil, errors.New("not found"))

	req, _ := http.NewRequest("GET", "/users/notfound", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestUserController_GetAll_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := mocks.NewMockUserServiceInterface(ctrl)
	mockTokenRepo := mocks.NewMockTokenRepositoryInterface(ctrl)
	controller := NewUserController(mockService, mockTokenRepo)
	router := setupGinUser()
	router.GET("/users", controller.GetAll)

	users := []models.User{
		{ID: "1", Name: "A"},
		{ID: "2", Name: "B"},
	}
	mockService.EXPECT().GetAll().Return(users, nil)

	req, _ := http.NewRequest("GET", "/users", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestUserController_GetAll_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := mocks.NewMockUserServiceInterface(ctrl)
	mockTokenRepo := mocks.NewMockTokenRepositoryInterface(ctrl)
	controller := NewUserController(mockService, mockTokenRepo)
	router := setupGinUser()
	router.GET("/users", controller.GetAll)

	mockService.EXPECT().GetAll().Return(nil, errors.New("db error"))

	req, _ := http.NewRequest("GET", "/users", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestUserController_Update_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := mocks.NewMockUserServiceInterface(ctrl)
	mockTokenRepo := mocks.NewMockTokenRepositoryInterface(ctrl)
	controller := NewUserController(mockService, mockTokenRepo)
	router := setupGinUser()
	router.PUT("/users/:id", controller.Update)

	expectedUser := &models.User{
		ID:    "11111111-1111-1111-1111-111111111111",
		Name:  "Updated User",
		Email: "updated@example.com",
	}
	mockService.EXPECT().Update("11111111-1111-1111-1111-111111111111", &models.UserUpdate{Name: "Updated User"}).Return(expectedUser, nil)

	update := models.UserUpdate{Name: "Updated User"}
	body, _ := json.Marshal(update)
	req, _ := http.NewRequest("PUT", "/users/11111111-1111-1111-1111-111111111111", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestUserController_Update_BadRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := mocks.NewMockUserServiceInterface(ctrl)
	mockTokenRepo := mocks.NewMockTokenRepositoryInterface(ctrl)
	controller := NewUserController(mockService, mockTokenRepo)
	router := setupGinUser()
	router.PUT("/users/:id", controller.Update)

	// Invalid JSON
	req, _ := http.NewRequest("PUT", "/users/11111111-1111-1111-1111-111111111111", bytes.NewBuffer([]byte("{invalid")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUserController_Update_ServiceError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := mocks.NewMockUserServiceInterface(ctrl)
	mockTokenRepo := mocks.NewMockTokenRepositoryInterface(ctrl)
	controller := NewUserController(mockService, mockTokenRepo)
	router := setupGinUser()
	router.PUT("/users/:id", controller.Update)

	mockService.EXPECT().Update("11111111-1111-1111-1111-111111111111", &models.UserUpdate{Name: "Updated User"}).Return(nil, errors.New("update error"))

	update := models.UserUpdate{Name: "Updated User"}
	body, _ := json.Marshal(update)
	req, _ := http.NewRequest("PUT", "/users/11111111-1111-1111-1111-111111111111", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUserController_Delete_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := mocks.NewMockUserServiceInterface(ctrl)
	mockTokenRepo := mocks.NewMockTokenRepositoryInterface(ctrl)
	controller := NewUserController(mockService, mockTokenRepo)
	router := setupGinUser()
	router.DELETE("/users/:id", controller.Delete)

	mockService.EXPECT().Delete("11111111-1111-1111-1111-111111111111").Return(nil)

	req, _ := http.NewRequest("DELETE", "/users/11111111-1111-1111-1111-111111111111", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestUserController_Delete_ServiceError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := mocks.NewMockUserServiceInterface(ctrl)
	mockTokenRepo := mocks.NewMockTokenRepositoryInterface(ctrl)
	controller := NewUserController(mockService, mockTokenRepo)
	router := setupGinUser()
	router.DELETE("/users/:id", controller.Delete)

	mockService.EXPECT().Delete("11111111-1111-1111-1111-111111111111").Return(errors.New("delete error"))

	req, _ := http.NewRequest("DELETE", "/users/11111111-1111-1111-1111-111111111111", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUserController_Login_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := mocks.NewMockUserServiceInterface(ctrl)
	mockTokenRepo := mocks.NewMockTokenRepositoryInterface(ctrl)
	controller := NewUserController(mockService, mockTokenRepo)
	router := setupGinUser()
	router.POST("/users/login", controller.Login)

	reqBody := models.UserLogin{
		Email:    "test@example.com",
		Password: "password123",
	}
	expectedResp := &models.LoginResponse{
		AccessToken:  "access-token",
		RefreshToken: "refresh-token",
		ExpiresAt:    1234567890,
	}
	mockService.EXPECT().Login(&reqBody, mockTokenRepo).Return(expectedResp, nil)

	body, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/users/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestUserController_Login_BadRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := mocks.NewMockUserServiceInterface(ctrl)
	mockTokenRepo := mocks.NewMockTokenRepositoryInterface(ctrl)
	controller := NewUserController(mockService, mockTokenRepo)
	router := setupGinUser()
	router.POST("/users/login", controller.Login)

	// Invalid JSON
	req, _ := http.NewRequest("POST", "/users/login", bytes.NewBuffer([]byte("{invalid")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUserController_Login_ServiceError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := mocks.NewMockUserServiceInterface(ctrl)
	mockTokenRepo := mocks.NewMockTokenRepositoryInterface(ctrl)
	controller := NewUserController(mockService, mockTokenRepo)
	router := setupGinUser()
	router.POST("/users/login", controller.Login)

	reqBody := models.UserLogin{
		Email:    "test@example.com",
		Password: "password123",
	}
	mockService.EXPECT().Login(&reqBody, mockTokenRepo).Return(nil, errors.New("login error"))

	body, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/users/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUserController_GetByID_Forbidden(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := mocks.NewMockUserServiceInterface(ctrl)
	mockTokenRepo := mocks.NewMockTokenRepositoryInterface(ctrl)
	controller := NewUserController(mockService, mockTokenRepo)
	router := setupGinUser()
	// Setup middleware RequireAdminOrSameUser
	router.GET("/users/:id", func(c *gin.Context) {
		c.Set("user_id", "other-user-id")
		c.Set("role", "USER")
	}, 
	middleware.RequireAdminOrSameUser(), controller.GetByID)

	// No need mockService.EXPECT().GetByID, middleware reject this request before it hit controller.

	req, _ := http.NewRequest("GET", "/users/target-user-id", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
}
