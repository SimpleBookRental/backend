package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

/*
	Refactored response helpers to return data directly (no wrapping in "data", "success", "message").
	For error responses, return a simple JSON with "error" field.
*/

// Created returns a 201 Created response with data directly
func Created(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, data)
}

// OK returns a 200 OK response with data directly
func OK(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, data)
}

// BadRequest returns a 400 Bad Request response with error message
func BadRequest(c *gin.Context, err interface{}) {
	c.JSON(http.StatusBadRequest, gin.H{"error": err})
}

// NotFound returns a 404 Not Found response with error message
func NotFound(c *gin.Context, err interface{}) {
	c.JSON(http.StatusNotFound, gin.H{"error": err})
}

// InternalServerError returns a 500 Internal Server Error response with error message
func InternalServerError(c *gin.Context, err interface{}) {
	c.JSON(http.StatusInternalServerError, gin.H{"error": err})
}

// Forbidden returns a 403 Forbidden response with error message
func Forbidden(c *gin.Context, err interface{}) {
	c.JSON(http.StatusForbidden, gin.H{"error": err})
}
