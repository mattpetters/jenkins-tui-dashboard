package ui

import (
	"encoding/json"
	"fmt"
	"os"
)

// saveState persists the current builds to disk
func (m Model) saveState() error {
	if m.configPath == "" {
		return nil // No config path set, skip saving
	}

	data, err := json.MarshalIndent(m.state.Builds, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(m.configPath, data, 0644)
}

// LoadPersistedBuilds loads builds from disk (public method for startup)
func (m *Model) LoadPersistedBuilds() error {
	if m.configPath == "" {
		return nil // No config path set, skip loading
	}

	// Check if file exists
	if _, err := os.Stat(m.configPath); os.IsNotExist(err) {
		return nil // File doesn't exist yet, not an error
	}

	data, err := os.ReadFile(m.configPath)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(data, &m.state.Builds); err != nil {
		return err
	}

	// Update status message if builds were loaded
	if len(m.state.Builds) > 0 {
		m.statusMessage = fmt.Sprintf("Loaded %d saved build(s). Press 'a' to add more.", len(m.state.Builds))
	}

	return nil
}

