package api

import (
	"strconv"
	"time"

	"github.com/SimpleBookRental/backend/internal/domain"
	"github.com/SimpleBookRental/backend/pkg/auth"
	"github.com/SimpleBookRental/backend/pkg/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// RentalHandler handles rental requests
type RentalHandler struct {
	rentalService domain.RentalService
	jwtService    *auth.JWTService
	logger        *logger.Logger
}

// NewRentalHandler creates a new RentalHandler
func NewRentalHandler(rentalService domain.RentalService, jwtService *auth.JWTService, logger *logger.Logger) *RentalHandler {
	return &RentalHandler{
		rentalService: rentalService,
		jwtService:    jwtService,
		logger:        logger,
	}
}

// RentalRequest represents a rental request
type RentalRequest struct {
	BookID     int64     `json:"book_id" binding:"required"`
	RentalDate time.Time `json:"rental_date"`
	DueDate    time.Time `json:"due_date"`
}

// ExtendRentalRequest represents a rental extension request
type ExtendRentalRequest struct {
	Days int `json:"days" binding:"required,min=1"`
}

// GetByID handles getting a rental by ID
func (h *RentalHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		h.logger.Error("Invalid rental ID", zap.Error(err))
		SendError(c, domain.NewInvalidInputError("invalid rental ID"))
		return
	}

	rental, err := h.rentalService.GetByID(id)
	if err != nil {
		h.logger.Error("Failed to get rental by ID", zap.Int64("id", id), zap.Error(err))
		SendError(c, err)
		return
	}

	// Check if user is requesting their own rental or is an admin/librarian
	userID, exists := c.Get("userID")
	if !exists {
		SendError(c, domain.ErrUnauthorized)
		return
	}

	userRole, _ := c.Get("userRole")
	role := domain.UserRole(userRole.(string))

	if userID.(int64) != rental.UserID && !auth.IsLibrarian(role) {
		SendError(c, domain.ErrForbidden)
		return
	}

	SendSuccess(c, rental, "Rental retrieved successfully")
}

// List handles listing rentals with pagination
func (h *RentalHandler) List(c *gin.Context) {
	// Only admins and librarians can list all rentals
	userRole, exists := c.Get("userRole")
	if !exists {
		SendError(c, domain.ErrUnauthorized)
		return
	}

	role := domain.UserRole(userRole.(string))
	if !auth.IsLibrarian(role) {
		SendError(c, domain.ErrForbidden)
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	rentals, err := h.rentalService.List(int32(limit), int32(offset))
	if err != nil {
		h.logger.Error("Failed to list rentals", zap.Error(err))
		SendError(c, err)
		return
	}

	SendPaginated(c, rentals, int64(len(rentals)), int32(limit), int32(offset), "Rentals retrieved successfully")
}

// ListByUser handles listing rentals for a specific user with pagination
func (h *RentalHandler) ListByUser(c *gin.Context) {
	userID, err := strconv.ParseInt(c.Param("userId"), 10, 64)
	if err != nil {
		h.logger.Error("Invalid user ID", zap.Error(err))
		SendError(c, domain.NewInvalidInputError("invalid user ID"))
		return
	}

	// Check if user is requesting their own rentals or is an admin/librarian
	currentUserID, exists := c.Get("userID")
	if !exists {
		SendError(c, domain.ErrUnauthorized)
		return
	}

	userRole, _ := c.Get("userRole")
	role := domain.UserRole(userRole.(string))

	if currentUserID.(int64) != userID && !auth.IsLibrarian(role) {
		SendError(c, domain.ErrForbidden)
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	rentals, err := h.rentalService.ListByUser(userID, int32(limit), int32(offset))
	if err != nil {
		h.logger.Error("Failed to list rentals by user", zap.Int64("userID", userID), zap.Error(err))
		SendError(c, err)
		return
	}

	SendPaginated(c, rentals, int64(len(rentals)), int32(limit), int32(offset), "Rentals retrieved successfully")
}

// Create handles creating a rental
func (h *RentalHandler) Create(c *gin.Context) {
	// Get user ID from context
	userID, exists := c.Get("userID")
	if !exists {
		SendError(c, domain.ErrUnauthorized)
		return
	}

	var req RentalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid request body", zap.Error(err))
		SendError(c, domain.NewInvalidInputError(err.Error()))
		return
	}

	rental := &domain.Rental{
		UserID:     userID.(int64),
		BookID:     req.BookID,
		RentalDate: req.RentalDate,
		DueDate:    req.DueDate,
		Status:     domain.RentalStatusActive,
	}

	createdRental, err := h.rentalService.Create(rental)
	if err != nil {
		h.logger.Error("Failed to create rental", zap.Error(err))
		SendError(c, err)
		return
	}

	SendCreated(c, createdRental, "Rental created successfully")
}

// Return handles returning a rental
func (h *RentalHandler) Return(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		h.logger.Error("Invalid rental ID", zap.Error(err))
		SendError(c, domain.NewInvalidInputError("invalid rental ID"))
		return
	}

	// Get rental to check ownership
	rental, err := h.rentalService.GetByID(id)
	if err != nil {
		h.logger.Error("Failed to get rental by ID", zap.Int64("id", id), zap.Error(err))
		SendError(c, err)
		return
	}

	// Check if user is returning their own rental or is an admin/librarian
	userID, exists := c.Get("userID")
	if !exists {
		SendError(c, domain.ErrUnauthorized)
		return
	}

	userRole, _ := c.Get("userRole")
	role := domain.UserRole(userRole.(string))

	if userID.(int64) != rental.UserID && !auth.IsLibrarian(role) {
		SendError(c, domain.ErrForbidden)
		return
	}

	returnedRental, err := h.rentalService.Return(id)
	if err != nil {
		h.logger.Error("Failed to return rental", zap.Int64("id", id), zap.Error(err))
		SendError(c, err)
		return
	}

	// Calculate late fee if any
	lateFee, err := h.rentalService.CalculateLateFee(returnedRental)
	if err != nil {
		h.logger.Error("Failed to calculate late fee", zap.Int64("id", id), zap.Error(err))
		// Continue anyway, just log the error
	}

	response := gin.H{
		"rental":   returnedRental,
		"late_fee": lateFee,
	}

	SendSuccess(c, response, "Rental returned successfully")
}

// Extend handles extending a rental
func (h *RentalHandler) Extend(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		h.logger.Error("Invalid rental ID", zap.Error(err))
		SendError(c, domain.NewInvalidInputError("invalid rental ID"))
		return
	}

	// Get rental to check ownership
	rental, err := h.rentalService.GetByID(id)
	if err != nil {
		h.logger.Error("Failed to get rental by ID", zap.Int64("id", id), zap.Error(err))
		SendError(c, err)
		return
	}

	// Check if user is extending their own rental or is an admin/librarian
	userID, exists := c.Get("userID")
	if !exists {
		SendError(c, domain.ErrUnauthorized)
		return
	}

	userRole, _ := c.Get("userRole")
	role := domain.UserRole(userRole.(string))

	if userID.(int64) != rental.UserID && !auth.IsLibrarian(role) {
		SendError(c, domain.ErrForbidden)
		return
	}

	var req ExtendRentalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid request body", zap.Error(err))
		SendError(c, domain.NewInvalidInputError(err.Error()))
		return
	}

	extendedRental, err := h.rentalService.Extend(id, req.Days)
	if err != nil {
		h.logger.Error("Failed to extend rental", zap.Int64("id", id), zap.Error(err))
		SendError(c, err)
		return
	}

	SendSuccess(c, extendedRental, "Rental extended successfully")
}
