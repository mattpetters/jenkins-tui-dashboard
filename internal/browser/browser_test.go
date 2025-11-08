package browser

import (
	"testing"
)

// Test: RED - OpenURL function
func TestOpenURL(t *testing.T) {
	// We can't really test the browser opening in unit tests
	// but we can verify the function exists and doesn't panic

	testURL := "https://github.com/test"
	err := OpenURL(testURL)

	// Function should execute without panic
	// Error is acceptable (no browser in test environment)
	_ = err
}

func TestOpenURL_EmptyURL(t *testing.T) {
	err := OpenURL("")
	if err == nil {
		t.Error("OpenURL with empty string should return error")
	}
}
