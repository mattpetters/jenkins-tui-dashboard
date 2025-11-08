package github

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Test: Fetch PR check status from GitHub
func TestFetchPRCheckStatus(t *testing.T) {
	// Test that function exists and returns check info
	status := FetchPRCheckStatus("", "identity-manage/account", "3934")
	
	// Function should not panic
	_ = status
}

// TestCheckStatusSummary tests the check status summary generation logic
func TestCheckStatusSummary(t *testing.T) {
	tests := []struct {
		name          string
		checkRuns     []map[string]string
		totalCount    int
		expectedMsg   string
		expectedPassed int
		expectedFailed int
	}{
		{
			name: "no checks at all",
			checkRuns: []map[string]string{},
			totalCount: 0,
			expectedMsg: "no checks",
			expectedPassed: 0,
			expectedFailed: 0,
		},
		{
			name: "all checks passed",
			checkRuns: []map[string]string{
				{"status": "completed", "conclusion": "success"},
				{"status": "completed", "conclusion": "success"},
				{"status": "completed", "conclusion": "success"},
			},
			totalCount: 3,
			expectedMsg: "all passed",
			expectedPassed: 3,
			expectedFailed: 0,
		},
		{
			name: "some checks in progress - 4 passed out of 6",
			checkRuns: []map[string]string{
				{"status": "completed", "conclusion": "success"},
				{"status": "completed", "conclusion": "success"},
				{"status": "completed", "conclusion": "success"},
				{"status": "completed", "conclusion": "success"},
				{"status": "in_progress", "conclusion": ""},
				{"status": "queued", "conclusion": ""},
			},
			totalCount: 6,
			expectedMsg: "4/6 passing",
			expectedPassed: 4,
			expectedFailed: 0,
		},
		{
			name: "one failure out of 6 checks",
			checkRuns: []map[string]string{
				{"status": "completed", "conclusion": "success"},
				{"status": "completed", "conclusion": "success"},
				{"status": "completed", "conclusion": "success"},
				{"status": "completed", "conclusion": "success"},
				{"status": "completed", "conclusion": "success"},
				{"status": "completed", "conclusion": "failure"},
			},
			totalCount: 6,
			expectedMsg: "1 failed",
			expectedPassed: 5,
			expectedFailed: 1,
		},
		{
			name: "multiple failures",
			checkRuns: []map[string]string{
				{"status": "completed", "conclusion": "success"},
				{"status": "completed", "conclusion": "failure"},
				{"status": "completed", "conclusion": "failure"},
			},
			totalCount: 3,
			expectedMsg: "2 failed",
			expectedPassed: 1,
			expectedFailed: 2,
		},
		{
			name: "mixed state - some passed, some failed, some in progress",
			checkRuns: []map[string]string{
				{"status": "completed", "conclusion": "success"},
				{"status": "completed", "conclusion": "success"},
				{"status": "completed", "conclusion": "failure"},
				{"status": "in_progress", "conclusion": ""},
				{"status": "queued", "conclusion": ""},
			},
			totalCount: 5,
			expectedMsg: "1 failed, 3/5 done",
			expectedPassed: 2,
			expectedFailed: 1,
		},
		{
			name: "all checks in progress - none completed yet",
			checkRuns: []map[string]string{
				{"status": "in_progress", "conclusion": ""},
				{"status": "queued", "conclusion": ""},
				{"status": "queued", "conclusion": ""},
			},
			totalCount: 3,
			expectedMsg: "0/3 done",
			expectedPassed: 0,
			expectedFailed: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a mock server
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Mock PR endpoint
				if r.URL.Path == "/repos/test-owner/test-repo/pulls/123" {
					response := map[string]interface{}{
						"head": map[string]string{
							"sha": "abc123",
						},
					}
					json.NewEncoder(w).Encode(response)
					return
				}
				
				// Mock check runs endpoint
				if r.URL.Path == "/repos/test-owner/test-repo/commits/abc123/check-runs" {
					response := map[string]interface{}{
						"total_count": tt.totalCount,
						"check_runs": tt.checkRuns,
					}
					json.NewEncoder(w).Encode(response)
					return
				}
				
				http.NotFound(w, r)
			}))
			defer server.Close()

			// Temporarily replace the GitHub API base
			oldBase := githubAPIBase
			githubAPIBase = server.URL
			defer func() { githubAPIBase = oldBase }()

			// Fetch check status
			status := FetchPRCheckStatus("test-token", "test-owner/test-repo", "123")

			// Verify results
			if status.Summary != tt.expectedMsg {
				t.Errorf("Expected summary '%s', got '%s'", tt.expectedMsg, status.Summary)
			}
			if status.PassedChecks != tt.expectedPassed {
				t.Errorf("Expected %d passed checks, got %d", tt.expectedPassed, status.PassedChecks)
			}
			if status.FailedChecks != tt.expectedFailed {
				t.Errorf("Expected %d failed checks, got %d", tt.expectedFailed, status.FailedChecks)
			}
			if status.TotalChecks != tt.totalCount {
				t.Errorf("Expected %d total checks, got %d", tt.totalCount, status.TotalChecks)
			}
		})
	}
}

