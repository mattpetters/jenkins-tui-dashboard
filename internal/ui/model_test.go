package ui

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/mpetters/jenkins-dash/internal/models"
)

// Test 11: RED - Bubbletea model initialization
func TestModel_Init(t *testing.T) {
	m := NewModel()

	// Model should initialize with default values
	if m.state == nil {
		t.Error("state should be initialized")
	}
	if m.state.GridColumns == 0 {
		t.Error("GridColumns should be initialized to default value")
	}
	if m.inputMode {
		t.Error("inputMode should be false initially")
	}
	if m.statusMessage == "" {
		t.Error("statusMessage should have default help text")
	}

	// Init should return a valid command
	cmd := m.Init()
	if cmd == nil {
		t.Error("Init() should return a command")
	}
}

func TestModel_Update_AddBuildKey(t *testing.T) {
	m := NewModel()

	// Press 'a' to add build
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}}
	newModel, _ := m.Update(msg)

	// Should enter input mode
	updatedModel := newModel.(Model)
	if !updatedModel.inputMode {
		t.Error("Pressing 'a' should enter input mode")
	}
	if updatedModel.inputValue != "" {
		t.Error("Input value should be empty initially")
	}
	if updatedModel.statusMessage == "" {
		t.Error("Status message should guide user in input mode")
	}
}

func TestModel_Update_NavigationKeys(t *testing.T) {
	m := NewModel()

	// Add some test builds
	m.state.AddBuild(models.Build{PRNumber: "1", Status: models.StatusSuccess})
	m.state.AddBuild(models.Build{PRNumber: "2", Status: models.StatusSuccess})
	m.state.AddBuild(models.Build{PRNumber: "3", Status: models.StatusSuccess})
	m.state.GridColumns = 3
	m.state.SelectedIndex = 0

	// Press right arrow
	msg := tea.KeyMsg{Type: tea.KeyRight}
	newModel, _ := m.Update(msg)
	updatedModel := newModel.(Model)

	if updatedModel.state.SelectedIndex != 1 {
		t.Errorf("Right arrow should move selection from 0 to 1, got %d", updatedModel.state.SelectedIndex)
	}
}

func TestModel_Update_QuitKey(t *testing.T) {
	m := NewModel()

	// Press 'q' to quit
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}}
	_, cmd := m.Update(msg)

	// Should return quit command
	if cmd == nil {
		t.Error("Quit should return a command")
	}
	// Note: We can't easily test if cmd == tea.Quit without executing it
}

// CRITICAL BUG FIX TEST: Verify build is actually added when Enter is pressed
func TestModel_Update_InputSubmit_ActuallyAddsBuild(t *testing.T) {
	m := NewModel()
	
	// Initial state - no builds
	if len(m.state.Builds) != 0 {
		t.Fatalf("Expected 0 builds initially, got %d", len(m.state.Builds))
	}
	
	// Enter input mode by pressing 'a'
	newModel, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
	m = newModel.(Model)
	
	if !m.inputMode {
		t.Fatal("Should be in input mode after pressing 'a'")
	}
	
	// Type "3859"
	for _, r := range "3859" {
		newModel, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}})
		m = newModel.(Model)
	}
	
	if m.inputValue != "3859" {
		t.Fatalf("Input value should be '3859', got '%s'", m.inputValue)
	}
	
	// Press Enter to submit
	newModel, _ = m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = newModel.(Model)
	
	// CRITICAL: Build should now be in state
	if len(m.state.Builds) != 1 {
		t.Errorf("Expected 1 build after Enter, got %d", len(m.state.Builds))
	}
	
	if len(m.state.Builds) > 0 && m.state.Builds[0].PRNumber != "3859" {
		t.Errorf("Expected PR-3859, got PR-%s", m.state.Builds[0].PRNumber)
	}
	
	// Should exit input mode
	if m.inputMode {
		t.Error("Should exit input mode after Enter")
	}
}

// Test 12: RED - View rendering
func TestModel_View(t *testing.T) {
	m := NewModel()

	// View should not be empty
	view := m.View()
	if view == "" {
		t.Error("View should render content, not empty string")
	}

	// View should contain status message
	if !strings.Contains(view, m.statusMessage) {
		t.Error("View should contain status message")
	}
}

func TestModel_View_WithBuilds(t *testing.T) {
	m := NewModel()
	m.state.AddBuild(models.Build{PRNumber: "3859", Status: models.StatusSuccess, JobName: "test-job", BuildNumber: 42})
	m.state.AddBuild(models.Build{PRNumber: "3860", Status: models.StatusFailure, JobName: "test-job2", BuildNumber: 43})

	view := m.View()

	// View should contain build information
	if !strings.Contains(view, "PR-3859") {
		t.Error("View should contain PR-3859")
	}
	if !strings.Contains(view, "PR-3860") {
		t.Error("View should contain PR-3860")
	}
	if !strings.Contains(view, "#42") {
		t.Error("View should contain build number")
	}
}

func TestModel_View_InputMode(t *testing.T) {
	m := NewModel()
	m.inputMode = true
	m.inputValue = "3859"

	view := m.View()

	// View should show input field
	if !strings.Contains(view, "3859") {
		t.Error("View should show input value")
	}
}
