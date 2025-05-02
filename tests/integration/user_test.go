package integration

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
)

// TestUserCRUD tests the user management API endpoints
func TestUserCRUD(t *testing.T) {
	// Test create user (admin only)
	createURL := fmt.Sprintf("%s/api/v1/users", baseURL)
	userData := map[string]interface{}{
		"email":     "test.user@example.com",
		"password":  "TestPassword123!",
		"firstName": "Test",
		"lastName":  "User",
		"role":      "member",
	}
	
	// Create with admin token
	resp, err := makeAuthenticatedRequest("POST", createURL, userData, adminToken)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}
	defer resp.Body.Close()
	
	checkStatusCode(t, resp, http.StatusCreated)
	
	// Extract user ID from response
	var createResp map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&createResp); err != nil {
		t.Fatalf("Failed to decode create response: %v", err)
	}
	
	data, ok := createResp["data"].(map[string]interface{})
	if !ok {
		t.Fatalf("Failed to extract data from response")
	}
	
	userID, ok := data["id"].(float64)
	if !ok {
		t.Fatalf("Failed to extract user ID from response")
	}
	
	// Store user ID for other tests
	testUserID = fmt.Sprintf("%.0f", userID)
	
	// Test get user by ID
	getURL := fmt.Sprintf("%s/api/v1/users/%s", baseURL, testUserID)
	resp, err = makeAuthenticatedRequest("GET", getURL, nil, adminToken)
	if err != nil {
		t.Fatalf("Failed to get user: %v", err)
	}
	defer resp.Body.Close()
	
	checkStatusCode(t, resp, http.StatusOK)
	
	// Test update user
	updateURL := fmt.Sprintf("%s/api/v1/users/%s", baseURL, testUserID)
	updateData := map[string]interface{}{
		"firstName": "Updated",
		"lastName":  "TestUser",
		"role":      "librarian", // Promote user to librarian
	}
	
	resp, err = makeAuthenticatedRequest("PUT", updateURL, updateData, adminToken)
	if err != nil {
		t.Fatalf("Failed to update user: %v", err)
	}
	defer resp.Body.Close()
	
	checkStatusCode(t, resp, http.StatusOK)
	
	// Test get user after update
	resp, err = makeAuthenticatedRequest("GET", getURL, nil, adminToken)
	if err != nil {
		t.Fatalf("Failed to get updated user: %v", err)
	}
	defer resp.Body.Close()
	
	checkStatusCode(t, resp, http.StatusOK)
	
	var getResp map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&getResp); err != nil {
		t.Fatalf("Failed to decode get response: %v", err)
	}
	
	getData, ok := getResp["data"].(map[string]interface{})
	if !ok {
		t.Fatalf("Failed to extract data from get response")
	}
	
	// Verify user was updated
	if firstName, ok := getData["firstName"].(string); !ok || firstName != "Updated" {
		t.Errorf("Expected firstName 'Updated', got %v", firstName)
	}
	
	if role, ok := getData["role"].(string); !ok || role != "librarian" {
		t.Errorf("Expected role 'librarian', got %v", role)
	}
	
	// Test list users
	listURL := fmt.Sprintf("%s/api/v1/users", baseURL)
	resp, err = makeAuthenticatedRequest("GET", listURL, nil, adminToken)
	if err != nil {
		t.Fatalf("Failed to list users: %v", err)
	}
	defer resp.Body.Close()
	
	checkStatusCode(t, resp, http.StatusOK)
	
	// Test delete user (if supported)
	deleteURL := fmt.Sprintf("%s/api/v1/users/%s", baseURL, testUserID)
	resp, err = makeAuthenticatedRequest("DELETE", deleteURL, nil, adminToken)
	if err != nil {
		t.Fatalf("Failed to delete user: %v", err)
	}
	defer resp.Body.Close()
	
	// Note: API might not allow actual deletion (soft delete)
	// Both 200 OK and 204 No Content are acceptable
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		t.Errorf("Expected status OK or NoContent for delete, got %d", resp.StatusCode)
	}
	
	// Verify user was deleted or deactivated
	resp, err = makeAuthenticatedRequest("GET", getURL, nil, adminToken)
	if err != nil {
		t.Fatalf("Failed to get deleted user: %v", err)
	}
	defer resp.Body.Close()
	
	// Either 404 Not Found or 200 OK with isActive=false is acceptable
	if resp.StatusCode != http.StatusNotFound && resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status NotFound or OK for deleted user, got %d", resp.StatusCode)
	}
	
	if resp.StatusCode == http.StatusOK {
		// If the API uses soft delete, verify the user is marked as inactive
		var deleteVerifyResp map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&deleteVerifyResp); err != nil {
			t.Fatalf("Failed to decode deleted user response: %v", err)
		}
		
		deleteVerifyData, ok := deleteVerifyResp["data"].(map[string]interface{})
		if !ok {
			t.Fatalf("Failed to extract data from deleted user response")
		}
		
		isActive, ok := deleteVerifyData["isActive"].(bool)
		if !ok || isActive {
			t.Errorf("Expected isActive to be false for deleted user")
		}
	}
}

