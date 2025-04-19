package database

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/SimpleBookRental/backend/internal/config"
	"github.com/SimpleBookRental/backend/internal/models"
	"github.com/SimpleBookRental/backend/pkg/utils"
)

// DB is the database connection
var DB *gorm.DB

// Connect connects to the database with retry mechanism
func Connect(config *config.Config) (*gorm.DB, error) {
	var err error
	var db *gorm.DB

	dsn := config.Database.GetDSN()
	log.Printf("Attempting to connect to database with DSN: %s\n", dsn)

	// Retry parameters
	maxRetries := 5
	retryDelay := time.Second * 5

	// Retry loop
	for i := 0; i < maxRetries; i++ {
		log.Printf("Database connection attempt %d of %d\n", i+1, maxRetries)

		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		})

		if err == nil {
			log.Println("Successfully connected to database")
			db = db.Debug()
			DB = db
			return db, nil
		}

		log.Printf("Failed to connect to database: %v. Retrying in %v...\n", err, retryDelay)
		time.Sleep(retryDelay)
	}

	return nil, fmt.Errorf("failed to connect to database after %d attempts: %w", maxRetries, err)
}

// Migrate runs database migrations
func Migrate(db *gorm.DB) error {
	log.Println("Running database migrations")

	err := db.AutoMigrate(
		&models.User{},
		&models.Book{},
		&models.IssuedToken{},
	)

	if err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}

	// Create default admin user if it doesn't exist
	err = createDefaultAdminUser(db)
	if err != nil {
		return fmt.Errorf("failed to create default admin user: %w", err)
	}

	log.Println("Database migrations completed")
	return nil
}

// createDefaultAdminUser creates a default admin user if it doesn't exist
func createDefaultAdminUser(db *gorm.DB) error {
	// Check if admin user already exists
	var count int64
	db.Model(&models.User{}).Where("email = ?", "admin@system.com").Count(&count)

	if count > 0 {
		log.Println("Default admin user already exists")
		return nil
	}

	// Hash password
	hashedPassword, err := utils.HashPassword("admin")
	if err != nil {
		return fmt.Errorf("error hashing password: %w", err)
	}

	// Create admin user
	adminUser := &models.User{
		Name:     "Admin",
		Email:    "admin@system.com",
		Password: hashedPassword,
		Role:     models.AdminRole,
	}

	result := db.Create(adminUser)
	if result.Error != nil {
		return fmt.Errorf("error creating admin user: %w", result.Error)
	}

	log.Println("Default admin user created successfully")
	return nil
}
