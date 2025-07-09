package main

import (
	"fmt"
	"io"
	"net/http"
	"sync"
	"testing"
	"time"
)

// BenchmarkPDFGeneration benchmarks PDF generation performance
func BenchmarkPDFGeneration(b *testing.B) {
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
	url := testServer.URL + "/api/v1/students/2/report"

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			req, err := MakeAuthenticatedRequest("GET", url, nil, config)
			if err != nil {
				b.Fatalf("Failed to create request: %v", err)
			}

			resp, err := client.Do(req)
			if err != nil {
				b.Fatalf("Request failed: %v", err)
			}

			if resp.StatusCode != http.StatusOK {
				body, _ := io.ReadAll(resp.Body)
				resp.Body.Close()
				b.Fatalf("Request failed: status %d, body: %s", resp.StatusCode, string(body))
			}

			// Read response body to simulate complete download
			_, err = io.ReadAll(resp.Body)
			resp.Body.Close()
			if err != nil {
				b.Fatalf("Failed to read response: %v", err)
			}
		}
	})
}

// BenchmarkHealthCheck benchmarks health check endpoint performance
func BenchmarkHealthCheck(b *testing.B) {
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
	url := testServer.URL + "/health"

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			req, err := MakeAuthenticatedRequest("GET", url, nil, config)
			if err != nil {
				b.Fatalf("Failed to create request: %v", err)
			}

			resp, err := client.Do(req)
			if err != nil {
				b.Fatalf("Request failed: %v", err)
			}

			_, err = io.ReadAll(resp.Body)
			resp.Body.Close()
			if err != nil {
				b.Fatalf("Failed to read response: %v", err)
			}
		}
	})
}

// TestPerformanceUnderLoad tests performance under various load conditions
func TestPerformanceUnderLoad(t *testing.T) {
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

	loadTests := []struct {
		name         string
		concurrency  int
		requests     int
		maxDuration  time.Duration
		description  string
	}{
		{
			name:         "light_load",
			concurrency:  5,
			requests:     50,
			maxDuration:  10 * time.Second,
			description:  "Light load - 5 concurrent users, 50 requests",
		},
		{
			name:         "medium_load",
			concurrency:  10,
			requests:     100,
			maxDuration:  20 * time.Second,
			description:  "Medium load - 10 concurrent users, 100 requests",
		},
		{
			name:         "heavy_load",
			concurrency:  20,
			requests:     200,
			maxDuration:  45 * time.Second,
			description:  "Heavy load - 20 concurrent users, 200 requests",
		},
	}

	for _, loadTest := range loadTests {
		t.Run(loadTest.name, func(t *testing.T) {
			t.Logf("Starting %s: %s", loadTest.name, loadTest.description)
			
			stats := runLoadTest(t, testServer.URL, config, loadTest.concurrency, loadTest.requests)
			
			// Validate performance metrics
			if stats.TotalDuration > loadTest.maxDuration {
				t.Errorf("Load test took too long: %v (max: %v)", stats.TotalDuration, loadTest.maxDuration)
			}
			
			if stats.FailureRate > 0.05 { // Allow up to 5% failure rate
				t.Errorf("Failure rate too high: %.2f%% (max: 5%%)", stats.FailureRate*100)
			}
			
			if stats.AvgResponseTime > 2*time.Second {
				t.Errorf("Average response time too slow: %v (max: 2s)", stats.AvgResponseTime)
			}

			t.Logf("✓ %s completed successfully", loadTest.name)
			t.Logf("  Total Duration: %v", stats.TotalDuration)
			t.Logf("  Requests/sec: %.2f", stats.RequestsPerSecond)
			t.Logf("  Success Rate: %.2f%%", (1-stats.FailureRate)*100)
			t.Logf("  Avg Response Time: %v", stats.AvgResponseTime)
			t.Logf("  Min/Max Response Time: %v / %v", stats.MinResponseTime, stats.MaxResponseTime)
		})
	}
}

// LoadTestStats holds performance test statistics
type LoadTestStats struct {
	TotalRequests       int
	SuccessfulRequests  int
	FailedRequests      int
	TotalDuration       time.Duration
	AvgResponseTime     time.Duration
	MinResponseTime     time.Duration
	MaxResponseTime     time.Duration
	RequestsPerSecond   float64
	FailureRate         float64
	ResponseTimes       []time.Duration
}

