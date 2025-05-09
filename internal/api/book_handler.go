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

// BookHandler handles book requests
type BookHandler struct {
	bookService domain.BookService
	jwtService  *auth.JWTService
	logger      *logger.Logger
}

// NewBookHandler creates a new BookHandler
func NewBookHandler(bookService domain.BookService, jwtService *auth.JWTService, logger *logger.Logger) *BookHandler {
	return &BookHandler{
		bookService: bookService,
		jwtService:  jwtService,
		logger:      logger,
	}
}

// BookRequest represents a book request
type BookRequest struct {
	Title         string `json:"title" binding:"required"`
	Author        string `json:"author" binding:"required"`
	ISBN          string `json:"isbn" binding:"required"`
	Description   string `json:"description"`
	PublishedYear int32  `json:"published_year"`
	Publisher     string `json:"publisher"`
	TotalCopies   int32  `json:"total_copies" binding:"required,min=1"`
	CategoryID    int64  `json:"category_id"`
}

// BookCopiesRequest represents a book copies update request
type BookCopiesRequest struct {
	TotalCopies     int32 `json:"total_copies" binding:"required,min=1"`
	AvailableCopies int32 `json:"available_copies" binding:"required,min=0"`
}

// BookSearchRequest represents a book search request
type BookSearchRequest struct {
	Title         string `form:"title"`
	Author        string `form:"author"`
	ISBN          string `form:"isbn"`
	PublishedYear int32  `form:"published_year"`
	CategoryID    int64  `form:"category_id"`
	Available     bool   `form:"available"`
	Limit         int32  `form:"limit,default=10"`
	Offset        int32  `form:"offset,default=0"`
}

// GetByID handles getting a book by ID
// @Summary      Get a book by ID
// @Description  Retrieve a single book by its ID
// @Tags         books
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Book ID"
// @Success      200  {object}  domain.Book
// @Failure      400  {object}  domain.ErrorResponse
// @Failure      404  {object}  domain.ErrorResponse
// @Failure      500  {object}  domain.ErrorResponse
// @Router       /books/{id} [get]
func (h *BookHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		h.logger.Error("Invalid book ID", zap.Error(err))
		SendError(c, domain.NewInvalidInputError("invalid book ID"))
		return
	}

	book, err := h.bookService.GetByID(id)
	if err != nil {
		h.logger.Error("Failed to get book by ID", zap.Int64("id", id), zap.Error(err))
		SendError(c, err)
		return
	}

	SendSuccess(c, book, "Book retrieved successfully")
}

// List handles listing books with pagination
// @Summary      List all books
// @Description  Get a paginated list of all books
// @Tags         books
// @Accept       json
// @Produce      json
// @Param        limit  query    int     false  "Limit"  default(10)
// @Param        offset query    int     false  "Offset" default(0)
// @Success      200    {object} PaginatedResponse{data=[]domain.Book}
// @Failure      400    {object} domain.ErrorResponse
// @Failure      500    {object} domain.ErrorResponse
// @Router       /books [get]
func (h *BookHandler) List(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	books, err := h.bookService.List(int32(limit), int32(offset))
	if err != nil {
		h.logger.Error("Failed to list books", zap.Error(err))
		SendError(c, err)
		return
	}

	SendPaginated(c, books, int64(len(books)), int32(limit), int32(offset), "Books retrieved successfully")
}

