package controllers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/SimpleBookRental/backend/internal/mocks"
	"github.com/SimpleBookRental/backend/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func setupTokenController(t *testing.T) (*gin.Engine, *mocks.MockTokenServiceInterface) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	ctrl := gomock.NewController(t)
	mockService := mocks.NewMockTokenServiceInterface(ctrl)
	controller := &TokenController{
		tokenService: mockService,
	}

	// Setup routes
	v1 := router.Group("/api/v1")
	{
		v1.POST("/refresh-token", controller.RefreshToken)
		v1.POST("/logout", controller.Logout)
	}

	return router, mockService
}

func TestTokenController_RefreshToken(t *testing.T) {
	// Setup
	router, mockService := setupTokenController(t)

	refreshTokenRequest := &models.RefreshTokenRequest{
		RefreshToken: "refresh-token",
	}

	refreshTokenResponse := &models.RefreshTokenResponse{
		AccessToken:  "new-access-token",
		RefreshToken: "new-refresh-token",
		ExpiresAt:    time.Now().Unix(),
	}

	// Expectations
	mockService.EXPECT().RefreshToken(gomock.AssignableToTypeOf(&models.RefreshTokenRequest{})).Return(refreshTokenResponse, nil)

	// Create request
	body, _ := json.Marshal(refreshTokenRequest)
	req, _ := http.NewRequest("POST", "/api/v1/refresh-token", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Perform request
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, true, response["success"])
	assert.Equal(t, "Token refreshed successfully", response["message"])

	// Verify expectations handled by gomock controller
}

func TestTokenController_RefreshToken_InvalidRequest(t *testing.T) {
	// Setup
	router, _ := setupTokenController(t)

	// Invalid request (missing required fields)
	refreshTokenRequest := map[string]interface{}{
		// Missing refresh_token
	}

	// Create request
	body, _ := json.Marshal(refreshTokenRequest)
	req, _ := http.NewRequest("POST", "/api/v1/refresh-token", bytes.NewBuffer(body))
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

	// No expectations to verify as the service should not be called
}

func TestTokenController_RefreshToken_ServiceError(t *testing.T) {
	// Setup
	router, mockService := setupTokenController(t)

	refreshTokenRequest := &models.RefreshTokenRequest{
		RefreshToken: "refresh-token",
	}

	// Expectations
	mockService.EXPECT().RefreshToken(gomock.AssignableToTypeOf(&models.RefreshTokenRequest{})).Return(nil, errors.New("service error"))

	// Create request
	body, _ := json.Marshal(refreshTokenRequest)
	req, _ := http.NewRequest("POST", "/api/v1/refresh-token", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Perform request
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, false, response["success"])
	assert.Equal(t, "Failed to refresh token", response["message"])

	// Verify expectations handled by gomock controller
}

func TestTokenController_Logout(t *testing.T) {
	// Setup
	router, mockService := setupTokenController(t)

	logoutRequest := &models.LogoutRequest{
		RefreshToken: "refresh-token",
	}

	// Expectations
	mockService.EXPECT().Logout(gomock.AssignableToTypeOf(&models.LogoutRequest{})).Return(nil)

	// Create request
	body, _ := json.Marshal(logoutRequest)
	req, _ := http.NewRequest("POST", "/api/v1/logout", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Perform request
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, true, response["success"])
	assert.Equal(t, "Logged out successfully", response["message"])

	// Verify expectations handled by gomock controller
}

func TestTokenController_Logout_InvalidRequest(t *testing.T) {
	// Setup
	router, _ := setupTokenController(t)

	// Invalid request (missing required fields)
	logoutRequest := map[string]interface{}{
		// Missing refresh_token
	}

	// Create request
	body, _ := json.Marshal(logoutRequest)
	req, _ := http.NewRequest("POST", "/api/v1/logout", bytes.NewBuffer(body))
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

	// No expectations to verify as the service should not be called
}

func TestTokenController_Logout_ServiceError(t *testing.T) {
	// Setup
	router, mockService := setupTokenController(t)

	logoutRequest := &models.LogoutRequest{
		RefreshToken: "refresh-token",
	}

	// Expectations
	mockService.EXPECT().Logout(gomock.AssignableToTypeOf(&models.LogoutRequest{})).Return(errors.New("service error"))

	// Create request
	body, _ := json.Marshal(logoutRequest)
	req, _ := http.NewRequest("POST", "/api/v1/logout", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Perform request
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, false, response["success"])
	assert.Equal(t, "Failed to logout", response["message"])

	// Verify expectations handled by gomock controller
}
