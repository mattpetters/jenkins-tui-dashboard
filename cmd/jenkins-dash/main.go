package main

import (
	"fmt"
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/joho/godotenv"
	"github.com/mpetters/jenkins-dash/internal/jenkins"
	"github.com/mpetters/jenkins-dash/internal/ui"
)

func main() {
	// Load environment variables from .env file
	_ = godotenv.Load() // Ignore error if .env doesn't exist

	// Get Jenkins credentials (uses Basic Auth with username:token)
	username := os.Getenv("JENKINS_USER")
	jenkinsToken := os.Getenv("JENKINS_TOKEN")
	githubToken := os.Getenv("GITHUB_TOKEN")
	
	if username == "" || jenkinsToken == "" {
		fmt.Println("⚠️  Warning: JENKINS_USER or JENKINS_TOKEN not set in .env file")
		fmt.Println("    Set both to fetch real Jenkins data")
		fmt.Println()
	}
	
	if githubToken == "" {
		fmt.Println("⚠️  Warning: GITHUB_TOKEN not set in .env file")
		fmt.Println("    Set it to auto-fetch Git branch names")
		fmt.Println()
	}

	// Create Jenkins client
	jenkinsClient := jenkins.NewClient(username, jenkinsToken)

	// Get config file path
	configPath := getConfigPath()

	// Create the model with Jenkins client and config path
	m := ui.NewModelWithClient(jenkinsClient, configPath)

	// Load persisted builds
	if err := m.LoadPersistedBuilds(); err != nil {
		fmt.Printf("Warning: Could not load saved builds: %v\n", err)
	}

	// Run the program
	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v\n", err)
		os.Exit(1)
	}
}

func getConfigPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return ".jenkins-dash-builds.json"
	}
	return filepath.Join(homeDir, ".jenkins-dash-builds.json")
}
