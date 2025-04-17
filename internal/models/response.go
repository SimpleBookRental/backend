// Package models contains API models and response types.
package models

// ErrorResponse represents a standard error response.
// @Description Error response
type ErrorResponse struct {
	Error string `json:"error" example:"Validation failed"` // Error details
}
