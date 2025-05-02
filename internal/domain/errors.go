package domain

import (
	"errors"
	"fmt"
)

// Common errors
var (
	ErrNotFound          = errors.New("resource not found")
	ErrInvalidInput      = errors.New("invalid input")
	ErrUnauthorized      = errors.New("unauthorized")
	ErrForbidden         = errors.New("forbidden")
	ErrInternalServer    = errors.New("internal server error")
	ErrConflict          = errors.New("resource conflict")
	ErrResourceExhausted = errors.New("resource exhausted")
)

// User errors
var (
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidPassword   = errors.New("invalid password")
)

// Book errors
var (
	ErrBookNotFound      = errors.New("book not found")
	ErrBookAlreadyExists = errors.New("book already exists")
	ErrBookNotAvailable  = errors.New("book not available")
)

// Category errors
var (
	ErrCategoryNotFound      = errors.New("category not found")
	ErrCategoryAlreadyExists = errors.New("category already exists")
)

// Rental errors
var (
	ErrRentalNotFound      = errors.New("rental not found")
	ErrRentalAlreadyExists = errors.New("rental already exists")
	ErrRentalNotActive     = errors.New("rental not active")
	ErrRentalOverdue       = errors.New("rental is overdue")
)

// Payment errors
var (
	ErrPaymentNotFound      = errors.New("payment not found")
	ErrPaymentAlreadyExists = errors.New("payment already exists")
	ErrPaymentFailed        = errors.New("payment failed")
	ErrInvalidPaymentStatus = errors.New("invalid payment status")
)

// ErrorResponse represents an error response for API
type ErrorResponse struct {
	Success bool   `json:"success" example:"false"`
	Error   string `json:"error" example:"Resource not found"`
}

// AppError represents an application error
type AppError struct {
	Err     error
	Message string
	Code    int
}

// Error returns the error message
func (e *AppError) Error() string {
	if e.Message != "" {
		return e.Message
	}
	return e.Err.Error()
}

// Unwrap returns the wrapped error
func (e *AppError) Unwrap() error {
	return e.Err
}

// NewAppError creates a new AppError
func NewAppError(err error, message string, code int) *AppError {
	return &AppError{
		Err:     err,
		Message: message,
		Code:    code,
	}
}

// NewNotFoundError creates a new not found error
func NewNotFoundError(resource string, id interface{}) *AppError {
	return &AppError{
		Err:     ErrNotFound,
		Message: fmt.Sprintf("%s with ID %v not found", resource, id),
		Code:    404,
	}
}

// NewConflictError creates a new conflict error
func NewConflictError(resource string, field string, value interface{}) *AppError {
	return &AppError{
		Err:     ErrConflict,
		Message: fmt.Sprintf("%s with %s %v already exists", resource, field, value),
		Code:    409,
	}
}

// NewUnauthorizedError creates a new unauthorized error
func NewUnauthorizedError(message string) *AppError {
	return &AppError{
		Err:     ErrUnauthorized,
		Message: message,
		Code:    401,
	}
}

// NewForbiddenError creates a new forbidden error
func NewForbiddenError(message string) *AppError {
	return &AppError{
		Err:     ErrForbidden,
		Message: message,
		Code:    403,
	}
}

// NewInvalidInputError creates a new invalid input error
func NewInvalidInputError(message string) *AppError {
	return &AppError{
		Err:     ErrInvalidInput,
		Message: message,
		Code:    400,
	}
}

// NewInternalServerError creates a new internal server error
func NewInternalServerError(err error) *AppError {
	return &AppError{
		Err:     ErrInternalServer,
		Message: "internal server error",
		Code:    500,
	}
}
