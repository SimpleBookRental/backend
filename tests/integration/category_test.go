package integration

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
)

// TestCategoryCRUD tests the category API endpoints
func TestCategoryCRUD(t *testing.T) {
	// Test create category (admin/librarian only)
	createURL := fmt.Sprintf("%s/api/v1/categories", baseURL)
	categoryData := map[string]interface{}{
		"name":        "Test Category",
		"description": "Category for integration testing",
	}
	
	// Create with librarian token
	resp, err := makeAuthenticatedRequest("POST", createURL, categoryData, librianToken)
	if err != nil {
		t.Fatalf("Failed to create category: %v", err)
	}
	defer resp.Body.Close()
	
	checkStatusCode(t, resp, http.StatusCreated)
	
	// Extract category ID from response
	var createResp map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&createResp); err != nil {
		t.Fatalf("Failed to decode create response: %v", err)
	}
	
	data, ok := createResp["data"].(map[string]interface{})
	if !ok {
		t.Fatalf("Failed to extract data from response")
	}
	
	categoryID, ok := data["id"].(float64)
	if !ok {
		t.Fatalf("Failed to extract category ID from response")
	}
	
	// Store category ID for other tests
	testCategoryID = fmt.Sprintf("%.0f", categoryID)
	
	// Test get category by ID
	getURL := fmt.Sprintf("%s/api/v1/categories/%s", baseURL, testCategoryID)
	resp, err = makeAuthenticatedRequest("GET", getURL, nil, memberToken) // Anyone should be able to get categories
	if err != nil {
		t.Fatalf("Failed to get category: %v", err)
	}
	defer resp.Body.Close()
	
	checkStatusCode(t, resp, http.StatusOK)
	
	// Test update category
	updateURL := fmt.Sprintf("%s/api/v1/categories/%s", baseURL, testCategoryID)
	updateData := map[string]interface{}{
		"name":        "Updated Test Category",
		"description": "Updated description for integration testing",
	}
	
	resp, err = makeAuthenticatedRequest("PUT", updateURL, updateData, adminToken)
	if err != nil {
		t.Fatalf("Failed to update category: %v", err)
	}
	defer resp.Body.Close()
	
	checkStatusCode(t, resp, http.StatusOK)
	
	// Test get category after update
	resp, err = makeAuthenticatedRequest("GET", getURL, nil, memberToken)
	if err != nil {
		t.Fatalf("Failed to get updated category: %v", err)
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
	
	// Verify category was updated
	if name, ok := getData["name"].(string); !ok || name != "Updated Test Category" {
		t.Errorf("Expected name 'Updated Test Category', got %v", name)
	}
	
	// Test list categories
	listURL := fmt.Sprintf("%s/api/v1/categories", baseURL)
	resp, err = makeAuthenticatedRequest("GET", listURL, nil, memberToken)
	if err != nil {
		t.Fatalf("Failed to list categories: %v", err)
	}
	defer resp.Body.Close()
	
	checkStatusCode(t, resp, http.StatusOK)
	
	// Verify our test category is in the list
	var listResp map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&listResp); err != nil {
		t.Fatalf("Failed to decode list response: %v", err)
	}
	
	listData, ok := listResp["data"].(map[string]interface{})
	if !ok {
		t.Fatalf("Failed to extract data from list response")
	}
	
	categories, ok := listData["categories"].([]interface{})
	if !ok {
		t.Fatalf("Failed to extract categories from list response")
	}
	
	found := false
	for _, cat := range categories {
		category, ok := cat.(map[string]interface{})
		if !ok {
			continue
		}
		
		id, ok := category["id"].(float64)
		if !ok {
			continue
		}
		
		if fmt.Sprintf("%.0f", id) == testCategoryID {
			found = true
			break
		}
	}
	
	if !found {
		t.Errorf("Expected to find category with ID %s in list, but it was not found", testCategoryID)
	}
}

