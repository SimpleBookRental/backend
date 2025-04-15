package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/SimpleBookRental/backend/internal/mocks"
	"github.com/SimpleBookRental/backend/internal/models"
)

func setupAuthMiddleware(t *testing.T) (*gin.Engine, *mocks.MockTokenRepositoryInterface) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	ctrl := gomock.NewController(t)
	mockRepo := mocks.NewMockTokenRepositoryInterface(ctrl)

	// Setup routes
	router.GET("/protected", AuthMiddleware(mockRepo), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Protected route accessed successfully",
			"data":    nil,
		})
	})

	return router, mockRepo
}

func TestAuthMiddleware_MissingAuthHeader(t *testing.T) {
	// Setup
	router, _ := setupAuthMiddleware(t)

	// Create request without Authorization header
	req, _ := http.NewRequest("GET", "/protected", nil)
	w := httptest.NewRecorder()

	// Perform request
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "Authorization header is required")
}

func TestAuthMiddleware_InvalidAuthHeaderFormat(t *testing.T) {
	// Setup
	router, _ := setupAuthMiddleware(t)

	// Create request with invalid Authorization header format
	req, _ := http.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "InvalidFormat")
	w := httptest.NewRecorder()

	// Perform request
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "Authorization header format must be Bearer {token}")
}

func TestAuthMiddleware_TokenNotFound(t *testing.T) {
	// Setup
	router, mockRepo := setupAuthMiddleware(t)

	token := "invalid-token"

	// Expectations
	mockRepo.EXPECT().FindTokenByValue(token).Return(nil, nil)

	// Create request
	req, _ := http.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	// Perform request
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid token")

	// Verify expectations handled by gomock controller
}

func TestAuthMiddleware_TokenRevoked(t *testing.T) {
	// Setup
	router, mockRepo := setupAuthMiddleware(t)

	token := "revoked-token"
	userID := "user-id"

	issuedToken := &models.IssuedToken{
		ID:        "token-id",
		UserID:    userID,
		Token:     token,
		TokenType: string(models.AccessToken),
		ExpiresAt: time.Now().Add(time.Hour),
		IsRevoked: true, // Token is revoked
	}

	// Expectations
	mockRepo.EXPECT().FindTokenByValue(token).Return(issuedToken, nil)

	// Create request
	req, _ := http.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	// Perform request
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "Token has been revoked")

	// Verify expectations handled by gomock controller
}

func TestAuthMiddleware_TokenExpired(t *testing.T) {
	// Setup
	router, mockRepo := setupAuthMiddleware(t)

	token := "expired-token"
	userID := "user-id"

	issuedToken := &models.IssuedToken{
		ID:        "token-id",
		UserID:    userID,
		Token:     token,
		TokenType: string(models.AccessToken),
		ExpiresAt: time.Now().Add(-time.Hour), // Token is expired
		IsRevoked: false,
	}

	// Expectations
	mockRepo.EXPECT().FindTokenByValue(token).Return(issuedToken, nil)

	// Create request
	req, _ := http.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	// Perform request
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "Token has expired")

	// Verify expectations handled by gomock controller
}

func TestAuthMiddleware_InvalidTokenType(t *testing.T) {
	// Setup
	router, mockRepo := setupAuthMiddleware(t)

	token := "refresh-token"
	userID := "user-id"

	issuedToken := &models.IssuedToken{
		ID:        "token-id",
		UserID:    userID,
		Token:     token,
		TokenType: string(models.RefreshToken), // Wrong token type
		ExpiresAt: time.Now().Add(time.Hour),
		IsRevoked: false,
	}

	// Expectations
	mockRepo.EXPECT().FindTokenByValue(token).Return(issuedToken, nil)

	// Create request
	req, _ := http.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	// Perform request
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid token type")

	// Verify expectations handled by gomock controller
}
