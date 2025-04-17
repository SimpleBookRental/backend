// Unit tests for role-based middleware using testify.
package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/SimpleBookRental/backend/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupGinWithRole(mw gin.HandlerFunc, inject func(*gin.Context)) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		if inject != nil {
			inject(c)
		}
		c.Next()
	})
	r.Use(mw)
	r.GET("/protected", func(c *gin.Context) {
		c.JSON(200, gin.H{"success": true})
	})
	return r
}

func TestRequireRole_Success(t *testing.T) {
	router := setupGinWithRole(RequireRole(models.AdminRole), func(c *gin.Context) {
		c.Set("role", models.AdminRole)
	})
	req, _ := http.NewRequest("GET", "/protected", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
}

func TestRequireRole_Forbidden(t *testing.T) {
	router := setupGinWithRole(RequireRole(models.AdminRole), func(c *gin.Context) {
		c.Set("role", models.UserRole)
	})
	req, _ := http.NewRequest("GET", "/protected", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, 403, w.Code)
	assert.Contains(t, w.Body.String(), "insufficient permissions")
}

func TestRequireRole_MissingRole(t *testing.T) {
	router := setupGinWithRole(RequireRole(models.AdminRole), nil)
	req, _ := http.NewRequest("GET", "/protected", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, 401, w.Code)
	assert.Contains(t, w.Body.String(), "role not found")
}

func TestRequireAdmin_Success(t *testing.T) {
	router := setupGinWithRole(RequireAdmin(), func(c *gin.Context) {
		c.Set("role", models.AdminRole)
	})
	req, _ := http.NewRequest("GET", "/protected", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
}

func TestRequireUser_Success(t *testing.T) {
	router := setupGinWithRole(RequireUser(), func(c *gin.Context) {
		c.Set("role", models.UserRole)
	})
	req, _ := http.NewRequest("GET", "/protected", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
}

func TestRequireAdminOrSameUser_Admin(t *testing.T) {
	router := setupGinWithRole(RequireAdminOrSameUser(), func(c *gin.Context) {
		c.Set("role", models.AdminRole)
		c.Set("user_id", "user-1")
	})
	req, _ := http.NewRequest("GET", "/protected", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
}

func TestRequireAdminOrSameUser_SameUser(t *testing.T) {
	router := setupGinWithRole(RequireAdminOrSameUser(), func(c *gin.Context) {
		c.Set("role", models.UserRole)
		c.Set("user_id", "user-1")
		c.Params = gin.Params{{Key: "id", Value: "user-1"}}
	})
	req, _ := http.NewRequest("GET", "/protected?id=user-1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
}

func TestRequireAdminOrSameUser_Forbidden(t *testing.T) {
	router := setupGinWithRole(RequireAdminOrSameUser(), func(c *gin.Context) {
		c.Set("role", models.UserRole)
		c.Set("user_id", "user-1")
		c.Params = gin.Params{{Key: "id", Value: "user-2"}}
	})
	req, _ := http.NewRequest("GET", "/protected?id=user-2", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, 403, w.Code)
	assert.Contains(t, w.Body.String(), "insufficient permissions")
}

func TestRequireAdminOrSameUser_MissingContext(t *testing.T) {
	router := setupGinWithRole(RequireAdminOrSameUser(), nil)
	req, _ := http.NewRequest("GET", "/protected", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, 401, w.Code)
	assert.Contains(t, w.Body.String(), "user information not found")
}
