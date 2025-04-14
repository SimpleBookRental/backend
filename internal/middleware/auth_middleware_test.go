package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/SimpleBookRental/backend/internal/models"
	"github.com/SimpleBookRental/backend/internal/repositories"
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

// Ensure MockTokenRepository implements TokenRepositoryInterface
var _ repositories.TokenRepositoryInterface = (*MockTokenRepository)(nil)

func setupAuthMiddleware() (*gin.Engine, *MockTokenRepository) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	mockRepo := new(MockTokenRepository)

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
	router, _ := setupAuthMiddleware()

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
	router, _ := setupAuthMiddleware()

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
	router, mockRepo := setupAuthMiddleware()

	token := "invalid-token"

	// Expectations
	mockRepo.On("FindTokenByValue", token).Return(nil, nil)

	// Create request
	req, _ := http.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	// Perform request
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid token")

	// Verify expectations
	mockRepo.AssertExpectations(t)
}

func TestAuthMiddleware_TokenRevoked(t *testing.T) {
	// Setup
	router, mockRepo := setupAuthMiddleware()

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
	mockRepo.On("FindTokenByValue", token).Return(issuedToken, nil)

	// Create request
	req, _ := http.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	// Perform request
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "Token has been revoked")

	// Verify expectations
	mockRepo.AssertExpectations(t)
}

func TestAuthMiddleware_TokenExpired(t *testing.T) {
	// Setup
	router, mockRepo := setupAuthMiddleware()

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
	mockRepo.On("FindTokenByValue", token).Return(issuedToken, nil)

	// Create request
	req, _ := http.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	// Perform request
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "Token has expired")

	// Verify expectations
	mockRepo.AssertExpectations(t)
}

func TestAuthMiddleware_InvalidTokenType(t *testing.T) {
	// Setup
	router, mockRepo := setupAuthMiddleware()

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
	mockRepo.On("FindTokenByValue", token).Return(issuedToken, nil)

	// Create request
	req, _ := http.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	// Perform request
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid token type")

	// Verify expectations
	mockRepo.AssertExpectations(t)
}
