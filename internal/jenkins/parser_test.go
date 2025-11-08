package jenkins

import (
	"testing"

	"github.com/mpetters/jenkins-dash/internal/models"
)

// Test: Better stage extraction from real Jenkins data
func TestExtractStageInfo(t *testing.T) {
	// Real data structure from Jenkins
	stages := []interface{}{
		map[string]interface{}{"name": "BUILD:", "status": "SUCCESS"},
		map[string]interface{}{"name": "Run Unit Tests", "status": "SUCCESS"},
		map[string]interface{}{"name": "Run Integration Tests", "status": "IN_PROGRESS"},
		map[string]interface{}{"name": "QAL:", "status": "NOT_EXECUTED"},
	}

	phase, jobs := ExtractStageInfo(stages, models.StatusRunning)

	// Phase should be the last phase label seen before current activity
	if phase != "BUILD:" {
		t.Errorf("Expected phase 'BUILD:', got '%s'", phase)
	}

	// Jobs should be the currently running tasks
	if jobs != "Run Integration Tests" {
		t.Errorf("Expected jobs 'Run Integration Tests', got '%s'", jobs)
	}
}

func TestExtractStageInfo_ParallelStages(t *testing.T) {
	// Multiple tasks running in parallel
	stages := []interface{}{
		map[string]interface{}{"name": "Test:", "status": "SUCCESS"},
		map[string]interface{}{"name": "Run Unit Tests", "status": "IN_PROGRESS"},
		map[string]interface{}{"name": "Run Integration Tests", "status": "IN_PROGRESS"},
		map[string]interface{}{"name": "Deploy:", "status": "NOT_EXECUTED"},
	}

	phase, jobs := ExtractStageInfo(stages, models.StatusRunning)

	if phase != "Test:" {
		t.Errorf("Expected phase 'Test:', got '%s'", phase)
	}

	// Should show both parallel tasks
	if !contains(jobs, "Run Unit Tests") || !contains(jobs, "Run Integration Tests") {
		t.Errorf("Expected both parallel tasks, got '%s'", jobs)
	}
}

func TestExtractStageInfo_RunningWithActiveTasks(t *testing.T) {
	stages := []interface{}{
		map[string]interface{}{"name": "BUILD:", "status": "SUCCESS"},
		map[string]interface{}{"name": "Compile", "status": "SUCCESS"},
		map[string]interface{}{"name": "Test:", "status": "IN_PROGRESS"},
		map[string]interface{}{"name": "Run Tests", "status": "IN_PROGRESS"},
	}

	phase, jobs := ExtractStageInfo(stages, models.StatusRunning)

	if phase != "Test:" {
		t.Errorf("Expected phase 'Test:', got '%s'", phase)
	}

	if jobs != "Run Tests" {
		t.Errorf("Expected jobs 'Run Tests', got '%s'", jobs)
	}
}

func TestExtractStageInfo_CompletedSuccess(t *testing.T) {
	stages := []interface{}{
		map[string]interface{}{"name": "BUILD:", "status": "SUCCESS"},
		map[string]interface{}{"name": "Compile", "status": "SUCCESS"},
		map[string]interface{}{"name": "Test:", "status": "SUCCESS"},
		map[string]interface{}{"name": "Run Tests", "status": "SUCCESS"},
	}

	phase, jobs := ExtractStageInfo(stages, models.StatusSuccess)

	if phase != "Passed" {
		t.Errorf("Expected phase 'Passed' for successful build, got '%s'", phase)
	}

	if jobs != "Passed" {
		t.Errorf("Expected jobs 'Passed' for successful build, got '%s'", jobs)
	}
}

func TestExtractStageInfo_CompletedFailure(t *testing.T) {
	stages := []interface{}{
		map[string]interface{}{"name": "BUILD:", "status": "SUCCESS"},
		map[string]interface{}{"name": "Compile", "status": "SUCCESS"},
		map[string]interface{}{"name": "Test:", "status": "FAILED"},
		map[string]interface{}{"name": "Run Tests", "status": "FAILED"},
	}

	phase, jobs := ExtractStageInfo(stages, models.StatusFailure)

	if phase != "Failed" {
		t.Errorf("Expected phase 'Failed' for failed build, got '%s'", phase)
	}

	if jobs != "Failed" {
		t.Errorf("Expected jobs 'Failed' for failed build, got '%s'", jobs)
	}
}

func contains(s, substr string) bool {
	return len(s) > 0 && len(substr) > 0 && (s == substr || len(s) > len(substr) && (s[:len(substr)] == substr || s[len(s)-len(substr):] == substr || containsMiddle(s, substr)))
}

func containsMiddle(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

