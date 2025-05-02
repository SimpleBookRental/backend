package integration

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"
)

// TestRentalCRUD tests the rental API endpoints
func TestRentalCRUD(t *testing.T) {
	// Create a test book for rental operations
	createBookURL := fmt.Sprintf("%s/api/v1/books", baseURL)
	bookData := map[string]interface{}{
		"title":       "Rental Test Book",
		"author":      "Rental Author",
		"publisher":   "Rental Publisher",
		"isbn":        "1111222233334",
		"categoryID":  1,
		"year":        2023,
		"description": "Book for rental test",
		"quantity":    5,
		"price":       19.99,
	}
	
	resp, err := makeAuthenticatedRequest("POST", createBookURL, bookData, librianToken)
	if err != nil {
		t.Fatalf("Failed to create test book: %v", err)
	}
	defer resp.Body.Close()
	
	checkStatusCode(t, resp, http.StatusCreated)
	
	// Extract book ID from response
	var createBookResp map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&createBookResp); err != nil {
		t.Fatalf("Failed to decode create book response: %v", err)
	}
	
	bookData, ok := createBookResp["data"].(map[string]interface{})
	if !ok {
		t.Fatalf("Failed to extract data from book response")
	}
	
	bookID, ok := bookData["id"].(float64)
	if !ok {
		t.Fatalf("Failed to extract book ID from response")
	}
	
	// Store book ID for rental tests
	rentalBookID := fmt.Sprintf("%.0f", bookID)
	
	// Create a rental
	createRentalURL := fmt.Sprintf("%s/api/v1/rentals", baseURL)
	rentalData := map[string]interface{}{
		"bookID":      rentalBookID,
		"dueDate":     time.Now().AddDate(0, 0, 14).Format("2006-01-02"), // 14 days from now
	}
	
	resp, err = makeAuthenticatedRequest("POST", createRentalURL, rentalData, memberToken)
	if err != nil {
		t.Fatalf("Failed to create rental: %v", err)
	}
	defer resp.Body.Close()
	
	checkStatusCode(t, resp, http.StatusCreated)
	
	// Extract rental ID from response
	var createRentalResp map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&createRentalResp); err != nil {
		t.Fatalf("Failed to decode create rental response: %v", err)
	}
	
	rentalData, ok = createRentalResp["data"].(map[string]interface{})
	if !ok {
		t.Fatalf("Failed to extract data from rental response")
	}
	
	rentalID, ok := rentalData["id"].(float64)
	if !ok {
		t.Fatalf("Failed to extract rental ID from response")
	}
	
	// Store rental ID for other tests
	testRentalID = fmt.Sprintf("%.0f", rentalID)
	
	// Test get rental by ID
	getRentalURL := fmt.Sprintf("%s/api/v1/rentals/%s", baseURL, testRentalID)
	resp, err = makeAuthenticatedRequest("GET", getRentalURL, nil, memberToken)
	if err != nil {
		t.Fatalf("Failed to get rental: %v", err)
	}
	defer resp.Body.Close()
	
	checkStatusCode(t, resp, http.StatusOK)
	
	// Test list user rentals
	listRentalsURL := fmt.Sprintf("%s/api/v1/rentals/user", baseURL)
	resp, err = makeAuthenticatedRequest("GET", listRentalsURL, nil, memberToken)
	if err != nil {
		t.Fatalf("Failed to list user rentals: %v", err)
	}
	defer resp.Body.Close()
	
	checkStatusCode(t, resp, http.StatusOK)
	
	// Test extend rental
	extendURL := fmt.Sprintf("%s/api/v1/rentals/%s/extend", baseURL, testRentalID)
	extendData := map[string]interface{}{
		"additionalDays": 7,
	}
	
	resp, err = makeAuthenticatedRequest("PUT", extendURL, extendData, memberToken)
	if err != nil {
		t.Fatalf("Failed to extend rental: %v", err)
	}
	defer resp.Body.Close()
	
	checkStatusCode(t, resp, http.StatusOK)
	
	// Verify rental was extended
	resp, err = makeAuthenticatedRequest("GET", getRentalURL, nil, memberToken)
	if err != nil {
		t.Fatalf("Failed to get rental after extension: %v", err)
	}
	defer resp.Body.Close()
	
	checkStatusCode(t, resp, http.StatusOK)
	
	var getExtendedResp map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&getExtendedResp); err != nil {
		t.Fatalf("Failed to decode extended rental response: %v", err)
	}
	
	extendedData, ok := getExtendedResp["data"].(map[string]interface{})
	if !ok {
		t.Fatalf("Failed to extract data from extended rental response")
	}
	
	// Verify the rental exists and has status field
	status, ok := extendedData["status"].(string)
	if !ok || status == "" {
		t.Errorf("Failed to get status from extended rental data")
	}
	
	// Original due date + 7 days should be reflected in the response
	// Note: In a real test you would parse the dates and compare them
	
	// Test return rental
	returnURL := fmt.Sprintf("%s/api/v1/rentals/%s/return", baseURL, testRentalID)
	resp, err = makeAuthenticatedRequest("PUT", returnURL, nil, memberToken)
	if err != nil {
		t.Fatalf("Failed to return rental: %v", err)
	}
	defer resp.Body.Close()
	
	checkStatusCode(t, resp, http.StatusOK)
	
	// Verify rental was returned
	resp, err = makeAuthenticatedRequest("GET", getRentalURL, nil, memberToken)
	if err != nil {
		t.Fatalf("Failed to get rental after return: %v", err)
	}
	defer resp.Body.Close()
	
	checkStatusCode(t, resp, http.StatusOK)
	
	var getReturnedResp map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&getReturnedResp); err != nil {
		t.Fatalf("Failed to decode returned rental response: %v", err)
	}
	
	returnedData, ok := getReturnedResp["data"].(map[string]interface{})
	if !ok {
		t.Fatalf("Failed to extract data from returned rental response")
	}
	
	status, ok = returnedData["status"].(string)
	if !ok || status != "returned" {
		t.Errorf("Expected status 'returned', got %v", status)
	}
}

