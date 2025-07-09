package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"go-service/pkg/models"
)

// NodejsClient handles communication with the Node.js backend API
type NodejsClient struct {
	BaseURL    string
	HTTPClient *http.Client
	// Authentication tokens will be added here
	AccessToken string
	CSRFToken   string
}

// NewNodejsClient creates a new client for the Node.js backend API
func NewNodejsClient(baseURL string) *NodejsClient {
	return &NodejsClient{
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Timeout: time.Second * 30,
		},
	}
}

// SetAuthTokens sets the authentication tokens for API requests
func (c *NodejsClient) SetAuthTokens(accessToken, csrfToken string) {
	c.AccessToken = accessToken
	c.CSRFToken = csrfToken
}

// GetStudent fetches a single student by ID from the Node.js API
func (c *NodejsClient) GetStudent(studentID string) (*models.Student, error) {
	url := fmt.Sprintf("%s/api/v1/students/%s", c.BaseURL, studentID)
	
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add authentication headers if tokens are available
	if c.AccessToken != "" {
		req.Header.Set("Cookie", fmt.Sprintf("accessToken=%s", c.AccessToken))
	}
	if c.CSRFToken != "" {
		req.Header.Set("X-CSRF-Token", c.CSRFToken)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var student models.Student
	if err := json.Unmarshal(body, &student); err != nil {
		return nil, fmt.Errorf("failed to unmarshal student data: %w", err)
	}

	return &student, nil
}

// GetStudents fetches all students from the Node.js API (optional, for future use)
func (c *NodejsClient) GetStudents() (models.StudentList, error) {
	url := fmt.Sprintf("%s/api/v1/students", c.BaseURL)
	
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add authentication headers if tokens are available
	if c.AccessToken != "" {
		req.Header.Set("Cookie", fmt.Sprintf("accessToken=%s", c.AccessToken))
	}
	if c.CSRFToken != "" {
		req.Header.Set("X-CSRF-Token", c.CSRFToken)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var students models.StudentList
	if err := json.Unmarshal(body, &students); err != nil {
		return nil, fmt.Errorf("failed to unmarshal students data: %w", err)
	}

	return students, nil
}

// HealthCheck verifies the Node.js API is accessible
func (c *NodejsClient) HealthCheck() error {
	url := fmt.Sprintf("%s/api/v1/dashboard", c.BaseURL)
	
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create health check request: %w", err)
	}

	// Add authentication headers if tokens are available
	if c.AccessToken != "" {
		req.Header.Set("Cookie", fmt.Sprintf("accessToken=%s", c.AccessToken))
	}
	if c.CSRFToken != "" {
		req.Header.Set("X-CSRF-Token", c.CSRFToken)
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("health check request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Node.js API health check failed with status: %d", resp.StatusCode)
	}

	return nil
} 