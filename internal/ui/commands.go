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
			// Get GitHub repo from env or use default
			githubRepo := os.Getenv("GITHUB_REPO")
			if githubRepo == "" {
				githubRepo = "identity-manage/account"
			}

			// Fetch PR information (branch, author, repository)
			if build.GitBranch == "" || build.PRAuthor == "" || build.Repository == "" {
				prInfo, err := github.FetchPRBranch(token, githubRepo, prNumber)
				if err == nil {
					if build.GitBranch == "" && prInfo.BranchName != "" {
						build.GitBranch = prInfo.BranchName
					}
					if build.PRAuthor == "" && prInfo.Author != "" {
						build.PRAuthor = prInfo.Author
					}
					if build.Repository == "" && prInfo.Repository != "" {
						build.Repository = prInfo.Repository
					}
				}
			}

			// Fetch PR check status
			checkStatus := github.FetchPRCheckStatus(token, githubRepo, prNumber)
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

// fetchBuildCmd is used for refresh - preserves Git branch but refreshes PR check status
func fetchBuildCmd(client Client, prNumber string, index int, existingGitBranch, existingPRCheckStatus string) tea.Cmd {
	return func() tea.Msg {
		jobPath := jenkins.InferJobPath(prNumber)
		branch := "PR-" + prNumber
		build, err := client.GetBuildStatus(jobPath, branch, 0)

		if build != nil {
			// Preserve the Git branch (doesn't change often)
			if existingGitBranch != "" {
				build.GitBranch = existingGitBranch
			}

			// Re-fetch PR check status for real-time updates
			token := os.Getenv("GITHUB_TOKEN")
			if token != "" {
				githubRepo := os.Getenv("GITHUB_REPO")
				if githubRepo == "" {
					githubRepo = "identity-manage/account"
				}

				checkStatus := github.FetchPRCheckStatus(token, githubRepo, prNumber)
				build.PRCheckStatus = checkStatus.Summary
			}
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
