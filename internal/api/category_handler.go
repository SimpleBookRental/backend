package api

import (
	"net/http"
	"strconv"

	"github.com/SimpleBookRental/backend/internal/domain"
	"github.com/SimpleBookRental/backend/pkg/auth"
	"github.com/SimpleBookRental/backend/pkg/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// CategoryHandler handles category requests
type CategoryHandler struct {
	categoryService domain.CategoryService
	jwtService      *auth.JWTService
	logger          *logger.Logger
}

// NewCategoryHandler creates a new CategoryHandler
func NewCategoryHandler(categoryService domain.CategoryService, jwtService *auth.JWTService, logger *logger.Logger) *CategoryHandler {
	return &CategoryHandler{
		categoryService: categoryService,
		jwtService:      jwtService,
		logger:          logger,
	}
}

// CategoryRequest represents a category request
type CategoryRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

// GetByID handles getting a category by ID
func (h *CategoryHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		h.logger.Error("Invalid category ID", zap.Error(err))
		SendError(c, domain.NewInvalidInputError("invalid category ID"))
		return
	}

	category, err := h.categoryService.GetByID(id)
	if err != nil {
		h.logger.Error("Failed to get category by ID", zap.Int64("id", id), zap.Error(err))
		SendError(c, err)
		return
	}

	SendSuccess(c, category, "Category retrieved successfully")
}

// List handles listing categories with pagination
func (h *CategoryHandler) List(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	categories, err := h.categoryService.List(int32(limit), int32(offset))
	if err != nil {
		h.logger.Error("Failed to list categories", zap.Error(err))
		SendError(c, err)
		return
	}

	SendPaginated(c, categories, int64(len(categories)), int32(limit), int32(offset), "Categories retrieved successfully")
}

// ListAll handles listing all categories
func (h *CategoryHandler) ListAll(c *gin.Context) {
	categories, err := h.categoryService.ListAll()
	if err != nil {
		h.logger.Error("Failed to list all categories", zap.Error(err))
		SendError(c, err)
		return
	}

	SendSuccess(c, categories, "All categories retrieved successfully")
}

// Create handles creating a category
func (h *CategoryHandler) Create(c *gin.Context) {
	// Only admins and librarians can create categories
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

	var req CategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid request body", zap.Error(err))
		SendError(c, domain.NewInvalidInputError(err.Error()))
		return
	}

	category := &domain.Category{
		Name:        req.Name,
		Description: req.Description,
	}

	createdCategory, err := h.categoryService.Create(category)
	if err != nil {
		h.logger.Error("Failed to create category", zap.Error(err))
		SendError(c, err)
		return
	}

	SendCreated(c, createdCategory, "Category created successfully")
}

// Update handles updating a category
func (h *CategoryHandler) Update(c *gin.Context) {
	// Only admins and librarians can update categories
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

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		h.logger.Error("Invalid category ID", zap.Error(err))
		SendError(c, domain.NewInvalidInputError("invalid category ID"))
		return
	}

	var req CategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid request body", zap.Error(err))
		SendError(c, domain.NewInvalidInputError(err.Error()))
		return
	}

	// Get existing category
	existingCategory, err := h.categoryService.GetByID(id)
	if err != nil {
		h.logger.Error("Failed to get category by ID", zap.Int64("id", id), zap.Error(err))
		SendError(c, err)
		return
	}

	// Update category fields
	existingCategory.Name = req.Name
	existingCategory.Description = req.Description

	updatedCategory, err := h.categoryService.Update(existingCategory)
	if err != nil {
		h.logger.Error("Failed to update category", zap.Int64("id", id), zap.Error(err))
		SendError(c, err)
		return
	}

	SendSuccess(c, updatedCategory, "Category updated successfully")
}

// Delete handles deleting a category
func (h *CategoryHandler) Delete(c *gin.Context) {
	// Only admins and librarians can delete categories
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

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		h.logger.Error("Invalid category ID", zap.Error(err))
		SendError(c, domain.NewInvalidInputError("invalid category ID"))
		return
	}

	err = h.categoryService.Delete(id)
	if err != nil {
		h.logger.Error("Failed to delete category", zap.Int64("id", id), zap.Error(err))
		SendError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Category deleted successfully"})
}
