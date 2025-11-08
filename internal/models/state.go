package models

// DashboardState manages the list of builds and selection
type DashboardState struct {
	Builds        []Build
	SelectedIndex int
	GridColumns   int // Number of columns in the grid layout
}

// AddBuild adds a build to the list
func (s *DashboardState) AddBuild(build Build) {
	s.Builds = append(s.Builds, build)
}

// RemoveBuild removes a build at the specified index
// Returns true if successful, false if index is invalid
func (s *DashboardState) RemoveBuild(index int) bool {
	if index < 0 || index >= len(s.Builds) {
		return false
	}

	s.Builds = append(s.Builds[:index], s.Builds[index+1:]...)

	// Adjust selected index if needed
	if s.SelectedIndex >= len(s.Builds) && len(s.Builds) > 0 {
		s.SelectedIndex = len(s.Builds) - 1
	}

	return true
}

// GetSelectedBuild returns the currently selected build
// Returns nil if no build is selected or list is empty
func (s *DashboardState) GetSelectedBuild() *Build {
	if s.SelectedIndex < 0 || s.SelectedIndex >= len(s.Builds) {
		return nil
	}
	return &s.Builds[s.SelectedIndex]
}

// MoveSelection moves the selection in the specified direction
// Directions: "up", "down", "left", "right"
func (s *DashboardState) MoveSelection(direction string) {
	if len(s.Builds) == 0 || s.GridColumns == 0 {
		return
	}

	currentRow := s.SelectedIndex / s.GridColumns
	currentCol := s.SelectedIndex % s.GridColumns
	totalRows := (len(s.Builds) + s.GridColumns - 1) / s.GridColumns

	var newRow, newCol int

	switch direction {
	case "up":
		newRow = currentRow - 1
		if newRow < 0 {
			return // Stay at current position
		}
		newCol = currentCol

	case "down":
		newRow = currentRow + 1
		if newRow >= totalRows {
			return // Stay at current position
		}
		newCol = currentCol

	case "left":
		if currentCol == 0 {
			return // Stay at current position
		}
		newRow = currentRow
		newCol = currentCol - 1

	case "right":
		if currentCol >= s.GridColumns-1 {
			return // Stay at current position
		}
		newRow = currentRow
		newCol = currentCol + 1

	default:
		return
	}

	// Calculate new index
	newIndex := newRow*s.GridColumns + newCol

	// Verify new index is valid
	if newIndex >= 0 && newIndex < len(s.Builds) {
		s.SelectedIndex = newIndex
	}
}
