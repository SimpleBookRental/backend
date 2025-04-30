package api

import (
	"strconv"

	"github.com/SimpleBookRental/backend/internal/domain"
	"github.com/SimpleBookRental/backend/internal/service"
	"github.com/SimpleBookRental/backend/pkg/auth"
	"github.com/SimpleBookRental/backend/pkg/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// ReportHandler handles report requests
type ReportHandler struct {
	reportService service.ReportService
	jwtService    *auth.JWTService
	logger        *logger.Logger
}

// NewReportHandler creates a new ReportHandler
func NewReportHandler(reportService service.ReportService, jwtService *auth.JWTService, logger *logger.Logger) *ReportHandler {
	return &ReportHandler{
		reportService: reportService,
		jwtService:    jwtService,
		logger:        logger,
	}
}

// RevenueReportRequest represents a revenue report request
type RevenueReportRequest struct {
	StartDate string `form:"start_date" binding:"required"`
	EndDate   string `form:"end_date" binding:"required"`
}

// GetPopularBooks handles getting popular books
func (h *ReportHandler) GetPopularBooks(c *gin.Context) {
	// Only admins and librarians can access reports
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

	books, err := h.reportService.GetPopularBooks(int32(limit), int32(offset))
	if err != nil {
		h.logger.Error("Failed to get popular books", zap.Error(err))
		SendError(c, err)
		return
	}

	SendPaginated(c, books, int64(len(books)), int32(limit), int32(offset), "Popular books retrieved successfully")
}

// GetRevenueReport handles getting a revenue report
func (h *ReportHandler) GetRevenueReport(c *gin.Context) {
	// Only admins can access revenue reports
	userRole, exists := c.Get("userRole")
	if !exists {
		SendError(c, domain.ErrUnauthorized)
		return
	}

	role := domain.UserRole(userRole.(string))
	if role != domain.RoleAdmin {
		SendError(c, domain.ErrForbidden)
		return
	}

	var req RevenueReportRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		h.logger.Error("Invalid request parameters", zap.Error(err))
		SendError(c, domain.NewInvalidInputError(err.Error()))
		return
	}

	report, err := h.reportService.GetRevenueReport(req.StartDate, req.EndDate)
	if err != nil {
		h.logger.Error("Failed to get revenue report", zap.Error(err))
		SendError(c, err)
		return
	}

	SendSuccess(c, report, "Revenue report retrieved successfully")
}

// GetOverdueBooks handles getting overdue books
func (h *ReportHandler) GetOverdueBooks(c *gin.Context) {
	// Only admins and librarians can access reports
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

	rentals, err := h.reportService.GetOverdueBooks(int32(limit), int32(offset))
	if err != nil {
		h.logger.Error("Failed to get overdue books", zap.Error(err))
		SendError(c, err)
		return
	}

	SendPaginated(c, rentals, int64(len(rentals)), int32(limit), int32(offset), "Overdue books retrieved successfully")
}
