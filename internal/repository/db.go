package repository

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq" // PostgreSQL driver

	"github.com/SimpleBookRental/backend/pkg/config"
	"github.com/SimpleBookRental/backend/pkg/logger"
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

// Close closes the database connection
func (c *DBConn) Close() error {
	if c.DB != nil {
		c.Logger.Info("Closing database connection")
		return c.DB.Close()
	}
	return nil
}
