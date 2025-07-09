package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"go-service/internal/api"
	"go-service/pkg/models"
)

// TestConfig holds configuration for test runs
type TestConfig struct {
	NodejsAPIURL     string
	GoServicePort    string
	TestAccessToken  string
	TestCSRFToken    string
	TestStudentID    string
	UseRealBackend   bool
}

// DefaultTestConfig returns default test configuration
func DefaultTestConfig() *TestConfig {
	return &TestConfig{
		NodejsAPIURL:    "http://localhost:5007",
		GoServicePort:   "8080",
		TestAccessToken: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6MSwicm9sZSI6ImFkbWluIiwicm9sZUlkIjoxLCJjc3JmX2htYWMiOiI4MTU1NTA5YWRjZjJhZjIwNzA0ZmUyNWVmYmUzMTBhZDk1MmE2NjBkZjNjYmFmZGExYWNhNTQzZjg3ZDA5NGI4IiwiaWF0IjoxNzUyMDA4NTM0LCJleHAiOjE3NTIwMDk0MzR9.BLorB5VRlhWh6HlUP9-obcAHgzCNalIyGNjjMFGbdew",
		TestCSRFToken:   "32175c1f-5df7-418b-a9a4-24eadf5d7526",
		TestStudentID:   "2",
		UseRealBackend:  false,
	}
}

// MockNodejsServer creates a mock Node.js server for testing
func MockNodejsServer() *httptest.Server {
	mux := http.NewServeMux()

	// Mock student endpoint
	mux.HandleFunc("/api/v1/students/", func(w http.ResponseWriter, r *http.Request) {
		// Extract student ID from path
		path := strings.TrimPrefix(r.URL.Path, "/api/v1/students/")
		studentID := strings.TrimSuffix(path, "/")
		
		// Check authentication
		authCookie := r.Header.Get("Cookie")
		csrfToken := r.Header.Get("X-CSRF-Token")
		
		if !strings.Contains(authCookie, "accessToken=") {
			http.Error(w, `{"error":"Authentication required"}`, http.StatusUnauthorized)
			return
		}
		
		if csrfToken == "" {
			http.Error(w, `{"error":"CSRF token required"}`, http.StatusForbidden)
			return
		}

		// Mock student data
		switch studentID {
		case "2":
			student := models.Student{
				ID:                 2,
				Name:               "Alice Johnson",
				Email:              "alice.johnson@school.edu",
				SystemAccess:       true,
				Phone:              "555-0102",
				Gender:             "Female",
				DOB:                time.Date(2005, 8, 15, 0, 0, 0, 0, time.UTC),
				Class:              "Grade 10",
				Section:            "A",
				Roll:               2,
				FatherName:         "Robert Johnson",
				FatherPhone:        "555-0103",
				MotherName:         "Sarah Johnson",
				MotherPhone:        "555-0104",
				GuardianName:       "Robert Johnson",
				GuardianPhone:      "555-0103",
				RelationOfGuardian: "Father",
				CurrentAddress:     "456 Oak Ave, Springfield, IL 62701",
				PermanentAddress:   "456 Oak Ave, Springfield, IL 62701",
				AdmissionDate:      time.Date(2020, 9, 1, 0, 0, 0, 0, time.UTC),
				ReporterName:       "Mrs. Smith",
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(student)
		case "999":
			http.Error(w, `{"error":"Student not found"}`, http.StatusNotFound)
		default:
			// Return a generic student for other IDs
			student := models.Student{
				ID:                 1,
				Name:               "Test Student",
				Email:              "test@school.edu",
				SystemAccess:       true,
				Phone:              "555-0001",
				Gender:             "Male",
				DOB:                time.Date(2005, 1, 1, 0, 0, 0, 0, time.UTC),
				Class:              "Grade 10",
				Section:            "A",
				Roll:               1,
				FatherName:         "Test Father",
				FatherPhone:        "555-0002",
				MotherName:         "Test Mother",
				MotherPhone:        "555-0003",
				GuardianName:       "Test Father",
				GuardianPhone:      "555-0002",
				RelationOfGuardian: "Father",
				CurrentAddress:     "123 Test St, Test City, TC 12345",
				PermanentAddress:   "123 Test St, Test City, TC 12345",
				AdmissionDate:      time.Date(2020, 9, 1, 0, 0, 0, 0, time.UTC),
				ReporterName:       "Test Teacher",
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(student)
		}
	})

	// Mock dashboard endpoint for health check
	mux.HandleFunc("/api/v1/dashboard", func(w http.ResponseWriter, r *http.Request) {
		// Check authentication
		authCookie := r.Header.Get("Cookie")
		if !strings.Contains(authCookie, "accessToken=") {
			http.Error(w, `{"error":"Authentication required"}`, http.StatusUnauthorized)
			return
		}
		
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok","data":{}}`))
	})

	return httptest.NewServer(mux)
}

// SetupTestEnvironment prepares the test environment
func SetupTestEnvironment(config *TestConfig) func() {
	// Set environment variables for testing
	os.Setenv("AUTH_MODE", "test")
	if !config.UseRealBackend {
		os.Setenv("NODEJS_API_URL", config.NodejsAPIURL)
	}

	// Return cleanup function
	return func() {
		os.Unsetenv("AUTH_MODE")
		os.Unsetenv("NODEJS_API_URL")
	}
}

// CreateTestServer creates a test HTTP server with the Go service router
func CreateTestServer() *httptest.Server {
	router := api.NewRouter()
	return httptest.NewServer(router)
}

// MakeAuthenticatedRequest creates an HTTP request with authentication tokens
func MakeAuthenticatedRequest(method, url string, body io.Reader, config *TestConfig) (*http.Request, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	// Add authentication cookies
	req.AddCookie(&http.Cookie{
		Name:  "accessToken",
		Value: config.TestAccessToken,
	})
	req.AddCookie(&http.Cookie{
		Name:  "csrfToken",
		Value: config.TestCSRFToken,
	})

	// Add CSRF header
	req.Header.Set("X-CSRF-Token", config.TestCSRFToken)

	return req, nil
}

// MakeUnauthenticatedRequest creates an HTTP request without authentication
func MakeUnauthenticatedRequest(method, url string, body io.Reader) (*http.Request, error) {
	return http.NewRequest(method, url, body)
}

// ValidatePDFResponse checks if the response contains a valid PDF
func ValidatePDFResponse(t *testing.T, resp *http.Response) []byte {
	t.Helper()

	// Check content type
	contentType := resp.Header.Get("Content-Type")
	if contentType != "application/pdf" {
		t.Errorf("Expected Content-Type: application/pdf, got: %s", contentType)
	}

	// Check content disposition
	contentDisposition := resp.Header.Get("Content-Disposition")
	if !strings.Contains(contentDisposition, "attachment") || !strings.Contains(contentDisposition, ".pdf") {
		t.Errorf("Invalid Content-Disposition header: %s", contentDisposition)
	}

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	// Check PDF magic bytes
	if len(body) < 4 || !bytes.HasPrefix(body, []byte("%PDF")) {
		t.Error("Response body does not appear to be a valid PDF")
	}

	// Check minimum PDF size (should be at least a few KB)
	if len(body) < 1000 {
		t.Errorf("PDF file seems too small: %d bytes", len(body))
	}

	return body
}

// ValidateHealthResponse checks if the health response is valid
func ValidateHealthResponse(t *testing.T, resp *http.Response, expectedHealthy bool) {
	t.Helper()

	// Check content type
	contentType := resp.Header.Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Expected Content-Type: application/json, got: %s", contentType)
	}

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	// Parse JSON response
	var healthResp map[string]interface{}
	if err := json.Unmarshal(body, &healthResp); err != nil {
		t.Fatalf("Failed to parse health response: %v", err)
	}

	// Check response structure
	if service, ok := healthResp["service"]; !ok || service != "go-pdf-service" {
		t.Error("Health response missing or invalid service field")
	}

	if expectedHealthy {
		if status, ok := healthResp["status"]; !ok || status != "healthy" {
			t.Errorf("Expected healthy status, got: %v", healthResp["status"])
		}
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got: %d", resp.StatusCode)
		}
	} else {
		if status, ok := healthResp["status"]; !ok || status != "unhealthy" {
			t.Errorf("Expected unhealthy status, got: %v", healthResp["status"])
		}
		if resp.StatusCode != http.StatusServiceUnavailable {
			t.Errorf("Expected status 503, got: %d", resp.StatusCode)
		}
	}
}

