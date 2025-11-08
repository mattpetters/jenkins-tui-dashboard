package persistence

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/mpetters/jenkins-dash/internal/models"
)

const configFileName = ".jenkins-dash-builds.json"

// GetConfigPath returns the path to the config file in the user's home directory
func GetConfigPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		// Fallback to current directory if can't get home
		return configFileName
	}
	return filepath.Join(homeDir, configFileName)
}

// SaveBuilds saves the builds list to a JSON file
func SaveBuilds(filePath string, builds []models.Build) error {
	// Create directory if it doesn't exist
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	// Convert to JSON
	data, err := json.MarshalIndent(builds, "", "  ")
	if err != nil {
		return err
	}

	// Write to file
	return os.WriteFile(filePath, data, 0644)
}

// LoadBuilds loads builds from a JSON file
// Returns empty list if file doesn't exist (not an error)
func LoadBuilds(filePath string) ([]models.Build, error) {
	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return []models.Build{}, nil
	}

	// Read file
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	// Parse JSON
	var builds []models.Build
	if err := json.Unmarshal(data, &builds); err != nil {
		return nil, err
	}

	return builds, nil
}

