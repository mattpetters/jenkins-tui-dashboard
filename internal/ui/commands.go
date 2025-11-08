package ui

import (
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
		
		// Fetch Git branch from GitHub (non-blocking, best effort)
		if build != nil {
			// Try to fetch Git branch, but don't fail if it doesn't work
			token := os.Getenv("GITHUB_TOKEN")
			if token != "" {
				if gitBranch, _ := github.FetchPRBranch(token, "identity-manage/account", prNumber); gitBranch != "" {
					build.GitBranch = gitBranch
				}
			}
			// If GitHub fetch failed, GitBranch stays empty and tile shows "PR-3934" as fallback
		}
		
		return buildFetchedMsg{
			index: index,
			build: build,
			err:   err,
		}
	}
}

// fetchBuildCmd is a simpler version for refresh (doesn't re-fetch Git branch)
func fetchBuildCmd(client Client, prNumber string, index int) tea.Cmd {
	return func() tea.Msg {
		jobPath := jenkins.InferJobPath(prNumber)
		branch := "PR-" + prNumber
		build, err := client.GetBuildStatus(jobPath, branch, 0)
		
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
