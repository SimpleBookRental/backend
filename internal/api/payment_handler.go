package api

import (
	"strconv"

	"github.com/SimpleBookRental/backend/internal/domain"
	"github.com/SimpleBookRental/backend/pkg/auth"
	"github.com/SimpleBookRental/backend/pkg/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// PaymentHandler handles payment requests
type PaymentHandler struct {
	paymentService domain.PaymentService
	jwtService     *auth.JWTService
	logger         *logger.Logger
}

// NewPaymentHandler creates a new PaymentHandler
func NewPaymentHandler(paymentService domain.PaymentService, jwtService *auth.JWTService, logger *logger.Logger) *PaymentHandler {
	return &PaymentHandler{
		paymentService: paymentService,
		jwtService:     jwtService,
		logger:         logger,
	}
}

// PaymentRequest represents a payment request
type PaymentRequest struct {
	RentalID      *int64  `json:"rental_id" example:"1"`
	Amount        float64 `json:"amount" binding:"required,gt=0" example:"15.50"`
	PaymentMethod string  `json:"payment_method" binding:"required" example:"credit_card"`
}

// GetByID handles getting a payment by ID
// @Summary      Get a payment by ID
// @Description  Retrieve a single payment by its ID. Users can only view their own payments unless they are admins/librarians.
// @Tags         payments
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Payment ID"
// @Success      200  {object}  domain.Payment
// @Failure      400  {object}  domain.ErrorResponse
// @Failure      401  {object}  domain.ErrorResponse
// @Failure      403  {object}  domain.ErrorResponse
// @Failure      404  {object}  domain.ErrorResponse
// @Failure      500  {object}  domain.ErrorResponse
// @Security     Bearer
// @Router       /payments/{id} [get]
func (h *PaymentHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		h.logger.Error("Invalid payment ID", zap.Error(err))
		SendError(c, domain.NewInvalidInputError("invalid payment ID"))
		return
	}

	payment, err := h.paymentService.GetByID(id)
	if err != nil {
		h.logger.Error("Failed to get payment by ID", zap.Int64("id", id), zap.Error(err))
		SendError(c, err)
		return
	}

	// Check if user is requesting their own payment or is an admin/librarian
	userID, exists := c.Get("userID")
	if !exists {
		SendError(c, domain.ErrUnauthorized)
		return
	}

	userRole, _ := c.Get("userRole")
	role := domain.UserRole(userRole.(string))

	if userID.(int64) != payment.UserID && !auth.IsLibrarian(role) {
		SendError(c, domain.ErrForbidden)
		return
	}

	SendSuccess(c, payment, "Payment retrieved successfully")
}

// List handles listing payments with pagination
// @Summary      List all payments
// @Description  Get a paginated list of all payments. Only admins can access this endpoint.
// @Tags         payments
// @Accept       json
// @Produce      json
// @Param        limit  query    int     false  "Limit"  default(10)
// @Param        offset query    int     false  "Offset" default(0)
// @Success      200    {object} PaginatedResponse{data=[]domain.Payment}
// @Failure      401    {object} domain.ErrorResponse
// @Failure      403    {object} domain.ErrorResponse
// @Failure      500    {object} domain.ErrorResponse
// @Security     Bearer
// @Router       /payments [get]
func (h *PaymentHandler) List(c *gin.Context) {
	// Only admins and librarians can list all payments
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

	payments, err := h.paymentService.List(int32(limit), int32(offset))
	if err != nil {
		h.logger.Error("Failed to list payments", zap.Error(err))
		SendError(c, err)
		return
	}

	SendPaginated(c, payments, int64(len(payments)), int32(limit), int32(offset), "Payments retrieved successfully")
}

// ListByUser handles listing payments for a specific user with pagination
// @Summary      List user payments
// @Description  Get a paginated list of payments for a specific user. Users can only view their own payments unless they are admins/librarians.
// @Tags         payments
// @Accept       json
// @Produce      json
// @Param        userId path     int     true   "User ID"
// @Param        limit  query    int     false  "Limit"  default(10)
// @Param        offset query    int     false  "Offset" default(0)
// @Success      200    {object} PaginatedResponse{data=[]domain.Payment}
// @Failure      400    {object} domain.ErrorResponse
// @Failure      401    {object} domain.ErrorResponse
// @Failure      403    {object} domain.ErrorResponse
// @Failure      500    {object} domain.ErrorResponse
// @Security     Bearer
// @Router       /payments/user/{userId} [get]
func (h *PaymentHandler) ListByUser(c *gin.Context) {
	userID, err := strconv.ParseInt(c.Param("userId"), 10, 64)
	if err != nil {
		h.logger.Error("Invalid user ID", zap.Error(err))
		SendError(c, domain.NewInvalidInputError("invalid user ID"))
		return
	}

	// Check if user is requesting their own payments or is an admin/librarian
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

	payments, err := h.paymentService.ListByUser(userID, int32(limit), int32(offset))
	if err != nil {
		h.logger.Error("Failed to list payments by user", zap.Int64("userID", userID), zap.Error(err))
		SendError(c, err)
		return
	}

	SendPaginated(c, payments, int64(len(payments)), int32(limit), int32(offset), "Payments retrieved successfully")
}

