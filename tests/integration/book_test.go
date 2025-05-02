package integration

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
)

// TestBookCRUD tests the book API endpoints
func TestBookCRUD(t *testing.T) {
	// Use librarian token for book operations
	token := librianToken
	
	// Test create book
	createURL := fmt.Sprintf("%s/api/v1/books", baseURL)
	bookData := map[string]interface{}{
		"title":       "Test Book",
		"author":      "Test Author",
		"publisher":   "Test Publisher",
		"isbn":        "1234567890123",
		"categoryID":  1, // Assume category ID 1 exists
		"year":        2023,
		"description": "Test book description",
		"quantity":    10,
		"price":       19.99,
	}
	
	resp, err := makeAuthenticatedRequest("POST", createURL, bookData, token)
	if err != nil {
		t.Fatalf("Failed to create book: %v", err)
	}
	defer resp.Body.Close()
	
	checkStatusCode(t, resp, http.StatusCreated)
	
	// Extract book ID from response
	var createResp map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&createResp); err != nil {
		t.Fatalf("Failed to decode create response: %v", err)
	}
	
	data, ok := createResp["data"].(map[string]interface{})
	if !ok {
		t.Fatalf("Failed to extract data from response")
	}
	
	bookID, ok := data["id"].(float64)
	if !ok {
		t.Fatalf("Failed to extract book ID from response")
	}
	
	// Store book ID for other tests
	testBookID = fmt.Sprintf("%.0f", bookID)
	
	// Test get book by ID
	getURL := fmt.Sprintf("%s/api/v1/books/%s", baseURL, testBookID)
	resp, err = makeAuthenticatedRequest("GET", getURL, nil, token)
	if err != nil {
		t.Fatalf("Failed to get book: %v", err)
	}
	defer resp.Body.Close()
	
	checkStatusCode(t, resp, http.StatusOK)
	
	// Test update book
	updateURL := fmt.Sprintf("%s/api/v1/books/%s", baseURL, testBookID)
	updateData := map[string]interface{}{
		"title":      "Updated Test Book",
		"author":     "Updated Test Author",
		"quantity":   5,
		"price":      29.99,
	}
	
	resp, err = makeAuthenticatedRequest("PUT", updateURL, updateData, token)
	if err != nil {
		t.Fatalf("Failed to update book: %v", err)
	}
	defer resp.Body.Close()
	
	checkStatusCode(t, resp, http.StatusOK)
	
	// Test get book after update
	resp, err = makeAuthenticatedRequest("GET", getURL, nil, token)
	if err != nil {
		t.Fatalf("Failed to get updated book: %v", err)
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
	
	// Verify book was updated
	if title, ok := getData["title"].(string); !ok || title != "Updated Test Book" {
		t.Errorf("Expected title 'Updated Test Book', got %v", title)
	}
	
	// Test list books
	listURL := fmt.Sprintf("%s/api/v1/books", baseURL)
	resp, err = makeAuthenticatedRequest("GET", listURL, nil, token)
	if err != nil {
		t.Fatalf("Failed to list books: %v", err)
	}
	defer resp.Body.Close()
	
	checkStatusCode(t, resp, http.StatusOK)
	
	// Test delete book
	deleteURL := fmt.Sprintf("%s/api/v1/books/%s", baseURL, testBookID)
	resp, err = makeAuthenticatedRequest("DELETE", deleteURL, nil, token)
	if err != nil {
		t.Fatalf("Failed to delete book: %v", err)
	}
	defer resp.Body.Close()
	
	checkStatusCode(t, resp, http.StatusOK)
	
	// Verify book was deleted
	resp, err = makeAuthenticatedRequest("GET", getURL, nil, token)
	if err != nil {
		t.Fatalf("Failed to get book after deletion: %v", err)
	}
	defer resp.Body.Close()
	
	checkStatusCode(t, resp, http.StatusNotFound)
}

