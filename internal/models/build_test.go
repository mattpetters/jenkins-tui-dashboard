package models

import (
	"testing"
)

// Test 1: BuildStatus enum (COMPLETE)
func TestBuildStatus_String(t *testing.T) {
	tests := []struct {
		name   string
		status BuildStatus
		want   string
	}{
		{
			name:   "Pending status",
			status: StatusPending,
			want:   "pending",
		},
		{
			name:   "Running status",
			status: StatusRunning,
			want:   "running",
		},
		{
			name:   "Success status",
			status: StatusSuccess,
			want:   "success",
		},
		{
			name:   "Failure status",
			status: StatusFailure,
			want:   "failure",
		},
		{
			name:   "Error status",
			status: StatusError,
			want:   "error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.status.String(); got != tt.want {
				t.Errorf("BuildStatus.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Test 2: RED - Build struct creation and basic fields
func TestBuild_Creation(t *testing.T) {
	build := Build{
		PRNumber:        "3859",
		Status:          StatusSuccess,
		Stage:           "Deploy",
		JobName:         "maven-build",
		JobPath:         "intuit-auth/job/auth-service/job/pr-ci",
		BuildNumber:     142,
		BuildURL:        "https://build.intuit.com/job/PR-3859/142",
		PRURL:           "https://github.com/IntuitDeveloper/auth/pull/3859",
		DurationSeconds: 323,
		Timestamp:       1699564800,
		ErrorMessage:    "",
	}

	if build.PRNumber != "3859" {
		t.Errorf("PRNumber = %v, want 3859", build.PRNumber)
	}
	if build.Status != StatusSuccess {
		t.Errorf("Status = %v, want StatusSuccess", build.Status)
	}
	if build.BuildNumber != 142 {
		t.Errorf("BuildNumber = %v, want 142", build.BuildNumber)
	}
	if build.DurationSeconds != 323 {
		t.Errorf("DurationSeconds = %v, want 323", build.DurationSeconds)
	}
}

// Test 3: RED - Build status helper methods
func TestBuild_StatusMethods(t *testing.T) {
	tests := []struct {
		name      string
		status    BuildStatus
		isRunning bool
		isSuccess bool
		isFailure bool
	}{
		{
			name:      "Running build",
			status:    StatusRunning,
			isRunning: true,
			isSuccess: false,
			isFailure: false,
		},
		{
			name:      "Successful build",
			status:    StatusSuccess,
			isRunning: false,
			isSuccess: true,
			isFailure: false,
		},
		{
			name:      "Failed build",
			status:    StatusFailure,
			isRunning: false,
			isSuccess: false,
			isFailure: true,
		},
		{
			name:      "Pending build",
			status:    StatusPending,
			isRunning: false,
			isSuccess: false,
			isFailure: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			build := Build{Status: tt.status}

			if got := build.IsRunning(); got != tt.isRunning {
				t.Errorf("IsRunning() = %v, want %v", got, tt.isRunning)
			}
			if got := build.IsSuccess(); got != tt.isSuccess {
				t.Errorf("IsSuccess() = %v, want %v", got, tt.isSuccess)
			}
			if got := build.IsFailure(); got != tt.isFailure {
				t.Errorf("IsFailure() = %v, want %v", got, tt.isFailure)
			}
		})
	}
}

// Test 4: RED - FormatDuration method
func TestBuild_FormatDuration(t *testing.T) {
	tests := []struct {
		name    string
		seconds int
		want    string
	}{
		{
			name:    "Zero duration",
			seconds: 0,
			want:    "0s",
		},
		{
			name:    "Seconds only",
			seconds: 45,
			want:    "45s",
		},
		{
			name:    "Minutes and seconds",
			seconds: 125, // 2m 5s
			want:    "2m 5s",
		},
		{
			name:    "Exactly 1 minute",
			seconds: 60,
			want:    "1m 0s",
		},
		{
			name:    "Hours, minutes, and seconds",
			seconds: 3725, // 1h 2m 5s
			want:    "1h 2m 5s",
		},
		{
			name:    "Multiple hours",
			seconds: 7323, // 2h 2m 3s
			want:    "2h 2m 3s",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			build := Build{DurationSeconds: tt.seconds}
			if got := build.FormatDuration(); got != tt.want {
				t.Errorf("FormatDuration() = %v, want %v", got, tt.want)
			}
		})
	}
}
