package integration

import (
	"net/http"
	"testing"
)

// TestAuthRegister tests the registration endpoint
func TestAuthRegister(t *testing.T) {
	// Test successful registration
	registerURL := baseURL + "/api/v1/auth/register"
	registerData := map[string]interface{}{
		"email":     "test_register@example.com",
		"password":  "Password123!",
		"firstName": "TestFirst",
		"lastName":  "TestLast",
		"role":      "member",
	}
	
	resp, err := makeAuthenticatedRequest("POST", registerURL, registerData, "")
	if err != nil {
		t.Fatalf("Failed to make register request: %v", err)
	}
	defer resp.Body.Close()
	
	checkStatusCode(t, resp, http.StatusCreated)
	
	// Test duplicate email
	resp, err = makeAuthenticatedRequest("POST", registerURL, registerData, "")
	if err != nil {
		t.Fatalf("Failed to make duplicate register request: %v", err)
	}
	defer resp.Body.Close()
	
	checkStatusCode(t, resp, http.StatusConflict)
	
	// Test invalid email format
	invalidEmailData := map[string]interface{}{
		"email":     "invalid-email",
		"password":  "Password123!",
		"firstName": "TestFirst",
		"lastName":  "TestLast",
		"role":      "member",
	}
	
	resp, err = makeAuthenticatedRequest("POST", registerURL, invalidEmailData, "")
	if err != nil {
		t.Fatalf("Failed to make invalid email register request: %v", err)
	}
	defer resp.Body.Close()
	
	checkStatusCode(t, resp, http.StatusBadRequest)
	
	// Test weak password
	weakPasswordData := map[string]interface{}{
		"email":     "test_weak@example.com",
		"password":  "weak",
		"firstName": "TestFirst",
		"lastName":  "TestLast",
		"role":      "member",
	}
	
	resp, err = makeAuthenticatedRequest("POST", registerURL, weakPasswordData, "")
	if err != nil {
		t.Fatalf("Failed to make weak password register request: %v", err)
	}
	defer resp.Body.Close()
	
	checkStatusCode(t, resp, http.StatusBadRequest)
}

// TestAuthLogin tests the login endpoint
func TestAuthLogin(t *testing.T) {
	// Test successful login
	loginURL := baseURL + "/api/v1/auth/login"
	loginData := map[string]interface{}{
		"email":    "member@example.com",
		"password": "Member123!",
	}
	
	resp, err := makeAuthenticatedRequest("POST", loginURL, loginData, "")
	if err != nil {
		t.Fatalf("Failed to make login request: %v", err)
	}
	defer resp.Body.Close()
	
	checkStatusCode(t, resp, http.StatusOK)
	
	// Test invalid credentials
	invalidLoginData := map[string]interface{}{
		"email":    "member@example.com",
		"password": "WrongPassword123!",
	}
	
	resp, err = makeAuthenticatedRequest("POST", loginURL, invalidLoginData, "")
	if err != nil {
		t.Fatalf("Failed to make invalid login request: %v", err)
	}
	defer resp.Body.Close()
	
	checkStatusCode(t, resp, http.StatusUnauthorized)
	
	// Test non-existent user
	nonExistentData := map[string]interface{}{
		"email":    "nonexistent@example.com",
		"password": "Password123!",
	}
	
	resp, err = makeAuthenticatedRequest("POST", loginURL, nonExistentData, "")
	if err != nil {
		t.Fatalf("Failed to make non-existent user login request: %v", err)
	}
	defer resp.Body.Close()
	
	checkStatusCode(t, resp, http.StatusUnauthorized)
}

// TestAuthLogout tests the logout endpoint
func TestAuthLogout(t *testing.T) {
	// Test successful logout
	logoutURL := baseURL + "/api/v1/auth/logout"
	
	resp, err := makeAuthenticatedRequest("POST", logoutURL, nil, memberToken)
	if err != nil {
		t.Fatalf("Failed to make logout request: %v", err)
	}
	defer resp.Body.Close()
	
	checkStatusCode(t, resp, http.StatusOK)
	
	// Test unauthenticated logout
	resp, err = makeAuthenticatedRequest("POST", logoutURL, nil, "")
	if err != nil {
		t.Fatalf("Failed to make unauthenticated logout request: %v", err)
	}
	defer resp.Body.Close()
	
	checkStatusCode(t, resp, http.StatusUnauthorized)
}

// TestAuthRefresh tests the token refresh endpoint
func TestAuthRefresh(t *testing.T) {
	// Test successful token refresh
	refreshURL := baseURL + "/api/v1/auth/refresh"
	
	resp, err := makeAuthenticatedRequest("POST", refreshURL, nil, memberToken)
	if err != nil {
		t.Fatalf("Failed to make refresh request: %v", err)
	}
	defer resp.Body.Close()
	
	checkStatusCode(t, resp, http.StatusOK)
	
	// Test unauthenticated refresh
	resp, err = makeAuthenticatedRequest("POST", refreshURL, nil, "")
	if err != nil {
		t.Fatalf("Failed to make unauthenticated refresh request: %v", err)
	}
	defer resp.Body.Close()
	
	checkStatusCode(t, resp, http.StatusUnauthorized)
	
	// Test with invalid token
	resp, err = makeAuthenticatedRequest("POST", refreshURL, nil, "invalid-token")
	if err != nil {
		t.Fatalf("Failed to make invalid token refresh request: %v", err)
	}
	defer resp.Body.Close()
	
	checkStatusCode(t, resp, http.StatusUnauthorized)
}

// TestAuthMe tests the get current user endpoint
func TestAuthMe(t *testing.T) {
	// Test successful get current user
	meURL := baseURL + "/api/v1/auth/me"
	
	resp, err := makeAuthenticatedRequest("GET", meURL, nil, memberToken)
	if err != nil {
		t.Fatalf("Failed to make me request: %v", err)
	}
	defer resp.Body.Close()
	
	checkStatusCode(t, resp, http.StatusOK)
	
	// Test unauthenticated get current user
	resp, err = makeAuthenticatedRequest("GET", meURL, nil, "")
	if err != nil {
		t.Fatalf("Failed to make unauthenticated me request: %v", err)
	}
	defer resp.Body.Close()
	
	checkStatusCode(t, resp, http.StatusUnauthorized)
}
