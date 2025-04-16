package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/SimpleBookRental/backend/internal/config"
	"github.com/SimpleBookRental/backend/internal/controllers"
	"github.com/SimpleBookRental/backend/internal/repositories"
	"github.com/SimpleBookRental/backend/internal/routes"
	"github.com/SimpleBookRental/backend/internal/services"
	"github.com/SimpleBookRental/backend/pkg/database"
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

	// Configure CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Setup routes
	routes.SetupRoutes(router, userController, bookController, tokenController, tokenRepo, userRepo)

	// Start server
	serverAddr := fmt.Sprintf(":%s", cfg.Server.Port)
	log.Printf("Server starting on %s", serverAddr)
	if err := router.Run(serverAddr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
