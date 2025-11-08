package ui

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/mpetters/jenkins-dash/internal/models"
)

// Aesthetically pleasing pastel color palette
var (
	// Success - soft green (pastel, easy on eyes)
	colorSuccessBg = lipgloss.Color("#98C379")  // Soft green
	colorSuccessFg = lipgloss.Color("#1A1A1A")  // Dark text for contrast
	
	// Failure - soft red (pastel, not harsh)
	colorFailureBg = lipgloss.Color("#E06C75")  // Soft red
	colorFailureFg = lipgloss.Color("#1A1A1A")  // Dark text
	
	// Running - soft blue (calming pastel)
	colorRunningBg = lipgloss.Color("#61AFEF")  // Soft blue
	colorRunningFg = lipgloss.Color("#1A1A1A")  // Dark text
	
	// Pending - soft yellow/amber (warm pastel)
	colorPendingBg = lipgloss.Color("#E5C07B")  // Soft yellow/amber
	colorPendingFg = lipgloss.Color("#1A1A1A")  // Dark text
	
	// Error - same as failure but could be different
	colorErrorBg = lipgloss.Color("#E06C75")  // Soft red
	colorErrorFg = lipgloss.Color("#1A1A1A")  // Dark text
)

// GetTileColors returns the background and foreground colors for a build status
func GetTileColors(status models.BuildStatus) (bg lipgloss.Color, fg lipgloss.Color) {
	switch status {
	case models.StatusSuccess:
		return colorSuccessBg, colorSuccessFg
	case models.StatusFailure:
		return colorFailureBg, colorFailureFg
	case models.StatusRunning:
		return colorRunningBg, colorRunningFg
	case models.StatusPending:
		return colorPendingBg, colorPendingFg
	case models.StatusError:
		return colorErrorBg, colorErrorFg
	default:
		return colorPendingBg, colorPendingFg
	}
}

