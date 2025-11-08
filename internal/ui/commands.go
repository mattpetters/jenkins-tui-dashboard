package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/mpetters/jenkins-dash/internal/browser"
	"github.com/mpetters/jenkins-dash/internal/jenkins"
	"github.com/mpetters/jenkins-dash/internal/models"
)

// buildFetchedMsg is sent when a build fetch completes (success or error)
type buildFetchedMsg struct {
	index int
	build *models.Build
	err   error
}

// fetchBuildCmd creates a command to fetch build data from Jenkins
// Also fetches Git branch name from GitHub if available
func fetchBuildCmd(client Client, prNumber string, index int) tea.Cmd {
	return func() tea.Msg {
		jobPath := jenkins.InferJobPath(prNumber)
		branch := "PR-" + prNumber
		build, err := client.GetBuildStatus(jobPath, branch, 0)
		
		// TODO: Fetch Git branch name from GitHub API and set build.GitBranch
		// For now, Git branch comes from Jenkins API (might be "master" after merge)
		
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
