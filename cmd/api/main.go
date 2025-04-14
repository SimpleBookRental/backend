package main

import (
	"fmt"
	"log"

	"github.com/SimpleBookRental/backend/internal/config"
	"github.com/SimpleBookRental/backend/internal/controllers"
	"github.com/SimpleBookRental/backend/internal/repositories"
	"github.com/SimpleBookRental/backend/internal/routes"
	"github.com/SimpleBookRental/backend/internal/services"
	"github.com/SimpleBookRental/backend/pkg/database"
	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Connect to database
	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Run migrations
	err = database.Migrate(db)
	if err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Initialize repositories
	userRepo := repositories.NewUserRepository(db)
	bookRepo := repositories.NewBookRepository(db)
	tokenRepo := repositories.NewTokenRepository(db)

	// Initialize services
	userService := services.NewUserService(userRepo)
	bookService := services.NewBookService(bookRepo, userRepo)
	tokenService := services.NewTokenService(tokenRepo, userRepo)

	// Initialize controllers
	userController := controllers.NewUserController(userService, tokenRepo)
	bookController := controllers.NewBookController(bookService)
	tokenController := controllers.NewTokenController(tokenService)

	// Initialize router
	router := gin.Default()

	// Setup routes
	routes.SetupRoutes(router, userController, bookController, tokenController, tokenRepo)

	// Start server
	serverAddr := fmt.Sprintf(":%s", cfg.Server.Port)
	log.Printf("Server starting on %s", serverAddr)
	if err := router.Run(serverAddr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
