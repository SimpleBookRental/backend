// Unit tests for TokenController using gomock and testify.
package controllers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/SimpleBookRental/backend/internal/mocks"
	"github.com/SimpleBookRental/backend/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func setupGinToken() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.New()
}

func TestTokenController_RefreshToken_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := mocks.NewMockTokenServiceInterface(ctrl)
	controller := NewTokenController(mockService)
	router := setupGinToken()
	router.POST("/token/refresh", controller.RefreshToken)

	reqBody := models.RefreshTokenRequest{
		RefreshToken: "valid-refresh-token",
	}
	expectedResp := &models.RefreshTokenResponse{
		AccessToken:  "access-token",
		RefreshToken: "refresh-token",
		ExpiresAt:    1234567890,
	}
	mockService.EXPECT().RefreshToken(&reqBody).Return(expectedResp, nil)

	body, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/token/refresh", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestTokenController_RefreshToken_BadRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := mocks.NewMockTokenServiceInterface(ctrl)
	controller := NewTokenController(mockService)
	router := setupGinToken()
	router.POST("/token/refresh", controller.RefreshToken)

	// Invalid JSON
	req, _ := http.NewRequest("POST", "/token/refresh", bytes.NewBuffer([]byte("{invalid")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestTokenController_RefreshToken_ServiceError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := mocks.NewMockTokenServiceInterface(ctrl)
	controller := NewTokenController(mockService)
	router := setupGinToken()
	router.POST("/token/refresh", controller.RefreshToken)

	reqBody := models.RefreshTokenRequest{
		RefreshToken: "invalid-refresh-token",
	}
	mockService.EXPECT().RefreshToken(&reqBody).Return(nil, errors.New("invalid token"))

	body, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/token/refresh", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestTokenController_Logout_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := mocks.NewMockTokenServiceInterface(ctrl)
	controller := NewTokenController(mockService)
	router := setupGinToken()
	router.POST("/token/logout", controller.Logout)

	reqBody := models.LogoutRequest{
		RefreshToken: "valid-refresh-token",
	}
	mockService.EXPECT().Logout(&reqBody).Return(nil)

	body, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/token/logout", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestTokenController_Logout_BadRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := mocks.NewMockTokenServiceInterface(ctrl)
	controller := NewTokenController(mockService)
	router := setupGinToken()
	router.POST("/token/logout", controller.Logout)

	// Invalid JSON
	req, _ := http.NewRequest("POST", "/token/logout", bytes.NewBuffer([]byte("{invalid")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestTokenController_Logout_ServiceError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := mocks.NewMockTokenServiceInterface(ctrl)
	controller := NewTokenController(mockService)
	router := setupGinToken()
	router.POST("/token/logout", controller.Logout)

	reqBody := models.LogoutRequest{
		RefreshToken: "invalid-refresh-token",
	}
	mockService.EXPECT().Logout(&reqBody).Return(errors.New("logout error"))

	body, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/token/logout", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
