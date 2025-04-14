package services

import (
	"errors"
	"testing"
	"time"

	"github.com/SimpleBookRental/backend/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestTokenService_RefreshToken(t *testing.T) {
	// Setup
	mockTokenRepo := new(MockTokenRepository)
	mockUserRepo := new(MockUserRepository)
	service := NewTokenService(mockTokenRepo, mockUserRepo)

	refreshToken := "refresh-token"
	request := &models.RefreshTokenRequest{
		RefreshToken: refreshToken,
	}

	user := &models.User{
		ID:    "user-id",
		Email: "test@example.com",
	}

	issuedToken := &models.IssuedToken{
		ID:        "token-id",
		UserID:    user.ID,
		Token:     refreshToken,
		TokenType: string(models.RefreshToken),
		ExpiresAt: time.Now().Add(time.Hour),
		IsRevoked: false,
	}

	// Expectations
	mockTokenRepo.On("FindTokenByValue", refreshToken).Return(issuedToken, nil)
	mockUserRepo.On("FindByID", user.ID).Return(user, nil)
	mockTokenRepo.On("CreateToken", mock.AnythingOfType("*models.IssuedToken")).Return(nil).Times(2)
	mockTokenRepo.On("RevokeToken", issuedToken).Return(nil)

	// Test
	response, err := service.RefreshToken(request)
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.NotEmpty(t, response.AccessToken)
	assert.NotEmpty(t, response.RefreshToken)
	assert.NotZero(t, response.ExpiresAt)

	// Verify expectations
	mockTokenRepo.AssertExpectations(t)
	mockUserRepo.AssertExpectations(t)
}

func TestTokenService_RefreshToken_TokenNotFound(t *testing.T) {
	// Setup
	mockTokenRepo := new(MockTokenRepository)
	mockUserRepo := new(MockUserRepository)
	service := NewTokenService(mockTokenRepo, mockUserRepo)

	refreshToken := "refresh-token"
	request := &models.RefreshTokenRequest{
		RefreshToken: refreshToken,
	}

	// Expectations
	mockTokenRepo.On("FindTokenByValue", refreshToken).Return(nil, nil)

	// Test
	response, err := service.RefreshToken(request)
	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Contains(t, err.Error(), "token not found")

	// Verify expectations
	mockTokenRepo.AssertExpectations(t)
	mockUserRepo.AssertExpectations(t)
}

func TestTokenService_RefreshToken_TokenRevoked(t *testing.T) {
	// Setup
	mockTokenRepo := new(MockTokenRepository)
	mockUserRepo := new(MockUserRepository)
	service := NewTokenService(mockTokenRepo, mockUserRepo)

	refreshToken := "refresh-token"
	request := &models.RefreshTokenRequest{
		RefreshToken: refreshToken,
	}

	issuedToken := &models.IssuedToken{
		ID:        "token-id",
		UserID:    "user-id",
		Token:     refreshToken,
		TokenType: string(models.RefreshToken),
		ExpiresAt: time.Now().Add(time.Hour),
		IsRevoked: true,
	}

	// Expectations
	mockTokenRepo.On("FindTokenByValue", refreshToken).Return(issuedToken, nil)

	// Test
	response, err := service.RefreshToken(request)
	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Contains(t, err.Error(), "token has been revoked")

	// Verify expectations
	mockTokenRepo.AssertExpectations(t)
	mockUserRepo.AssertExpectations(t)
}

func TestTokenService_RefreshToken_TokenExpired(t *testing.T) {
	// Setup
	mockTokenRepo := new(MockTokenRepository)
	mockUserRepo := new(MockUserRepository)
	service := NewTokenService(mockTokenRepo, mockUserRepo)

	refreshToken := "refresh-token"
	request := &models.RefreshTokenRequest{
		RefreshToken: refreshToken,
	}

	issuedToken := &models.IssuedToken{
		ID:        "token-id",
		UserID:    "user-id",
		Token:     refreshToken,
		TokenType: string(models.RefreshToken),
		ExpiresAt: time.Now().Add(-time.Hour), // Expired
		IsRevoked: false,
	}

	// Expectations
	mockTokenRepo.On("FindTokenByValue", refreshToken).Return(issuedToken, nil)

	// Test
	response, err := service.RefreshToken(request)
	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Contains(t, err.Error(), "token has expired")

	// Verify expectations
	mockTokenRepo.AssertExpectations(t)
	mockUserRepo.AssertExpectations(t)
}

func TestTokenService_RefreshToken_InvalidTokenType(t *testing.T) {
	// Setup
	mockTokenRepo := new(MockTokenRepository)
	mockUserRepo := new(MockUserRepository)
	service := NewTokenService(mockTokenRepo, mockUserRepo)

	refreshToken := "refresh-token"
	request := &models.RefreshTokenRequest{
		RefreshToken: refreshToken,
	}

	issuedToken := &models.IssuedToken{
		ID:        "token-id",
		UserID:    "user-id",
		Token:     refreshToken,
		TokenType: string(models.AccessToken), // Wrong token type
		ExpiresAt: time.Now().Add(time.Hour),
		IsRevoked: false,
	}

	// Expectations
	mockTokenRepo.On("FindTokenByValue", refreshToken).Return(issuedToken, nil)

	// Test
	response, err := service.RefreshToken(request)
	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Contains(t, err.Error(), "invalid token type")

	// Verify expectations
	mockTokenRepo.AssertExpectations(t)
	mockUserRepo.AssertExpectations(t)
}

