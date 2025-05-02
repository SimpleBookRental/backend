package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/SimpleBookRental/backend/internal/api"
	"github.com/SimpleBookRental/backend/internal/repository"
	"github.com/SimpleBookRental/backend/internal/service"
	"github.com/SimpleBookRental/backend/pkg/auth"
	"github.com/SimpleBookRental/backend/pkg/config"
	"github.com/SimpleBookRental/backend/pkg/logger"
	"github.com/gin-gonic/gin"
	_ "github.com/SimpleBookRental/backend/docs"
)

var (
	testServer     *httptest.Server
	testClient     *http.Client
	testDB         *repository.DBConn
	adminToken     string
	librianToken   string
	memberToken    string
	baseURL        string
	testUserID     string
	testBookID     string
	testCategoryID string
	testRentalID   string
	testPaymentID  string
)

// TestMain sets up the test environment
func TestMain(m *testing.M) {
	// Set up test environment
	setup()
	
	// Run tests
	code := m.Run()
	
	// Tear down test environment
	teardown()
	
	os.Exit(code)
}

// setup initializes the test environment
func setup() {
	gin.SetMode(gin.TestMode)
	
	// Initialize test database
	setupTestDatabase()
	
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
	
	// Initialize JWT service
	jwtService := auth.NewJWTService(&cfg.JWT)
	
	// Initialize repositories
	repos := repository.NewRepository(testDB)
	
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
	
	// Add a ping endpoint for testing
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})
	
	// Register routes
	handlers.RegisterRoutes(router, middleware)
	
	// Create test server
	testServer = httptest.NewServer(router)
	baseURL = testServer.URL
	
	testClient = &http.Client{
		Timeout: time.Second * 10,
	}
	
	// Create test users and get tokens
	createTestUsers()
}

// teardown cleans up the test environment
func teardown() {
	testServer.Close()
	cleanupTestDatabase()
}

// setupTestDatabase initializes a test database
func setupTestDatabase() {
	// In a real implementation, this would connect to a test database
	// This would typically use environment variables specific for testing
	
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	
	// Initialize logger for testing
	appLogger, err := logger.NewLogger(cfg.Logger)
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	
	// Initialize test database
	db, err := repository.NewPostgresDB(cfg.Database)
	if err != nil {
		log.Printf("WARNING: Failed to initialize database: %v", err)
		log.Printf("Some tests may fail if they require database access")
		
		// Create mock DBConn for tests that don't require real DB
		testDB = &repository.DBConn{
			Logger: appLogger,
		}
		return
	}
	
	// Create DBConn
	testDB = &repository.DBConn{
		DB:     db,
		Logger: appLogger,
	}
}

// cleanupTestDatabase cleans up the test database
func cleanupTestDatabase() {
	if testDB != nil && testDB.DB != nil {
		testDB.DB.Close()
	}
}

// createTestUsers creates test users for different roles
func createTestUsers() {
	if testDB == nil || testDB.DB == nil {
		// If database isn't available, use mock tokens
		adminToken = "mock-admin-token"
		librianToken = "mock-librarian-token"
		memberToken = "mock-member-token"
		log.Printf("Using mock tokens for tests")
		return
	}
	
	// Create admin user
	adminToken = createUserAndGetToken("admin@example.com", "Admin123!", "admin")
	
	// Create librarian user
	librianToken = createUserAndGetToken("librarian@example.com", "Librarian123!", "librarian")
	
	// Create member user
	memberToken = createUserAndGetToken("member@example.com", "Member123!", "member")
}

// createUserAndGetToken creates a user and returns the authentication token
func createUserAndGetToken(email, password, role string) string {
	// Register user
	registerURL := fmt.Sprintf("%s/api/v1/auth/register", baseURL)
	registerData := map[string]interface{}{
		"email":     email,
		"password":  password,
		"firstName": "Test",
		"lastName":  "User",
		"role":      role,
	}
	
	registerBody, _ := json.Marshal(registerData)
	registerReq, _ := http.NewRequest("POST", registerURL, bytes.NewBuffer(registerBody))
	registerReq.Header.Set("Content-Type", "application/json")
	
	registerResp, err := testClient.Do(registerReq)
	if err != nil {
		log.Fatalf("Failed to register user: %v", err)
	}
	registerResp.Body.Close()
	
	// Login user
	loginURL := fmt.Sprintf("%s/api/v1/auth/login", baseURL)
	loginData := map[string]string{
		"email":    email,
		"password": password,
	}
	
	loginBody, _ := json.Marshal(loginData)
	loginReq, _ := http.NewRequest("POST", loginURL, bytes.NewBuffer(loginBody))
	loginReq.Header.Set("Content-Type", "application/json")
	
	loginResp, err := testClient.Do(loginReq)
	if err != nil {
		log.Fatalf("Failed to login user: %v", err)
	}
	defer loginResp.Body.Close()
	
	var loginResult map[string]interface{}
	json.NewDecoder(loginResp.Body).Decode(&loginResult)
	
	// Extract token from response
	token, ok := loginResult["token"].(string)
	if !ok {
		log.Fatalf("Failed to extract token from login response")
	}
	
	return token
}

// TestPing tests the ping endpoint to verify the server is running
func TestPing(t *testing.T) {
	pingURL := fmt.Sprintf("%s/ping", baseURL)
	req, _ := http.NewRequest("GET", pingURL, nil)
	
	resp, err := testClient.Do(req)
	if err != nil {
		t.Fatalf("Failed to ping server: %v", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status OK; got %v", resp.Status)
	}
	
	var result map[string]string
	json.NewDecoder(resp.Body).Decode(&result)
	
	if message, exists := result["message"]; !exists || message != "pong" {
		t.Errorf("Expected message to be 'pong'; got %v", message)
	}
}

// Helper function to make authenticated requests
func makeAuthenticatedRequest(method, url string, body interface{}, token string) (*http.Response, error) {
	var bodyReader *bytes.Buffer
	
	if body != nil {
		bodyData, _ := json.Marshal(body)
		bodyReader = bytes.NewBuffer(bodyData)
	} else {
		bodyReader = bytes.NewBuffer(nil)
	}
	
	req, _ := http.NewRequest(method, url, bodyReader)
	req.Header.Set("Content-Type", "application/json")
	
	if token != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	}
	
	return testClient.Do(req)
}

// Helper function to check if status code matches expected
func checkStatusCode(t *testing.T, resp *http.Response, expected int) {
	if resp.StatusCode != expected {
		t.Errorf("Expected status %d; got %d", expected, resp.StatusCode)
	}
}
