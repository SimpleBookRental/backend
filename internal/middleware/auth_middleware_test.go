// Unit tests for AuthMiddleware using gomock and testify.
package middleware

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/SimpleBookRental/backend/internal/models"
	"github.com/SimpleBookRental/backend/internal/mocks"
	"github.com/SimpleBookRental/backend/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func setupGinWithAuth(tokenRepo *mocks.MockTokenRepositoryInterface, userRepo *mocks.MockUserRepositoryInterface) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(AuthMiddleware(tokenRepo, userRepo))
	r.GET("/protected", func(c *gin.Context) {
		c.JSON(200, gin.H{"success": true})
	})
	return r
}

func TestAuthMiddleware_MissingHeader(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	tokenRepo := mocks.NewMockTokenRepositoryInterface(ctrl)
	userRepo := mocks.NewMockUserRepositoryInterface(ctrl)
	router := setupGinWithAuth(tokenRepo, userRepo)

	req, _ := http.NewRequest("GET", "/protected", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, 401, w.Code)
	assert.Contains(t, w.Body.String(), "Authorization header is required")
}

func TestAuthMiddleware_InvalidHeaderFormat(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	tokenRepo := mocks.NewMockTokenRepositoryInterface(ctrl)
	userRepo := mocks.NewMockUserRepositoryInterface(ctrl)
	router := setupGinWithAuth(tokenRepo, userRepo)

	req, _ := http.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Token abc")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, 401, w.Code)
	assert.Contains(t, w.Body.String(), "Authorization header format must be Bearer")
}

func TestAuthMiddleware_InvalidToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	tokenRepo := mocks.NewMockTokenRepositoryInterface(ctrl)
	userRepo := mocks.NewMockUserRepositoryInterface(ctrl)
	router := setupGinWithAuth(tokenRepo, userRepo)

	req, _ := http.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer invalidtoken")
	// Patch utils.ValidateToken to always return error for this test
	origValidateToken := utils.ValidateToken
	utils.ValidateToken = func(token string, secret []byte) (*utils.Claims, error) {
		return nil, errors.New("invalid")
	}
	defer func() { utils.ValidateToken = origValidateToken }()

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, 401, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid or expired token")
}

func TestAuthMiddleware_TokenNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	tokenRepo := mocks.NewMockTokenRepositoryInterface(ctrl)
	userRepo := mocks.NewMockUserRepositoryInterface(ctrl)
	router := setupGinWithAuth(tokenRepo, userRepo)

	claims := &utils.Claims{UserID: "user-1", Email: "test@example.com"}
	origValidateToken := utils.ValidateToken
	utils.ValidateToken = func(token string, secret []byte) (*utils.Claims, error) {
		return claims, nil
	}
	defer func() { utils.ValidateToken = origValidateToken }()

	tokenRepo.EXPECT().FindTokenByValue("validtoken").Return(nil, nil)

	req, _ := http.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer validtoken")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, 401, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid token")
}

func TestAuthMiddleware_TokenRevoked(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	tokenRepo := mocks.NewMockTokenRepositoryInterface(ctrl)
	userRepo := mocks.NewMockUserRepositoryInterface(ctrl)
	router := setupGinWithAuth(tokenRepo, userRepo)

	claims := &utils.Claims{UserID: "user-1", Email: "test@example.com"}
	origValidateToken := utils.ValidateToken
	utils.ValidateToken = func(token string, secret []byte) (*utils.Claims, error) {
		return claims, nil
	}
	defer func() { utils.ValidateToken = origValidateToken }()

	issuedToken := &models.IssuedToken{
		Token:     "validtoken",
		IsRevoked: true,
		ExpiresAt: time.Now().Add(time.Hour),
		TokenType: string(models.AccessToken),
	}
	tokenRepo.EXPECT().FindTokenByValue("validtoken").Return(issuedToken, nil)

	req, _ := http.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer validtoken")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, 401, w.Code)
	assert.Contains(t, w.Body.String(), "Token has been revoked")
}

