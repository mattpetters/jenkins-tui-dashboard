package jenkins

import (
	"fmt"
	"regexp"
	"strings"
)

const (
	// Updated to match your actual Jenkins job structure
	defaultJobPath = "identity/job/identity-manage/job/account/job/account-eks"
	jenkinsBaseURL = "https://build.intuit.com"
	githubRepo     = "IntuitDeveloper/authentication-service"  // TODO: Update to correct repo
	githubBaseURL  = "https://github.com"
)

// InferJobPath returns the Jenkins job path for a given PR number
// Currently returns a hardcoded path but could be made configurable
func InferJobPath(prNumber string) string {
	return defaultJobPath
}

// BuildPRURL constructs the Blue Ocean PR URL for the given PR number
func BuildPRURL(prNumber string) string {
	// Blue Ocean URL format: /blue/organizations/jenkins/{job-path}/detail/PR-{number}/{build}/pipeline
	// For "view all builds" we can omit the build number
	jobPathEncoded := strings.ReplaceAll(defaultJobPath, "/job/", "%2F")
	return fmt.Sprintf("%s/blue/organizations/jenkins/%s/detail/PR-%s/", jenkinsBaseURL, jobPathEncoded, prNumber)
}

// BuildJenkinsURL constructs the full Jenkins build URL
// If buildNumber is 0, uses "lastBuild" instead
func BuildJenkinsURL(jobPath, branch string, buildNumber int) string {
	buildRef := "lastBuild"
	if buildNumber > 0 {
		buildRef = fmt.Sprintf("%d", buildNumber)
	}
	return fmt.Sprintf("%s/%s/job/%s/%s", jenkinsBaseURL, jobPath, branch, buildRef)
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