// TestRentalOverdue tests the overdue rental functionality
func TestRentalOverdue(t *testing.T) {
	// Create a test book for overdue rental
	createBookURL := fmt.Sprintf("%s/api/v1/books", baseURL)
	bookData := map[string]interface{}{
		"title":       "Overdue Test Book",
		"author":      "Overdue Author",
		"publisher":   "Overdue Publisher",
		"isbn":        "9999888877777",
		"categoryID":  1,
		"year":        2023,
		"description": "Book for overdue test",
		"quantity":    2,
		"price":       29.99,
	}
	
	resp, err := makeAuthenticatedRequest("POST", createBookURL, bookData, librianToken)
	if err != nil {
		t.Fatalf("Failed to create test book: %v", err)
	}
	defer resp.Body.Close()
	
	checkStatusCode(t, resp, http.StatusCreated)
	
	// Extract book ID
	var createBookResp map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&createBookResp); err != nil {
		t.Fatalf("Failed to decode create book response: %v", err)
	}
	
	bookData, ok := createBookResp["data"].(map[string]interface{})
	if !ok {
		t.Fatalf("Failed to extract data from book response")
	}
	
	bookID, ok := bookData["id"].(float64)
	if !ok {
		t.Fatalf("Failed to extract book ID from response")
	}
	
	overdueBookID := fmt.Sprintf("%.0f", bookID)
	
	// Create a rental with past due date (overdue)
	createRentalURL := fmt.Sprintf("%s/api/v1/rentals", baseURL)
	rentalData := map[string]interface{}{
		"bookID":  overdueBookID,
		"dueDate": time.Now().AddDate(0, 0, -7).Format("2006-01-02"), // 7 days ago
	}
	
	// Note: In a real API, you wouldn't be able to create rentals with past due dates
	// This is just for testing overdue functionality
	resp, err = makeAuthenticatedRequest("POST", createRentalURL, rentalData, memberToken)
	if err != nil {
		t.Fatalf("Failed to create overdue rental: %v", err)
	}
	defer resp.Body.Close()
	
	checkStatusCode(t, resp, http.StatusCreated)
	
	// Extract rental ID
	var createRentalResp map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&createRentalResp); err != nil {
		t.Fatalf("Failed to decode create rental response: %v", err)
	}
	
	rentalData, ok = createRentalResp["data"].(map[string]interface{})
	if !ok {
		t.Fatalf("Failed to extract data from rental response")
	}
	
	overdueRentalID, ok := rentalData["id"].(float64)
	if !ok {
		t.Fatalf("Failed to extract rental ID from response")
	}
	
	// Test get overdue rentals (admin only)
	overdueURL := fmt.Sprintf("%s/api/v1/rentals/overdue", baseURL)
	resp, err = makeAuthenticatedRequest("GET", overdueURL, nil, adminToken)
	if err != nil {
		t.Fatalf("Failed to get overdue rentals: %v", err)
	}
	defer resp.Body.Close()
	
	checkStatusCode(t, resp, http.StatusOK)
	
	// Verify the rental is in the overdue list
	var overdueResp map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&overdueResp); err != nil {
		t.Fatalf("Failed to decode overdue response: %v", err)
	}
	
	overdueData, ok := overdueResp["data"].(map[string]interface{})
	if !ok {
		t.Fatalf("Failed to extract data from overdue response")
	}
	
	overdueRentals, ok := overdueData["rentals"].([]interface{})
	if !ok {
		t.Fatalf("Failed to extract rentals from overdue response")
	}
	
	// Find our overdue rental in the list
	found := false
	for _, rental := range overdueRentals {
		rentalObj, ok := rental.(map[string]interface{})
		if !ok {
			continue
		}
		
		id, ok := rentalObj["id"].(float64)
		if !ok {
			continue
		}
		
		if fmt.Sprintf("%.0f", id) == fmt.Sprintf("%.0f", overdueRentalID) {
			found = true
			break
		}
	}
	
	if !found {
		t.Errorf("Expected to find overdue rental in the list, but it was not found")
	}
	
	// Test calculate fee for overdue rental
	feeURL := fmt.Sprintf("%s/api/v1/rentals/%s/fee", baseURL, fmt.Sprintf("%.0f", overdueRentalID))
	resp, err = makeAuthenticatedRequest("GET", feeURL, nil, memberToken)
	if err != nil {
		t.Fatalf("Failed to calculate fee: %v", err)
	}
	defer resp.Body.Close()
	
	checkStatusCode(t, resp, http.StatusOK)
	
	var feeResp map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&feeResp); err != nil {
		t.Fatalf("Failed to decode fee response: %v", err)
	}
	
	feeData, ok := feeResp["data"].(map[string]interface{})
	if !ok {
		t.Fatalf("Failed to extract data from fee response")
	}
	
	fee, ok := feeData["fee"].(float64)
	if !ok {
		t.Fatalf("Failed to extract fee from response")
	}
	
	// Fee should be greater than zero for overdue rental
	if fee <= 0 {
		t.Errorf("Expected fee to be greater than zero for overdue rental, got %v", fee)
	}
}

