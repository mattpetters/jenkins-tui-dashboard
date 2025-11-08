package testdata

import (
	"time"

	"github.com/mpetters/jenkins-dash/internal/models"
)

// CreateTestBuilds returns sample builds for testing the UI
func CreateTestBuilds() []models.Build {
	now := time.Now().Unix()

	return []models.Build{
		{
			PRNumber:        "3859",
			Status:          models.StatusSuccess,
			Stage:           "Deploy",
			JobName:         "maven-build",
			JobPath:         "intuit-authentication/job/auth-service-security/job/authentication-service-pr-ci",
			BuildNumber:     142,
			BuildURL:        "https://build.intuit.com/intuit-authentication/job/auth-service-security/job/authentication-service-pr-ci/job/PR-3859/142",
			PRURL:           "https://github.com/IntuitDeveloper/authentication-service/pull/3859",
			DurationSeconds: 323,
			Timestamp:       now - 600,
		},
		{
			PRNumber:        "3860",
			Status:          models.StatusFailure,
			Stage:           "Test",
			JobName:         "maven-test",
			JobPath:         "intuit-authentication/job/auth-service-security/job/authentication-service-pr-ci",
			BuildNumber:     143,
			BuildURL:        "https://build.intuit.com/intuit-authentication/job/auth-service-security/job/authentication-service-pr-ci/job/PR-3860/143",
			PRURL:           "https://github.com/IntuitDeveloper/authentication-service/pull/3860",
			DurationSeconds: 180,
			Timestamp:       now - 300,
		},
		{
			PRNumber:        "3861",
			Status:          models.StatusRunning,
			Stage:           "Build",
			JobName:         "maven-build",
			JobPath:         "intuit-authentication/job/auth-service-security/job/authentication-service-pr-ci",
			BuildNumber:     144,
			BuildURL:        "https://build.intuit.com/intuit-authentication/job/auth-service-security/job/authentication-service-pr-ci/job/PR-3861/144",
			PRURL:           "https://github.com/IntuitDeveloper/authentication-service/pull/3861",
			DurationSeconds: 45,
			Timestamp:       now - 45,
		},
		{
			PRNumber:        "3862",
			Status:          models.StatusPending,
			Stage:           "",
			JobName:         "",
			JobPath:         "intuit-authentication/job/auth-service-security/job/authentication-service-pr-ci",
			BuildNumber:     0,
			BuildURL:        "",
			PRURL:           "https://github.com/IntuitDeveloper/authentication-service/pull/3862",
			DurationSeconds: 0,
			Timestamp:       now,
		},
	}
}