func TestAuthMiddleware_TokenExpired(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	tokenRepo := mocks.NewMockTokenRepositoryInterface(ctrl)
	userRepo := mocks.NewMockUserRepositoryInterface(ctrl)
	router := setupGinWithAuth(tokenRepo, userRepo)

	claims := &utils.Claims{UserID: "user-1", Email: "test@example.com"}
	origValidateToken := utils.ValidateToken
	utils.ValidateToken = func(token string, secret []byte) (*utils.Claims, error) {
		return claims, nil
	}
	defer func() { utils.ValidateToken = origValidateToken }()

	issuedToken := &models.IssuedToken{
		Token:     "validtoken",
		IsRevoked: false,
		ExpiresAt: time.Now().Add(-time.Hour),
		TokenType: string(models.AccessToken),
	}
	tokenRepo.EXPECT().FindTokenByValue("validtoken").Return(issuedToken, nil)

	req, _ := http.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer validtoken")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, 401, w.Code)
	assert.Contains(t, w.Body.String(), "Token has expired")
}

func TestAuthMiddleware_InvalidTokenType(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	tokenRepo := mocks.NewMockTokenRepositoryInterface(ctrl)
	userRepo := mocks.NewMockUserRepositoryInterface(ctrl)
	router := setupGinWithAuth(tokenRepo, userRepo)

	claims := &utils.Claims{UserID: "user-1", Email: "test@example.com"}
	origValidateToken := utils.ValidateToken
	utils.ValidateToken = func(token string, secret []byte) (*utils.Claims, error) {
		return claims, nil
	}
	defer func() { utils.ValidateToken = origValidateToken }()

	issuedToken := &models.IssuedToken{
		Token:     "validtoken",
		IsRevoked: false,
		ExpiresAt: time.Now().Add(time.Hour),
		TokenType: string(models.RefreshToken),
	}
	tokenRepo.EXPECT().FindTokenByValue("validtoken").Return(issuedToken, nil)

	req, _ := http.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer validtoken")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, 401, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid token type")
}

func TestAuthMiddleware_UserNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	tokenRepo := mocks.NewMockTokenRepositoryInterface(ctrl)
	userRepo := mocks.NewMockUserRepositoryInterface(ctrl)
	router := setupGinWithAuth(tokenRepo, userRepo)

	claims := &utils.Claims{UserID: "user-1", Email: "test@example.com"}
	origValidateToken := utils.ValidateToken
	utils.ValidateToken = func(token string, secret []byte) (*utils.Claims, error) {
		return claims, nil
	}
	defer func() { utils.ValidateToken = origValidateToken }()

	issuedToken := &models.IssuedToken{
		Token:     "validtoken",
		IsRevoked: false,
		ExpiresAt: time.Now().Add(time.Hour),
		TokenType: string(models.AccessToken),
	}
	tokenRepo.EXPECT().FindTokenByValue("validtoken").Return(issuedToken, nil)
	userRepo.EXPECT().FindByID("user-1").Return(nil, nil)

	req, _ := http.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer validtoken")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, 401, w.Code)
	assert.Contains(t, w.Body.String(), "User not found")
}

func TestAuthMiddleware_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	tokenRepo := mocks.NewMockTokenRepositoryInterface(ctrl)
	userRepo := mocks.NewMockUserRepositoryInterface(ctrl)
	router := setupGinWithAuth(tokenRepo, userRepo)

	claims := &utils.Claims{UserID: "user-1", Email: "test@example.com"}
	origValidateToken := utils.ValidateToken
	utils.ValidateToken = func(token string, secret []byte) (*utils.Claims, error) {
		return claims, nil
	}
	defer func() { utils.ValidateToken = origValidateToken }()

	issuedToken := &models.IssuedToken{
		Token:     "validtoken",
		IsRevoked: false,
		ExpiresAt: time.Now().Add(time.Hour),
		TokenType: string(models.AccessToken),
	}
	tokenRepo.EXPECT().FindTokenByValue("validtoken").Return(issuedToken, nil)
	userRepo.EXPECT().FindByID("user-1").Return(&models.User{ID: "user-1", Role: "ADMIN"}, nil)

	req, _ := http.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer validtoken")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), "success")
}
