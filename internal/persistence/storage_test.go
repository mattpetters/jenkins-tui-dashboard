package persistence

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/mpetters/jenkins-dash/internal/models"
)

// Test: RED - Save builds to file
func TestSaveBuilds(t *testing.T) {
	// Create temp directory for test
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test-builds.json")

	builds := []models.Build{
		{PRNumber: "3859", Status: models.StatusSuccess, JobName: "test-job", BuildNumber: 42},
		{PRNumber: "3860", Status: models.StatusFailure, JobName: "test-job2", BuildNumber: 43},
	}

	err := SaveBuilds(testFile, builds)
	if err != nil {
		t.Fatalf("SaveBuilds() error = %v", err)
	}

	// Verify file was created
	if _, err := os.Stat(testFile); os.IsNotExist(err) {
		t.Error("File should have been created")
	}
}

// Test: RED - Load builds from file
func TestLoadBuilds(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test-builds.json")

	// Save some builds first
	originalBuilds := []models.Build{
		{PRNumber: "3859", Status: models.StatusSuccess, JobName: "job1", BuildNumber: 100},
		{PRNumber: "3860", Status: models.StatusRunning, JobName: "job2", BuildNumber: 101},
	}
	err := SaveBuilds(testFile, originalBuilds)
	if err != nil {
		t.Fatalf("Setup failed: %v", err)
	}

	// Load them back
	loadedBuilds, err := LoadBuilds(testFile)
	if err != nil {
		t.Fatalf("LoadBuilds() error = %v", err)
	}

	// Verify data matches
	if len(loadedBuilds) != 2 {
		t.Errorf("Expected 2 builds, got %d", len(loadedBuilds))
	}

	if loadedBuilds[0].PRNumber != "3859" {
		t.Errorf("Expected PR-3859, got PR-%s", loadedBuilds[0].PRNumber)
	}

	if loadedBuilds[1].BuildNumber != 101 {
		t.Errorf("Expected build #101, got #%d", loadedBuilds[1].BuildNumber)
	}
}

// Test: Load from non-existent file should return empty list
func TestLoadBuilds_FileNotExist(t *testing.T) {
	tmpDir := t.TempDir()
	nonExistentFile := filepath.Join(tmpDir, "does-not-exist.json")

	builds, err := LoadBuilds(nonExistentFile)
	if err != nil {
		t.Errorf("Loading non-existent file should not error, got: %v", err)
	}

	if len(builds) != 0 {
		t.Errorf("Expected empty list for non-existent file, got %d builds", len(builds))
	}
}

// Test: Get default config path
func TestGetConfigPath(t *testing.T) {
	path := GetConfigPath()

	if path == "" {
		t.Error("Config path should not be empty")
	}

	// Should end with .jenkins-dash-builds.json
	if filepath.Base(path) != ".jenkins-dash-builds.json" {
		t.Errorf("Expected .jenkins-dash-builds.json, got %s", filepath.Base(path))
	}
}

