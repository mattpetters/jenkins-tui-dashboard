package jenkins

import (
	"testing"

	"github.com/mpetters/jenkins-dash/internal/models"
)

// Test that extractGitBranch doesn't return master/main since those are base branches
func TestExtractGitBranch_IgnoresMasterBranches(t *testing.T) {
	tests := []struct {
		name     string
		data     map[string]interface{}
		expected string
	}{
		{
			name: "returns empty for master branch",
			data: map[string]interface{}{
				"actions": []interface{}{
					map[string]interface{}{
						"lastBuiltRevision": map[string]interface{}{
							"branch": []interface{}{
								map[string]interface{}{
									"name": "origin/master",
								},
							},
						},
					},
				},
			},
			expected: "", // Should return empty, not "master"
		},
		{
			name: "returns empty for main branch",
			data: map[string]interface{}{
				"actions": []interface{}{
					map[string]interface{}{
						"lastBuiltRevision": map[string]interface{}{
							"branch": []interface{}{
								map[string]interface{}{
									"name": "origin/main",
								},
							},
						},
					},
				},
			},
			expected: "", // Should return empty, not "main"
		},
		{
			name: "returns empty for refs/remotes/origin/master",
			data: map[string]interface{}{
				"actions": []interface{}{
					map[string]interface{}{
						"lastBuiltRevision": map[string]interface{}{
							"branch": []interface{}{
								map[string]interface{}{
									"name": "refs/remotes/origin/master",
								},
							},
						},
					},
				},
			},
			expected: "", // Should return empty
		},
		{
			name: "returns feature branch",
			data: map[string]interface{}{
				"actions": []interface{}{
					map[string]interface{}{
						"lastBuiltRevision": map[string]interface{}{
							"branch": []interface{}{
								map[string]interface{}{
									"name": "origin/feature/add-auth",
								},
							},
						},
					},
				},
			},
			expected: "feature/add-auth",
		},
		{
			name: "returns bugfix branch",
			data: map[string]interface{}{
				"actions": []interface{}{
					map[string]interface{}{
						"lastBuiltRevision": map[string]interface{}{
							"branch": []interface{}{
								map[string]interface{}{
									"name": "origin/bugfix/fix-login",
								},
							},
						},
					},
				},
			},
			expected: "bugfix/fix-login",
		},
		{
			name:     "returns empty for no branch data",
			data:     map[string]interface{}{},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractGitBranch(tt.data)
			if result != tt.expected {
				t.Errorf("extractGitBranch() = %q, want %q", result, tt.expected)
			}
		})
	}
}

// Test that ParseBuildResponse doesn't overwrite GitBranch with master/main
func TestParseBuildResponse_PreservesFeatureBranches(t *testing.T) {
	data := map[string]interface{}{
		"building": false,
		"result":   "SUCCESS",
		"number":   float64(142),
		"duration": float64(120000),
		"timestamp": float64(1234567890000),
		"actions": []interface{}{
			map[string]interface{}{
				"lastBuiltRevision": map[string]interface{}{
					"branch": []interface{}{
						map[string]interface{}{
							"name": "origin/master", // Jenkins reports master
						},
					},
				},
			},
		},
	}

	build := ParseBuildResponse(data, "PR-3859", "test/job/path")

	// GitBranch should be empty (not "master") so GitHub branch can be used
	if build.GitBranch != "" {
		t.Errorf("GitBranch should be empty when Jenkins reports master, got %q", build.GitBranch)
	}
}

// Test that ParseBuildResponse preserves feature branches from Jenkins
func TestParseBuildResponse_UsesJenkinsFeatureBranch(t *testing.T) {
	data := map[string]interface{}{
		"building": false,
		"result":   "SUCCESS",
		"number":   float64(142),
		"duration": float64(120000),
		"timestamp": float64(1234567890000),
		"actions": []interface{}{
			map[string]interface{}{
				"lastBuiltRevision": map[string]interface{}{
					"branch": []interface{}{
						map[string]interface{}{
							"name": "origin/feature/new-feature",
						},
					},
				},
			},
		},
	}

	build := ParseBuildResponse(data, "PR-3859", "test/job/path")

	// GitBranch should be the feature branch from Jenkins
	if build.GitBranch != "feature/new-feature" {
		t.Errorf("GitBranch = %q, want %q", build.GitBranch, "feature/new-feature")
	}
}

func TestParseBuildResponse_BasicFields(t *testing.T) {
	data := map[string]interface{}{
		"building": false,
		"result":   "SUCCESS",
		"number":   float64(142),
		"duration": float64(120000),
		"timestamp": float64(1234567890000),
	}

	build := ParseBuildResponse(data, "PR-3859", "test/job/path")

	if build.PRNumber != "3859" {
		t.Errorf("PRNumber = %s, want 3859", build.PRNumber)
	}
	if build.Status != models.StatusSuccess {
		t.Errorf("Status = %v, want StatusSuccess", build.Status)
	}
	if build.BuildNumber != 142 {
		t.Errorf("BuildNumber = %d, want 142", build.BuildNumber)
	}
	if build.DurationSeconds != 120 {
		t.Errorf("DurationSeconds = %d, want 120", build.DurationSeconds)
	}
}

func TestExtractJobName_ValidInput(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"auth-service » PR-3859 » #142", "auth-service"},
		{"my-job » PR-123 » #1", "my-job"},
		{"single-job", "single-job"},
		{"", "Unknown"},
	}

	for _, tt := range tests {
		result := ExtractJobName(tt.input)
		if result != tt.expected {
			t.Errorf("ExtractJobName(%q) = %q, want %q", tt.input, result, tt.expected)
		}
	}
}

func TestExtractStageInfo_CompletedBuild(t *testing.T) {
	stages := []interface{}{
		map[string]interface{}{"name": "BUILD:", "status": "SUCCESS"},
		map[string]interface{}{"name": "Test", "status": "SUCCESS"},
	}

	phase, job := ExtractStageInfo(stages, models.StatusSuccess)
	if phase != "Passed" || job != "Passed" {
		t.Errorf("Expected Passed/Passed for success, got %s/%s", phase, job)
	}

	phase, job = ExtractStageInfo(stages, models.StatusFailure)
	if phase != "Failed" || job != "Failed" {
		t.Errorf("Expected Failed/Failed for failure, got %s/%s", phase, job)
	}
}

func TestExtractStageInfo_RunningBuild(t *testing.T) {
	stages := []interface{}{
		map[string]interface{}{"name": "BUILD:", "status": "SUCCESS"},
		map[string]interface{}{"name": "Compile", "status": "SUCCESS"},
		map[string]interface{}{"name": "QAL:", "status": "IN_PROGRESS"},
		map[string]interface{}{"name": "Integration Tests", "status": "IN_PROGRESS"},
	}

	phase, job := ExtractStageInfo(stages, models.StatusRunning)
	
	if phase != "QAL:" {
		t.Errorf("Expected phase 'QAL:', got %q", phase)
	}
	if job != "Integration Tests" {
		t.Errorf("Expected job 'Integration Tests', got %q", job)
	}
}
