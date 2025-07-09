package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"
)

// TestEndToEndPDFGenerationWorkflow tests the complete user workflow
func TestEndToEndPDFGenerationWorkflow(t *testing.T) {
	// Start mock Node.js server
	mockServer := MockNodejsServer()
	defer mockServer.Close()

	// Configure test to use mock server
	config := DefaultTestConfig()
	config.NodejsAPIURL = mockServer.URL
	config.UseRealBackend = false
	
	// Set up environment
	cleanup := SetupTestEnvironment(config)
	defer cleanup()

	// Start Go service test server
	testServer := CreateTestServer()
	defer testServer.Close()

	t.Run("complete_student_report_workflow", func(t *testing.T) {
		// Step 1: Health Check - Verify service is ready
		t.Log("Step 1: Checking service health...")
		req, err := MakeAuthenticatedRequest("GET", testServer.URL+"/health", nil, config)
		if err != nil {
			t.Fatalf("Failed to create health check request: %v", err)
		}

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("Health check failed: %v", err)
		}
		resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Fatalf("Service not healthy: status %d", resp.StatusCode)
		}
		t.Log("✓ Service is healthy")

		// Step 2: Generate PDF for multiple students
		studentIDs := []string{"1", "2"}
		
		for _, studentID := range studentIDs {
			t.Logf("Step 2.%s: Generating PDF for student %s...", studentID, studentID)
			
			url := fmt.Sprintf("%s/api/v1/students/%s/report", testServer.URL, studentID)
			req, err := MakeAuthenticatedRequest("GET", url, nil, config)
			if err != nil {
				t.Fatalf("Failed to create PDF request for student %s: %v", studentID, err)
			}

			resp, err := client.Do(req)
			if err != nil {
				t.Fatalf("PDF generation failed for student %s: %v", studentID, err)
			}
			
			if resp.StatusCode != http.StatusOK {
				body, _ := io.ReadAll(resp.Body)
				resp.Body.Close()
				t.Fatalf("PDF generation failed for student %s: status %d, body: %s", studentID, resp.StatusCode, string(body))
			}

			// Validate PDF
			pdfBytes := ValidatePDFResponse(t, resp)
			resp.Body.Close()

			// Verify PDF size is reasonable
			if len(pdfBytes) < 1000 {
				t.Errorf("PDF for student %s seems too small: %d bytes", studentID, len(pdfBytes))
			}

			t.Logf("✓ Generated PDF for student %s (%d bytes)", studentID, len(pdfBytes))
		}

		t.Log("✓ Complete workflow successful")
	})
}

// TestEndToEndErrorRecovery tests how the system handles and recovers from errors
func TestEndToEndErrorRecovery(t *testing.T) {
	// Start mock Node.js server
	mockServer := MockNodejsServer()
	defer mockServer.Close()

	// Configure test to use mock server
	config := DefaultTestConfig()
	config.NodejsAPIURL = mockServer.URL
	config.UseRealBackend = false
	
	// Set up environment
	cleanup := SetupTestEnvironment(config)
	defer cleanup()

	// Start Go service test server
	testServer := CreateTestServer()
	defer testServer.Close()

	client := &http.Client{}

	t.Run("error_recovery_workflow", func(t *testing.T) {
		// Step 1: Try invalid student ID
		t.Log("Step 1: Testing error handling with invalid student ID...")
		
		url := testServer.URL + "/api/v1/students/999/report"
		req, err := MakeAuthenticatedRequest("GET", url, nil, config)
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("Request failed: %v", err)
		}
		
		if resp.StatusCode == http.StatusOK {
			resp.Body.Close()
			t.Fatal("Expected error for invalid student ID, but got success")
		}
		
		ValidateErrorResponse(t, resp, http.StatusNotFound, "Student not found")
		resp.Body.Close()
		t.Log("✓ Proper error handling for invalid student ID")

		// Step 2: Try without authentication
		t.Log("Step 2: Testing error handling without authentication...")
		
		url = testServer.URL + "/api/v1/students/2/report"
		req, err = MakeUnauthenticatedRequest("GET", url, nil)
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		resp, err = client.Do(req)
		if err != nil {
			t.Fatalf("Request failed: %v", err)
		}
		
		if resp.StatusCode == http.StatusOK {
			resp.Body.Close()
			t.Fatal("Expected authentication error, but got success")
		}
		resp.Body.Close()
		t.Log("✓ Proper error handling for missing authentication")

		// Step 3: Successful recovery with valid request
		t.Log("Step 3: Testing successful recovery...")
		
		req, err = MakeAuthenticatedRequest("GET", url, nil, config)
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		resp, err = client.Do(req)
		if err != nil {
			t.Fatalf("Request failed: %v", err)
		}
		
		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			resp.Body.Close()
			t.Fatalf("Expected success after recovery, got status %d: %s", resp.StatusCode, string(body))
		}
		
		ValidatePDFResponse(t, resp)
		resp.Body.Close()
		t.Log("✓ Successful recovery with valid request")

		t.Log("✓ Error recovery workflow completed")
	})
}

