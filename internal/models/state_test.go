package models

import (
	"testing"
)

// Test 5: RED - DashboardState build list management
func TestDashboardState_AddBuild(t *testing.T) {
	state := DashboardState{}

	// Initially empty
	if len(state.Builds) != 0 {
		t.Errorf("Expected empty builds list, got %d", len(state.Builds))
	}

	// Add first build
	build1 := Build{PRNumber: "3859", Status: StatusSuccess}
	state.AddBuild(build1)

	if len(state.Builds) != 1 {
		t.Errorf("Expected 1 build, got %d", len(state.Builds))
	}
	if state.Builds[0].PRNumber != "3859" {
		t.Errorf("Expected PR-3859, got PR-%s", state.Builds[0].PRNumber)
	}

	// Add second build
	build2 := Build{PRNumber: "3860", Status: StatusFailure}
	state.AddBuild(build2)

	if len(state.Builds) != 2 {
		t.Errorf("Expected 2 builds, got %d", len(state.Builds))
	}
}

func TestDashboardState_RemoveBuild(t *testing.T) {
	state := DashboardState{}
	state.AddBuild(Build{PRNumber: "3859"})
	state.AddBuild(Build{PRNumber: "3860"})
	state.AddBuild(Build{PRNumber: "3861"})

	// Remove middle build
	if !state.RemoveBuild(1) {
		t.Error("RemoveBuild(1) should return true")
	}

	if len(state.Builds) != 2 {
		t.Errorf("Expected 2 builds after removal, got %d", len(state.Builds))
	}

	// Verify correct build was removed
	if state.Builds[0].PRNumber != "3859" {
		t.Errorf("First build should be 3859, got %s", state.Builds[0].PRNumber)
	}
	if state.Builds[1].PRNumber != "3861" {
		t.Errorf("Second build should be 3861, got %s", state.Builds[1].PRNumber)
	}

	// Try to remove invalid index
	if state.RemoveBuild(5) {
		t.Error("RemoveBuild(5) should return false for invalid index")
	}
}

func TestDashboardState_GetSelectedBuild(t *testing.T) {
	state := DashboardState{}
	state.AddBuild(Build{PRNumber: "3859", Status: StatusSuccess})
	state.AddBuild(Build{PRNumber: "3860", Status: StatusFailure})
	state.SelectedIndex = 1

	build := state.GetSelectedBuild()
	if build == nil {
		t.Fatal("GetSelectedBuild() should not return nil")
	}
	if build.PRNumber != "3860" {
		t.Errorf("Expected PR-3860, got PR-%s", build.PRNumber)
	}

	// Test with empty builds
	emptyState := DashboardState{}
	if build := emptyState.GetSelectedBuild(); build != nil {
		t.Error("GetSelectedBuild() should return nil for empty state")
	}
}

// Test 6: RED - Grid navigation (up/down/left/right)
func TestDashboardState_MoveSelection(t *testing.T) {
	// Create a state with 6 builds (2 rows x 3 columns)
	state := DashboardState{GridColumns: 3}
	for i := 0; i < 6; i++ {
		state.AddBuild(Build{PRNumber: string(rune('0' + i))})
	}

	// Grid layout:
	// [0] [1] [2]
	// [3] [4] [5]

	// Test moving right
	state.SelectedIndex = 0
	state.MoveSelection("right")
	if state.SelectedIndex != 1 {
		t.Errorf("After moving right from 0, expected 1, got %d", state.SelectedIndex)
	}

	// Test moving down
	state.SelectedIndex = 1
	state.MoveSelection("down")
	if state.SelectedIndex != 4 {
		t.Errorf("After moving down from 1, expected 4, got %d", state.SelectedIndex)
	}

	// Test moving left
	state.SelectedIndex = 4
	state.MoveSelection("left")
	if state.SelectedIndex != 3 {
		t.Errorf("After moving left from 4, expected 3, got %d", state.SelectedIndex)
	}

	// Test moving up
	state.SelectedIndex = 3
	state.MoveSelection("up")
	if state.SelectedIndex != 0 {
		t.Errorf("After moving up from 3, expected 0, got %d", state.SelectedIndex)
	}

	// Test edge cases - should stay at current position
	state.SelectedIndex = 0
	state.MoveSelection("up")
	if state.SelectedIndex != 0 {
		t.Errorf("Moving up from top should stay at 0, got %d", state.SelectedIndex)
	}

	state.SelectedIndex = 2
	state.MoveSelection("right")
	if state.SelectedIndex != 2 {
		t.Errorf("Moving right from edge should stay at 2, got %d", state.SelectedIndex)
	}
}