// TestBookSearch tests the book search endpoint
func TestBookSearch(t *testing.T) {
	// Create a test book for searching
	createURL := fmt.Sprintf("%s/api/v1/books", baseURL)
	searchBookData := map[string]interface{}{
		"title":       "Search Test Book",
		"author":      "Search Author",
		"publisher":   "Search Publisher",
		"isbn":        "9876543210987",
		"categoryID":  1,
		"year":        2023,
		"description": "Book for search test",
		"quantity":    5,
		"price":       24.99,
	}
	
	resp, err := makeAuthenticatedRequest("POST", createURL, searchBookData, librianToken)
	if err != nil {
		t.Fatalf("Failed to create search test book: %v", err)
	}
	defer resp.Body.Close()
	
	checkStatusCode(t, resp, http.StatusCreated)
	
	// Test search by title
	searchURL := fmt.Sprintf("%s/api/v1/books/search?title=Search", baseURL)
	resp, err = makeAuthenticatedRequest("GET", searchURL, nil, memberToken)
	if err != nil {
		t.Fatalf("Failed to search books by title: %v", err)
	}
	defer resp.Body.Close()
	
	checkStatusCode(t, resp, http.StatusOK)
	
	// Test search by author
	searchURL = fmt.Sprintf("%s/api/v1/books/search?author=Search", baseURL)
	resp, err = makeAuthenticatedRequest("GET", searchURL, nil, memberToken)
	if err != nil {
		t.Fatalf("Failed to search books by author: %v", err)
	}
	defer resp.Body.Close()
	
	checkStatusCode(t, resp, http.StatusOK)
	
	// Test search by ISBN
	searchURL = fmt.Sprintf("%s/api/v1/books/search?isbn=9876543210987", baseURL)
	resp, err = makeAuthenticatedRequest("GET", searchURL, nil, memberToken)
	if err != nil {
		t.Fatalf("Failed to search books by ISBN: %v", err)
	}
	defer resp.Body.Close()
	
	checkStatusCode(t, resp, http.StatusOK)
	
	// Test search with pagination
	searchURL = fmt.Sprintf("%s/api/v1/books/search?limit=10&page=1", baseURL)
	resp, err = makeAuthenticatedRequest("GET", searchURL, nil, memberToken)
	if err != nil {
		t.Fatalf("Failed to search books with pagination: %v", err)
	}
	defer resp.Body.Close()
	
	checkStatusCode(t, resp, http.StatusOK)
	
	// Test search with no results
	searchURL = fmt.Sprintf("%s/api/v1/books/search?title=NonExistentBook", baseURL)
	resp, err = makeAuthenticatedRequest("GET", searchURL, nil, memberToken)
	if err != nil {
		t.Fatalf("Failed to search books with no results: %v", err)
	}
	defer resp.Body.Close()
	
	checkStatusCode(t, resp, http.StatusOK)
	
	var searchResp map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&searchResp); err != nil {
		t.Fatalf("Failed to decode search response: %v", err)
	}
	
	// Verify search returned empty results
	data, ok := searchResp["data"].(map[string]interface{})
	if !ok {
		t.Fatalf("Failed to extract data from search response")
	}
	
	books, ok := data["books"].([]interface{})
	if !ok {
		t.Fatalf("Failed to extract books from search response")
	}
	
	if len(books) != 0 {
		t.Errorf("Expected 0 books for non-existent title search, got %d", len(books))
	}
}

// TestBookAvailability tests the book availability endpoint
func TestBookAvailability(t *testing.T) {
	// Create a test book for availability check
	createURL := fmt.Sprintf("%s/api/v1/books", baseURL)
	availBookData := map[string]interface{}{
		"title":       "Availability Test Book",
		"author":      "Availability Author",
		"publisher":   "Availability Publisher",
		"isbn":        "5555555555555",
		"categoryID":  1,
		"year":        2023,
		"description": "Book for availability test",
		"quantity":    3,
		"price":       14.99,
	}
	
	resp, err := makeAuthenticatedRequest("POST", createURL, availBookData, librianToken)
	if err != nil {
		t.Fatalf("Failed to create availability test book: %v", err)
	}
	defer resp.Body.Close()
	
	checkStatusCode(t, resp, http.StatusCreated)
	
	// Extract book ID from response
	var createResp map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&createResp); err != nil {
		t.Fatalf("Failed to decode create response: %v", err)
	}
	
	data, ok := createResp["data"].(map[string]interface{})
	if !ok {
		t.Fatalf("Failed to extract data from response")
	}
	
	bookID, ok := data["id"].(float64)
	if !ok {
		t.Fatalf("Failed to extract book ID from response")
	}
	
	// Store book ID for availability test
	availBookID := fmt.Sprintf("%.0f", bookID)
	
	// Test availability endpoint
	availURL := fmt.Sprintf("%s/api/v1/books/%s/availability", baseURL, availBookID)
	resp, err = makeAuthenticatedRequest("GET", availURL, nil, memberToken)
	if err != nil {
		t.Fatalf("Failed to check book availability: %v", err)
	}
	defer resp.Body.Close()
	
	checkStatusCode(t, resp, http.StatusOK)
	
	var availResp map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&availResp); err != nil {
		t.Fatalf("Failed to decode availability response: %v", err)
	}
	
	// Verify availability data
	availData, ok := availResp["data"].(map[string]interface{})
	if !ok {
		t.Fatalf("Failed to extract data from availability response")
	}
	
	available, ok := availData["available"].(bool)
	if !ok {
		t.Fatalf("Failed to extract available flag from response")
	}
	
	if !available {
		t.Errorf("Expected book to be available, but it's not")
	}
	
	availableCount, ok := availData["available_count"].(float64)
	if !ok {
		t.Fatalf("Failed to extract available count from response")
	}
	
	if int(availableCount) != 3 {
		t.Errorf("Expected 3 available books, got %d", int(availableCount))
	}
}
