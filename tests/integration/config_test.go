package integration

import (
	"os"
	"testing"

	"github.com/SimpleBookRental/backend/pkg/config"
)

// TestConfig verifies that configuration is loaded correctly
func TestConfig(t *testing.T) {
	// Force loading config
	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}
	
	// Verify database config
	if cfg.Database.Host != "127.0.0.1" {
		t.Errorf("Expected DB host to be 127.0.0.1, got %s", cfg.Database.Host)
	}
	
	if cfg.Database.Name != "book_rental_test" {
		t.Errorf("Expected DB name to be book_rental_test, got %s", cfg.Database.Name)
	}
	
	// Verify JWT config
	if cfg.JWT.Secret != "test_jwt_secret_key_for_integration_tests" {
		t.Errorf("Expected JWT secret to be test_jwt_secret_key_for_integration_tests, got %s", cfg.JWT.Secret)
	}
	
	// Verify environment is test
	if cfg.Server.Env != "test" {
		t.Errorf("Expected environment to be test, got %s", cfg.Server.Env)
	}
}

// TestConfigEnvOverrides tests that environment variables override config file
func TestConfigEnvOverrides(t *testing.T) {
	// Set environment variable
	os.Setenv("DB_HOST", "env_override_host")
	defer os.Unsetenv("DB_HOST")
	
	// Load config with environment override
	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}
	
	// Verify environment variable override
	if cfg.Database.Host != "env_override_host" {
		t.Errorf("Expected DB host to be env_override_host, got %s", cfg.Database.Host)
	}
}
