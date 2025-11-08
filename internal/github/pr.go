package github

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

var (
	githubAPIBase = "https://github.intuit.com/api/v3"
	defaultRepo   = "identity-manage/account"
)

// PRInfo contains PR information from GitHub
type PRInfo struct {
	BranchName string
	Author     string
	Repository string
	Title      string
}

// FetchPRBranch fetches PR information including branch name, author, and repository from GitHub
func FetchPRBranch(token, repo, prNumber string) (PRInfo, error) {
	if repo == "" {
		repo = defaultRepo
	}

	client := &http.Client{Timeout: 3 * time.Second} // Short timeout
	
	// GitHub API: GET /repos/{owner}/{repo}/pulls/{pull_number}
	// For identity-manage/account, owner=identity-manage, repo=account
	url := fmt.Sprintf("%s/repos/%s/pulls/%s", githubAPIBase, repo, prNumber)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return PRInfo{Repository: repo}, err
	}

	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := client.Do(req)
	if err != nil {
		return PRInfo{Repository: repo}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		// Return empty PRInfo with repo instead of error - PR data is optional
		return PRInfo{Repository: repo}, nil
	}

	var pr struct {
		Head struct {
			Ref string `json:"ref"` // This is the branch name
		} `json:"head"`
		User struct {
			Login string `json:"login"` // PR author username
		} `json:"user"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&pr); err != nil {
		return PRInfo{Repository: repo}, nil
	}

	return PRInfo{
		BranchName: pr.Head.Ref,
		Author:     pr.User.Login,
		Repository: repo,
	}, nil
}

