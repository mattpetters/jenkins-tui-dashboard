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
	
	// Fetch Check Runs (GitHub Actions, GitHub Apps)
	checkRunsURL := fmt.Sprintf("%s/repos/%s/commits/%s/check-runs", githubAPIBase, repo, sha)
	checkRunsReq, err := http.NewRequest("GET", checkRunsURL, nil)
	if err != nil {
		return CheckStatus{}, err
	}
	if token != "" {
		checkRunsReq.Header.Set("Authorization", "Bearer "+token)
	}
	checkRunsReq.Header.Set("Accept", "application/vnd.github.v3+json")
	
	checkRunsResp, err := client.Do(checkRunsReq)
	if err != nil {
		return CheckStatus{}, err
	}
	defer checkRunsResp.Body.Close()
	
	var checkRunsResult struct {
		TotalCount int `json:"total_count"`
		CheckRuns  []struct {
			Status     string `json:"status"`
			Conclusion string `json:"conclusion"`
		} `json:"check_runs"`
	}
	
	if checkRunsResp.StatusCode == http.StatusOK {
		if err := json.NewDecoder(checkRunsResp.Body).Decode(&checkRunsResult); err != nil {
			return CheckStatus{}, err
		}
	}
	
	// Fetch Commit Statuses (traditional CI/CD status checks)
	statusesURL := fmt.Sprintf("%s/repos/%s/commits/%s/statuses", githubAPIBase, repo, sha)
	statusesReq, err := http.NewRequest("GET", statusesURL, nil)
	if err != nil {
		return CheckStatus{}, err
	}
	if token != "" {
		statusesReq.Header.Set("Authorization", "Bearer "+token)
	}
	statusesReq.Header.Set("Accept", "application/vnd.github.v3+json")
	
	statusesResp, err := client.Do(statusesReq)
	if err != nil {
		return CheckStatus{}, err
	}
	defer statusesResp.Body.Close()
	
	var statuses []struct {
		State   string `json:"state"`
		Context string `json:"context"`
	}
	
	if statusesResp.StatusCode == http.StatusOK {
		if err := json.NewDecoder(statusesResp.Body).Decode(&statuses); err != nil {
			return CheckStatus{}, err
		}
	}
	
	// Combine and count all checks
	total := 0
	passed := 0
	failed := 0
	completed := 0
	
	// Count Check Runs
	for _, check := range checkRunsResult.CheckRuns {
		total++
		if check.Status == "completed" {
			completed++
			if check.Conclusion == "success" {
				passed++
			} else if check.Conclusion == "failure" {
				failed++
			}
		}
	}
	
	// Count Commit Statuses (deduplicate by context - only latest per context)
	seenContexts := make(map[string]bool)
	for _, status := range statuses {
		// Skip if we've already seen this context (API returns newest first)
		if seenContexts[status.Context] {
			continue
		}
		seenContexts[status.Context] = true
		
		total++
		// Commit statuses have different state values: success, pending, failure, error
		if status.State == "success" {
			completed++
			passed++
		} else if status.State == "failure" || status.State == "error" {
			completed++
			failed++
		}
		// "pending" means not completed, so don't increment completed
	}
	
	// Build summary based on state
	var summary string
	if total == 0 {
		summary = "no checks"
	} else if completed == total {
		// All checks are completed
		if failed == 0 {
			summary = "all passed"
		} else if failed == 1 {
			summary = "1 failed"
		} else {
			summary = fmt.Sprintf("%d failed", failed)
		}
	} else {
		// Some checks still in progress
		if failed > 0 {
			// Has failures and still running
			summary = fmt.Sprintf("%d failed, %d/%d done", failed, completed, total)
		} else if passed > 0 {
			// Has passes, no failures yet, still running
			summary = fmt.Sprintf("%d/%d passing", passed, total)
		} else {
			// Nothing completed yet
			summary = fmt.Sprintf("%d/%d done", completed, total)
		}
	}
	
	return CheckStatus{
		TotalChecks:  total,
		PassedChecks: passed,
		FailedChecks: failed,
		Summary:      summary,
	}, nil
}

