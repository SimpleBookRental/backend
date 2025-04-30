package repository

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	"github.com/SimpleBookRental/backend/pkg/config"
	"github.com/SimpleBookRental/backend/pkg/logger"
	_ "github.com/lib/pq" // PostgreSQL driver
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file" // file source for migrations
)

// DBConn represents a database connection
type DBConn struct {
	DB     *sql.DB
	Logger *logger.Logger
}

// NewDBConn creates a new database connection
func NewDBConn(cfg *config.DatabaseConfig, logger *logger.Logger) (*DBConn, error) {
	dsn := cfg.GetDSN()
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Set connection pool settings
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)

	logger.Info("Connected to database")

	return &DBConn{
		DB:     db,
		Logger: logger,
	}, nil
}

// NewPostgresDB creates a new PostgreSQL database connection
func NewPostgresDB(cfg config.DatabaseConfig) (*sql.DB, error) {
	dsn := cfg.GetDSN()
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Set connection pool settings
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)

	return db, nil
}

// RunMigrations runs database migrations
func RunMigrations(cfg config.DatabaseConfig) error {
	db, err := sql.Open("postgres", cfg.GetDSN())
	if err != nil {
		return fmt.Errorf("failed to open database connection: %w", err)
	}
	defer db.Close()

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("failed to create migration driver: %w", err)
	}

	// Get the current working directory
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current working directory: %w", err)
	}

	// Construct the migrations path
	migrationsPath := filepath.Join(cwd, "migrations")
	migrationsURL := fmt.Sprintf("file://%s", migrationsPath)

	m, err := migrate.NewWithDatabaseInstance(migrationsURL, "postgres", driver)
	if err != nil {
		return fmt.Errorf("failed to create migration instance: %w", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	return nil
}

// Close closes the database connection
func (c *DBConn) Close() error {
	if c.DB != nil {
		c.Logger.Info("Closing database connection")
		return c.DB.Close()
	}
	return nil
}