// TestUserPermissions tests authorization for user management
func TestUserPermissions(t *testing.T) {
	// Test member cannot create users
	createURL := fmt.Sprintf("%s/api/v1/users", baseURL)
	userData := map[string]interface{}{
		"email":     "perm.test@example.com",
		"password":  "Password123!",
		"firstName": "Permission",
		"lastName":  "Test",
		"role":      "member",
	}
	
	resp, err := makeAuthenticatedRequest("POST", createURL, userData, memberToken)
	if err != nil {
		t.Fatalf("Failed to make create request as member: %v", err)
	}
	defer resp.Body.Close()
	
	checkStatusCode(t, resp, http.StatusForbidden)
	
	// Test librarian create permissions (may or may not be allowed)
	resp, err = makeAuthenticatedRequest("POST", createURL, userData, librianToken)
	if err != nil {
		t.Fatalf("Failed to make create request as librarian: %v", err)
	}
	defer resp.Body.Close()
	
	// Either forbidden or created is acceptable depending on implementation
	if resp.StatusCode != http.StatusForbidden && resp.StatusCode != http.StatusCreated {
		t.Errorf("Expected status Forbidden or Created, got %d", resp.StatusCode)
	}
	
	// If librarian can create users, extract the ID for further tests
	var libCreateUserID string
	if resp.StatusCode == http.StatusCreated {
		var createResp map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&createResp); err != nil {
			t.Fatalf("Failed to decode librarian create response: %v", err)
		}
		
		data, ok := createResp["data"].(map[string]interface{})
		if !ok {
			t.Fatalf("Failed to extract data from librarian create response")
		}
		
		userID, ok := data["id"].(float64)
		if !ok {
			t.Fatalf("Failed to extract user ID from librarian create response")
		}
		
		libCreateUserID = fmt.Sprintf("%.0f", userID)
	} else {
		// If librarian cannot create users, create a test user with admin token
		resp, err = makeAuthenticatedRequest("POST", createURL, userData, adminToken)
		if err != nil {
			t.Fatalf("Failed to create test user as admin: %v", err)
		}
		defer resp.Body.Close()
		
		checkStatusCode(t, resp, http.StatusCreated)
		
		var createResp map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&createResp); err != nil {
			t.Fatalf("Failed to decode admin create response: %v", err)
		}
		
		data, ok := createResp["data"].(map[string]interface{})
		if !ok {
			t.Fatalf("Failed to extract data from admin create response")
		}
		
		userID, ok := data["id"].(float64)
		if !ok {
			t.Fatalf("Failed to extract user ID from admin create response")
		}
		
		libCreateUserID = fmt.Sprintf("%.0f", userID)
	}
	
	// Test member cannot view other users' details
	getURL := fmt.Sprintf("%s/api/v1/users/%s", baseURL, libCreateUserID)
	resp, err = makeAuthenticatedRequest("GET", getURL, nil, memberToken)
	if err != nil {
		t.Fatalf("Failed to make get request as member: %v", err)
	}
	defer resp.Body.Close()
	
	checkStatusCode(t, resp, http.StatusForbidden)
	
	// Test member cannot update other users
	updateURL := fmt.Sprintf("%s/api/v1/users/%s", baseURL, libCreateUserID)
	updateData := map[string]interface{}{
		"firstName": "Hacked",
		"lastName":  "ByMember",
	}
	
	resp, err = makeAuthenticatedRequest("PUT", updateURL, updateData, memberToken)
	if err != nil {
		t.Fatalf("Failed to make update request as member: %v", err)
	}
	defer resp.Body.Close()
	
	checkStatusCode(t, resp, http.StatusForbidden)
	
	// Test librarian access to normal member user
	resp, err = makeAuthenticatedRequest("GET", getURL, nil, librianToken)
	if err != nil {
		t.Fatalf("Failed to make get request as librarian: %v", err)
	}
	defer resp.Body.Close()
	
	// Librarian should be able to view member details
	checkStatusCode(t, resp, http.StatusOK)
	
	// Clean up - delete test user
	deleteURL := fmt.Sprintf("%s/api/v1/users/%s", baseURL, libCreateUserID)
	resp, err = makeAuthenticatedRequest("DELETE", deleteURL, nil, adminToken)
	if err != nil {
		t.Fatalf("Failed to delete test user: %v", err)
	}
	resp.Body.Close()
}

