package routes

import (
	"github.com/SimpleBookRental/backend/internal/controllers"
	"github.com/gin-gonic/gin"
)

// SetupRoutes sets up all the routes for the application
func SetupRoutes(
	router *gin.Engine,
	userController *controllers.UserController,
	bookController *controllers.BookController,
) {
	// API v1 group
	v1 := router.Group("/api/v1")
	{
		// User routes
		users := v1.Group("/users")
		{
			users.POST("", userController.Create)
			users.GET("", userController.GetAll)
			users.GET("/:id", userController.GetByID)
			users.PUT("/:id", userController.Update)
			users.DELETE("/:id", userController.Delete)
		}

		// Book routes
		books := v1.Group("/books")
		{
			books.POST("", bookController.Create)
			books.GET("", bookController.GetAll)
			books.GET("/:id", bookController.GetByID)
			books.PUT("/:id", bookController.Update)
			books.DELETE("/:id", bookController.Delete)
		}

		// User's books routes
		v1.GET("/user-books/:user_id", bookController.GetByUserID)
	}
}
