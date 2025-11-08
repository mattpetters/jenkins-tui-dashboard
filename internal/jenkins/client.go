package jenkins

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/mpetters/jenkins-dash/internal/models"
)

// Client handles communication with Jenkins API
type Client struct {
	baseURL    string
	username   string
	token      string
	httpClient *http.Client
}

// NewClient creates a new Jenkins API client
// Jenkins uses Basic Auth with username:token
func NewClient(username, token string) *Client {
	return &Client{
		baseURL:    jenkinsBaseURL,
		username:   username,
		token:      token,
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
}

// GetBuildStatus fetches build status from Jenkins API
func (c *Client) GetBuildStatus(jobPath, branch string, buildNum int) (*models.Build, error) {
	// Build API URL
	url := BuildJenkinsURL(jobPath, branch, buildNum) + "/api/json"

	// Create request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	// Add Basic Auth (Jenkins uses username:token, not Bearer)
	if c.username != "" && c.token != "" {
		req.SetBasicAuth(c.username, c.token)
	}
	req.Header.Set("Accept", "application/json")

	// Make request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("making request: %w", err)
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Jenkins API returned status %d for URL: %s", resp.StatusCode, url)
	}

	// Parse JSON response
	var data map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("parsing JSON: %w", err)
	}

	// Convert to Build struct
	build := ParseBuildResponse(data, branch, jobPath)
	return &build, nil
}
