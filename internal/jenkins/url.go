package jenkins

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

// Default configuration (can be overridden with environment variables)
const (
	defaultJobPath     = "identity/job/identity-manage/job/account/job/account-eks"
	defaultJenkinsURL  = "https://build.intuit.com"
	defaultGitHubRepo  = "identity-manage/account"
	defaultGitHubURL   = "https://github.intuit.com"
)

var (
	jenkinsBaseURL string
	githubRepo     string
	githubBaseURL  string
)

func init() {
	// Load from environment or use defaults
	jenkinsBaseURL = getEnvOrDefault("JENKINS_BASE_URL", defaultJenkinsURL)
	githubRepo = getEnvOrDefault("GITHUB_REPO", defaultGitHubRepo)
	githubBaseURL = getEnvOrDefault("GITHUB_BASE_URL", defaultGitHubURL)
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// InferJobPath returns the Jenkins job path for a given PR number
// Reads from JENKINS_JOB_PATH environment variable or uses default
func InferJobPath(prNumber string) string {
	if jobPath := os.Getenv("JENKINS_JOB_PATH"); jobPath != "" {
		return jobPath
	}
	return defaultJobPath
}

// BuildPRURL constructs the GitHub PR URL for the given PR number
func BuildPRURL(prNumber string) string {
	// Intuit GitHub URL
	return fmt.Sprintf("%s/%s/pull/%s", githubBaseURL, githubRepo, prNumber)
}

// BuildJenkinsURL constructs the full Jenkins build URL (classic view)
// Used for API calls only
func BuildJenkinsURL(jobPath, branch string, buildNumber int) string {
	buildRef := "lastBuild"
	if buildNumber > 0 {
		buildRef = fmt.Sprintf("%d", buildNumber)
	}
	return fmt.Sprintf("%s/%s/job/%s/%s", jenkinsBaseURL, jobPath, branch, buildRef)
}

// BuildBlueOceanBuildURL constructs the Blue Ocean pipeline view URL for a specific build
func BuildBlueOceanBuildURL(jobPath, branch string, buildNumber int) string {
	// Blue Ocean format: https://build.intuit.com/{first-segment}/blue/organizations/jenkins/{rest-of-path}/detail/{branch}/{build}/pipeline
	// Example: identity/job/identity-manage/job/account/job/account-eks
	// Becomes: https://build.intuit.com/identity/blue/organizations/jenkins/identity-manage%2Faccount%2Faccount-eks/detail/PR-3934/8/pipeline
	
	// Split job path to extract first segment
	parts := strings.Split(jobPath, "/job/")
	if len(parts) < 2 {
		// Fallback if job path doesn't have expected structure
		return fmt.Sprintf("%s/%s", jenkinsBaseURL, jobPath)
	}
	
	firstSegment := parts[0]  // e.g., "identity"
	restOfPath := strings.Join(parts[1:], "%2F")  // e.g., "identity-manage%2Faccount%2Faccount-eks"
	
	buildRef := fmt.Sprintf("%d", buildNumber)
	if buildNumber == 0 {
		buildRef = "lastBuild"
	}
	
	return fmt.Sprintf("%s/%s/blue/organizations/jenkins/%s/detail/%s/%s/pipeline", 
		jenkinsBaseURL, firstSegment, restOfPath, branch, buildRef)
}

// ParsePRNumber parses and validates a PR number from user input
// Removes "PR-" prefix if present, trims whitespace, and validates it's numeric
func ParsePRNumber(input string) (string, error) {
	// Trim whitespace and convert to uppercase
	cleaned := strings.TrimSpace(strings.ToUpper(input))

	// Remove "PR-" prefix if present
	cleaned = strings.TrimPrefix(cleaned, "PR-")

	// Validate it's not empty
	if cleaned == "" {
		return "", fmt.Errorf("PR number cannot be empty")
	}

	// Validate it's numeric
	matched, err := regexp.MatchString(`^\d+$`, cleaned)
	if err != nil {
		return "", err
	}
	if !matched {
		return "", fmt.Errorf("PR number must contain only digits, got: %s", input)
	}

	return cleaned, nil
}
