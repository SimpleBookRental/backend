package integration

import (
	"fmt"
	"net/http"
	"testing"
)

// TestReportAccess tests access to various report endpoints
func TestReportAccess(t *testing.T) {
	// Test popular books report - should be accessible to admins and librarians
	popularBooksURL := fmt.Sprintf("%s/api/v1/reports/popular-books", baseURL)
	
	// Admin should have access
	resp, err := makeAuthenticatedRequest("GET", popularBooksURL, nil, adminToken)
	if err != nil {
		t.Fatalf("Failed to access popular books report as admin: %v", err)
	}
	defer resp.Body.Close()
	
	checkStatusCode(t, resp, http.StatusOK)
	
	// Librarian should have access
	resp, err = makeAuthenticatedRequest("GET", popularBooksURL, nil, librianToken)
	if err != nil {
		t.Fatalf("Failed to access popular books report as librarian: %v", err)
	}
	defer resp.Body.Close()
	
	checkStatusCode(t, resp, http.StatusOK)
	
	// Member should NOT have access
	resp, err = makeAuthenticatedRequest("GET", popularBooksURL, nil, memberToken)
	if err != nil {
		t.Fatalf("Failed to make request to popular books report as member: %v", err)
	}
	defer resp.Body.Close()
	
	checkStatusCode(t, resp, http.StatusForbidden)
}

// TestReportQueries tests various query parameters for reports
func TestReportQueries(t *testing.T) {
	// Test report with time period - last 30 days
	thirtyDayReportURL := fmt.Sprintf("%s/api/v1/reports/popular-books?period=30days", baseURL)
	resp, err := makeAuthenticatedRequest("GET", thirtyDayReportURL, nil, adminToken)
	if err != nil {
		t.Fatalf("Failed to access 30-day report: %v", err)
	}
	defer resp.Body.Close()
	
	checkStatusCode(t, resp, http.StatusOK)
	
	// Test report with custom date range
	customRangeURL := fmt.Sprintf("%s/api/v1/reports/popular-books?startDate=2023-01-01&endDate=2023-12-31", baseURL)
	resp, err = makeAuthenticatedRequest("GET", customRangeURL, nil, adminToken)
	if err != nil {
		t.Fatalf("Failed to access custom date range report: %v", err)
	}
	defer resp.Body.Close()
	
	checkStatusCode(t, resp, http.StatusOK)
	
	// Test report with limit parameter
	limitedReportURL := fmt.Sprintf("%s/api/v1/reports/popular-books?limit=5", baseURL)
	resp, err = makeAuthenticatedRequest("GET", limitedReportURL, nil, adminToken)
	if err != nil {
		t.Fatalf("Failed to access limited report: %v", err)
	}
	defer resp.Body.Close()
	
	checkStatusCode(t, resp, http.StatusOK)
	
	// Test with invalid period parameter - should return 400 Bad Request
	invalidPeriodURL := fmt.Sprintf("%s/api/v1/reports/popular-books?period=invalid", baseURL)
	resp, err = makeAuthenticatedRequest("GET", invalidPeriodURL, nil, adminToken)
	if err != nil {
		t.Fatalf("Failed to make request with invalid period: %v", err)
	}
	defer resp.Body.Close()
	
	checkStatusCode(t, resp, http.StatusBadRequest)
}

// TestReportTypes tests different types of reports
func TestReportTypes(t *testing.T) {
	// Test revenue report
	revenueReportURL := fmt.Sprintf("%s/api/v1/reports/revenue", baseURL)
	resp, err := makeAuthenticatedRequest("GET", revenueReportURL, nil, adminToken)
	if err != nil {
		t.Fatalf("Failed to access revenue report: %v", err)
	}
	defer resp.Body.Close()
	
	checkStatusCode(t, resp, http.StatusOK)
	
	// Test overdue books report
	overdueReportURL := fmt.Sprintf("%s/api/v1/reports/overdue-books", baseURL)
	resp, err = makeAuthenticatedRequest("GET", overdueReportURL, nil, adminToken)
	if err != nil {
		t.Fatalf("Failed to access overdue books report: %v", err)
	}
	defer resp.Body.Close()
	
	checkStatusCode(t, resp, http.StatusOK)
	
	// Test user activity report
	userActivityURL := fmt.Sprintf("%s/api/v1/reports/user-activity", baseURL)
	resp, err = makeAuthenticatedRequest("GET", userActivityURL, nil, adminToken)
	if err != nil {
		t.Fatalf("Failed to access user activity report: %v", err)
	}
	defer resp.Body.Close()
	
	checkStatusCode(t, resp, http.StatusOK)
	
	// Test book utilization report
	bookUtilizationURL := fmt.Sprintf("%s/api/v1/reports/book-utilization", baseURL)
	resp, err = makeAuthenticatedRequest("GET", bookUtilizationURL, nil, adminToken)
	if err != nil {
		t.Fatalf("Failed to access book utilization report: %v", err)
	}
	defer resp.Body.Close()
	
	checkStatusCode(t, resp, http.StatusOK)
}

