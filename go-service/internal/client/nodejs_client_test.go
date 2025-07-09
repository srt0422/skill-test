package client

import (
	"testing"
)

// TestNewNodejsClient tests the creation of a new Node.js client
func TestNewNodejsClient(t *testing.T) {
	baseURL := "http://localhost:5007"
	client := NewNodejsClient(baseURL)

	if client == nil {
		t.Error("Expected client to be created, got nil")
	}

	if client.BaseURL != baseURL {
		t.Errorf("Expected BaseURL to be %s, got %s", baseURL, client.BaseURL)
	}

	if client.HTTPClient == nil {
		t.Error("Expected HTTPClient to be created, got nil")
	}
}

// TestSetAuthTokens tests setting authentication tokens
func TestSetAuthTokens(t *testing.T) {
	client := NewNodejsClient("http://localhost:5007")
	
	accessToken := "test-access-token"
	csrfToken := "test-csrf-token"
	
	client.SetAuthTokens(accessToken, csrfToken)
	
	if client.AccessToken != accessToken {
		t.Errorf("Expected AccessToken to be %s, got %s", accessToken, client.AccessToken)
	}
	
	if client.CSRFToken != csrfToken {
		t.Errorf("Expected CSRFToken to be %s, got %s", csrfToken, client.CSRFToken)
	}
}

// Note: Integration tests would require the Node.js backend to be running
// For now, we'll test the basic functionality without actual HTTP calls 