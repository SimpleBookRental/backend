package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/SimpleBookRental/backend/internal/api"
	"github.com/SimpleBookRental/backend/internal/repository"
	"github.com/SimpleBookRental/backend/internal/service"
	"github.com/SimpleBookRental/backend/pkg/auth"
	"github.com/SimpleBookRental/backend/pkg/config"
	"github.com/SimpleBookRental/backend/pkg/logger"
	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize logger
	appLogger, err := logger.NewLogger(cfg.Logger)
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer appLogger.Sync()

	// Set Gin mode
	if cfg.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize database
	db, err := repository.NewPostgresDB(cfg.Database)
	if err != nil {
		appLogger.Fatal("Failed to initialize database", err)
	}
	defer db.Close()

	// Create DBConn
	dbConn := &repository.DBConn{
		DB:     db,
		Logger: appLogger,
	}

	// Initialize repositories
	repos := repository.NewRepository(dbConn)

	// Initialize JWT service
	jwtService := auth.NewJWTService(&cfg.JWT)

	// Initialize services
	services := service.NewService(repos, cfg, jwtService, appLogger)

	// Initialize handlers
	handlers := api.NewHandler(services, cfg, jwtService, appLogger)

	// Initialize middleware
	middleware := api.NewMiddleware(jwtService, appLogger)

	// Initialize router
	router := gin.New()

	// Apply middleware
	router.Use(middleware.LoggerMiddleware())
	router.Use(middleware.RecoveryMiddleware())
	router.Use(middleware.CORSMiddleware())

	// Register routes
	handlers.RegisterRoutes(router, middleware)

	// Create HTTP server
	server := &http.Server{
		Addr:           fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:        router,
		ReadTimeout:    cfg.Server.ReadTimeout,
		WriteTimeout:   cfg.Server.WriteTimeout,
		MaxHeaderBytes: 1 << 20, // 1 MB
	}

	// Start server in a goroutine
	go func() {
		appLogger.Info(fmt.Sprintf("Starting server on port %d", cfg.Server.Port))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			appLogger.Fatal("Failed to start server", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	appLogger.Info("Shutting down server...")

	// Create a deadline to wait for
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Shutdown server
	if err := server.Shutdown(ctx); err != nil {
		appLogger.Fatal("Server forced to shutdown", err)
	}

	appLogger.Info("Server exiting")
}
