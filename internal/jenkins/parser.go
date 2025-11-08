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
		// No stages data - show status-based text
		switch status {
		case models.StatusSuccess:
			stage, jobName = "Passed", "Passed"
		case models.StatusFailure:
			stage, jobName = "Failed", "Failed"
		case models.StatusRunning:
			stage, jobName = "Running", "In Progress"
		case models.StatusPending:
			stage, jobName = "Pending", "Queued"
		default:
			stage, jobName = "Unknown", "Unknown"
		}
	}
	
	// Extract Git branch from actions if available
	gitBranch := extractGitBranch(data)

	// Build Blue Ocean URL for better pipeline visualization
	blueOceanURL := BuildBlueOceanBuildURL(jobPath, prBranch, buildNumber)
	
	return models.Build{
		PRNumber:        prNumber,
		GitBranch:       gitBranch,
		Status:          status,
		Stage:           stage,
		JobName:         jobName,
		JobPath:         jobPath,
		BuildNumber:     buildNumber,
		BuildURL:        blueOceanURL, // Use Blue Ocean URL instead of classic
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
// Stage = outer phase label (e.g., "BUILD:", "QAL:")
// Job = nested task name (e.g., "Podman Multi-Stage Build(NO Tests)")
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

	// For running/pending builds, find the nested structure
	var currentPhase string      // Last "LABEL:" seen
	var phaseForActiveTasks string // Phase label where active tasks are
	var activeTasks []string      // Tasks that are IN_PROGRESS

	for _, stageData := range stages {
		stageMap, ok := stageData.(map[string]interface{})
		if !ok {
			continue
		}

		stageName, _ := stageMap["name"].(string)
		stageStatus, _ := stageMap["status"].(string)

		// Phase labels end with ":" (e.g., "BUILD:", "QAL:")
		if len(stageName) > 0 && stageName[len(stageName)-1] == ':' {
			currentPhase = stageName
			continue
		}

		// Track IN_PROGRESS tasks and remember which phase they belong to
		if stageStatus == "IN_PROGRESS" {
			activeTasks = append(activeTasks, stageName)
			if phaseForActiveTasks == "" {
				phaseForActiveTasks = currentPhase // Remember the phase label
			}
		}
	}

	// Stage = the outer phase label (BUILD:, QAL:, etc.)
	if phaseForActiveTasks != "" {
		phase = phaseForActiveTasks
	} else if currentPhase != "" {
		phase = currentPhase
	} else {
		phase = "Starting"
	}

	// Job = the nested task name(s)
	if len(activeTasks) > 0 {
		jobs = strings.Join(activeTasks, ", ")
	} else {
		jobs = "Starting..."
	}

	return phase, jobs
}

// extractGitBranch attempts to extract the Git branch name from Jenkins actions
// Filters out master/main branches since those are base branches, not PR branches
func extractGitBranch(data map[string]interface{}) string {
	// Try to get from actions array
	if actions, ok := data["actions"].([]interface{}); ok {
		for _, action := range actions {
			actionMap, ok := action.(map[string]interface{})
			if !ok {
				continue
			}

			// Look for branch name in various possible fields
			if branchName, ok := actionMap["lastBuiltRevision"].(map[string]interface{}); ok {
				if branch, ok := branchName["branch"].([]interface{}); ok && len(branch) > 0 {
					if branchMap, ok := branch[0].(map[string]interface{}); ok {
						if name, ok := branchMap["name"].(string); ok {
							// Extract short branch name (e.g., "origin/feature/auth" -> "feature/auth")
							parts := strings.Split(name, "/")
							var shortName string
							if len(parts) > 1 {
								shortName = strings.Join(parts[1:], "/")
							} else {
								shortName = name
							}
							
							// Filter out master/main branches - those are base branches, not PR branches
							// We want the feature branch name from GitHub, not the base branch
							if shortName == "master" || shortName == "main" || 
							   shortName == "remotes/origin/master" || shortName == "remotes/origin/main" {
								return "" // Return empty so GitHub branch is preferred
							}
							
							return shortName
						}
					}
				}
			}
		}
	}

	return "" // Will fallback to GitHub branch or PR number in tile
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
