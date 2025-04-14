package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateUUID(t *testing.T) {
	// Test case 1: Generate UUID
	uuid := GenerateUUID()
	assert.NotEmpty(t, uuid)
	assert.Len(t, uuid, 36) // UUID v4 is 36 characters long

	// Test case 2: Generate multiple UUIDs and check they are different
	uuid1 := GenerateUUID()
	uuid2 := GenerateUUID()
	assert.NotEqual(t, uuid1, uuid2)
}

func TestIsValidUUID(t *testing.T) {
	// Test case 1: Valid UUID
	uuid := GenerateUUID()
	isValid := IsValidUUID(uuid)
	assert.True(t, isValid)

	// Test case 2: Invalid UUID
	isValid = IsValidUUID("invalid-uuid")
	assert.False(t, isValid)

	// Test case 3: Empty UUID
	isValid = IsValidUUID("")
	assert.False(t, isValid)

	// Test case 4: UUID with wrong format
	isValid = IsValidUUID("123e4567-e89b-12d3-a456-426655440000")
	assert.True(t, isValid)

	// Test case 5: UUID with wrong characters
	isValid = IsValidUUID("123e4567-e89b-12d3-a456-42665544000g")
	assert.False(t, isValid)
}
