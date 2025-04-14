package database

import (
	"fmt"
	"log"
	"time"

	"github.com/SimpleBookRental/backend/internal/config"
	"github.com/SimpleBookRental/backend/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
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
	)

	if err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}

	log.Println("Database migrations completed")
	return nil
}
