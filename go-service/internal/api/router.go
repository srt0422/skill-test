package api

import (
	"os"
	
	"github.com/gorilla/mux"
)

// NewRouter creates and configures the API router
func NewRouter() *mux.Router {
	// Initialize service with dependencies
	service := NewService()
	
	// For development/testing, set test tokens if AUTH_MODE=test
	if os.Getenv("AUTH_MODE") == "test" {
		service.SetTestTokens()
	}
	
	router := mux.NewRouter()

	// API v1 routes
	api := router.PathPrefix("/api/v1").Subrouter()
	
	// Students routes with authentication middleware
	api.HandleFunc("/students/{id}/report", service.AuthMiddleware(service.HandleStudentReport)).Methods("GET")
	
	// Health check endpoint (no auth required)
	router.HandleFunc("/health", service.HandleHealth).Methods("GET")

	return router
} 