// TestCheckStatusCombinesCheckRunsAndStatuses tests that we fetch from both APIs
func TestCheckStatusCombinesCheckRunsAndStatuses(t *testing.T) {
	// Create a mock server that returns both check runs and commit statuses
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Mock PR endpoint
		if r.URL.Path == "/repos/test-owner/test-repo/pulls/123" {
			response := map[string]interface{}{
				"head": map[string]string{
					"sha": "abc123",
				},
			}
			json.NewEncoder(w).Encode(response)
			return
		}
		
		// Mock check runs endpoint (e.g., GitHub Actions)
		if r.URL.Path == "/repos/test-owner/test-repo/commits/abc123/check-runs" {
			response := map[string]interface{}{
				"total_count": 3,
				"check_runs": []map[string]string{
					{"status": "completed", "conclusion": "success"},
					{"status": "completed", "conclusion": "success"},
					{"status": "in_progress", "conclusion": ""},
				},
			}
			json.NewEncoder(w).Encode(response)
			return
		}
		
		// Mock commit statuses endpoint (e.g., Jenkins, CircleCI external checks)
		if r.URL.Path == "/repos/test-owner/test-repo/commits/abc123/statuses" {
			response := []map[string]string{
				{"state": "success", "context": "ci/jenkins"},
				{"state": "failure", "context": "ci/external"},
				{"state": "pending", "context": "ci/codecov"},
			}
			json.NewEncoder(w).Encode(response)
			return
		}
		
		http.NotFound(w, r)
	}))
	defer server.Close()

	// Temporarily replace the GitHub API base
	oldBase := githubAPIBase
	githubAPIBase = server.URL
	defer func() { githubAPIBase = oldBase }()

	// Fetch check status (should combine both APIs)
	status := FetchPRCheckStatus("test-token", "test-owner/test-repo", "123")

	// Verify results: 3 check runs + 3 statuses = 6 total
	expectedTotal := 6
	if status.TotalChecks != expectedTotal {
		t.Errorf("Expected %d total checks (3 check runs + 3 statuses), got %d", expectedTotal, status.TotalChecks)
	}

	// 2 passed from check runs + 1 passed from statuses = 3 passed
	expectedPassed := 3
	if status.PassedChecks != expectedPassed {
		t.Errorf("Expected %d passed checks, got %d", expectedPassed, status.PassedChecks)
	}

	// 1 failed from statuses
	expectedFailed := 1
	if status.FailedChecks != expectedFailed {
		t.Errorf("Expected %d failed checks, got %d", expectedFailed, status.FailedChecks)
	}

	// Should show "1 failed, 4/6 done" (2 passed + 1 passed + 1 failed = 4 completed out of 6)
	expectedSummary := "1 failed, 4/6 done"
	if status.Summary != expectedSummary {
		t.Errorf("Expected summary '%s', got '%s'", expectedSummary, status.Summary)
	}

	t.Logf("✓ Successfully combined check runs and commit statuses: %s", status.Summary)
}

