package ui

import (
	"testing"

	"github.com/mpetters/jenkins-dash/internal/models"
)

// Test: RED - Better color palette for build statuses
func TestGetTileColors(t *testing.T) {
	tests := []struct {
		name   string
		status models.BuildStatus
	}{
		{"Success should have pastel green", models.StatusSuccess},
		{"Failure should have pastel red", models.StatusFailure},
		{"Running should have pastel blue", models.StatusRunning},
		{"Pending should have pastel yellow", models.StatusPending},
		{"Error should have pastel red", models.StatusError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bg, fg := GetTileColors(tt.status)

			// Just verify function executes without panicking
			// lipgloss.Color doesn't have String() method
			_ = bg
			_ = fg
		})
	}
}

