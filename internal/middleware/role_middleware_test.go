package middleware

import (
	"net/http/httptest"
	"testing"

	"github.com/SimpleBookRental/backend/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestRequireRole(t *testing.T) {
	// Test with required role matching
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("role", "ADMIN")

	middleware := RequireRole("ADMIN")
	middleware(c)

	// If middleware calls c.Next(), it means it passed
	assert.False(t, c.IsAborted())

	// Test with required role not matching
	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)
	c.Set("role", "USER")

	middleware = RequireRole("ADMIN")
	middleware(c)

	// If middleware calls c.Abort(), it means it failed
	assert.True(t, c.IsAborted())

	// Test with role not in context
	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	middleware = RequireRole("ADMIN")
	middleware(c)

	// If middleware calls c.Abort(), it means it failed
	assert.True(t, c.IsAborted())
}

func TestRequireAdmin(t *testing.T) {
	// Test with admin role
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("role", models.AdminRole)

	middleware := RequireAdmin()
	middleware(c)

	// If middleware calls c.Next(), it means it passed
	assert.False(t, c.IsAborted())

	// Test with non-admin role
	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)
	c.Set("role", models.UserRole)

	middleware = RequireAdmin()
	middleware(c)

	// If middleware calls c.Abort(), it means it failed
	assert.True(t, c.IsAborted())
}

func TestRequireUser(t *testing.T) {
	// Test with user role
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("role", models.UserRole)

	middleware := RequireUser()
	middleware(c)

	// If middleware calls c.Next(), it means it passed
	assert.False(t, c.IsAborted())

	// Test with non-user role
	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)
	c.Set("role", models.AdminRole)

	middleware = RequireUser()
	middleware(c)

	// If middleware calls c.Abort(), it means it failed
	assert.True(t, c.IsAborted())
}

func TestRequireAdminOrSameUser(t *testing.T) {
	// Test with admin role
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("role", models.AdminRole)
	c.Set("user_id", "123")

	middleware := RequireAdminOrSameUser()
	middleware(c)

	// If middleware calls c.Next(), it means it passed
	assert.False(t, c.IsAborted())

	// Test with same user
	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)
	c.Set("role", models.UserRole)
	c.Set("user_id", "123")
	c.Params = []gin.Param{{Key: "id", Value: "123"}}

	middleware = RequireAdminOrSameUser()
	middleware(c)

	// If middleware calls c.Next(), it means it passed
	assert.False(t, c.IsAborted())

	// Test with different user
	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)
	c.Set("role", models.UserRole)
	c.Set("user_id", "456")
	c.Params = []gin.Param{{Key: "id", Value: "123"}}

	middleware = RequireAdminOrSameUser()
	middleware(c)

	// If middleware calls c.Abort(), it means it failed
	assert.True(t, c.IsAborted())

	// Test with missing role
	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)
	c.Set("user_id", "123")

	middleware = RequireAdminOrSameUser()
	middleware(c)

	// If middleware calls c.Abort(), it means it failed
	assert.True(t, c.IsAborted())

	// Test with missing user_id
	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)
	c.Set("role", models.UserRole)

	middleware = RequireAdminOrSameUser()
	middleware(c)

	// If middleware calls c.Abort(), it means it failed
	assert.True(t, c.IsAborted())
}
