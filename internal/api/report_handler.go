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
	StartDate string `form:"start_date" binding:"required" example:"2025-01-01"`
	EndDate   string `form:"end_date" binding:"required" example:"2025-05-01"`
}

// GetPopularBooks handles getting popular books
// @Summary      Get popular books
// @Description  Retrieve a list of popular books based on rental frequency. Only accessible by admins and librarians.
// @Tags         reports
// @Accept       json
// @Produce      json
// @Param        limit  query    int     false  "Limit"  default(10)
// @Param        offset query    int     false  "Offset" default(0)
// @Success      200    {object} PaginatedResponse{data=[]domain.Book}
// @Failure      401    {object} domain.ErrorResponse
// @Failure      403    {object} domain.ErrorResponse
// @Failure      500    {object} domain.ErrorResponse
// @Security     Bearer
// @Router       /reports/books/popular [get]
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
// @Summary      Get revenue report
// @Description  Retrieve a revenue report for a specific date range. Only accessible by admins.
// @Tags         reports
// @Accept       json
// @Produce      json
// @Param        start_date query    string  true   "Start date (YYYY-MM-DD)"
// @Param        end_date   query    string  true   "End date (YYYY-MM-DD)"
// @Success      200        {object} domain.RevenueReport
// @Failure      400        {object} domain.ErrorResponse
// @Failure      401        {object} domain.ErrorResponse
// @Failure      403        {object} domain.ErrorResponse
// @Failure      500        {object} domain.ErrorResponse
// @Security     Bearer
// @Router       /reports/revenue [get]
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
// @Summary      Get overdue books
// @Description  Retrieve a list of books that are currently overdue. Only accessible by admins and librarians.
// @Tags         reports
// @Accept       json
// @Produce      json
// @Param        limit  query    int     false  "Limit"  default(10)
// @Param        offset query    int     false  "Offset" default(0)
// @Success      200    {object} PaginatedResponse{data=[]domain.Rental}
// @Failure      401    {object} domain.ErrorResponse
// @Failure      403    {object} domain.ErrorResponse
// @Failure      500    {object} domain.ErrorResponse
// @Security     Bearer
// @Router       /reports/overdue [get]
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