// ValidateErrorResponse checks if the error response is correctly formatted
func ValidateErrorResponse(t *testing.T, resp *http.Response, expectedStatusCode int, expectedErrorSubstring string) {
	t.Helper()

	if resp.StatusCode != expectedStatusCode {
		t.Errorf("Expected status code %d, got: %d", expectedStatusCode, resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	// Parse JSON error response
	var errorResp map[string]interface{}
	if err := json.Unmarshal(body, &errorResp); err != nil {
		t.Fatalf("Failed to parse error response: %v", err)
	}

	if errorMsg, ok := errorResp["error"]; !ok {
		t.Error("Error response missing 'error' field")
	} else if errorStr, ok := errorMsg.(string); !ok {
		t.Error("Error field is not a string")
	} else if !strings.Contains(errorStr, expectedErrorSubstring) {
		t.Errorf("Expected error message to contain '%s', got: %s", expectedErrorSubstring, errorStr)
	}
}

// LoginToRealBackend attempts to login to the real Node.js backend and return fresh tokens
func LoginToRealBackend(config *TestConfig) (accessToken, csrfToken string, err error) {
	loginURL := fmt.Sprintf("%s/api/v1/auth/login", config.NodejsAPIURL)
	
	loginData := map[string]string{
		"email":    "admin@school-admin.com",
		"password": "3OU4zn3q6Zh9",
	}
	
	loginJSON, err := json.Marshal(loginData)
	if err != nil {
		return "", "", fmt.Errorf("failed to marshal login data: %w", err)
	}

	resp, err := http.Post(loginURL, "application/json", bytes.NewBuffer(loginJSON))
	if err != nil {
		return "", "", fmt.Errorf("failed to login: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", "", fmt.Errorf("login failed with status %d: %s", resp.StatusCode, string(body))
	}

	// Extract tokens from cookies
	for _, cookie := range resp.Cookies() {
		switch cookie.Name {
		case "accessToken":
			accessToken = cookie.Value
		case "csrfToken":
			csrfToken = cookie.Value
		}
	}

	if accessToken == "" || csrfToken == "" {
		return "", "", fmt.Errorf("failed to extract authentication tokens from login response")
	}

	return accessToken, csrfToken, nil
}

// SkipIfNoBackend skips the test if the real backend is not available
func SkipIfNoBackend(t *testing.T, config *TestConfig) {
	if !config.UseRealBackend {
		return
	}

	// Test if backend is accessible
	resp, err := http.Get(fmt.Sprintf("%s/health", config.NodejsAPIURL))
	if err != nil || resp.StatusCode != http.StatusOK {
		t.Skip("Real Node.js backend not available, skipping test")
	}
	if resp != nil {
		resp.Body.Close()
	}
} 