package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/SimpleBookRental/backend/internal/controllers"
	"github.com/SimpleBookRental/backend/internal/middleware"
	"github.com/SimpleBookRental/backend/internal/repositories"
)

// SetupRoutes sets up all the routes for the application
func SetupRoutes(
	router *gin.Engine,
	userController *controllers.UserController,
	bookController *controllers.BookController,
	tokenController *controllers.TokenController,
	tokenRepo *repositories.TokenRepository,
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
		auth.Use(middleware.AuthMiddleware(tokenRepo))
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
				books.POST("", bookController.Create)
				books.GET("", bookController.GetAll)
				books.GET("/:id", bookController.GetByID)
				books.PUT("/:id", bookController.Update)
				books.DELETE("/:id", bookController.Delete)
			}

			// User's books routes
			auth.GET("/user-books/:user_id", bookController.GetByUserID)
		}
	}
}
