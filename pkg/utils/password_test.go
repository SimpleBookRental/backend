package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHashPassword(t *testing.T) {
	// Test case 1: Valid password
	password := "password123"
	hashedPassword, err := HashPassword(password)
	assert.NoError(t, err)
	assert.NotEmpty(t, hashedPassword)
	assert.NotEqual(t, password, hashedPassword)

	// Test case 2: Check if the hash is valid
	isValid := CheckPasswordHash(password, hashedPassword)
	assert.True(t, isValid)

	// Test case 3: Check with wrong password
	isValid = CheckPasswordHash("wrongpassword", hashedPassword)
	assert.False(t, isValid)

	// Test case 4: Empty password
	_, err = HashPassword("")
	assert.NoError(t, err) // bcrypt can hash empty strings
}

func TestCheckPasswordHash(t *testing.T) {
	// Test case 1: Valid password and hash
	password := "password123"
	hashedPassword, _ := HashPassword(password)
	isValid := CheckPasswordHash(password, hashedPassword)
	assert.True(t, isValid)

	// Test case 2: Invalid password
	isValid = CheckPasswordHash("wrongpassword", hashedPassword)
	assert.False(t, isValid)

	// Test case 3: Invalid hash
	isValid = CheckPasswordHash(password, "invalid_hash")
	assert.False(t, isValid)

	// Test case 4: Empty password and hash
	isValid = CheckPasswordHash("", "")
	assert.False(t, isValid)
}
