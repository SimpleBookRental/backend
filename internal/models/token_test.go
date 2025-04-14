package models

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestIssuedTokenModel(t *testing.T) {
	// Test case 1: Create a new issued token
	now := time.Now()
	token := &IssuedToken{
		ID:        "test-id",
		UserID:    "user-id",
		Token:     "token-value",
		TokenType: string(AccessToken),
		ExpiresAt: now.Add(time.Hour),
		IsRevoked: false,
		CreatedAt: now,
		UpdatedAt: now,
	}

	// Assert
	assert.Equal(t, "test-id", token.ID)
	assert.Equal(t, "user-id", token.UserID)
	assert.Equal(t, "token-value", token.Token)
	assert.Equal(t, string(AccessToken), token.TokenType)
	assert.False(t, token.IsRevoked)
	assert.Nil(t, token.RevokedAt)
	assert.Equal(t, now, token.CreatedAt)
	assert.Equal(t, now, token.UpdatedAt)
}

func TestIssuedTokenBeforeCreate(t *testing.T) {
	// Test case 1: BeforeCreate with empty ID
	token := &IssuedToken{
		UserID:    "user-id",
		Token:     "token-value",
		TokenType: string(AccessToken),
		ExpiresAt: time.Now().Add(time.Hour),
	}

	err := token.BeforeCreate(&gorm.DB{})
	assert.NoError(t, err)
	assert.NotEmpty(t, token.ID)

	// Test case 2: BeforeCreate with existing ID
	token = &IssuedToken{
		ID:        "test-id",
		UserID:    "user-id",
		Token:     "token-value",
		TokenType: string(AccessToken),
		ExpiresAt: time.Now().Add(time.Hour),
	}

	err = token.BeforeCreate(&gorm.DB{})
	assert.NoError(t, err)
	assert.Equal(t, "test-id", token.ID)
}

func TestIssuedTokenTableName(t *testing.T) {
	// Test case 1: Check table name
	token := IssuedToken{}
	assert.Equal(t, "br_issued_token", token.TableName())
}

func TestRefreshTokenRequest(t *testing.T) {
	// Test case 1: Valid refresh token request
	request := &RefreshTokenRequest{
		RefreshToken: "refresh-token",
	}

	// Assert
	assert.Equal(t, "refresh-token", request.RefreshToken)
}

func TestRefreshTokenResponse(t *testing.T) {
	// Test case 1: Valid refresh token response
	response := &RefreshTokenResponse{
		AccessToken:  "access-token",
		RefreshToken: "refresh-token",
		ExpiresAt:    time.Now().Unix(),
	}

	// Assert
	assert.Equal(t, "access-token", response.AccessToken)
	assert.Equal(t, "refresh-token", response.RefreshToken)
	assert.NotZero(t, response.ExpiresAt)
}

func TestLogoutRequest(t *testing.T) {
	// Test case 1: Valid logout request
	request := &LogoutRequest{
		RefreshToken: "refresh-token",
	}

	// Assert
	assert.Equal(t, "refresh-token", request.RefreshToken)
}

func TestTokenType(t *testing.T) {
	// Test case 1: Check token types
	assert.Equal(t, TokenType("access"), AccessToken)
	assert.Equal(t, TokenType("refresh"), RefreshToken)
}
