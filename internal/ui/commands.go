package ui

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/mpetters/jenkins-dash/internal/browser"
	"github.com/mpetters/jenkins-dash/internal/github"
	"github.com/mpetters/jenkins-dash/internal/jenkins"
	"github.com/mpetters/jenkins-dash/internal/models"
)

// buildFetchedMsg is sent when a build fetch completes (success or error)
type buildFetchedMsg struct {
	index int
	build *models.Build
	err   error
}

// fetchBuildAndBranchCmd fetches both Jenkins build data and GitHub branch name
func fetchBuildAndBranchCmd(client Client, prNumber string, index int) tea.Cmd {
	return func() tea.Msg {
		jobPath := jenkins.InferJobPath(prNumber)
		branch := "PR-" + prNumber
		build, err := client.GetBuildStatus(jobPath, branch, 0)
		
		if err != nil {
			return buildFetchedMsg{index: index, build: nil, err: err}
		}
		
		if build == nil {
			return buildFetchedMsg{index: index, build: nil, err: fmt.Errorf("no build data returned")}
		}
		
		// Fetch Git branch and PR check status from GitHub
		if build.GitBranch == "" || build.PRCheckStatus == "" {
			token := os.Getenv("GITHUB_TOKEN")
			if token != "" {
				// Fetch Git branch
				if build.GitBranch == "" {
					gitBranch, err := github.FetchPRBranch(token, "identity-manage/account", prNumber)
					if err == nil && gitBranch != "" {
						build.GitBranch = gitBranch
					}
				}
				
				// Fetch PR check status
				checkStatus := github.FetchPRCheckStatus(token, "identity-manage/account", prNumber)
				build.PRCheckStatus = checkStatus.Summary
			}
		}
		
		return buildFetchedMsg{
			index: index,
			build: build,
			err:   nil,
		}
	}
}

// fetchBuildCmd is used for refresh - MUST preserve existing Git branch
func fetchBuildCmd(client Client, prNumber string, index int, existingGitBranch string) tea.Cmd {
	return func() tea.Msg {
		jobPath := jenkins.InferJobPath(prNumber)
		branch := "PR-" + prNumber
		build, err := client.GetBuildStatus(jobPath, branch, 0)
		
		// Preserve the Git branch from persistence file
		// Don't re-fetch on every refresh - it doesn't change
		if build != nil && existingGitBranch != "" {
			build.GitBranch = existingGitBranch
		}
		
		return buildFetchedMsg{
			index: index,
			build: build,
			err:   err,
		}
	}
}

// urlOpenedMsg is sent after attempting to open a URL
type urlOpenedMsg struct {
	url string
	err error
}

// openURLCmd creates a command to open a URL in the browser
func openURLCmd(url string) tea.Cmd {
	return func() tea.Msg {
		err := browser.OpenURL(url)
		return urlOpenedMsg{
			url: url,
			err: err,
		}
	}
}
