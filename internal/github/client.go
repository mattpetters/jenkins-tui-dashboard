package github

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const githubAPIBase = "https://api.github.com"

// Client handles GitHub API calls
type Client struct {
	token      string
	httpClient *http.Client
}

// NewClient creates a new GitHub API client
func NewClient(token string) *Client {
	return &Client{
		token:      token,
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}
}

// GetPRBranchName fetches the branch name for a PR from GitHub
func (c *Client) GetPRBranchName(repo, prNumber string) (string, error) {
	// GitHub API: GET /repos/{owner}/{repo}/pulls/{pull_number}
	url := fmt.Sprintf("%s/repos/%s/pulls/%s", githubAPIBase, repo, prNumber)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("GitHub API returned status %d", resp.StatusCode)
	}

	var pr struct {
		Head struct {
			Ref string `json:"ref"` // This is the branch name
		} `json:"head"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&pr); err != nil {
		return "", err
	}

	return pr.Head.Ref, nil
}

// GetPRBranchName is a convenience function that uses default repo
func GetPRBranchName(repo, prNumber string) string {
	// For now, return empty - will be fetched async
	// This is a placeholder for the interface
	return ""
}

