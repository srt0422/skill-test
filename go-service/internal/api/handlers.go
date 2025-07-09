package api

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"go-service/internal/client"
	"go-service/internal/pdf"

	"github.com/gorilla/mux"
)

// Service holds the dependencies for handlers
type Service struct {
	NodejsClient *client.NodejsClient
}

// NewService creates a new service with initialized dependencies
func NewService() *Service {
	// Get Node.js backend URL from environment or use default
	nodejsURL := os.Getenv("NODEJS_API_URL")
	if nodejsURL == "" {
		nodejsURL = "http://localhost:5007"
	}

	return &Service{
		NodejsClient: client.NewNodejsClient(nodejsURL),
	}
}

// HandleStudentReport generates and returns a PDF report for a student
func (s *Service) HandleStudentReport(w http.ResponseWriter, r *http.Request) {
	// Extract student ID from URL
	vars := mux.Vars(r)
	studentID := vars["id"]
	
	if studentID == "" {
		http.Error(w, `{"error":"Student ID is required"}`, http.StatusBadRequest)
		return
	}

	// TODO: In next task, we'll add authentication handling
	// For now, we'll make the request without authentication
	
	// Fetch student data from Node.js API
	student, err := s.NodejsClient.GetStudent(studentID)
	if err != nil {
		// Log the error for debugging
		fmt.Printf("Error fetching student %s: %v\n", studentID, err)
		
		// Return appropriate error response based on status code
		errorMsg := err.Error()
		if strings.Contains(errorMsg, "status 404") {
			http.Error(w, `{"error":"Student not found"}`, http.StatusNotFound)
			return
		}
		
		http.Error(w, `{"error":"Failed to fetch student data"}`, http.StatusInternalServerError)
		return
	}

	// Generate PDF report
	generator := pdf.NewGenerator()
	pdfBytes, err := generator.GenerateStudentReport(student)
	if err != nil {
		fmt.Printf("Error generating PDF for student %s: %v\n", studentID, err)
		http.Error(w, `{"error":"Failed to generate PDF report"}`, http.StatusInternalServerError)
		return
	}

	// Set response headers for PDF download
	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=student_%s_report.pdf", studentID))
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(pdfBytes)))

	// Write PDF to response
	_, err = w.Write(pdfBytes)
	if err != nil {
		fmt.Printf("Error writing PDF response for student %s: %v\n", studentID, err)
		return
	}

	fmt.Printf("Successfully generated PDF report for student %s\n", studentID)
}

// HandleHealth provides a health check endpoint
func (s *Service) HandleHealth(w http.ResponseWriter, r *http.Request) {
	// Check if Node.js API is accessible
	err := s.NodejsClient.HealthCheck()
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Write([]byte(`{"status":"unhealthy","service":"go-pdf-service","error":"Node.js API unavailable"}`))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"healthy","service":"go-pdf-service","nodejs_api":"connected"}`))
} 