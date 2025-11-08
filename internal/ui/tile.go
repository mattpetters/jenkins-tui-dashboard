package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/mpetters/jenkins-dash/internal/models"
)

const (
	tileWidth = 32
)

// RenderTile renders a build tile with proper styling
func RenderTile(build models.Build, isSelected bool) string {
	// Get aesthetically pleasing pastel colors
	bgColor, fgColor := GetTileColors(build.Status)

	// Build the tile content
	lines := make([]string, 0, 8)

	// Top border
	lines = append(lines, "┌"+strings.Repeat("─", tileWidth-2)+"┐")

	// PR number (centered)
	prText := fmt.Sprintf("PR-%s", build.PRNumber)
	padding := (tileWidth - 4 - len(prText)) / 2
	prLine := fmt.Sprintf("│ %s%s%s │",
		strings.Repeat(" ", padding),
		prText,
		strings.Repeat(" ", tileWidth-4-len(prText)-padding))
	lines = append(lines, prLine)
	
	// Git branch name (centered, smaller text)
	branchText := build.GitBranch
	if branchText == "" {
		branchText = "PR-" + build.PRNumber // Fallback to PR number
	}
	if len(branchText) > tileWidth-4 {
		branchText = branchText[:tileWidth-4]
	}
	branchPadding := (tileWidth - 4 - len(branchText)) / 2
	branchLine := fmt.Sprintf("│ %s%s%s │",
		strings.Repeat(" ", branchPadding),
		branchText,
		strings.Repeat(" ", tileWidth-4-len(branchText)-branchPadding))
	lines = append(lines, branchLine)

	// Separator
	lines = append(lines, "├"+strings.Repeat("─", tileWidth-2)+"┤")

	// Stage
	stageText := build.Stage
	if stageText == "" {
		if build.Status == models.StatusPending {
			stageText = "Loading..."
		} else {
			stageText = "Unknown"
		}
	}
	if len(stageText) > 18 {
		stageText = stageText[:18]
	}
	stageLine := fmt.Sprintf("│ Stage: %-18s │", stageText)
	lines = append(lines, stageLine)

	// Job
	jobText := build.JobName
	if jobText == "" {
		if build.Status == models.StatusPending {
			jobText = "Fetching data..."
		} else {
			jobText = "Unknown"
		}
	}
	if len(jobText) > 18 {
		jobText = jobText[:18]
	}
	jobLine := fmt.Sprintf("│ Job: %-20s │", jobText)
	lines = append(lines, jobLine)

	// Duration
	durationText := build.FormatDuration()
	timeLine := fmt.Sprintf("│ Time: %-20s │", durationText)
	lines = append(lines, timeLine)

	// Build number (bottom right)
	buildNumText := fmt.Sprintf("#%d", build.BuildNumber)
	if build.BuildNumber == 0 {
		buildNumText = "..."
	}
	buildNumLine := fmt.Sprintf("│ %s%s │",
		strings.Repeat(" ", tileWidth-4-len(buildNumText)),
		buildNumText)
	lines = append(lines, buildNumLine)

	// Bottom border
	lines = append(lines, "└"+strings.Repeat("─", tileWidth-2)+"┘")

	// Join all lines
	content := strings.Join(lines, "\n")

	// Apply styling
	style := lipgloss.NewStyle().
		Foreground(fgColor).
		Background(bgColor)

	if isSelected {
		style = style.Bold(true)
	}

	return style.Render(content)
}
