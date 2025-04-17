// Package utils provides utility functions for JWT handling.
package utils

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestGetAccessSecret checks if GetAccessSecret returns the correct secret.
func TestGetAccessSecret(t *testing.T) {
	// Test with environment variable set
	os.Setenv("JWT_ACCESS_SECRET", "test_secret")
	secret := GetAccessSecret()
	assert.Equal(t, []byte("test_secret"), secret)

	// Test with environment variable unset
	os.Unsetenv("JWT_ACCESS_SECRET")
	secret = GetAccessSecret()
	assert.Equal(t, []byte("access_token_secret_key"), secret)
}

// TestGenerateTokenPairAndValidateToken tests token generation and validation.
func TestGenerateTokenPairAndValidateToken(t *testing.T) {
	os.Setenv("JWT_ACCESS_SECRET", "access_secret")
	os.Setenv("JWT_REFRESH_SECRET", "refresh_secret")
	os.Setenv("JWT_ACCESS_EXPIRATION", "1h")
	os.Setenv("JWT_REFRESH_EXPIRATION", "2h")

	userID := "123"
	email := "test@example.com"

	// Generate token pair
	tokens, err := GenerateTokenPair(userID, email)
	assert.NoError(t, err)
	assert.NotEmpty(t, tokens.AccessToken)
	assert.NotEmpty(t, tokens.RefreshToken)

	// Validate access token
	claims, err := ValidateToken(tokens.AccessToken, GetAccessSecret())
	assert.NoError(t, err)
	assert.Equal(t, userID, claims.UserID)
	assert.Equal(t, email, claims.Email)

	// Validate refresh token
	refreshClaims, err := ValidateToken(tokens.RefreshToken, []byte("refresh_secret"))
	assert.NoError(t, err)
	assert.Equal(t, userID, refreshClaims.UserID)
	assert.Equal(t, email, refreshClaims.Email)
}

// TestValidateToken_InvalidToken tests validation with an invalid token.
func TestValidateToken_InvalidToken(t *testing.T) {
	invalidToken := "invalid.token.string"
	_, err := ValidateToken(invalidToken, GetAccessSecret())
	assert.Error(t, err)
}

// TestValidateToken_WrongSecret tests validation with a wrong secret.
func TestValidateToken_WrongSecret(t *testing.T) {
	os.Setenv("JWT_ACCESS_SECRET", "access_secret")
	userID := "123"
	email := "test@example.com"
	tokens, err := GenerateTokenPair(userID, email)
	assert.NoError(t, err)

	// Validate with wrong secret
	_, err = ValidateToken(tokens.AccessToken, []byte("wrong_secret"))
	assert.Error(t, err)
}

// TestRefreshTokens tests refreshing tokens using a valid refresh token.
func TestRefreshTokens(t *testing.T) {
	os.Setenv("JWT_ACCESS_SECRET", "access_secret")
	os.Setenv("JWT_REFRESH_SECRET", "refresh_secret")
	os.Setenv("JWT_ACCESS_EXPIRATION", "1h")
	os.Setenv("JWT_REFRESH_EXPIRATION", "2h")

	userID := "456"
	email := "refresh@example.com"
	tokens, err := GenerateTokenPair(userID, email)
	assert.NoError(t, err)

	newTokens, err := RefreshTokens(tokens.RefreshToken)
	assert.NoError(t, err)
	assert.NotEmpty(t, newTokens.AccessToken)
	assert.NotEmpty(t, newTokens.RefreshToken)
}

// TestRefreshTokens_InvalidRefreshToken tests refreshing with an invalid refresh token.
func TestRefreshTokens_InvalidRefreshToken(t *testing.T) {
	invalidToken := "invalid.token.string"
	_, err := RefreshTokens(invalidToken)
	assert.Error(t, err)
}

// TestGenerateTokenPair_ErrorHandling tests error handling in GenerateTokenPair.
func TestGenerateTokenPair_ErrorHandling(t *testing.T) {
	// Patch generateToken to return error (simulate signing error)
	// Not possible directly, but we can test with empty secret
	os.Setenv("JWT_ACCESS_SECRET", "")
	os.Setenv("JWT_REFRESH_SECRET", "")

	// Unlikely to fail, but for coverage, try with empty userID/email
	tokens, err := GenerateTokenPair("", "")
	assert.NoError(t, err)
	assert.NotEmpty(t, tokens.AccessToken)
	assert.NotEmpty(t, tokens.RefreshToken)
}

// TestGenerateToken_Expiration tests token expiration logic.
func TestGenerateToken_Expiration(t *testing.T) {
	secret := []byte("test_secret")
	userID := "789"
	email := "expire@example.com"
	tokenStr, err := generateToken(userID, email, secret, time.Millisecond)
	assert.NoError(t, err)
	time.Sleep(2 * time.Millisecond)
	_, err = ValidateToken(tokenStr, secret)
	assert.Error(t, err)
}
