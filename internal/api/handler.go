package api

import (
	"github.com/SimpleBookRental/backend/internal/domain"
	"github.com/SimpleBookRental/backend/internal/service"
	"github.com/SimpleBookRental/backend/pkg/auth"
	"github.com/SimpleBookRental/backend/pkg/config"
	"github.com/SimpleBookRental/backend/pkg/logger"
	"github.com/gin-gonic/gin"
)

// Handler is a factory for all API handlers
type Handler struct {
	AuthHandler     *AuthHandler
	UserHandler     *UserHandler
	CategoryHandler *CategoryHandler
	BookHandler     *BookHandler
	RentalHandler   *RentalHandler
	PaymentHandler  *PaymentHandler
	ReportHandler   *ReportHandler
	Logger          *logger.Logger
}

// NewHandler creates a new handler factory
func NewHandler(services *service.Service, cfg *config.Config, jwtService *auth.JWTService, logger *logger.Logger) *Handler {
	handlerLogger := logger.Named("handler")

	return &Handler{
		AuthHandler:     NewAuthHandler(services.Auth, jwtService, handlerLogger.Named("auth")),
		UserHandler:     NewUserHandler(services.User, jwtService, handlerLogger.Named("user")),
		CategoryHandler: NewCategoryHandler(services.Category, jwtService, handlerLogger.Named("category")),
		BookHandler:     NewBookHandler(services.Book, jwtService, handlerLogger.Named("book")),
		RentalHandler:   NewRentalHandler(services.Rental, jwtService, handlerLogger.Named("rental")),
		PaymentHandler:  NewPaymentHandler(services.Payment, jwtService, handlerLogger.Named("payment")),
		ReportHandler:   NewReportHandler(services.Report, jwtService, handlerLogger.Named("report")),
		Logger:          handlerLogger,
	}
}

// RegisterRoutes registers all API routes
func (h *Handler) RegisterRoutes(router *gin.Engine, middleware *Middleware) {
	// API v1 group
	v1 := router.Group("/api/v1")
	{
		// Auth routes - public endpoints
		auth := v1.Group("/auth")
		{
			auth.POST("/register", h.AuthHandler.Register)
			auth.POST("/login", h.AuthHandler.Login)
			auth.POST("/refresh", h.AuthHandler.RefreshToken)
			auth.POST("/logout", middleware.AuthMiddleware(), h.AuthHandler.Logout)
		}

		// User routes - protected endpoints
		users := v1.Group("/users")
		users.Use(middleware.AuthMiddleware())
		{
			users.GET("", middleware.RoleMiddleware(domain.RoleAdmin), h.UserHandler.List)
			users.GET("/:id", h.UserHandler.GetByID) // Handler checks if user is requesting their own profile or is admin
			users.PUT("/:id", h.UserHandler.Update)  // Handler checks if user is updating their own profile or is admin
			users.DELETE("/:id", middleware.RoleMiddleware(domain.RoleAdmin), h.UserHandler.Delete)
		}

		// Category routes
		categories := v1.Group("/categories")
		{
			// Public endpoints for browsing categories
			categories.GET("", h.CategoryHandler.List)
			categories.GET("/all", h.CategoryHandler.ListAll)
			categories.GET("/:id", h.CategoryHandler.GetByID)
			
			// Protected endpoints for managing categories
			categoriesProtected := categories.Group("")
			categoriesProtected.Use(middleware.AuthMiddleware(), middleware.RoleMiddleware(domain.RoleLibrarian))
			{
				categoriesProtected.POST("", h.CategoryHandler.Create)
				categoriesProtected.PUT("/:id", h.CategoryHandler.Update)
				categoriesProtected.DELETE("/:id", h.CategoryHandler.Delete)
			}
		}

		// Book routes
		books := v1.Group("/books")
		{
			// Public endpoints for browsing books
			books.GET("", h.BookHandler.List)
			books.GET("/search", h.BookHandler.Search)
			books.GET("/category/:id", h.BookHandler.ListByCategory)
			books.GET("/:id", h.BookHandler.GetByID)
			
			// Protected endpoints for managing books
			booksProtected := books.Group("")
			booksProtected.Use(middleware.AuthMiddleware(), middleware.RoleMiddleware(domain.RoleLibrarian))
			{
				booksProtected.POST("", h.BookHandler.Create)
				booksProtected.PUT("/:id", h.BookHandler.Update)
				booksProtected.PUT("/:id/copies", h.BookHandler.UpdateCopies)
				booksProtected.DELETE("/:id", h.BookHandler.Delete)
			}
		}

		// Rental routes - all require authentication
		rentals := v1.Group("/rentals")
		rentals.Use(middleware.AuthMiddleware())
		{
			// Admin/Librarian endpoints
			rentals.GET("", middleware.RoleMiddleware(domain.RoleLibrarian), h.RentalHandler.List)
			
			// Member endpoints (handlers check if user is requesting their own rentals or is admin/librarian)
			rentals.GET("/user/:userId", h.RentalHandler.ListByUser)
			rentals.GET("/:id", h.RentalHandler.GetByID)
			rentals.POST("", h.RentalHandler.Create)
			rentals.PUT("/:id/return", h.RentalHandler.Return)
			rentals.PUT("/:id/extend", h.RentalHandler.Extend)
		}

		// Payment routes - all require authentication
		payments := v1.Group("/payments")
		payments.Use(middleware.AuthMiddleware())
		{
			// Admin/Librarian endpoints
			payments.GET("", middleware.RoleMiddleware(domain.RoleAdmin), h.PaymentHandler.List)
			
			// Member endpoints (handlers check if user is requesting their own payments or is admin/librarian)
			payments.GET("/user/:userId", h.PaymentHandler.ListByUser)
			payments.GET("/:id", h.PaymentHandler.GetByID)
			payments.POST("", h.PaymentHandler.Create)
			payments.POST("/process", h.PaymentHandler.Process)
			
			// Admin/Librarian endpoints
			payments.PUT("/:id/refund", middleware.RoleMiddleware(domain.RoleLibrarian), h.PaymentHandler.Refund)
		}

		// Report routes - all require authentication and appropriate roles
		reports := v1.Group("/reports")
		reports.Use(middleware.AuthMiddleware())
		{
			reports.GET("/books/popular", middleware.RoleMiddleware(domain.RoleLibrarian), h.ReportHandler.GetPopularBooks)
			reports.GET("/revenue", middleware.RoleMiddleware(domain.RoleAdmin), h.ReportHandler.GetRevenueReport)
			reports.GET("/overdue", middleware.RoleMiddleware(domain.RoleLibrarian), h.ReportHandler.GetOverdueBooks)
		}
	}
}