// TestReportFormats tests different output formats for reports
func TestReportFormats(t *testing.T) {
	// Test JSON format (default)
	jsonReportURL := fmt.Sprintf("%s/api/v1/reports/popular-books", baseURL)
	resp, err := makeAuthenticatedRequest("GET", jsonReportURL, nil, adminToken)
	if err != nil {
		t.Fatalf("Failed to access JSON report: %v", err)
	}
	defer resp.Body.Close()
	
	checkStatusCode(t, resp, http.StatusOK)
	
	// Test CSV format
	csvReportURL := fmt.Sprintf("%s/api/v1/reports/popular-books?format=csv", baseURL)
	resp, err = makeAuthenticatedRequest("GET", csvReportURL, nil, adminToken)
	if err != nil {
		t.Fatalf("Failed to access CSV report: %v", err)
	}
	defer resp.Body.Close()
	
	checkStatusCode(t, resp, http.StatusOK)
	
	// Verify content type for CSV
	contentType := resp.Header.Get("Content-Type")
	if contentType != "text/csv" {
		t.Errorf("Expected Content-Type 'text/csv', got '%s'", contentType)
	}
	
	// Test PDF format (if supported)
	pdfReportURL := fmt.Sprintf("%s/api/v1/reports/popular-books?format=pdf", baseURL)
	resp, err = makeAuthenticatedRequest("GET", pdfReportURL, nil, adminToken)
	if err != nil {
		t.Fatalf("Failed to make PDF report request: %v", err)
	}
	defer resp.Body.Close()
	
	// Note: If PDF is not supported, this might return 400 Bad Request
	// We're just checking that the server handles the request properly
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status OK or BadRequest for PDF format, got %d", resp.StatusCode)
	}
}

// TestReportCaching tests report caching behavior
func TestReportCaching(t *testing.T) {
	// Make initial request to generate cache
	cacheReportURL := fmt.Sprintf("%s/api/v1/reports/popular-books", baseURL)
	resp1, err := makeAuthenticatedRequest("GET", cacheReportURL, nil, adminToken)
	if err != nil {
		t.Fatalf("Failed to access report for cache test: %v", err)
	}
	
	// Get cache headers
	etag1 := resp1.Header.Get("ETag")
	lastModified1 := resp1.Header.Get("Last-Modified")
	resp1.Body.Close()
	
	// Make second request with If-None-Match header if ETag was provided
	if etag1 != "" {
		req, _ := http.NewRequest("GET", cacheReportURL, nil)
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", adminToken))
		req.Header.Set("If-None-Match", etag1)
		
		resp2, err := testClient.Do(req)
		if err != nil {
			t.Fatalf("Failed to make conditional request with ETag: %v", err)
		}
		defer resp2.Body.Close()
		
		// Should get 304 Not Modified if caching is implemented
		if resp2.StatusCode != http.StatusNotModified && resp2.StatusCode != http.StatusOK {
			t.Errorf("Expected status NotModified or OK for ETag request, got %d", resp2.StatusCode)
		}
	}
	
	// Make third request with If-Modified-Since if Last-Modified was provided
	if lastModified1 != "" {
		req, _ := http.NewRequest("GET", cacheReportURL, nil)
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", adminToken))
		req.Header.Set("If-Modified-Since", lastModified1)
		
		resp3, err := testClient.Do(req)
		if err != nil {
			t.Fatalf("Failed to make conditional request with Last-Modified: %v", err)
		}
		defer resp3.Body.Close()
		
		// Should get 304 Not Modified if caching is implemented
		if resp3.StatusCode != http.StatusNotModified && resp3.StatusCode != http.StatusOK {
			t.Errorf("Expected status NotModified or OK for Last-Modified request, got %d", resp3.StatusCode)
		}
	}
}
