package github

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// Test 3: FetchPRBranch should return PRInfo with author, branch, and repository
func TestFetchPRBranch_ReturnsPRInfo(t *testing.T) {
	// Mock GitHub API server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/repos/identity-manage/account/pulls/3859" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{
				"head": {
					"ref": "feature/add-auth"
				},
				"user": {
					"login": "john.doe"
				}
			}`))
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	// Temporarily override the API base
	oldAPIBase := githubAPIBase
	githubAPIBase = server.URL
	defer func() { githubAPIBase = oldAPIBase }()

	prInfo, err := FetchPRBranch("test-token", "identity-manage/account", "3859")

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if prInfo.BranchName != "feature/add-auth" {
		t.Errorf("Expected branch name 'feature/add-auth', got '%s'", prInfo.BranchName)
	}

	if prInfo.Author != "john.doe" {
		t.Errorf("Expected author 'john.doe', got '%s'", prInfo.Author)
	}

	if prInfo.Repository != "identity-manage/account" {
		t.Errorf("Expected repository 'identity-manage/account', got '%s'", prInfo.Repository)
	}
}

// Test 4: FetchPRBranch handles missing author gracefully
func TestFetchPRBranch_MissingAuthor(t *testing.T) {
	// Mock GitHub API server without user field
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"head": {
				"ref": "feature/add-auth"
			}
		}`))
	}))
	defer server.Close()

	oldAPIBase := githubAPIBase
	githubAPIBase = server.URL
	defer func() { githubAPIBase = oldAPIBase }()

	prInfo, err := FetchPRBranch("test-token", "identity-manage/account", "3859")

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Should have empty author but not fail
	if prInfo.Author != "" {
		t.Errorf("Expected empty author for missing user field, got '%s'", prInfo.Author)
	}

	// Should still have branch name
	if prInfo.BranchName != "feature/add-auth" {
		t.Errorf("Expected branch name 'feature/add-auth', got '%s'", prInfo.BranchName)
	}
}

// Test 5: Repository defaults to identity-manage/account when empty
func TestFetchPRBranch_DefaultRepository(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Should use the default repo
		if r.URL.Path == "/repos/identity-manage/account/pulls/3859" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{
				"head": {
					"ref": "feature/add-auth"
				},
				"user": {
					"login": "jane.smith"
				}
			}`))
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	oldAPIBase := githubAPIBase
	githubAPIBase = server.URL
	defer func() { githubAPIBase = oldAPIBase }()

	// Pass empty string for repo - should use default
	prInfo, err := FetchPRBranch("test-token", "", "3859")

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if prInfo.Repository != "identity-manage/account" {
		t.Errorf("Expected default repository 'identity-manage/account', got '%s'", prInfo.Repository)
	}
}

