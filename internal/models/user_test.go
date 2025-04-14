package models

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestUserModel(t *testing.T) {
	// Test case 1: Create a new user
	user := &User{
		ID:        "test-id",
		Name:      "Test User",
		Email:     "test@example.com",
		Password:  "password123",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Assert
	assert.Equal(t, "test-id", user.ID)
	assert.Equal(t, "Test User", user.Name)
	assert.Equal(t, "test@example.com", user.Email)
	assert.Equal(t, "password123", user.Password)
	assert.NotZero(t, user.CreatedAt)
	assert.NotZero(t, user.UpdatedAt)
}

func TestUserBeforeCreate(t *testing.T) {
	// Test case 1: BeforeCreate with empty ID
	user := &User{
		Name:     "Test User",
		Email:    "test@example.com",
		Password: "password123",
	}

	err := user.BeforeCreate(&gorm.DB{})
	assert.NoError(t, err)
	assert.NotEmpty(t, user.ID)

	// Test case 2: BeforeCreate with existing ID
	user = &User{
		ID:       "test-id",
		Name:     "Test User",
		Email:    "test@example.com",
		Password: "password123",
	}

	err = user.BeforeCreate(&gorm.DB{})
	assert.NoError(t, err)
	assert.Equal(t, "test-id", user.ID)
}

func TestUserTableName(t *testing.T) {
	// Test case 1: Check table name
	user := User{}
	assert.Equal(t, "br_user", user.TableName())
}

func TestUserCreate(t *testing.T) {
	// Test case 1: Valid user create
	userCreate := &UserCreate{
		Name:     "Test User",
		Email:    "test@example.com",
		Password: "password123",
	}

	// Assert
	assert.Equal(t, "Test User", userCreate.Name)
	assert.Equal(t, "test@example.com", userCreate.Email)
	assert.Equal(t, "password123", userCreate.Password)
}

func TestUserUpdate(t *testing.T) {
	// Test case 1: Valid user update
	userUpdate := &UserUpdate{
		Name:     "Updated User",
		Email:    "updated@example.com",
		Password: "newpassword123",
	}

	// Assert
	assert.Equal(t, "Updated User", userUpdate.Name)
	assert.Equal(t, "updated@example.com", userUpdate.Email)
	assert.Equal(t, "newpassword123", userUpdate.Password)
}

func TestUserLogin(t *testing.T) {
	// Test case 1: Valid user login
	userLogin := &UserLogin{
		Email:    "test@example.com",
		Password: "password123",
	}

	// Assert
	assert.Equal(t, "test@example.com", userLogin.Email)
	assert.Equal(t, "password123", userLogin.Password)
}

func TestLoginResponse(t *testing.T) {
	// Test case 1: Valid login response
	user := &User{
		ID:    "test-id",
		Name:  "Test User",
		Email: "test@example.com",
	}

	loginResponse := &LoginResponse{
		User:         user,
		AccessToken:  "access-token",
		RefreshToken: "refresh-token",
		ExpiresAt:    time.Now().Unix(),
	}

	// Assert
	assert.Equal(t, user, loginResponse.User)
	assert.Equal(t, "access-token", loginResponse.AccessToken)
	assert.Equal(t, "refresh-token", loginResponse.RefreshToken)
	assert.NotZero(t, loginResponse.ExpiresAt)
}