// ListByCategory handles listing books by category with pagination
// @Summary      List books by category
// @Description  Get a paginated list of books filtered by category ID
// @Tags         books
// @Accept       json
// @Produce      json
// @Param        id     path     int     true   "Category ID"
// @Param        limit  query    int     false  "Limit"  default(10)
// @Param        offset query    int     false  "Offset" default(0)
// @Success      200    {object} PaginatedResponse{data=[]domain.Book}
// @Failure      400    {object} domain.ErrorResponse
// @Failure      404    {object} domain.ErrorResponse
// @Failure      500    {object} domain.ErrorResponse
// @Router       /categories/{id}/books [get]
func (h *BookHandler) ListByCategory(c *gin.Context) {
	categoryID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		h.logger.Error("Invalid category ID", zap.Error(err))
		SendError(c, domain.NewInvalidInputError("invalid category ID"))
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	books, err := h.bookService.ListByCategory(categoryID, int32(limit), int32(offset))
	if err != nil {
		h.logger.Error("Failed to list books by category", zap.Int64("categoryID", categoryID), zap.Error(err))
		SendError(c, err)
		return
	}

	SendPaginated(c, books, int64(len(books)), int32(limit), int32(offset), "Books retrieved successfully")
}

// Search handles searching books
// @Summary      Search books
// @Description  Search books with various filters
// @Tags         books
// @Accept       json
// @Produce      json
// @Param        title         query    string  false  "Title"
// @Param        author        query    string  false  "Author"
// @Param        isbn          query    string  false  "ISBN"
// @Param        published_year query    int     false  "Published Year"
// @Param        category_id   query    int     false  "Category ID"
// @Param        available     query    bool    false  "Available"
// @Param        limit         query    int     false  "Limit"  default(10)
// @Param        offset        query    int     false  "Offset" default(0)
// @Success      200           {object} PaginatedResponse{data=[]domain.Book}
// @Failure      400           {object} domain.ErrorResponse
// @Failure      500           {object} domain.ErrorResponse
// @Router       /books/search [get]
func (h *BookHandler) Search(c *gin.Context) {
	var req BookSearchRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		h.logger.Error("Invalid search parameters", zap.Error(err))
		SendError(c, domain.NewInvalidInputError(err.Error()))
		return
	}

	params := domain.BookSearchParams{
		Title:         req.Title,
		Author:        req.Author,
		ISBN:          req.ISBN,
		PublishedYear: req.PublishedYear,
		CategoryID:    req.CategoryID,
		Available:     req.Available,
		Limit:         req.Limit,
		Offset:        req.Offset,
	}

	books, err := h.bookService.Search(params)
	if err != nil {
		h.logger.Error("Failed to search books", zap.Error(err))
		SendError(c, err)
		return
	}

	SendPaginated(c, books, int64(len(books)), req.Limit, req.Offset, "Books retrieved successfully")
}

// Create handles creating a book
// @Summary      Create a new book
// @Description  Create a new book in the catalog
// @Tags         books
// @Accept       json
// @Produce      json
// @Param        book  body      BookRequest  true  "Book object"
// @Success      201   {object}  domain.Book
// @Failure      400   {object}  domain.ErrorResponse
// @Failure      401   {object}  domain.ErrorResponse
// @Failure      403   {object}  domain.ErrorResponse
// @Failure      500   {object}  domain.ErrorResponse
// @Security     Bearer
// @Router       /books [post]
func (h *BookHandler) Create(c *gin.Context) {
	// Only admins and librarians can create books
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

	var req BookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid request body", zap.Error(err))
		SendError(c, domain.NewInvalidInputError(err.Error()))
		return
	}

	book := &domain.Book{
		Title:           req.Title,
		Author:          req.Author,
		ISBN:            req.ISBN,
		Description:     req.Description,
		PublishedYear:   req.PublishedYear,
		Publisher:       req.Publisher,
		TotalCopies:     req.TotalCopies,
		AvailableCopies: req.TotalCopies, // Initially all copies are available
		CategoryID:      req.CategoryID,
	}

	createdBook, err := h.bookService.Create(book)
	if err != nil {
		h.logger.Error("Failed to create book", zap.Error(err))
		SendError(c, err)
		return
	}

	SendCreated(c, createdBook, "Book created successfully")
}

