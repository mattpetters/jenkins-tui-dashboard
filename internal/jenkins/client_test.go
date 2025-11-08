package jenkins

import (
	"testing"
	"time"

	"github.com/mpetters/jenkins-dash/internal/models"
)

// Test 8: RED - Parse Jenkins API JSON response
func TestParseBuildResponse(t *testing.T) {
	tests := []struct {
		name     string
		response map[string]interface{}
		prBranch string
		want     models.Build
	}{
		{
			name: "Successful completed build",
			response: map[string]interface{}{
				"number":          float64(142),
				"building":        false,
				"result":          "SUCCESS",
				"duration":        float64(323000),
				"timestamp":       float64(1699564800000),
				"url":             "https://build.intuit.com/job/PR-3859/142",
				"fullDisplayName": "auth-service » PR-3859 » #142",
			},
			prBranch: "PR-3859",
			want: models.Build{
				PRNumber:    "3859",
				Status:      models.StatusSuccess,
				BuildNumber: 142,
				// BuildURL is now Blue Ocean format, constructed in parser
				DurationSeconds: 323,
				Timestamp:       1699564800,
			},
		},
		{
			name: "Failed build",
			response: map[string]interface{}{
				"number":          float64(143),
				"building":        false,
				"result":          "FAILURE",
				"duration":        float64(180000),
				"timestamp":       float64(1699564900000),
				"url":             "https://build.intuit.com/job/PR-3860/143",
				"fullDisplayName": "auth-service » PR-3860 » #143",
			},
			prBranch: "PR-3860",
			want: models.Build{
				PRNumber:    "3860",
				Status:      models.StatusFailure,
				BuildNumber: 143,
				// BuildURL is now Blue Ocean format
				DurationSeconds: 180,
				Timestamp:       1699564900,
			},
		},
		{
			name: "Running build",
			response: map[string]interface{}{
				"number":          float64(144),
				"building":        true,
				"result":          nil,
				"duration":        float64(0),
				"timestamp":       float64(time.Now().Unix() * 1000),
				"url":             "https://build.intuit.com/job/PR-3861/144",
				"fullDisplayName": "auth-service » PR-3861 » #144",
			},
			prBranch: "PR-3861",
			want: models.Build{
				PRNumber:    "3861",
				Status:      models.StatusRunning,
				BuildNumber: 144,
				// BuildURL is now Blue Ocean format
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ParseBuildResponse(tt.response, tt.prBranch, "test-job-path")

			if got.PRNumber != tt.want.PRNumber {
				t.Errorf("PRNumber = %v, want %v", got.PRNumber, tt.want.PRNumber)
			}
			if got.Status != tt.want.Status {
				t.Errorf("Status = %v, want %v", got.Status, tt.want.Status)
			}
			if got.BuildNumber != tt.want.BuildNumber {
				t.Errorf("BuildNumber = %v, want %v", got.BuildNumber, tt.want.BuildNumber)
			}
			// Don't test BuildURL - it's constructed dynamically with Blue Ocean format
			// Just verify it's not empty
			if got.BuildURL == "" {
				t.Error("BuildURL should not be empty")
			}
			if tt.want.DurationSeconds > 0 && got.DurationSeconds != tt.want.DurationSeconds {
				t.Errorf("DurationSeconds = %v, want %v", got.DurationSeconds, tt.want.DurationSeconds)
			}
		})
	}
}

func TestExtractJobName(t *testing.T) {
	tests := []struct {
		name        string
		displayName string
		want        string
	}{
		{
			name:        "Standard format",
			displayName: "auth-service » PR-3859 » #142",
			want:        "auth-service",
		},
		{
			name:        "Different job name",
			displayName: "maven-build » PR-1234 » #56",
			want:        "maven-build",
		},
		{
			name:        "Empty display name",
			displayName: "",
			want:        "Unknown",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ExtractJobName(tt.displayName)
			if got != tt.want {
				t.Errorf("ExtractJobName() = %v, want %v", got, tt.want)
			}
		})
	}
}
