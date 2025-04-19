// @title           Simple Book Rental API
// @version         1.0.0
// @description     RESTful API for book rental system (Go, Gin, GORM, Clean Architecture)
// @BasePath        /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/SimpleBookRental/backend/internal/config"
	"github.com/SimpleBookRental/backend/internal/controllers"
	"github.com/SimpleBookRental/backend/internal/models"
	"github.com/SimpleBookRental/backend/internal/repositories"
	"github.com/SimpleBookRental/backend/internal/routes"
	"github.com/SimpleBookRental/backend/internal/services"
	"github.com/SimpleBookRental/backend/pkg/database"
	"github.com/SimpleBookRental/backend/pkg/utils"
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
	txManager := repositories.NewTransactionManager(db)

	// Seed default admin user
	existingAdmin, err := userRepo.FindByEmail(cfg.Admin.Email)
	if err != nil {
		log.Fatalf("Failed to check default admin: %v", err)
	}
	if existingAdmin == nil {
		hashedPassword, err := utils.HashPassword(cfg.Admin.Password)
		if err != nil {
			log.Fatalf("Failed to hash default admin password: %v", err)
		}
		adminUser := &models.User{
			Name:     cfg.Admin.Name,
			Email:    cfg.Admin.Email,
			Password: hashedPassword,
			Role:     models.AdminRole,
		}
		if err := userRepo.Create(adminUser); err != nil {
			log.Fatalf("Failed to create default admin: %v", err)
		}
		log.Println("Default admin user created")
	}

	// Initialize services
	userService := services.NewUserService(userRepo, bookRepo, tokenRepo)
	bookService := services.NewBookService(bookRepo, userRepo)
	tokenService := services.NewTokenService(tokenRepo, userRepo)
	bookUserService := services.NewBookUserService(txManager, bookRepo, userRepo)

	// Initialize controllers
	userController := controllers.NewUserController(userService, tokenRepo)
	bookController := controllers.NewBookController(bookService)
	tokenController := controllers.NewTokenController(tokenService)
	bookUserController := controllers.NewBookUserController(bookUserService, userRepo)

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
	routes.SetupRoutes(router, userController, bookController, tokenController,
		bookUserController, tokenRepo, userRepo)

	// Start server
	serverAddr := fmt.Sprintf(":%s", cfg.Server.Port)
	log.Printf("Server starting on %s", serverAddr)
	if err := router.Run(serverAddr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
