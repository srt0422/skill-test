package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"go-service/internal/api"
)

const (
	DefaultPort = "8080"
)

func main() {
	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = DefaultPort
	}

	// Initialize API router
	router := api.NewRouter()

	// Configure server
	server := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	// Start server
	fmt.Printf("🚀 Go PDF Report Service starting on port %s\n", port)
	fmt.Printf("📊 Student Report Endpoint: http://localhost:%s/api/v1/students/{id}/report\n", port)
	
	log.Fatal(server.ListenAndServe())
} 