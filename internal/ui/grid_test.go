package ui

import (
	"strings"
	"testing"

	"github.com/mpetters/jenkins-dash/internal/models"
)

// Test 10: RED - Grid layout calculation and rendering
func TestCalculateGridColumns(t *testing.T) {
	tests := []struct {
		name          string
		terminalWidth int
		want          int
	}{
		{
			name:          "Very narrow terminal",
			terminalWidth: 40,
			want:          1,
		},
		{
			name:          "Medium terminal",
			terminalWidth: 120,
			want:          3,
		},
		{
			name:          "Wide terminal",
			terminalWidth: 160,
			want:          4,
		},
		{
			name:          "Extra wide terminal",
			terminalWidth: 200,
			want:          4, // Max out at 4 columns
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CalculateGridColumns(tt.terminalWidth)
			if got != tt.want {
				t.Errorf("CalculateGridColumns(%d) = %v, want %v", tt.terminalWidth, got, tt.want)
			}
		})
	}
}

func TestRenderGrid(t *testing.T) {
	builds := []models.Build{
		{PRNumber: "3859", Status: models.StatusSuccess, JobName: "job1", BuildNumber: 1},
		{PRNumber: "3860", Status: models.StatusFailure, JobName: "job2", BuildNumber: 2},
		{PRNumber: "3861", Status: models.StatusRunning, JobName: "job3", BuildNumber: 3},
	}

	result := RenderGrid(builds, 0, 3, false)

	// Grid should contain all builds
	if !strings.Contains(result, "PR-3859") {
		t.Error("Grid should contain PR-3859")
	}
	if !strings.Contains(result, "PR-3860") {
		t.Error("Grid should contain PR-3860")
	}
	if !strings.Contains(result, "PR-3861") {
		t.Error("Grid should contain PR-3861")
	}
}

func TestRenderGrid_Empty(t *testing.T) {
	builds := []models.Build{}

	result := RenderGrid(builds, 0, 3, false)

	// Empty grid should still render without crashing
	if result == "" {
		t.Error("Empty grid should render some content (help text or empty state)")
	}
}

func TestRenderGrid_Selection(t *testing.T) {
	builds := []models.Build{
		{PRNumber: "3859", Status: models.StatusSuccess},
		{PRNumber: "3860", Status: models.StatusFailure},
		{PRNumber: "3861", Status: models.StatusRunning},
	}

	// Render with selected index
	result := RenderGrid(builds, 1, 3, false)

	// Grid should render without errors
	if result == "" {
		t.Error("Grid should not be empty")
	}

	// Grid should contain all builds
	if !strings.Contains(result, "PR-3859") {
		t.Error("Grid should contain all builds including selected one")
	}

	// Note: Selection highlighting uses ANSI codes which aren't visible in tests
}
