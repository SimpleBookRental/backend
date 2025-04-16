package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/SimpleBookRental/backend/internal/mocks"
)

func setupAuthMiddleware(t *testing.T) (*gin.Engine, *mocks.MockTokenRepositoryInterface, *mocks.MockUserRepositoryInterface) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	ctrl := gomock.NewController(t)
	mockTokenRepo := mocks.NewMockTokenRepositoryInterface(ctrl)
	mockUserRepo := mocks.NewMockUserRepositoryInterface(ctrl)

	// Setup routes
	router.GET("/protected", AuthMiddleware(mockTokenRepo, mockUserRepo), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Protected route accessed successfully",
			"data":    nil,
		})
	})

	return router, mockTokenRepo, mockUserRepo
}

func TestAuthMiddleware_MissingAuthHeader(t *testing.T) {
	// Setup
	router, _, _ := setupAuthMiddleware(t)

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
	router, _, _ := setupAuthMiddleware(t)

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

func TestAuthMiddleware_InvalidToken(t *testing.T) {
	// Setup
	router, _, _ := setupAuthMiddleware(t)

	token := "invalid-token"

	// Create request
	req, _ := http.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	// Perform request
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid or expired token")
}
