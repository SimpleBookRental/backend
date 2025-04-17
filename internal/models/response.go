// Package models contains API models and response types.
package models

// ErrorResponse represents a standard error response.
// @Description Error response
type ErrorResponse struct {
	Success bool   `json:"success" example:"false"`         // Success flag
	Message string `json:"message" example:"Error message"` // Error message
	Error   string `json:"error" example:"Validation failed"` // Error details
}

// SuccessResponse represents a standard success response.
// @Description Success response
type SuccessResponse struct {
	Success bool        `json:"success" example:"true"`         // Success flag
	Message string      `json:"message" example:"Operation successful"` // Success message
	Data    interface{} `json:"data,omitempty"`                  // Optional data
}
