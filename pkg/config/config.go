package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

// Config holds all configuration for the application
type Config struct {
	Server    ServerConfig
	Database  DatabaseConfig
	JWT       JWTConfig
	Logger    LoggingConfig
	Rental    RentalConfig
	RateLimit RateLimitConfig
}

// ServerConfig holds server configuration
type ServerConfig struct {
	Host         string
	Port         int
	Env          string
	Mode         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Host          string
	Port          int
	User          string
	Password      string
	Name          string
	SSLMode       string
	RunMigrations bool
}

// JWTConfig holds JWT configuration
type JWTConfig struct {
	Secret            string
	ExpirationHours   time.Duration
	RefreshExpiration time.Duration
}

// LoggingConfig holds logging configuration
type LoggingConfig struct {
	Level  string
	Format string
}

// RentalConfig holds rental configuration
type RentalConfig struct {
	DefaultRentalDays      int
	MaxRentalExtensionDays int
	LateFeePerDay          float64
}

// RateLimitConfig holds rate limiting configuration
type RateLimitConfig struct {
	Requests int
	Duration time.Duration
}

// Load loads configuration from environment variables and .env file
func Load() (*Config, error) {
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	
	// Set defaults
	setDefaults()
	
	// Map environment variables (important to do this BEFORE reading config file)
	viper.SetEnvPrefix("")
	viper.AutomaticEnv()
	
	// Explicitly bind environment variables
	viper.BindEnv("DB_HOST")
	viper.BindEnv("DB_PORT")
	viper.BindEnv("DB_USER")
	viper.BindEnv("DB_PASSWORD")
	viper.BindEnv("DB_NAME")
	viper.BindEnv("DB_SSL_MODE")
	viper.BindEnv("SERVER_HOST")
	viper.BindEnv("SERVER_PORT")
	
	// Read config file
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
		// Config file not found, using defaults and env vars
		fmt.Println("Config file not found, using environment variables and defaults")
	}
	
	// Debug - print actual values being used
	fmt.Printf("DB_HOST: %s\n", viper.GetString("DB_HOST"))
	fmt.Printf("DB_PORT: %d\n", viper.GetInt("DB_PORT"))
	fmt.Printf("SERVER_HOST: %s\n", viper.GetString("SERVER_HOST"))
	fmt.Printf("SERVER_PORT: %d\n", viper.GetInt("SERVER_PORT"))

	// Parse config
	config := &Config{
		Server: ServerConfig{
			Host:         viper.GetString("SERVER_HOST"),
			Port:         viper.GetInt("SERVER_PORT"),
			Env:          viper.GetString("ENV"),
			Mode:         viper.GetString("SERVER_MODE"),
			ReadTimeout:  viper.GetDuration("SERVER_READ_TIMEOUT"),
			WriteTimeout: viper.GetDuration("SERVER_WRITE_TIMEOUT"),
		},
		Database: DatabaseConfig{
			Host:          viper.GetString("DB_HOST"),
			Port:          viper.GetInt("DB_PORT"),
			User:          viper.GetString("DB_USER"),
			Password:      viper.GetString("DB_PASSWORD"),
			Name:          viper.GetString("DB_NAME"),
			SSLMode:       viper.GetString("DB_SSL_MODE"),
			RunMigrations: viper.GetBool("DB_RUN_MIGRATIONS"),
		},
		JWT: JWTConfig{
			Secret:            viper.GetString("JWT_SECRET"),
			ExpirationHours:   viper.GetDuration("JWT_EXPIRATION"),
			RefreshExpiration: viper.GetDuration("JWT_REFRESH_EXPIRATION"),
		},
		Logger: LoggingConfig{
			Level:  viper.GetString("LOG_LEVEL"),
			Format: viper.GetString("LOG_FORMAT"),
		},
		Rental: RentalConfig{
			DefaultRentalDays:      viper.GetInt("DEFAULT_RENTAL_DAYS"),
			MaxRentalExtensionDays: viper.GetInt("MAX_RENTAL_EXTENSION_DAYS"),
			LateFeePerDay:          viper.GetFloat64("LATE_FEE_PER_DAY"),
		},
		RateLimit: RateLimitConfig{
			Requests: viper.GetInt("RATE_LIMIT_REQUESTS"),
			Duration: viper.GetDuration("RATE_LIMIT_DURATION"),
		},
	}

	return config, nil
}

// setDefaults sets default values for configuration
func setDefaults() {
	// Server defaults
	viper.SetDefault("SERVER_HOST", "0.0.0.0")
	viper.SetDefault("SERVER_PORT", 3000)
	viper.SetDefault("ENV", "development")
	viper.SetDefault("SERVER_MODE", "debug")
	viper.SetDefault("SERVER_READ_TIMEOUT", "10s")
	viper.SetDefault("SERVER_WRITE_TIMEOUT", "10s")

	// Database defaults
	viper.SetDefault("DB_HOST", "postgres")
	viper.SetDefault("DB_PORT", 5432)
	viper.SetDefault("DB_USER", "postgres")
	viper.SetDefault("DB_PASSWORD", "postgres")
	viper.SetDefault("DB_NAME", "book_rental")
	viper.SetDefault("DB_SSL_MODE", "disable")
	viper.SetDefault("DB_RUN_MIGRATIONS", true)

	// JWT defaults
	viper.SetDefault("JWT_SECRET", "your_jwt_secret_key_here")
	viper.SetDefault("JWT_EXPIRATION", "24h")
	viper.SetDefault("JWT_REFRESH_EXPIRATION", "168h")

	// Logging defaults
	viper.SetDefault("LOG_LEVEL", "debug")
	viper.SetDefault("LOG_FORMAT", "json")

	// Rental defaults
	viper.SetDefault("DEFAULT_RENTAL_DAYS", 14)
	viper.SetDefault("MAX_RENTAL_EXTENSION_DAYS", 7)
	viper.SetDefault("LATE_FEE_PER_DAY", 1.00)

	// Rate limiting defaults
	viper.SetDefault("RATE_LIMIT_REQUESTS", 100)
	viper.SetDefault("RATE_LIMIT_DURATION", "1m")
}

// GetDSN returns the database connection string
func (c *DatabaseConfig) GetDSN() string {
	// Use URL format which is more reliable
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		c.User, c.Password, c.Host, c.Port, c.Name, c.SSLMode)
	fmt.Printf("DEBUG: Using DSN: %s\n", dsn)
	return dsn
}

// GetServerAddress returns the server address
func (c *ServerConfig) GetServerAddress() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

// IsDevelopment returns true if the environment is development
func (c *ServerConfig) IsDevelopment() bool {
	return c.Env == "development"
}

// IsProduction returns true if the environment is production
func (c *ServerConfig) IsProduction() bool {
	return c.Env == "production"
}
