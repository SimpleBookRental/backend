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

// UserResponse represents a response containing a single user.
// @Description User response
type UserResponse struct {
	Success bool   `json:"success" example:"true"`
	Message string `json:"message" example:"User retrieved successfully"`
	Data    *User  `json:"data"`
}

// UsersResponse represents a response containing a list of users.
// @Description Users response
type UsersResponse struct {
	Success bool    `json:"success" example:"true"`
	Message string  `json:"message" example:"Users retrieved successfully"`
	Data    []*User `json:"data"`
}

// BookResponse represents a response containing a single book.
// @Description Book response
type BookResponse struct {
	Success bool   `json:"success" example:"true"`
	Message string `json:"message" example:"Book retrieved successfully"`
	Data    *Book  `json:"data"`
}

// BooksResponse represents a response containing a list of books.
// @Description Books response
type BooksResponse struct {
	Success bool    `json:"success" example:"true"`
	Message string  `json:"message" example:"Books retrieved successfully"`
	Data    []*Book `json:"data"`
}