// TestEndToEndConcurrentRequests tests concurrent PDF generation
func TestEndToEndConcurrentRequests(t *testing.T) {
	// Start mock Node.js server
	mockServer := MockNodejsServer()
	defer mockServer.Close()

	// Configure test to use mock server
	config := DefaultTestConfig()
	config.NodejsAPIURL = mockServer.URL
	config.UseRealBackend = false
	
	// Set up environment
	cleanup := SetupTestEnvironment(config)
	defer cleanup()

	// Start Go service test server
	testServer := CreateTestServer()
	defer testServer.Close()

	t.Run("concurrent_pdf_generation", func(t *testing.T) {
		concurrency := 5
		studentIDs := []string{"1", "2", "1", "2", "1"} // Mix of student IDs
		
		results := make(chan error, concurrency)
		client := &http.Client{Timeout: 30 * time.Second}

		t.Logf("Starting %d concurrent PDF generation requests...", concurrency)
		
		// Start concurrent requests
		for i := 0; i < concurrency; i++ {
			go func(index int, studentID string) {
				url := fmt.Sprintf("%s/api/v1/students/%s/report", testServer.URL, studentID)
				
				req, err := MakeAuthenticatedRequest("GET", url, nil, config)
				if err != nil {
					results <- fmt.Errorf("request %d: failed to create request: %w", index, err)
					return
				}

				resp, err := client.Do(req)
				if err != nil {
					results <- fmt.Errorf("request %d: request failed: %w", index, err)
					return
				}
				defer resp.Body.Close()

				if resp.StatusCode != http.StatusOK {
					body, _ := io.ReadAll(resp.Body)
					results <- fmt.Errorf("request %d: status %d, body: %s", index, resp.StatusCode, string(body))
					return
				}

				// Validate PDF
				pdfBytes := ValidatePDFResponse(t, resp)
				if len(pdfBytes) < 1000 {
					results <- fmt.Errorf("request %d: PDF too small: %d bytes", index, len(pdfBytes))
					return
				}

				t.Logf("✓ Request %d completed successfully (%d bytes)", index, len(pdfBytes))
				results <- nil
			}(i, studentIDs[i])
		}

		// Wait for all requests to complete
		var errors []error
		for i := 0; i < concurrency; i++ {
			if err := <-results; err != nil {
				errors = append(errors, err)
			}
		}

		if len(errors) > 0 {
			t.Errorf("Concurrent requests failed:")
			for _, err := range errors {
				t.Errorf("  - %v", err)
			}
		} else {
			t.Log("✓ All concurrent requests completed successfully")
		}
	})
}

