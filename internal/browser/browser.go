package browser

import (
	"fmt"
	"os/exec"
	"runtime"
)

// OpenURL opens a URL in the default browser
func OpenURL(url string) error {
	if url == "" {
		return fmt.Errorf("URL cannot be empty")
	}

	var cmd *exec.Cmd

	// Use appropriate command based on OS
	switch runtime.GOOS {
	case "darwin": // macOS
		cmd = exec.Command("open", url)
	case "linux":
		cmd = exec.Command("xdg-open", url)
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	default:
		return fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}

	return cmd.Start()
}
