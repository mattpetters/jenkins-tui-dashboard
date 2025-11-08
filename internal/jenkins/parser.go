package jenkins

import (
	"strings"
	"time"

	"github.com/mpetters/jenkins-dash/internal/models"
)

// ParseBuildResponse parses Jenkins API JSON response into a Build object
func ParseBuildResponse(data map[string]interface{}, prBranch, jobPath string) models.Build {
	// Extract PR number from branch (e.g., "PR-3859" -> "3859")
	prNumber := strings.TrimPrefix(prBranch, "PR-")

	// Determine status
	building, _ := data["building"].(bool)
	result, _ := data["result"].(string)

	var status models.BuildStatus
	if building {
		status = models.StatusRunning
	} else if result == "SUCCESS" {
		status = models.StatusSuccess
	} else if result == "FAILURE" {
		status = models.StatusFailure
	} else if result == "" || result == "null" {
		status = models.StatusPending
	} else {
		status = models.StatusError
	}

	// Extract basic fields
	buildNumber := int(getFloat(data, "number"))
	buildURL, _ := data["url"].(string)

	// Calculate duration
	durationMs := getFloat(data, "duration")
	durationSeconds := int(durationMs / 1000)

	// If still running, calculate from timestamp
	if status == models.StatusRunning {
		timestampMs := getFloat(data, "timestamp")
		currentTimeMs := float64(time.Now().Unix() * 1000)
		durationSeconds = int((currentTimeMs - timestampMs) / 1000)
	}

	// Extract timestamp
	timestampMs := getFloat(data, "timestamp")
	timestamp := int64(timestampMs / 1000)

	// Extract stage and job info from stages array
	var stage, jobName string
	if stagesData, ok := data["stages"].([]interface{}); ok && len(stagesData) > 0 {
		stage, jobName = ExtractStageInfo(stagesData, status)
	} else {
		// No stages data - use simple status-based text
		switch status {
		case models.StatusSuccess:
			stage, jobName = "Passed", "Passed"
		case models.StatusFailure:
			stage, jobName = "Failed", "Failed"
		case models.StatusRunning:
			stage, jobName = "Running", "In Progress"
		default:
			stage, jobName = "Pending", "Waiting"
		}
	}

	return models.Build{
		PRNumber:        prNumber,
		Status:          status,
		Stage:           stage,
		JobName:         jobName,
		JobPath:         jobPath,
		BuildNumber:     buildNumber,
		BuildURL:        buildURL,
		PRURL:           BuildPRURL(prNumber),
		DurationSeconds: durationSeconds,
		Timestamp:       timestamp,
	}
}

// ExtractJobName extracts the job name from Jenkins full display name
// Example: "auth-service » PR-3859 » #142" -> "auth-service"
func ExtractJobName(fullDisplayName string) string {
	if fullDisplayName == "" {
		return "Unknown"
	}

	parts := strings.Split(fullDisplayName, "»")
	if len(parts) > 0 {
		return strings.TrimSpace(parts[0])
	}

	return "Unknown"
}

// Helper function to safely extract float from map
func getFloat(data map[string]interface{}, key string) float64 {
	if val, ok := data[key].(float64); ok {
		return val
	}
	return 0
}

// ExtractStageInfo extracts phase and job information from stages
// Phase = high-level label (e.g., "BUILD:", "TEST:")
// Jobs = actual tasks (e.g., "Run unit tests, Run integration tests")
// For completed builds, returns simple "Passed"/"Failed"
func ExtractStageInfo(stages []interface{}, buildStatus models.BuildStatus) (phase string, jobs string) {
	if len(stages) == 0 {
		return "Unknown", "Unknown"
	}

	// For completed builds, show simple status
	if buildStatus == models.StatusSuccess {
		return "Passed", "Passed"
	}
	if buildStatus == models.StatusFailure || buildStatus == models.StatusError {
		return "Failed", "Failed"
	}

	// For running/pending builds, show detailed stage info
	var currentPhase string
	var phaseForActiveTasks string
	var activeTasks []string

	// Look for IN_PROGRESS stages
	for _, stageData := range stages {
		stageMap, ok := stageData.(map[string]interface{})
		if !ok {
			continue
		}

		stageName, _ := stageMap["name"].(string)
		stageStatus, _ := stageMap["status"].(string)

		// Phase labels end with ":"
		if len(stageName) > 0 && stageName[len(stageName)-1] == ':' {
			currentPhase = stageName
			continue
		}

		// Track IN_PROGRESS tasks
		if stageStatus == "IN_PROGRESS" {
			activeTasks = append(activeTasks, stageName)
			// Remember the phase for the first active task
			if phaseForActiveTasks == "" {
				phaseForActiveTasks = currentPhase
			}
		}
	}

	// Use the phase where active tasks were found
	if phaseForActiveTasks != "" {
		phase = phaseForActiveTasks
	} else if currentPhase != "" {
		phase = currentPhase
	} else {
		phase = "Running"
	}

	// Combine active tasks
	if len(activeTasks) > 0 {
		jobs = strings.Join(activeTasks, ", ")
	} else {
		jobs = "Starting..."
	}

	return phase, jobs
}

// Helper function to extract current stage from stages array (legacy)
func extractCurrentStage(data map[string]interface{}) string {
	stages, ok := data["stages"].([]interface{})
	if !ok || len(stages) == 0 {
		return "Unknown"
	}

	// Use new extraction logic
	phase, _ := ExtractStageInfo(stages, models.StatusRunning)
	return phase
}
