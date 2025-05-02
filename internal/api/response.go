package api

import (
	"errors"
	"net/http"

	"github.com/SimpleBookRental/backend/internal/domain"
	"github.com/gin-gonic/gin"
)

// Response represents a standard API response
type Response struct {
	Success bool        `json:"success" example:"true"`
	Message string      `json:"message,omitempty" example:"Operation successful"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// PaginatedResponse represents a paginated API response
type PaginatedResponse struct {
	Success bool        `json:"success" example:"true"`
	Message string      `json:"message,omitempty" example:"Data retrieved successfully"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
	Total   int64       `json:"total" example:"100"`
	Limit   int32       `json:"limit" example:"10"`
	Offset  int32       `json:"offset" example:"0"`
}

// NewSuccessResponse creates a new success response
func NewSuccessResponse(data interface{}, message string) Response {
	return Response{
		Success: true,
		Message: message,
		Data:    data,
	}
}

// NewErrorResponse creates a new error response
func NewErrorResponse(err error) Response {
	return Response{
		Success: false,
		Error:   err.Error(),
	}
}

// NewPaginatedResponse creates a new paginated response
func NewPaginatedResponse(data interface{}, total int64, limit, offset int32, message string) PaginatedResponse {
	return PaginatedResponse{
		Success: true,
		Message: message,
		Data:    data,
		Total:   total,
		Limit:   limit,
		Offset:  offset,
	}
}

// SendSuccess sends a success response
func SendSuccess(c *gin.Context, data interface{}, message string) {
	c.JSON(http.StatusOK, NewSuccessResponse(data, message))
}

// SendCreated sends a created response
func SendCreated(c *gin.Context, data interface{}, message string) {
	c.JSON(http.StatusCreated, NewSuccessResponse(data, message))
}

// SendError sends an error response
func SendError(c *gin.Context, err error) {
	var statusCode int
	var appErr *domain.AppError

	if errors.As(err, &appErr) {
		statusCode = appErr.Code
	} else {
		switch {
		case errors.Is(err, domain.ErrNotFound) || 
			 errors.Is(err, domain.ErrUserNotFound) || 
			 errors.Is(err, domain.ErrBookNotFound) || 
			 errors.Is(err, domain.ErrCategoryNotFound) || 
			 errors.Is(err, domain.ErrRentalNotFound) || 
			 errors.Is(err, domain.ErrPaymentNotFound):
			statusCode = http.StatusNotFound
		case errors.Is(err, domain.ErrInvalidInput) || 
			 errors.Is(err, domain.ErrInvalidCredentials) || 
			 errors.Is(err, domain.ErrInvalidPassword):
			statusCode = http.StatusBadRequest
		case errors.Is(err, domain.ErrUnauthorized):
			statusCode = http.StatusUnauthorized
		case errors.Is(err, domain.ErrForbidden):
			statusCode = http.StatusForbidden
		case errors.Is(err, domain.ErrConflict) || 
			 errors.Is(err, domain.ErrUserAlreadyExists) || 
			 errors.Is(err, domain.ErrBookAlreadyExists) || 
			 errors.Is(err, domain.ErrCategoryAlreadyExists) || 
			 errors.Is(err, domain.ErrRentalAlreadyExists) || 
			 errors.Is(err, domain.ErrPaymentAlreadyExists):
			statusCode = http.StatusConflict
		case errors.Is(err, domain.ErrResourceExhausted) || 
			 errors.Is(err, domain.ErrBookNotAvailable):
			statusCode = http.StatusTooManyRequests
		default:
			statusCode = http.StatusInternalServerError
		}
	}

	c.JSON(statusCode, NewErrorResponse(err))
}

// SendPaginated sends a paginated response
func SendPaginated(c *gin.Context, data interface{}, total int64, limit, offset int32, message string) {
	c.JSON(http.StatusOK, NewPaginatedResponse(data, total, limit, offset, message))
}
