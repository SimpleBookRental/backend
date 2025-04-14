package utils

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetAccessSecret(t *testing.T) {
	// Test case 1: With environment variable
	os.Setenv("JWT_ACCESS_SECRET", "test_access_secret")
	secret := GetAccessSecret()
	assert.Equal(t, []byte("test_access_secret"), secret)

	// Test case 2: Without environment variable
	os.Unsetenv("JWT_ACCESS_SECRET")
	secret = GetAccessSecret()
	assert.Equal(t, []byte("access_token_secret_key"), secret) // Default value
}

func TestGenerateTokenPair(t *testing.T) {
	// Test case 1: Valid user ID and email
	userID := "test-user-id"
	email := "test@example.com"
	tokenPair, err := GenerateTokenPair(userID, email)
	assert.NoError(t, err)
	assert.NotEmpty(t, tokenPair.AccessToken)
	assert.NotEmpty(t, tokenPair.RefreshToken)

	// Test case 2: Empty user ID and email
	tokenPair, err = GenerateTokenPair("", "")
	assert.NoError(t, err)
	assert.NotEmpty(t, tokenPair.AccessToken)
	assert.NotEmpty(t, tokenPair.RefreshToken)
}