// TestCategoryBooksRelation tests the relationship between categories and books
func TestCategoryBooksRelation(t *testing.T) {
	// First, create a new category
	createCategoryURL := fmt.Sprintf("%s/api/v1/categories", baseURL)
	categoryData := map[string]interface{}{
		"name":        "Fiction Category",
		"description": "Category for fiction books",
	}
	
	resp, err := makeAuthenticatedRequest("POST", createCategoryURL, categoryData, librianToken)
	if err != nil {
		t.Fatalf("Failed to create category: %v", err)
	}
	defer resp.Body.Close()
	
	checkStatusCode(t, resp, http.StatusCreated)
	
	// Extract category ID
	var createCategoryResp map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&createCategoryResp); err != nil {
		t.Fatalf("Failed to decode create category response: %v", err)
	}
	
	categoryData, ok := createCategoryResp["data"].(map[string]interface{})
	if !ok {
		t.Fatalf("Failed to extract data from category response")
	}
	
	categoryID, ok := categoryData["id"].(float64)
	if !ok {
		t.Fatalf("Failed to extract category ID from response")
	}
	
	// Create books in this category
	createBookURL := fmt.Sprintf("%s/api/v1/books", baseURL)
	for i := 1; i <= 3; i++ {
		bookData := map[string]interface{}{
			"title":       fmt.Sprintf("Fiction Book %d", i),
			"author":      fmt.Sprintf("Fiction Author %d", i),
			"publisher":   "Fiction Publisher",
			"isbn":        fmt.Sprintf("9999999%04d", i),
			"categoryID":  fmt.Sprintf("%.0f", categoryID),
			"year":        2023,
			"description": fmt.Sprintf("Fiction book %d for category testing", i),
			"quantity":    5,
			"price":       19.99,
		}
		
		resp, err = makeAuthenticatedRequest("POST", createBookURL, bookData, librianToken)
		if err != nil {
			t.Fatalf("Failed to create book: %v", err)
		}
		resp.Body.Close()
		
		checkStatusCode(t, resp, http.StatusCreated)
	}
	
	// Get books by category
	categoryBooksURL := fmt.Sprintf("%s/api/v1/categories/%s/books", baseURL, fmt.Sprintf("%.0f", categoryID))
	resp, err = makeAuthenticatedRequest("GET", categoryBooksURL, nil, memberToken)
	if err != nil {
		t.Fatalf("Failed to get books by category: %v", err)
	}
	defer resp.Body.Close()
	
	checkStatusCode(t, resp, http.StatusOK)
	
	var booksByCategoryResp map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&booksByCategoryResp); err != nil {
		t.Fatalf("Failed to decode books by category response: %v", err)
	}
	
	booksByCategoryData, ok := booksByCategoryResp["data"].(map[string]interface{})
	if !ok {
		t.Fatalf("Failed to extract data from books by category response")
	}
	
	books, ok := booksByCategoryData["books"].([]interface{})
	if !ok {
		t.Fatalf("Failed to extract books from books by category response")
	}
	
	// Verify we got the books we created
	if len(books) != 3 {
		t.Errorf("Expected 3 books in category, got %d", len(books))
	}
}