// Create handles creating a payment
// @Summary      Create a payment
// @Description  Create a new payment for the authenticated user
// @Tags         payments
// @Accept       json
// @Produce      json
// @Param        payment body      PaymentRequest true "Payment information"
// @Success      201    {object}   domain.Payment
// @Failure      400    {object}   domain.ErrorResponse
// @Failure      401    {object}   domain.ErrorResponse
// @Failure      500    {object}   domain.ErrorResponse
// @Security     Bearer
// @Router       /payments [post]
func (h *PaymentHandler) Create(c *gin.Context) {
	// Get user ID from context
	userID, exists := c.Get("userID")
	if !exists {
		SendError(c, domain.ErrUnauthorized)
		return
	}

	var req PaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid request body", zap.Error(err))
		SendError(c, domain.NewInvalidInputError(err.Error()))
		return
	}

	payment := &domain.Payment{
		UserID:        userID.(int64),
		RentalID:      req.RentalID,
		Amount:        req.Amount,
		PaymentMethod: req.PaymentMethod,
		Status:        domain.PaymentStatusPending,
	}

	createdPayment, err := h.paymentService.Create(payment)
	if err != nil {
		h.logger.Error("Failed to create payment", zap.Error(err))
		SendError(c, err)
		return
	}

	SendCreated(c, createdPayment, "Payment created successfully")
}

// Process handles processing a payment
// @Summary      Process a payment
// @Description  Process a payment transaction
// @Tags         payments
// @Accept       json
// @Produce      json
// @Param        payment body      PaymentRequest true "Payment information"
// @Success      200    {object}   domain.Payment
// @Failure      400    {object}   domain.ErrorResponse
// @Failure      401    {object}   domain.ErrorResponse
// @Failure      500    {object}   domain.ErrorResponse
// @Security     Bearer
// @Router       /payments/process [post]
func (h *PaymentHandler) Process(c *gin.Context) {
	// Get user ID from context
	userID, exists := c.Get("userID")
	if !exists {
		SendError(c, domain.ErrUnauthorized)
		return
	}

	var req PaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid request body", zap.Error(err))
		SendError(c, domain.NewInvalidInputError(err.Error()))
		return
	}

	payment := &domain.Payment{
		UserID:        userID.(int64),
		RentalID:      req.RentalID,
		Amount:        req.Amount,
		PaymentMethod: req.PaymentMethod,
		Status:        domain.PaymentStatusPending,
	}

	processedPayment, err := h.paymentService.ProcessPayment(payment)
	if err != nil {
		h.logger.Error("Failed to process payment", zap.Error(err))
		SendError(c, err)
		return
	}

	SendSuccess(c, processedPayment, "Payment processed successfully")
}

// Refund handles refunding a payment
// @Summary      Refund a payment
// @Description  Refund a processed payment. Only admins and librarians can refund payments.
// @Tags         payments
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Payment ID"
// @Success      200  {object}  domain.Payment
// @Failure      400  {object}  domain.ErrorResponse
// @Failure      401  {object}  domain.ErrorResponse
// @Failure      403  {object}  domain.ErrorResponse
// @Failure      404  {object}  domain.ErrorResponse
// @Failure      500  {object}  domain.ErrorResponse
// @Security     Bearer
// @Router       /payments/{id}/refund [put]
func (h *PaymentHandler) Refund(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		h.logger.Error("Invalid payment ID", zap.Error(err))
		SendError(c, domain.NewInvalidInputError("invalid payment ID"))
		return
	}

	// Only admins and librarians can refund payments
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

	refundedPayment, err := h.paymentService.RefundPayment(id)
	if err != nil {
		h.logger.Error("Failed to refund payment", zap.Int64("id", id), zap.Error(err))
		SendError(c, err)
		return
	}

	SendSuccess(c, refundedPayment, "Payment refunded successfully")
}