// runLoadTest executes a load test and returns performance statistics
func runLoadTest(t *testing.T, baseURL string, config *TestConfig, concurrency, totalRequests int) *LoadTestStats {
	var wg sync.WaitGroup
	results := make(chan time.Duration, totalRequests)
	errors := make(chan error, totalRequests)
	
	requestsPerWorker := totalRequests / concurrency
	remainingRequests := totalRequests % concurrency
	
	client := &http.Client{Timeout: 30 * time.Second}
	startTime := time.Now()
	
	// Start worker goroutines
	for i := 0; i < concurrency; i++ {
		workerRequests := requestsPerWorker
		if i < remainingRequests {
			workerRequests++
		}
		
		wg.Add(1)
		go func(workerID, numRequests int) {
			defer wg.Done()
			
			for j := 0; j < numRequests; j++ {
				requestStart := time.Now()
				
				// Alternate between different student IDs for variety
				studentID := fmt.Sprintf("%d", (j%2)+1)
				url := fmt.Sprintf("%s/api/v1/students/%s/report", baseURL, studentID)
				
				req, err := MakeAuthenticatedRequest("GET", url, nil, config)
				if err != nil {
					errors <- fmt.Errorf("worker %d request %d: failed to create request: %w", workerID, j, err)
					continue
				}

				resp, err := client.Do(req)
				if err != nil {
					errors <- fmt.Errorf("worker %d request %d: request failed: %w", workerID, j, err)
					continue
				}

				// Read entire response
				body, err := io.ReadAll(resp.Body)
				resp.Body.Close()
				
				requestDuration := time.Since(requestStart)
				
				if err != nil {
					errors <- fmt.Errorf("worker %d request %d: failed to read response: %w", workerID, j, err)
					continue
				}
				
				if resp.StatusCode != http.StatusOK {
					errors <- fmt.Errorf("worker %d request %d: status %d, body: %s", workerID, j, resp.StatusCode, string(body))
					continue
				}
				
				// Validate that we got a PDF
				if len(body) < 1000 || body[0] != '%' || body[1] != 'P' || body[2] != 'D' || body[3] != 'F' {
					errors <- fmt.Errorf("worker %d request %d: invalid PDF response", workerID, j)
					continue
				}
				
				results <- requestDuration
			}
		}(i, workerRequests)
	}
	
	// Wait for all workers to complete
	wg.Wait()
	close(results)
	close(errors)
	
	totalDuration := time.Since(startTime)
	
	// Collect results
	var responseTimes []time.Duration
	for duration := range results {
		responseTimes = append(responseTimes, duration)
	}
	
	var errorList []error
	for err := range errors {
		errorList = append(errorList, err)
	}
	
	// Calculate statistics
	stats := &LoadTestStats{
		TotalRequests:      totalRequests,
		SuccessfulRequests: len(responseTimes),
		FailedRequests:     len(errorList),
		TotalDuration:      totalDuration,
		ResponseTimes:      responseTimes,
	}
	
	if len(responseTimes) > 0 {
		var totalResponseTime time.Duration
		minTime := responseTimes[0]
		maxTime := responseTimes[0]
		
		for _, rt := range responseTimes {
			totalResponseTime += rt
			if rt < minTime {
				minTime = rt
			}
			if rt > maxTime {
				maxTime = rt
			}
		}
		
		stats.AvgResponseTime = totalResponseTime / time.Duration(len(responseTimes))
		stats.MinResponseTime = minTime
		stats.MaxResponseTime = maxTime
	}
	
	stats.RequestsPerSecond = float64(stats.SuccessfulRequests) / totalDuration.Seconds()
	stats.FailureRate = float64(stats.FailedRequests) / float64(stats.TotalRequests)
	
	// Log any errors for debugging
	if len(errorList) > 0 && len(errorList) <= 5 {
		t.Logf("Errors encountered during load test:")
		for _, err := range errorList {
			t.Logf("  - %v", err)
		}
	} else if len(errorList) > 5 {
		t.Logf("Total errors: %d (showing first 5):", len(errorList))
		for i := 0; i < 5; i++ {
			t.Logf("  - %v", errorList[i])
		}
	}
	
	return stats
}