func TestTokenService_RefreshToken_UserNotFound(t *testing.T) {
	// Setup
	mockTokenRepo := new(MockTokenRepository)
	mockUserRepo := new(MockUserRepository)
	service := NewTokenService(mockTokenRepo, mockUserRepo)

	refreshToken := "refresh-token"
	request := &models.RefreshTokenRequest{
		RefreshToken: refreshToken,
	}

	issuedToken := &models.IssuedToken{
		ID:        "token-id",
		UserID:    "user-id",
		Token:     refreshToken,
		TokenType: string(models.RefreshToken),
		ExpiresAt: time.Now().Add(time.Hour),
		IsRevoked: false,
	}

	// Expectations
	mockTokenRepo.On("FindTokenByValue", refreshToken).Return(issuedToken, nil)
	mockUserRepo.On("FindByID", issuedToken.UserID).Return(nil, nil)

	// Test
	response, err := service.RefreshToken(request)
	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Contains(t, err.Error(), "user not found")

	// Verify expectations
	mockTokenRepo.AssertExpectations(t)
	mockUserRepo.AssertExpectations(t)
}

func TestTokenService_Logout(t *testing.T) {
	// Setup
	mockTokenRepo := new(MockTokenRepository)
	mockUserRepo := new(MockUserRepository)
	service := NewTokenService(mockTokenRepo, mockUserRepo)

	refreshToken := "refresh-token"
	request := &models.LogoutRequest{
		RefreshToken: refreshToken,
	}

	issuedToken := &models.IssuedToken{
		ID:        "token-id",
		UserID:    "user-id",
		Token:     refreshToken,
		TokenType: string(models.RefreshToken),
		ExpiresAt: time.Now().Add(time.Hour),
		IsRevoked: false,
	}

	// Expectations
	mockTokenRepo.On("FindTokenByValue", refreshToken).Return(issuedToken, nil)
	mockTokenRepo.On("RevokeAllUserTokens", issuedToken.UserID).Return(nil)
	mockTokenRepo.On("CleanupExpiredTokens").Return(nil)

	// Test
	err := service.Logout(request)
	assert.NoError(t, err)

	// Verify expectations
	mockTokenRepo.AssertExpectations(t)
	mockUserRepo.AssertExpectations(t)
}

func TestTokenService_Logout_TokenNotFound(t *testing.T) {
	// Setup
	mockTokenRepo := new(MockTokenRepository)
	mockUserRepo := new(MockUserRepository)
	service := NewTokenService(mockTokenRepo, mockUserRepo)

	refreshToken := "refresh-token"
	request := &models.LogoutRequest{
		RefreshToken: refreshToken,
	}

	// Expectations
	mockTokenRepo.On("FindTokenByValue", refreshToken).Return(nil, nil)

	// Test
	err := service.Logout(request)
	assert.NoError(t, err) // Should not error if token not found
	
	// Verify expectations
	mockTokenRepo.AssertExpectations(t)
	mockUserRepo.AssertExpectations(t)
}

func TestTokenService_Logout_InvalidTokenType(t *testing.T) {
	// Setup
	mockTokenRepo := new(MockTokenRepository)
	mockUserRepo := new(MockUserRepository)
	service := NewTokenService(mockTokenRepo, mockUserRepo)

	refreshToken := "refresh-token"
	request := &models.LogoutRequest{
		RefreshToken: refreshToken,
	}

	issuedToken := &models.IssuedToken{
		ID:        "token-id",
		UserID:    "user-id",
		Token:     refreshToken,
		TokenType: string(models.AccessToken), // Wrong token type
		ExpiresAt: time.Now().Add(time.Hour),
		IsRevoked: false,
	}

	// Expectations
	mockTokenRepo.On("FindTokenByValue", refreshToken).Return(issuedToken, nil)

	// Test
	err := service.Logout(request)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid token type")

	// Verify expectations
	mockTokenRepo.AssertExpectations(t)
	mockUserRepo.AssertExpectations(t)
}

func TestTokenService_Logout_RevokeError(t *testing.T) {
	// Setup
	mockTokenRepo := new(MockTokenRepository)
	mockUserRepo := new(MockUserRepository)
	service := NewTokenService(mockTokenRepo, mockUserRepo)

	refreshToken := "refresh-token"
	request := &models.LogoutRequest{
		RefreshToken: refreshToken,
	}

	issuedToken := &models.IssuedToken{
		ID:        "token-id",
		UserID:    "user-id",
		Token:     refreshToken,
		TokenType: string(models.RefreshToken),
		ExpiresAt: time.Now().Add(time.Hour),
		IsRevoked: false,
	}

	// Expectations
	mockTokenRepo.On("FindTokenByValue", refreshToken).Return(issuedToken, nil)
	mockTokenRepo.On("RevokeAllUserTokens", issuedToken.UserID).Return(errors.New("revoke error"))

	// Test
	err := service.Logout(request)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error revoking user tokens")

	// Verify expectations
	mockTokenRepo.AssertExpectations(t)
	mockUserRepo.AssertExpectations(t)
}
