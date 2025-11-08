package ui

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/mpetters/jenkins-dash/internal/models"
)

const (
	tileWidthWithPadding = 40 // 36 + 4 for more spacing between tiles
	minColumns           = 1
	maxColumns           = 4
)

// CalculateGridColumns determines how many columns to use based on terminal width
func CalculateGridColumns(terminalWidth int) int {
	columns := terminalWidth / tileWidthWithPadding
	if columns < minColumns {
		return minColumns
	}
	if columns > maxColumns {
		return maxColumns
	}
	return columns
}

// RenderGrid renders all build tiles in a grid layout
func RenderGrid(builds []models.Build, selectedIndex int, columns int, blinkState bool) string {
	if len(builds) == 0 {
		return lipgloss.NewStyle().
			Foreground(lipgloss.Color("#888888")).
			Render("No builds yet. Press 'a' to add a PR build.")
	}

	var rows []string
	var currentRow []string

	for i, build := range builds {
		isSelected := i == selectedIndex
		tile := RenderTile(build, isSelected)

		currentRow = append(currentRow, tile)

		// Start new row after filling columns
		if len(currentRow) == columns || i == len(builds)-1 {
			// Join tiles in this row with spacing
			rowStr := lipgloss.JoinHorizontal(lipgloss.Top, currentRow...)
			rows = append(rows, rowStr)
			currentRow = []string{}
		}
	}

	// Join all rows vertically with spacing between rows
	return lipgloss.JoinVertical(lipgloss.Left, rows...)
}
