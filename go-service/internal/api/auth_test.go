package api

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestExtractTokens tests token extraction from different sources
func TestExtractTokens(t *testing.T) {
	// Test 1: Extract from cookies
	req1 := httptest.NewRequest("GET", "/test", nil)
	req1.AddCookie(&http.Cookie{Name: "accessToken", Value: "test-access-token"})
	req1.AddCookie(&http.Cookie{Name: "csrfToken", Value: "test-csrf-token"})
	
	accessToken1, csrfToken1 := extractTokens(req1)
	if accessToken1 != "test-access-token" {
		t.Errorf("Expected accessToken 'test-access-token', got '%s'", accessToken1)
	}
	if csrfToken1 != "test-csrf-token" {
		t.Errorf("Expected csrfToken 'test-csrf-token', got '%s'", csrfToken1)
	}
	
	// Test 2: Extract from headers
	req2 := httptest.NewRequest("GET", "/test", nil)
	req2.Header.Set("Authorization", "Bearer header-access-token")
	req2.Header.Set("X-CSRF-Token", "header-csrf-token")
	
	accessToken2, csrfToken2 := extractTokens(req2)
	if accessToken2 != "header-access-token" {
		t.Errorf("Expected accessToken 'header-access-token', got '%s'", accessToken2)
	}
	if csrfToken2 != "header-csrf-token" {
		t.Errorf("Expected csrfToken 'header-csrf-token', got '%s'", csrfToken2)
	}
	
	// Test 3: No tokens
	req3 := httptest.NewRequest("GET", "/test", nil)
	accessToken3, csrfToken3 := extractTokens(req3)
	if accessToken3 != "" {
		t.Errorf("Expected empty accessToken, got '%s'", accessToken3)
	}
	if csrfToken3 != "" {
		t.Errorf("Expected empty csrfToken, got '%s'", csrfToken3)
	}
}

// TestSetTestTokens tests the test token setting functionality
func TestSetTestTokens(t *testing.T) {
	service := NewService()
	service.SetTestTokens()
	
	if service.NodejsClient.AccessToken == "" {
		t.Error("Expected test access token to be set")
	}
	if service.NodejsClient.CSRFToken == "" {
		t.Error("Expected test CSRF token to be set")
	}
	
	// Verify the tokens are the expected test values
	expectedAccessToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6MSwicm9sZSI6ImFkbWluIiwicm9sZUlkIjoxLCJjc3JmX2htYWMiOiI4MTU1NTA5YWRjZjJhZjIwNzA0ZmUyNWVmYmUzMTBhZDk1MmE2NjBkZjNjYmFmZGExYWNhNTQzZjg3ZDA5NGI4IiwiaWF0IjoxNzUyMDA4NTM0LCJleHAiOjE3NTIwMDk0MzR9.BLorB5VRlhWh6HlUP9-obcAHgzCNalIyGNjjMFGbdew"
	expectedCSRFToken := "32175c1f-5df7-418b-a9a4-24eadf5d7526"
	
	if service.NodejsClient.AccessToken != expectedAccessToken {
		t.Error("Access token does not match expected test value")
	}
	if service.NodejsClient.CSRFToken != expectedCSRFToken {
		t.Error("CSRF token does not match expected test value")
	}
} 