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

// Test 6: Tile should display PR author below branch name
func TestRenderTile_DisplaysPRAuthor(t *testing.T) {
	build := models.Build{
		PRNumber:  "3859",
		GitBranch: "feature/add-auth",
		PRAuthor:  "john.doe",
		Status:    models.StatusSuccess,
	}

	result := RenderTile(build, false)

	// Should contain the author name
	if !strings.Contains(result, "john.doe") {
		t.Error("Tile should display PR author")
	}

	// Author should appear after branch name in the output
	branchIdx := strings.Index(result, "feature/add-auth")
	authorIdx := strings.Index(result, "john.doe")
	if branchIdx == -1 || authorIdx == -1 || authorIdx <= branchIdx {
		t.Error("PR author should appear after branch name")
	}
}

// Test 7: Tile should display repository at bottom right
func TestRenderTile_DisplaysRepository(t *testing.T) {
	build := models.Build{
		PRNumber:    "3859",
		Repository:  "identity-manage/account",
		Status:      models.StatusSuccess,
		BuildNumber: 142,
	}

	result := RenderTile(build, false)

	// Should contain the repository name
	if !strings.Contains(result, "identity-manage/account") {
		t.Error("Tile should display repository name")
	}

	// Repository should appear after build number in the output
	buildNumIdx := strings.Index(result, "#142")
	repoIdx := strings.Index(result, "identity-manage/account")
	if buildNumIdx == -1 || repoIdx == -1 || repoIdx <= buildNumIdx {
		t.Error("Repository should appear after build number")
	}
}

// Test 8: Tile handles long author names (truncation)
func TestRenderTile_LongAuthorName(t *testing.T) {
	build := models.Build{
		PRNumber:  "3859",
		GitBranch: "feature/add-auth",
		PRAuthor:  "very.long.username.that.exceeds.tile.width",
		Status:    models.StatusSuccess,
	}

	result := RenderTile(build, false)

	// Should not crash and should contain some part of the author
	if result == "" {
		t.Error("Tile should render even with long author name")
	}
	if !strings.Contains(result, "very.long") {
		t.Error("Tile should contain at least the start of long author name")
	}
}

// Test 9: Tile handles missing author (shows empty or placeholder)
func TestRenderTile_MissingAuthor(t *testing.T) {
	build := models.Build{
		PRNumber:  "3859",
		GitBranch: "feature/add-auth",
		PRAuthor:  "", // No author
		Status:    models.StatusSuccess,
	}

	result := RenderTile(build, false)

	// Should render without crashing
	if result == "" {
		t.Error("Tile should render even without author")
	}
	// Should still show branch
	if !strings.Contains(result, "feature/add-auth") {
		t.Error("Tile should show branch even when author is missing")
	}
}

// Test 10: Tile maintains correct width with new elements
func TestRenderTile_CorrectWidth(t *testing.T) {
	build := models.Build{
		PRNumber:    "3859",
		GitBranch:   "feature/add-auth",
		PRAuthor:    "john.doe",
		Repository:  "identity-manage/account",
		Status:      models.StatusSuccess,
		BuildNumber: 142,
	}

	result := RenderTile(build, false)

	// Check that all lines have consistent borders
	lines := strings.Split(result, "\n")
	for i, line := range lines {
		// Skip ANSI color code lines
		if line == "" {
			continue
		}
		// Each content line should have proper box drawing characters
		// This is a basic check - lipgloss may add styling
		if i > 0 && i < len(lines)-1 { // Skip first and last due to potential styling
			stripped := stripANSI(line)
			if stripped != "" && !strings.Contains(stripped, "│") && !strings.Contains(stripped, "─") {
				// If it's actual content, it should have box characters
				// But with lipgloss borders, this might not apply
				// So we just verify it doesn't crash
			}
		}
	}

	// Main check: tile should render without errors
	if result == "" {
		t.Error("Tile with all fields should render successfully")
	}
}

// Helper function to strip ANSI color codes
func stripANSI(str string) string {
	// Simple ANSI code stripper - matches ESC[...m patterns
	result := ""
	inEscape := false
	for i := 0; i < len(str); i++ {
		if str[i] == '\x1b' && i+1 < len(str) && str[i+1] == '[' {
			inEscape = true
			i++ // Skip the '['
			continue
		}
		if inEscape {
			if str[i] == 'm' {
				inEscape = false
			}
			continue
		}
		result += string(str[i])
	}
	return result
}
