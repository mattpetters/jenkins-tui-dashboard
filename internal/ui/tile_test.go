package ui

import (
	"strings"
	"testing"

	"github.com/mpetters/jenkins-dash/internal/models"
)

// Test 9: RED - Tile rendering
func TestRenderTile(t *testing.T) {
	build := models.Build{
		PRNumber:        "3859",
		Status:          models.StatusSuccess,
		Stage:           "Deploy",
		JobName:         "maven-build",
		BuildNumber:     142,
		DurationSeconds: 323,
	}

	result := RenderTile(build, false)

	// Verify tile contains key information
	if !strings.Contains(result, "PR-3859") {
		t.Error("Tile should contain PR number")
	}
	if !strings.Contains(result, "Deploy") {
		t.Error("Tile should contain stage")
	}
	if !strings.Contains(result, "maven-build") {
		t.Error("Tile should contain job name")
	}
	if !strings.Contains(result, "#142") {
		t.Error("Tile should contain build number")
	}
	if !strings.Contains(result, "5m 23s") {
		t.Error("Tile should contain formatted duration")
	}

	// Verify tile has box drawing characters
	if !strings.Contains(result, "┌") || !strings.Contains(result, "┐") {
		t.Error("Tile should have top border")
	}
	if !strings.Contains(result, "└") || !strings.Contains(result, "┘") {
		t.Error("Tile should have bottom border")
	}
	if !strings.Contains(result, "│") {
		t.Error("Tile should have side borders")
	}
}

func TestRenderTile_Selected(t *testing.T) {
	build := models.Build{
		PRNumber: "3859",
		Status:   models.StatusSuccess,
	}

	selectedTile := RenderTile(build, true)

	// Selected tile should render without errors and contain content
	if selectedTile == "" {
		t.Error("Selected tile should not be empty")
	}

	// Should still contain the PR number
	if !strings.Contains(selectedTile, "PR-3859") {
		t.Error("Selected tile should contain PR number")
	}

	// Note: lipgloss ANSI codes are only added when output is to a terminal,
	// so we can't reliably test for visual differences in unit tests
}

func TestRenderTile_DifferentStatuses(t *testing.T) {
	statuses := []models.BuildStatus{
		models.StatusPending,
		models.StatusRunning,
		models.StatusSuccess,
		models.StatusFailure,
		models.StatusError,
	}

	for _, status := range statuses {
		t.Run(status.String(), func(t *testing.T) {
			build := models.Build{
				PRNumber: "3859",
				Status:   status,
			}

			result := RenderTile(build, false)

			// Each status should produce different colored output
			// We can't test exact colors without inspecting ANSI codes,
			// but we can verify it doesn't crash and returns content
			if result == "" {
				t.Errorf("Render should produce output for status %s", status)
			}
			if !strings.Contains(result, "PR-3859") {
				t.Errorf("Tile should contain PR number for status %s", status)
			}
		})
	}
}

func TestRenderTile_LoadingState(t *testing.T) {
	build := models.Build{
		PRNumber: "3859",
		Status:   models.StatusPending,
		Stage:    "",
		JobName:  "",
	}

	result := RenderTile(build, false)

	// Loading state should show loading text
	if !strings.Contains(result, "Loading") || !strings.Contains(result, "Fetching") {
		t.Error("Loading state should indicate data is being fetched")
	}
}
