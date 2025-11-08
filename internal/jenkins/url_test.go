package jenkins

import (
	"testing"
)

// Test 7: RED - URL building utilities
func TestInferJobPath(t *testing.T) {
	tests := []struct {
		name     string
		prNumber string
		wantPath string
	}{
		{
			name:     "Standard PR",
			prNumber: "3934",
			wantPath: "identity/job/identity-manage/job/account/job/account-eks",
		},
		{
			name:     "Different PR number",
			prNumber: "1234",
			wantPath: "identity/job/identity-manage/job/account/job/account-eks",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := InferJobPath(tt.prNumber)
			if got != tt.wantPath {
				t.Errorf("InferJobPath() = %v, want %v", got, tt.wantPath)
			}
		})
	}
}

func TestBuildPRURL(t *testing.T) {
	tests := []struct {
		name     string
		prNumber string
		want     string
	}{
		{
			name:     "Standard PR - Blue Ocean",
			prNumber: "3934",
			want:     "https://build.intuit.com/blue/organizations/jenkins/identity%2Fidentity-manage%2Faccount%2Faccount-eks/detail/PR-3934/",
		},
		{
			name:     "Different PR - Blue Ocean",
			prNumber: "1234",
			want:     "https://build.intuit.com/blue/organizations/jenkins/identity%2Fidentity-manage%2Faccount%2Faccount-eks/detail/PR-1234/",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := BuildPRURL(tt.prNumber)
			if got != tt.want {
				t.Errorf("BuildPRURL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBuildJenkinsURL(t *testing.T) {
	tests := []struct {
		name        string
		jobPath     string
		branch      string
		buildNumber int
		want        string
	}{
		{
			name:        "Specific build number",
			jobPath:     "intuit-auth/job/pr-ci",
			branch:      "PR-3859",
			buildNumber: 142,
			want:        "https://build.intuit.com/intuit-auth/job/pr-ci/job/PR-3859/142",
		},
		{
			name:        "Last build (number 0)",
			jobPath:     "intuit-auth/job/pr-ci",
			branch:      "PR-3860",
			buildNumber: 0,
			want:        "https://build.intuit.com/intuit-auth/job/pr-ci/job/PR-3860/lastBuild",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := BuildJenkinsURL(tt.jobPath, tt.branch, tt.buildNumber)
			if got != tt.want {
				t.Errorf("BuildJenkinsURL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParsePRNumber(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		want      string
		wantError bool
	}{
		{
			name:      "Just number",
			input:     "3859",
			want:      "3859",
			wantError: false,
		},
		{
			name:      "With PR- prefix",
			input:     "PR-3859",
			want:      "3859",
			wantError: false,
		},
		{
			name:      "Lowercase with spaces",
			input:     "  pr-1234  ",
			want:      "1234",
			wantError: false,
		},
		{
			name:      "Invalid - contains letters",
			input:     "abc123",
			want:      "",
			wantError: true,
		},
		{
			name:      "Empty string",
			input:     "",
			want:      "",
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParsePRNumber(tt.input)
			if tt.wantError {
				if err == nil {
					t.Error("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if got != tt.want {
					t.Errorf("ParsePRNumber() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}