// TestUserValidation tests validation of user data
func TestUserValidation(t *testing.T) {
	createURL := fmt.Sprintf("%s/api/v1/users", baseURL)
	
	// Test invalid email format
	invalidEmailData := map[string]interface{}{
		"email":     "invalid-email",
		"password":  "Password123!",
		"firstName": "Invalid",
		"lastName":  "Email",
		"role":      "member",
	}
	
	resp, err := makeAuthenticatedRequest("POST", createURL, invalidEmailData, adminToken)
	if err != nil {
		t.Fatalf("Failed to make invalid email request: %v", err)
	}
	defer resp.Body.Close()
	
	checkStatusCode(t, resp, http.StatusBadRequest)
	
	// Test weak password
	weakPasswordData := map[string]interface{}{
		"email":     "valid@example.com",
		"password":  "weak",
		"firstName": "Weak",
		"lastName":  "Password",
		"role":      "member",
	}
	
	resp, err = makeAuthenticatedRequest("POST", createURL, weakPasswordData, adminToken)
	if err != nil {
		t.Fatalf("Failed to make weak password request: %v", err)
	}
	defer resp.Body.Close()
	
	checkStatusCode(t, resp, http.StatusBadRequest)
	
	// Test missing required fields
	missingFieldsData := map[string]interface{}{
		"email":    "missing@example.com",
		"password": "Password123!",
		// Missing firstName and lastName
	}
	
	resp, err = makeAuthenticatedRequest("POST", createURL, missingFieldsData, adminToken)
	if err != nil {
		t.Fatalf("Failed to make missing fields request: %v", err)
	}
	defer resp.Body.Close()
	
	checkStatusCode(t, resp, http.StatusBadRequest)
	
	// Test invalid role
	invalidRoleData := map[string]interface{}{
		"email":     "role@example.com",
		"password":  "Password123!",
		"firstName": "Invalid",
		"lastName":  "Role",
		"role":      "invalid_role",
	}
	
	resp, err = makeAuthenticatedRequest("POST", createURL, invalidRoleData, adminToken)
	if err != nil {
		t.Fatalf("Failed to make invalid role request: %v", err)
	}
	defer resp.Body.Close()
	
	checkStatusCode(t, resp, http.StatusBadRequest)
	
	// Test duplicate email
	// First create a user with a unique email
	uniqueEmailData := map[string]interface{}{
		"email":     "duplicate.test@example.com",
		"password":  "Password123!",
		"firstName": "Duplicate",
		"lastName":  "Test",
		"role":      "member",
	}
	
	resp, err = makeAuthenticatedRequest("POST", createURL, uniqueEmailData, adminToken)
	if err != nil {
		t.Fatalf("Failed to create user for duplicate test: %v", err)
	}
	resp.Body.Close()
	
	checkStatusCode(t, resp, http.StatusCreated)
	
	// Now try to create another user with the same email
	resp, err = makeAuthenticatedRequest("POST", createURL, uniqueEmailData, adminToken)
	if err != nil {
		t.Fatalf("Failed to make duplicate email request: %v", err)
	}
	defer resp.Body.Close()
	
	checkStatusCode(t, resp, http.StatusConflict)
}

// TestUserProfile tests user profile management
func TestUserProfile(t *testing.T) {
	// Test user can view their own profile
	profileURL := fmt.Sprintf("%s/api/v1/users/profile", baseURL)
	resp, err := makeAuthenticatedRequest("GET", profileURL, nil, memberToken)
	if err != nil {
		t.Fatalf("Failed to get own profile: %v", err)
	}
	defer resp.Body.Close()
	
	checkStatusCode(t, resp, http.StatusOK)
	
	// Test user can update their own profile
	updateProfileData := map[string]interface{}{
		"firstName": "Updated",
		"lastName":  "Profile",
	}
	
	resp, err = makeAuthenticatedRequest("PUT", profileURL, updateProfileData, memberToken)
	if err != nil {
		t.Fatalf("Failed to update own profile: %v", err)
	}
	defer resp.Body.Close()
	
	checkStatusCode(t, resp, http.StatusOK)
	
	// Test change password functionality
	changePasswordURL := fmt.Sprintf("%s/api/v1/users/change-password", baseURL)
	changePasswordData := map[string]interface{}{
		"currentPassword": "Member123!",
		"newPassword":     "NewPassword123!",
	}
	
	resp, err = makeAuthenticatedRequest("PUT", changePasswordURL, changePasswordData, memberToken)
	if err != nil {
		t.Fatalf("Failed to change password: %v", err)
	}
	defer resp.Body.Close()
	
	checkStatusCode(t, resp, http.StatusOK)
	
	// Change it back to avoid breaking other tests
	changeBackData := map[string]interface{}{
		"currentPassword": "NewPassword123!",
		"newPassword":     "Member123!",
	}
	
	resp, err = makeAuthenticatedRequest("PUT", changePasswordURL, changeBackData, memberToken)
	if err != nil {
		t.Fatalf("Failed to change password back: %v", err)
	}
	defer resp.Body.Close()
	
	checkStatusCode(t, resp, http.StatusOK)
}
