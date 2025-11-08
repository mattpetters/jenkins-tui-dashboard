package models

import (
	"fmt"
	"time"
)

// BuildStatus represents the current status of a build
type BuildStatus int

const (
	StatusPending BuildStatus = iota
	StatusRunning
	StatusSuccess
	StatusFailure
	StatusError
)

// String returns the string representation of the BuildStatus
func (s BuildStatus) String() string {
	return [...]string{"pending", "running", "success", "failure", "error"}[s]
}

// Build represents a Jenkins build for a PR
type Build struct {
	PRNumber        string
	GitBranch       string // Git branch name from GitHub (e.g., "feature/add-auth")
	Status          BuildStatus
	Stage           string
	JobName         string
	JobPath         string
	BuildNumber     int
	BuildURL        string
	PRURL           string
	DurationSeconds int
	Timestamp       int64
	ErrorMessage    string
}

// IsRunning returns true if the build is currently running
func (b Build) IsRunning() bool {
	return b.Status == StatusRunning
}

// IsSuccess returns true if the build was successful
func (b Build) IsSuccess() bool {
	return b.Status == StatusSuccess
}

// IsFailure returns true if the build failed
func (b Build) IsFailure() bool {
	return b.Status == StatusFailure
}

// GetCurrentDuration returns the current duration, accounting for running builds
func (b Build) GetCurrentDuration() int {
	if b.IsRunning() && b.Timestamp > 0 {
		// Calculate elapsed time from when build started
		return int(time.Now().Unix() - b.Timestamp)
	}
	return b.DurationSeconds
}

// FormatDuration returns a human-readable duration string
func (b Build) FormatDuration() string {
	duration := b.GetCurrentDuration()
	
	if duration == 0 {
		return "0s"
	}

	hours := duration / 3600
	minutes := (duration % 3600) / 60
	seconds := duration % 60

	if hours > 0 {
		return fmt.Sprintf("%dh %dm %ds", hours, minutes, seconds)
	} else if minutes > 0 {
		return fmt.Sprintf("%dm %ds", minutes, seconds)
	}
	return fmt.Sprintf("%ds", seconds)
}

// FormatCompletedTime returns the completion time for finished builds in PT
func (b Build) FormatCompletedTime() string {
	if b.IsRunning() || b.Status == StatusPending {
		return "" // Not completed yet
	}
	
	// Calculate end time: start + duration
	endTime := time.Unix(b.Timestamp+int64(b.DurationSeconds), 0)
	
	// Convert to Pacific Time
	pt, err := time.LoadLocation("America/Los_Angeles")
	if err != nil {
		// Fallback to local time if PT not available
		pt = time.Local
	}
	endTimePT := endTime.In(pt)
	
	// Format as "11/7 10:45pm"
	month := int(endTimePT.Month())
	day := endTimePT.Day()
	hour := endTimePT.Hour()
	minute := endTimePT.Minute()
	
	ampm := "am"
	displayHour := hour
	if hour >= 12 {
		ampm = "pm"
		if hour > 12 {
			displayHour = hour - 12
		}
	}
	if displayHour == 0 {
		displayHour = 12
	}
	
	return fmt.Sprintf("%d/%d %d:%02d%s", month, day, displayHour, minute, ampm)
}