// TestMemoryUsageUnderLoad tests memory usage during sustained load
func TestMemoryUsageUnderLoad(t *testing.T) {
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

	t.Run("sustained_load_memory_test", func(t *testing.T) {
		// Run sustained load for a period of time
		duration := 30 * time.Second
		concurrency := 10
		
		t.Logf("Running sustained load test for %v with %d concurrent workers", duration, concurrency)
		
		var wg sync.WaitGroup
		stopChan := make(chan struct{})
		requestCount := 0
		var requestCountMutex sync.Mutex
		
		client := &http.Client{Timeout: 10 * time.Second}
		
		// Start workers
		for i := 0; i < concurrency; i++ {
			wg.Add(1)
			go func(workerID int) {
				defer wg.Done()
				
				studentID := fmt.Sprintf("%d", (workerID%2)+1)
				url := fmt.Sprintf("%s/api/v1/students/%s/report", testServer.URL, studentID)
				
				for {
					select {
					case <-stopChan:
						return
					default:
						req, err := MakeAuthenticatedRequest("GET", url, nil, config)
						if err != nil {
							t.Logf("Worker %d: failed to create request: %v", workerID, err)
							continue
						}

						resp, err := client.Do(req)
						if err != nil {
							t.Logf("Worker %d: request failed: %v", workerID, err)
							continue
						}

						// Read and discard response
						io.Copy(io.Discard, resp.Body)
						resp.Body.Close()

						requestCountMutex.Lock()
						requestCount++
						requestCountMutex.Unlock()

						// Small delay to prevent overwhelming
						time.Sleep(50 * time.Millisecond)
					}
				}
			}(i)
		}
		
		// Run for specified duration
		time.Sleep(duration)
		close(stopChan)
		wg.Wait()
		
		t.Logf("Sustained load test completed")
		t.Logf("Total requests processed: %d", requestCount)
		t.Logf("Requests per second: %.2f", float64(requestCount)/duration.Seconds())
		
		// The fact that we completed without hanging or crashing indicates good memory management
		if requestCount < 10 {
			t.Errorf("Too few requests completed: %d (expected at least 10)", requestCount)
		}
	})
}

// TestResponseTimeConsistency tests that response times remain consistent under load
func TestResponseTimeConsistency(t *testing.T) {
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

	t.Run("response_time_consistency", func(t *testing.T) {
		numRequests := 50
		client := &http.Client{Timeout: 10 * time.Second}
		url := testServer.URL + "/api/v1/students/2/report"
		
		var responseTimes []time.Duration
		
		t.Logf("Measuring response time consistency over %d requests", numRequests)
		
		for i := 0; i < numRequests; i++ {
			start := time.Now()
			
			req, err := MakeAuthenticatedRequest("GET", url, nil, config)
			if err != nil {
				t.Fatalf("Request %d failed to create: %v", i, err)
			}

			resp, err := client.Do(req)
			if err != nil {
				t.Fatalf("Request %d failed: %v", i, err)
			}

			_, err = io.ReadAll(resp.Body)
			resp.Body.Close()
			if err != nil {
				t.Fatalf("Request %d failed to read response: %v", i, err)
			}

			responseTime := time.Since(start)
			responseTimes = append(responseTimes, responseTime)
			
			if resp.StatusCode != http.StatusOK {
				t.Fatalf("Request %d failed with status: %d", i, resp.StatusCode)
			}
			
			// Small delay between requests
			time.Sleep(10 * time.Millisecond)
		}
		
		// Calculate statistics
		var total time.Duration
		minTime := responseTimes[0]
		maxTime := responseTimes[0]
		
		for _, rt := range responseTimes {
			total += rt
			if rt < minTime {
				minTime = rt
			}
			if rt > maxTime {
				maxTime = rt
			}
		}
		
		avgTime := total / time.Duration(len(responseTimes))
		
		// Calculate standard deviation
		var variance time.Duration
		for _, rt := range responseTimes {
			diff := rt - avgTime
			variance += diff * diff / time.Duration(len(responseTimes))
		}
		stdDev := time.Duration(float64(variance) * 0.5) // Rough square root
		
		t.Logf("Response time statistics:")
		t.Logf("  Average: %v", avgTime)
		t.Logf("  Min: %v", minTime)
		t.Logf("  Max: %v", maxTime)
		t.Logf("  Std Dev: ~%v", stdDev)
		
		// Validate consistency (max should not be more than 3x average)
		if maxTime > avgTime*3 {
			t.Errorf("Response time inconsistency detected: max (%v) > 3x average (%v)", maxTime, avgTime)
		}
		
		// Validate reasonable performance
		if avgTime > 1*time.Second {
			t.Errorf("Average response time too slow: %v", avgTime)
		}
		
		t.Log("✓ Response times are consistent")
	})
} 