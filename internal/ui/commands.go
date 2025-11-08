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
		
		// Fetch Git branch from GitHub if not already set (best effort, non-blocking)
		if build.GitBranch == "" {
			token := os.Getenv("GITHUB_TOKEN")
			if token != "" {
				if gitBranch, _ := github.FetchPRBranch(token, "identity-manage/account", prNumber); gitBranch != "" {
					build.GitBranch = gitBranch
				}
			}
		}
		
		return buildFetchedMsg{
			index: index,
			build: build,
			err:   nil,
		}
	}
}

// fetchBuildCmd is used for refresh - preserves existing Git branch
func fetchBuildCmd(client Client, prNumber string, index int) tea.Cmd {
	return func() tea.Msg {
		jobPath := jenkins.InferJobPath(prNumber)
		branch := "PR-" + prNumber
		build, err := client.GetBuildStatus(jobPath, branch, 0)
		
		// Don't re-fetch Git branch on refresh - it's saved in the persistence file
		// The branch name doesn't change, so we keep the existing value
		
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