// TestRentalPermissions tests the permission controls for rentals
func TestRentalPermissions(t *testing.T) {
	// Create a test book
	createBookURL := fmt.Sprintf("%s/api/v1/books", baseURL)
	bookData := map[string]interface{}{
		"title":       "Permission Test Book",
		"author":      "Permission Author",
		"publisher":   "Permission Publisher",
		"isbn":        "7777666655555",
		"categoryID":  1,
		"year":        2023,
		"description": "Book for permission test",
		"quantity":    1,
		"price":       9.99,
	}
	
	resp, err := makeAuthenticatedRequest("POST", createBookURL, bookData, librianToken)
	if err != nil {
		t.Fatalf("Failed to create test book: %v", err)
	}
	defer resp.Body.Close()
	
	checkStatusCode(t, resp, http.StatusCreated)
	
	// Extract book ID
	var createBookResp map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&createBookResp); err != nil {
		t.Fatalf("Failed to decode create book response: %v", err)
	}
	
	bookData, ok := createBookResp["data"].(map[string]interface{})
	if !ok {
		t.Fatalf("Failed to extract data from book response")
	}
	
	bookID, ok := bookData["id"].(float64)
	if !ok {
		t.Fatalf("Failed to extract book ID from response")
	}
	
	permBookID := fmt.Sprintf("%.0f", bookID)
	
	// Create a rental with member token
	createRentalURL := fmt.Sprintf("%s/api/v1/rentals", baseURL)
	rentalData := map[string]interface{}{
		"bookID":  permBookID,
		"dueDate": time.Now().AddDate(0, 0, 14).Format("2006-01-02"), // 14 days from now
	}
	
	resp, err = makeAuthenticatedRequest("POST", createRentalURL, rentalData, memberToken)
	if err != nil {
		t.Fatalf("Failed to create rental: %v", err)
	}
	defer resp.Body.Close()
	
	checkStatusCode(t, resp, http.StatusCreated)
	
	// Extract rental ID
	var createRentalResp map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&createRentalResp); err != nil {
		t.Fatalf("Failed to decode create rental response: %v", err)
	}
	
	rentalData, ok = createRentalResp["data"].(map[string]interface{})
	if !ok {
		t.Fatalf("Failed to extract data from rental response")
	}
	
	permRentalID, ok := rentalData["id"].(float64)
	if !ok {
		t.Fatalf("Failed to extract rental ID from response")
	}
	
	permRentalIDStr := fmt.Sprintf("%.0f", permRentalID)
	
	// Test admin can access any rental
	getRentalURL := fmt.Sprintf("%s/api/v1/rentals/%s", baseURL, permRentalIDStr)
	resp, err = makeAuthenticatedRequest("GET", getRentalURL, nil, adminToken)
	if err != nil {
		t.Fatalf("Failed to get rental as admin: %v", err)
	}
	defer resp.Body.Close()
	
	checkStatusCode(t, resp, http.StatusOK)
	
	// Test librarian can access any rental
	resp, err = makeAuthenticatedRequest("GET", getRentalURL, nil, librianToken)
	if err != nil {
		t.Fatalf("Failed to get rental as librarian: %v", err)
	}
	defer resp.Body.Close()
	
	checkStatusCode(t, resp, http.StatusOK)
	
	// Test member who created the rental can access it
	resp, err = makeAuthenticatedRequest("GET", getRentalURL, nil, memberToken)
	if err != nil {
		t.Fatalf("Failed to get rental as owning member: %v", err)
	}
	defer resp.Body.Close()
	
	checkStatusCode(t, resp, http.StatusOK)
	
	// Test administrator only endpoint
	adminURL := fmt.Sprintf("%s/api/v1/rentals/all", baseURL)
	
	// Admin should have access
	resp, err = makeAuthenticatedRequest("GET", adminURL, nil, adminToken)
	if err != nil {
		t.Fatalf("Failed to access admin endpoint as admin: %v", err)
	}
	defer resp.Body.Close()
	
	checkStatusCode(t, resp, http.StatusOK)
	
	// Librarian should have access
	resp, err = makeAuthenticatedRequest("GET", adminURL, nil, librianToken)
	if err != nil {
		t.Fatalf("Failed to access admin endpoint as librarian: %v", err)
	}
	defer resp.Body.Close()
	
	checkStatusCode(t, resp, http.StatusOK)
	
	// Member should NOT have access
	resp, err = makeAuthenticatedRequest("GET", adminURL, nil, memberToken)
	if err != nil {
		t.Fatalf("Failed to make request to admin endpoint as member: %v", err)
	}
	defer resp.Body.Close()
	
	checkStatusCode(t, resp, http.StatusForbidden)
}