// TestCategoryPermissions tests permissions for category operations
func TestCategoryPermissions(t *testing.T) {
	createURL := fmt.Sprintf("%s/api/v1/categories", baseURL)
	categoryData := map[string]interface{}{
		"name":        "Permission Test Category",
		"description": "Category for permission testing",
	}
	
	// Test member cannot create categories
	resp, err := makeAuthenticatedRequest("POST", createURL, categoryData, memberToken)
	if err != nil {
		t.Fatalf("Failed to make create request as member: %v", err)
	}
	defer resp.Body.Close()
	
	checkStatusCode(t, resp, http.StatusForbidden)
	
	// Test unauthenticated user cannot create categories
	resp, err = makeAuthenticatedRequest("POST", createURL, categoryData, "")
	if err != nil {
		t.Fatalf("Failed to make create request without authentication: %v", err)
	}
	defer resp.Body.Close()
	
	checkStatusCode(t, resp, http.StatusUnauthorized)
	
	// Create a category as admin for update/delete tests
	resp, err = makeAuthenticatedRequest("POST", createURL, categoryData, adminToken)
	if err != nil {
		t.Fatalf("Failed to create category as admin: %v", err)
	}
	defer resp.Body.Close()
	
	checkStatusCode(t, resp, http.StatusCreated)
	
	// Extract category ID
	var createResp map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&createResp); err != nil {
		t.Fatalf("Failed to decode create response: %v", err)
	}
	
	data, ok := createResp["data"].(map[string]interface{})
	if !ok {
		t.Fatalf("Failed to extract data from response")
	}
	
	categoryID, ok := data["id"].(float64)
	if !ok {
		t.Fatalf("Failed to extract category ID from response")
	}
	
	permCategoryID := fmt.Sprintf("%.0f", categoryID)
	
	// Test member cannot update categories
	updateURL := fmt.Sprintf("%s/api/v1/categories/%s", baseURL, permCategoryID)
	updateData := map[string]interface{}{
		"name":        "Updated by Member",
		"description": "This should fail",
	}
	
	resp, err = makeAuthenticatedRequest("PUT", updateURL, updateData, memberToken)
	if err != nil {
		t.Fatalf("Failed to make update request as member: %v", err)
	}
	defer resp.Body.Close()
	
	checkStatusCode(t, resp, http.StatusForbidden)
	
	// Test librarian can update categories
	updateData = map[string]interface{}{
		"name":        "Updated by Librarian",
		"description": "This should succeed",
	}
	
	resp, err = makeAuthenticatedRequest("PUT", updateURL, updateData, librianToken)
	if err != nil {
		t.Fatalf("Failed to make update request as librarian: %v", err)
	}
	defer resp.Body.Close()
	
	checkStatusCode(t, resp, http.StatusOK)
	
	// Test delete category - only admin should be able to delete
	deleteURL := fmt.Sprintf("%s/api/v1/categories/%s", baseURL, permCategoryID)
	
	// Test member cannot delete
	resp, err = makeAuthenticatedRequest("DELETE", deleteURL, nil, memberToken)
	if err != nil {
		t.Fatalf("Failed to make delete request as member: %v", err)
	}
	defer resp.Body.Close()
	
	checkStatusCode(t, resp, http.StatusForbidden)
	
	// Test librarian cannot delete (depends on implementation)
	resp, err = makeAuthenticatedRequest("DELETE", deleteURL, nil, librianToken)
	if err != nil {
		t.Fatalf("Failed to make delete request as librarian: %v", err)
	}
	defer resp.Body.Close()
	
	// Either forbidden or success is acceptable depending on implementation
	if resp.StatusCode != http.StatusForbidden && resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status Forbidden or OK, got %d", resp.StatusCode)
	}
	
	// If librarian didn't delete it, admin should be able to
	if resp.StatusCode == http.StatusForbidden {
		resp, err = makeAuthenticatedRequest("DELETE", deleteURL, nil, adminToken)
		if err != nil {
			t.Fatalf("Failed to make delete request as admin: %v", err)
		}
		defer resp.Body.Close()
		
		checkStatusCode(t, resp, http.StatusOK)
	}
}

// TestCategoryValidation tests validation of category data
func TestCategoryValidation(t *testing.T) {
	createURL := fmt.Sprintf("%s/api/v1/categories", baseURL)
	
	// Test empty name
	emptyNameData := map[string]interface{}{
		"name":        "",
		"description": "Should fail validation",
	}
	
	resp, err := makeAuthenticatedRequest("POST", createURL, emptyNameData, adminToken)
	if err != nil {
		t.Fatalf("Failed to make create request with empty name: %v", err)
	}
	defer resp.Body.Close()
	
	checkStatusCode(t, resp, http.StatusBadRequest)
	
	// Test name too long (if there's a limit)
	longName := ""
	for i := 0; i < 256; i++ {
		longName += "a"
	}
	
	longNameData := map[string]interface{}{
		"name":        longName,
		"description": "Should fail validation if there's a length limit",
	}
	
	resp, err = makeAuthenticatedRequest("POST", createURL, longNameData, adminToken)
	if err != nil {
		t.Fatalf("Failed to make create request with long name: %v", err)
	}
	defer resp.Body.Close()
	
	// Note: May or may not fail depending on implementation
	// We're just testing that the server handles it properly
	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status Created or BadRequest for long name, got %d", resp.StatusCode)
	}
	
	// Test duplicate name (create with same name)
	dupNameData := map[string]interface{}{
		"name":        "Duplicate Category Name",
		"description": "First category with this name",
	}
	
	// Create first category
	resp, err = makeAuthenticatedRequest("POST", createURL, dupNameData, adminToken)
	if err != nil {
		t.Fatalf("Failed to make create request for first category: %v", err)
	}
	resp.Body.Close()
	
	checkStatusCode(t, resp, http.StatusCreated)
	
	// Create second category with same name
	resp, err = makeAuthenticatedRequest("POST", createURL, dupNameData, adminToken)
	if err != nil {
		t.Fatalf("Failed to make create request for duplicate category: %v", err)
	}
	defer resp.Body.Close()
	
	// Should fail with conflict status
	checkStatusCode(t, resp, http.StatusConflict)
}
