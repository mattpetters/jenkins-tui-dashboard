package models

import (
	"strings"
	"testing"
	"time"
)

// Test FormatCompletedTime with real timestamp
func TestBuild_FormatCompletedTime(t *testing.T) {
	// Create a build that finished at a known time
	// Nov 7, 2025, 10:45 PM PT = Nov 8, 2025, 6:45 AM UTC
	pt, _ := time.LoadLocation("America/Los_Angeles")
	completionTime := time.Date(2025, 11, 7, 22, 45, 0, 0, pt)
	
	// Build started 30 minutes earlier
	startTime := completionTime.Add(-30 * time.Minute)
	
	build := Build{
		Status:          StatusSuccess,
		Timestamp:       startTime.Unix(),
		DurationSeconds: 1800, // 30 minutes
	}
	
	result := build.FormatCompletedTime()
	
	// Should show "11/7 10:45pm"
	if !strings.Contains(result, "11/7") {
		t.Errorf("Should contain date 11/7, got %s", result)
	}
	if !strings.Contains(result, "10:45pm") {
		t.Errorf("Should contain time 10:45pm, got %s", result)
	}
}

func TestBuild_FormatCompletedTime_Running(t *testing.T) {
	build := Build{
		Status:          StatusRunning,
		Timestamp:       time.Now().Unix(),
		DurationSeconds: 100,
	}
	
	result := build.FormatCompletedTime()
	
	// Running builds should return empty string
	if result != "" {
		t.Errorf("Running build should return empty string, got %s", result)
	}
}

// Test 1: Build model should have PRAuthor field
func TestBuild_HasPRAuthorField(t *testing.T) {
	build := Build{
		PRNumber:  "3859",
		PRAuthor:  "john.doe",
		GitBranch: "feature/add-auth",
	}
	
	if build.PRAuthor != "john.doe" {
		t.Errorf("Expected PRAuthor to be 'john.doe', got '%s'", build.PRAuthor)
	}
}

// Test 2: Build model should have Repository field
func TestBuild_HasRepositoryField(t *testing.T) {
	build := Build{
		PRNumber:   "3859",
		Repository: "identity-manage/account",
	}
	
	if build.Repository != "identity-manage/account" {
		t.Errorf("Expected Repository to be 'identity-manage/account', got '%s'", build.Repository)
	}
}
