package utils

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestOK(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Test case 1: OK response with data
	data := map[string]string{"key": "value"}
	OK(c, "Success", data)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Success")
	assert.Contains(t, w.Body.String(), "key")
	assert.Contains(t, w.Body.String(), "value")
}

func TestBadRequest(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Test case 1: Bad request response with error
	BadRequest(c, "Bad request", "Invalid input")

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Bad request")
	assert.Contains(t, w.Body.String(), "Invalid input")
}

func TestNotFound(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Test case 1: Not found response
	NotFound(c, "Resource not found")

	// Assert
	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "Resource not found")
}

func TestInternalServerError(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Test case 1: Internal server error response
	InternalServerError(c, "Internal server error", "Something went wrong")

	// Assert
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "Internal server error")
	assert.Contains(t, w.Body.String(), "Something went wrong")
}
