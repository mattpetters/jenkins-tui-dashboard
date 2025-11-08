package github

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const (
	githubAPIBase = "https://github.intuit.com/api/v3"
	defaultRepo   = "identity-manage/account"
)

// PRInfo contains PR information from GitHub
type PRInfo struct {
	BranchName string
	Title      string
}

// FetchPRBranch fetches the Git branch name for a PR from GitHub
func FetchPRBranch(token, repo, prNumber string) (string, error) {
	if repo == "" {
		repo = defaultRepo
	}

	client := &http.Client{Timeout: 3 * time.Second} // Short timeout
	
	// GitHub API: GET /repos/{owner}/{repo}/pulls/{pull_number}
	url := fmt.Sprintf("%s/repos/intuit/%s/pulls/%s", githubAPIBase, repo, prNumber)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		// Return empty string instead of error - branch name is optional
		return "", nil
	}

	var pr struct {
		Head struct {
			Ref string `json:"ref"` // This is the branch name
		} `json:"head"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&pr); err != nil {
		return "", nil
	}

	return pr.Head.Ref, nil
}