// TestEndToEndRealWorldScenario tests a realistic usage scenario
func TestEndToEndRealWorldScenario(t *testing.T) {
	// Start mock Node.js server
	mockServer := MockNodejsServer()
	defer mockServer.Close()

	// Configure test to use mock server
	config := DefaultTestConfig()
	config.NodejsAPIURL = mockServer.URL
	config.UseRealBackend = false
	
	// Set up environment
	cleanup := SetupTestEnvironment(config)
	defer cleanup()

	// Start Go service test server
	testServer := CreateTestServer()
	defer testServer.Close()

	t.Run("realistic_usage_scenario", func(t *testing.T) {
		client := &http.Client{}

		// Scenario: A teacher wants to generate reports for their class
		t.Log("Scenario: Teacher generating reports for class...")

		// Step 1: Check if service is available
		t.Log("Step 1: Checking service availability...")
		req, err := MakeUnauthenticatedRequest("GET", testServer.URL+"/health", nil)
		if err != nil {
			t.Fatalf("Failed to create health check: %v", err)
		}

		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("Health check failed: %v", err)
		}
		resp.Body.Close()
		
		// Service might be unhealthy without auth, but should respond
		if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusServiceUnavailable {
			t.Fatalf("Unexpected health check status: %d", resp.StatusCode)
		}
		t.Log("✓ Service is responding")

		// Step 2: Teacher authenticates (simulated by using test tokens)
		t.Log("Step 2: Teacher authentication...")
		// In real scenario, teacher would login through frontend
		t.Log("✓ Authentication successful (simulated)")

		// Step 3: Generate reports for multiple students in the class
		classStudents := []string{"1", "2"}
		generatedReports := make(map[string][]byte)

		for _, studentID := range classStudents {
			t.Logf("Step 3.%s: Generating report for student %s...", studentID, studentID)
			
			url := fmt.Sprintf("%s/api/v1/students/%s/report", testServer.URL, studentID)
			req, err := MakeAuthenticatedRequest("GET", url, nil, config)
			if err != nil {
				t.Fatalf("Failed to create request for student %s: %v", studentID, err)
			}

			// Add some realistic headers that a browser might send
			req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; School-Management-System)")
			req.Header.Set("Accept", "application/pdf,*/*")

			resp, err := client.Do(req)
			if err != nil {
				t.Fatalf("Failed to generate report for student %s: %v", studentID, err)
			}

			if resp.StatusCode != http.StatusOK {
				body, _ := io.ReadAll(resp.Body)
				resp.Body.Close()
				t.Fatalf("Report generation failed for student %s: status %d, body: %s", studentID, resp.StatusCode, string(body))
			}

			pdfBytes := ValidatePDFResponse(t, resp)
			resp.Body.Close()

			generatedReports[studentID] = pdfBytes
			t.Logf("✓ Generated report for student %s (%d bytes)", studentID, len(pdfBytes))

			// Simulate brief delay between requests (realistic usage)
			time.Sleep(100 * time.Millisecond)
		}

		// Step 4: Verify all reports were generated successfully
		t.Log("Step 4: Verifying all reports...")
		for studentID, pdfBytes := range generatedReports {
			if len(pdfBytes) < 1000 {
				t.Errorf("Report for student %s seems invalid (too small): %d bytes", studentID, len(pdfBytes))
			}
		}

		t.Logf("✓ Successfully generated %d reports", len(generatedReports))
		t.Log("✓ Realistic usage scenario completed successfully")
	})
}

// TestEndToEndWithRealBackend tests the complete flow with real Node.js backend
func TestEndToEndWithRealBackend(t *testing.T) {
	config := DefaultTestConfig()
	config.UseRealBackend = true

	// Skip if real backend is not available
	SkipIfNoBackend(t, config)

	t.Run("real_backend_end_to_end", func(t *testing.T) {
		// Step 1: Login to get fresh tokens
		t.Log("Step 1: Authenticating with real backend...")
		accessToken, csrfToken, err := LoginToRealBackend(config)
		if err != nil {
			t.Skipf("Failed to login to real backend: %v", err)
		}

		config.TestAccessToken = accessToken
		config.TestCSRFToken = csrfToken
		t.Log("✓ Successfully authenticated with real backend")

		// Step 2: Set up environment for real backend
		cleanup := SetupTestEnvironment(config)
		defer cleanup()

		// Step 3: Start Go service
		testServer := CreateTestServer()
		defer testServer.Close()
		t.Log("✓ Go service started")

		client := &http.Client{}

		// Step 4: Health check with real backend
		t.Log("Step 4: Health check with real backend...")
		req, err := MakeAuthenticatedRequest("GET", testServer.URL+"/health", nil, config)
		if err != nil {
			t.Fatalf("Failed to create health check: %v", err)
		}

		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("Health check failed: %v", err)
		}

		ValidateHealthResponse(t, resp, true)
		resp.Body.Close()
		t.Log("✓ Health check passed with real backend")

		// Step 5: Generate PDF with real data
		t.Log("Step 5: Generating PDF with real student data...")
		url := testServer.URL + "/api/v1/students/2/report"
		req, err = MakeAuthenticatedRequest("GET", url, nil, config)
		if err != nil {
			t.Fatalf("Failed to create PDF request: %v", err)
		}

		resp, err = client.Do(req)
		if err != nil {
			t.Fatalf("PDF generation failed: %v", err)
		}

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			resp.Body.Close()
			t.Fatalf("PDF generation failed: status %d, body: %s", resp.StatusCode, string(body))
		}

		pdfBytes := ValidatePDFResponse(t, resp)
		resp.Body.Close()

		// Step 6: Save PDF for manual verification (optional)
		if os.Getenv("SAVE_TEST_PDFS") == "true" {
			testDir := "test_output"
			os.MkdirAll(testDir, 0755)
			
			filename := filepath.Join(testDir, "real_backend_student_2_report.pdf")
			if err := os.WriteFile(filename, pdfBytes, 0644); err != nil {
				t.Logf("Warning: Failed to save test PDF: %v", err)
			} else {
				t.Logf("Test PDF saved to: %s", filename)
			}
		}

		t.Logf("✓ Successfully generated PDF with real backend (%d bytes)", len(pdfBytes))
		t.Log("✓ End-to-end test with real backend completed successfully")
	})
} 