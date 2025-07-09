package api

import (
	"net/http"
	"strings"
)

// AuthMiddleware extracts authentication tokens from the request and adds them to the client
func (s *Service) AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract tokens from various sources
		accessToken, csrfToken := extractTokens(r)
		
		// Set tokens in the Node.js client
		s.NodejsClient.SetAuthTokens(accessToken, csrfToken)
		
		// Call the next handler
		next(w, r)
	}
}

// extractTokens extracts authentication tokens from the request
func extractTokens(r *http.Request) (accessToken, csrfToken string) {
	// Method 1: Extract from cookies (preferred method)
	if cookie, err := r.Cookie("accessToken"); err == nil {
		accessToken = cookie.Value
	}
	if cookie, err := r.Cookie("csrfToken"); err == nil {
		csrfToken = cookie.Value
	}
	
	// Method 2: Extract from Authorization header as fallback
	if accessToken == "" {
		authHeader := r.Header.Get("Authorization")
		if strings.HasPrefix(authHeader, "Bearer ") {
			accessToken = strings.TrimPrefix(authHeader, "Bearer ")
		}
	}
	
	// Method 3: Extract CSRF from header as fallback
	if csrfToken == "" {
		csrfToken = r.Header.Get("X-CSRF-Token")
	}
	
	// Method 4: Extract from custom headers if needed
	if accessToken == "" {
		accessToken = r.Header.Get("X-Access-Token")
	}
	if csrfToken == "" {
		csrfToken = r.Header.Get("X-CSRF-Token")
	}
	
	return accessToken, csrfToken
}

// SetTestTokens sets hardcoded tokens for testing (when authentication is not available)
func (s *Service) SetTestTokens() {
	// Use the tokens from login_cookies.txt for testing
	accessToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6MSwicm9sZSI6ImFkbWluIiwicm9sZUlkIjoxLCJjc3JmX2htYWMiOiI4MTU1NTA5YWRjZjJhZjIwNzA0ZmUyNWVmYmUzMTBhZDk1MmE2NjBkZjNjYmFmZGExYWNhNTQzZjg3ZDA5NGI4IiwiaWF0IjoxNzUyMDA4NTM0LCJleHAiOjE3NTIwMDk0MzR9.BLorB5VRlhWh6HlUP9-obcAHgzCNalIyGNjjMFGbdew"
	csrfToken := "32175c1f-5df7-418b-a9a4-24eadf5d7526"
	
	s.NodejsClient.SetAuthTokens(accessToken, csrfToken)
} 