package github

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// CheckStatus represents the status of PR checks
type CheckStatus struct {
	TotalChecks  int
	PassedChecks int
	FailedChecks int
	Summary      string // e.g., "5/8 checks" or "all passed"
}

// FetchPRCheckStatus fetches the check run status for a PR
func FetchPRCheckStatus(token, repo, prNumber string) CheckStatus {
	if repo == "" {
		repo = defaultRepo
	}

	// First, get the PR to find the head SHA
	pr, err := fetchPR(token, repo, prNumber)
	if err != nil {
		return CheckStatus{Summary: "unknown"}
	}

	// Then get check runs for that SHA
	checks, err := fetchCheckRuns(token, repo, pr.HeadSHA)
	if err != nil {
		return CheckStatus{Summary: "unknown"}
	}

	return checks
}

type prInfo struct {
	HeadSHA string
}

func fetchPR(token, repo, prNumber string) (*prInfo, error) {
	client := &http.Client{Timeout: 3 * time.Second}
	
	url := fmt.Sprintf("%s/repos/%s/pulls/%s", githubAPIBase, repo, prNumber)
	
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status %d", resp.StatusCode)
	}
	
	var pr struct {
		Head struct {
			SHA string `json:"sha"`
		} `json:"head"`
	}
	
	if err := json.NewDecoder(resp.Body).Decode(&pr); err != nil {
		return nil, err
	}
	
	return &prInfo{HeadSHA: pr.Head.SHA}, nil
}

func fetchCheckRuns(token, repo, sha string) (CheckStatus, error) {
	client := &http.Client{Timeout: 3 * time.Second}
	
	url := fmt.Sprintf("%s/repos/%s/commits/%s/check-runs", githubAPIBase, repo, sha)
	
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return CheckStatus{}, err
	}
	
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	
	resp, err := client.Do(req)
	if err != nil {
		return CheckStatus{}, err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return CheckStatus{}, fmt.Errorf("status %d", resp.StatusCode)
	}
	
	var result struct {
		TotalCount int `json:"total_count"`
		CheckRuns  []struct {
			Status     string `json:"status"`
			Conclusion string `json:"conclusion"`
		} `json:"check_runs"`
	}
	
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return CheckStatus{}, err
	}
	
	// Count passed/failed checks
	total := result.TotalCount
	passed := 0
	failed := 0
	
	for _, check := range result.CheckRuns {
		if check.Status == "completed" {
			if check.Conclusion == "success" {
				passed++
			} else if check.Conclusion == "failure" {
				failed++
			}
		}
	}
	
	// Build summary
	var summary string
	if total == 0 {
		summary = "no checks"
	} else if passed == total {
		summary = "all passed"
	} else if failed > 0 {
		summary = fmt.Sprintf("%d/%d failed", failed, total)
	} else {
		summary = fmt.Sprintf("%d/%d checks", passed, total)
	}
	
	return CheckStatus{
		TotalChecks:  total,
		PassedChecks: passed,
		FailedChecks: failed,
		Summary:      summary,
	}, nil
}

