package jenkins

import (
	"testing"

	"github.com/mpetters/jenkins-dash/internal/models"
)

// Test: Better stage extraction from real Jenkins data
func TestExtractStageInfo(t *testing.T) {
	// Real nested pipeline structure
	stages := []interface{}{
		map[string]interface{}{"name": "BUILD:", "status": "SUCCESS"},
		map[string]interface{}{"name": "Podman Multi-Stage Build(NO Tests)", "status": "SUCCESS"},
		map[string]interface{}{"name": "Run Unit Tests", "status": "IN_PROGRESS"},
		map[string]interface{}{"name": "QAL:", "status": "NOT_EXECUTED"},
	}

	phase, jobs := ExtractStageInfo(stages, models.StatusRunning)

	// Phase should be the outer label (BUILD:)
	if phase != "BUILD:" {
		t.Errorf("Expected phase 'BUILD:', got '%s'", phase)
	}

	// Job should be the nested task name
	if jobs != "Run Unit Tests" {
		t.Errorf("Expected jobs 'Run Unit Tests', got '%s'", jobs)
	}
}

func TestExtractStageInfo_ParallelStages(t *testing.T) {
	// Multiple tasks running in parallel within BUILD phase
	stages := []interface{}{
		map[string]interface{}{"name": "BUILD:", "status": "SUCCESS"},
		map[string]interface{}{"name": "Compile", "status": "SUCCESS"},
		map[string]interface{}{"name": "Run Unit Tests", "status": "IN_PROGRESS"},
		map[string]interface{}{"name": "Run Integration Tests", "status": "IN_PROGRESS"},
		map[string]interface{}{"name": "Deploy:", "status": "NOT_EXECUTED"},
	}

	phase, jobs := ExtractStageInfo(stages, models.StatusRunning)

	// Phase should be the outer label
	if phase != "BUILD:" {
		t.Errorf("Expected phase 'BUILD:', got '%s'", phase)
	}

	// Jobs should show both parallel tasks
	if !contains(jobs, "Run Unit Tests") || !contains(jobs, "Run Integration Tests") {
		t.Errorf("Expected both parallel tasks in jobs, got '%s'", jobs)
	}
}

func TestExtractStageInfo_RunningWithActiveTasks(t *testing.T) {
	stages := []interface{}{
		map[string]interface{}{"name": "BUILD:", "status": "SUCCESS"},
		map[string]interface{}{"name": "Compile", "status": "SUCCESS"},
		map[string]interface{}{"name": "TEST:", "status": "SUCCESS"},
		map[string]interface{}{"name": "Run Tests", "status": "IN_PROGRESS"},
	}

	phase, jobs := ExtractStageInfo(stages, models.StatusRunning)

	// Stage should be the outer phase label
	if phase != "TEST:" {
		t.Errorf("Expected phase 'TEST:', got '%s'", phase)
	}

	// Job should be the nested task
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