// Update handles updating a book
// @Summary      Update a book
// @Description  Update an existing book's details
// @Tags         books
// @Accept       json
// @Produce      json
// @Param        id    path      int          true  "Book ID"
// @Param        book  body      BookRequest  true  "Updated book object"
// @Success      200   {object}  domain.Book
// @Failure      400   {object}  domain.ErrorResponse
// @Failure      401   {object}  domain.ErrorResponse
// @Failure      403   {object}  domain.ErrorResponse
// @Failure      404   {object}  domain.ErrorResponse
// @Failure      500   {object}  domain.ErrorResponse
// @Security     Bearer
// @Router       /books/{id} [put]
func (h *BookHandler) Update(c *gin.Context) {
	// Only admins and librarians can update books
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
		h.logger.Error("Invalid book ID", zap.Error(err))
		SendError(c, domain.NewInvalidInputError("invalid book ID"))
		return
	}

	var req BookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid request body", zap.Error(err))
		SendError(c, domain.NewInvalidInputError(err.Error()))
		return
	}

	// Get existing book
	existingBook, err := h.bookService.GetByID(id)
	if err != nil {
		h.logger.Error("Failed to get book by ID", zap.Int64("id", id), zap.Error(err))
		SendError(c, err)
		return
	}

	// Update book fields
	existingBook.Title = req.Title
	existingBook.Author = req.Author
	existingBook.ISBN = req.ISBN
	existingBook.Description = req.Description
	existingBook.PublishedYear = req.PublishedYear
	existingBook.Publisher = req.Publisher
	existingBook.CategoryID = req.CategoryID

	updatedBook, err := h.bookService.Update(existingBook)
	if err != nil {
		h.logger.Error("Failed to update book", zap.Int64("id", id), zap.Error(err))
		SendError(c, err)
		return
	}

	SendSuccess(c, updatedBook, "Book updated successfully")
}

// UpdateCopies handles updating book copies
// @Summary      Update book copies
// @Description  Update the number of total and available copies of a book
// @Tags         books
// @Accept       json
// @Produce      json
// @Param        id    path      int                true  "Book ID"
// @Param        copies body     BookCopiesRequest  true  "Updated copies information"
// @Success      200   {object}  domain.Book
// @Failure      400   {object}  domain.ErrorResponse
// @Failure      401   {object}  domain.ErrorResponse
// @Failure      403   {object}  domain.ErrorResponse
// @Failure      404   {object}  domain.ErrorResponse
// @Failure      500   {object}  domain.ErrorResponse
// @Security     Bearer
// @Router       /books/{id}/copies [patch]
func (h *BookHandler) UpdateCopies(c *gin.Context) {
	// Only admins and librarians can update book copies
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
		h.logger.Error("Invalid book ID", zap.Error(err))
		SendError(c, domain.NewInvalidInputError("invalid book ID"))
		return
	}

	var req BookCopiesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid request body", zap.Error(err))
		SendError(c, domain.NewInvalidInputError(err.Error()))
		return
	}

	updatedBook, err := h.bookService.UpdateCopies(id, req.TotalCopies, req.AvailableCopies)
	if err != nil {
		h.logger.Error("Failed to update book copies", zap.Int64("id", id), zap.Error(err))
		SendError(c, err)
		return
	}

	SendSuccess(c, updatedBook, "Book copies updated successfully")
}

// Delete handles deleting a book
// @Summary      Delete a book
// @Description  Delete a book from the catalog
// @Tags         books
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Book ID"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  domain.ErrorResponse
// @Failure      401  {object}  domain.ErrorResponse
// @Failure      403  {object}  domain.ErrorResponse
// @Failure      404  {object}  domain.ErrorResponse
// @Failure      500  {object}  domain.ErrorResponse
// @Security     Bearer
// @Router       /books/{id} [delete]
func (h *BookHandler) Delete(c *gin.Context) {
	// Only admins and librarians can delete books
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
		h.logger.Error("Invalid book ID", zap.Error(err))
		SendError(c, domain.NewInvalidInputError("invalid book ID"))
		return
	}

	err = h.bookService.Delete(id)
	if err != nil {
		h.logger.Error("Failed to delete book", zap.Int64("id", id), zap.Error(err))
		SendError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Book deleted successfully"})
}
