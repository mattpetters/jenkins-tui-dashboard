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
// Makes TWO calls: /api/json for basic info, /wfapi/describe for stages
func (c *Client) GetBuildStatus(jobPath, branch string, buildNum int) (*models.Build, error) {
	baseURL := BuildJenkinsURL(jobPath, branch, buildNum)
	
	// Call 1: Get basic build info from standard API
	basicData, err := c.fetchJSON(baseURL + "/api/json")
	if err != nil {
		return nil, fmt.Errorf("fetching build info: %w", err)
	}

	// Call 2: Get stages from wfapi (best effort, don't fail if missing)
	stagesData, _ := c.fetchJSON(baseURL + "/wfapi/describe")
	
	// Merge stages into basic data
	if stagesData != nil {
		if stages, ok := stagesData["stages"]; ok {
			basicData["stages"] = stages
		}
	}

	// Convert to Build struct
	build := ParseBuildResponse(basicData, branch, jobPath)
	return &build, nil
}

// fetchJSON is a helper to fetch and parse JSON from Jenkins
func (c *Client) fetchJSON(url string) (map[string]interface{}, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	if c.username != "" && c.token != "" {
		req.SetBasicAuth(c.username, c.token)
	}
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status %d for URL: %s", resp.StatusCode, url)
	}

	var data map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	return data, nil
}
