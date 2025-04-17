package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/SimpleBookRental/backend/internal/cache"
	"github.com/SimpleBookRental/backend/internal/controllers"
	"github.com/SimpleBookRental/backend/internal/middleware"
	"github.com/SimpleBookRental/backend/internal/repositories"
)

func SetupRoutes(
	router *gin.Engine,
	userController *controllers.UserController,
	bookController *controllers.BookController,
	tokenController *controllers.TokenController,
	bookUserController *controllers.BookUserController,
	tokenRepo *repositories.TokenRepository,
	userRepo *repositories.UserRepository,
	redisCache cache.Cache,
	cacheTTL int,
) {
	// API v1 group
	v1 := router.Group("/api/v1")
	{
		// Public routes (no authentication required)
		// User creation
		v1.POST("/users", userController.Create)

		// Auth routes
		v1.POST("/login", userController.Login)
		v1.POST("/refresh-token", tokenController.RefreshToken)
		v1.POST("/logout", tokenController.Logout)

		// Protected routes (authentication required)
		// Apply auth middleware to all protected routes
		auth := v1.Group("/")
		auth.Use(middleware.AuthMiddleware(tokenRepo, userRepo))
		{
			// User routes
			users := auth.Group("users")
			{
				users.GET("", userController.GetAll)
				users.GET("/:id", userController.GetByID)
				users.PUT("/:id", userController.Update)
				users.DELETE("/:id", userController.Delete)
			}

			// Book routes
			books := auth.Group("books")
			{
				// Apply cache middleware to GET requests
				books.Use(middleware.CacheMiddleware(redisCache, cacheTTL))

				// All authenticated users can create books
				books.POST("", bookController.Create)

				// All authenticated users can get all books (filtered by role in controller)
				books.GET("", bookController.GetAll)

				// All authenticated users can get a book by ID (filtered by role in controller)
				books.GET("/:id", bookController.GetByID)

				// All authenticated users can update a book (filtered by role in controller)
				books.PUT("/:id", bookController.Update)

				// All authenticated users can delete a book (filtered by role in controller)
				books.DELETE("/:id", bookController.Delete)

				// Book-User operations
				books.POST("/:id/transfer", bookUserController.TransferBookOwnership)
			}

			// Book-User routes
			bookUsers := auth.Group("book-users")
			{
				// Apply cache middleware to GET requests if any (currently only POST)
				bookUsers.Use(middleware.CacheMiddleware(redisCache, cacheTTL))

				// Create a book with user
				bookUsers.POST("", bookUserController.CreateBookWithUser)
			}

		}
	}
}
