package ui

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mpetters/jenkins-dash/internal/models"
)

// Model represents the Bubbletea application state
type Model struct {
	state          *models.DashboardState
	jenkinsClient  Client
	configPath     string
	inputMode      bool
	inputValue     string
	statusMessage  string
	termWidth      int
	termHeight     int
	blinkState     bool
}

// Client is an interface to avoid import cycle with jenkins package
type Client interface {
	GetBuildStatus(jobPath, branch string, buildNum int) (*models.Build, error)
}

// NewModel creates a new Model with default values
func NewModel() Model {
	return NewModelWithClient(nil, "")
}

// NewModelWithClient creates a new Model with a Jenkins client
func NewModelWithClient(client Client, configPath string) Model {
	return Model{
		state: &models.DashboardState{
			Builds:        []models.Build{},
			SelectedIndex: 0,
			GridColumns:   3,
		},
		jenkinsClient: client,
		configPath:    configPath,
		inputMode:     false,
		inputValue:    "",
		statusMessage: "Press 'a' to add a PR build, arrow keys to navigate",
		termWidth:     120,
		termHeight:    40,
		blinkState:    false,
	}
}

// AddTestBuild adds a build to the model (for testing/demo purposes)
func (m *Model) AddTestBuild(build models.Build) {
	m.state.AddBuild(build)
}

// Init initializes the model
func (m Model) Init() tea.Cmd {
	return tea.Batch(
		tickCmd(),     // 10 second build refresh
		blinkCmd(),    // 800ms blink for running builds
		timeTickCmd(), // 1 second time tick for live clock
	)
}

// Update handles messages and updates the model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m.handleKeyPress(msg)

	case tea.WindowSizeMsg:
		m.termWidth = msg.Width
		m.termHeight = msg.Height
		m.state.GridColumns = CalculateGridColumns(m.termWidth)
		return m, nil

	case buildFetchedMsg:
		// Update build with fetched data
		if msg.index >= 0 && msg.index < len(m.state.Builds) {
			if msg.err != nil {
				m.state.Builds[msg.index].Status = models.StatusError
				m.state.Builds[msg.index].ErrorMessage = msg.err.Error()
				m.statusMessage = fmt.Sprintf("✗ Error fetching PR-%s: %v", m.state.Builds[msg.index].PRNumber, msg.err)
			} else {
				m.state.Builds[msg.index] = *msg.build
				m.statusMessage = fmt.Sprintf("✓ PR-%s: %s", msg.build.PRNumber, msg.build.Status.String())
			}
			// Save state after update
			_ = m.saveState()
		}
		return m, nil

	case tickMsg:
		// Refresh all builds every 10 seconds
		if m.jenkinsClient != nil {
			var cmds []tea.Cmd
			for i, build := range m.state.Builds {
				if build.Status != models.StatusPending {
					cmds = append(cmds, fetchBuildCmd(m.jenkinsClient, build.PRNumber, i))
				}
			}
			cmds = append(cmds, tickCmd())
			return m, tea.Batch(cmds...)
		}
		return m, tickCmd()

	case blinkMsg:
		// Toggle blink state for running builds
		m.blinkState = !m.blinkState
		return m, blinkCmd()

	case timeTickMsg:
		// Just trigger re-render for live time updates
		return m, timeTickCmd()

	case urlOpenedMsg:
		// Browser opened (or failed)
		if msg.err != nil {
			m.statusMessage = fmt.Sprintf("✗ Error opening URL: %v", msg.err)
		} else {
			m.statusMessage = "✓ Opened in browser"
		}
		return m, nil
	}

	return m, nil
}

// handleKeyPress processes keyboard input
func (m Model) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// Handle input mode separately
	if m.inputMode {
		return m.handleInputMode(msg)
	}

	// Handle normal mode keys
	switch msg.Type {
	case tea.KeyCtrlC, tea.KeyEsc:
		return m, tea.Quit

	case tea.KeyRunes:
		switch string(msg.Runes) {
		case "q":
			return m, tea.Quit
		case "a":
			m.inputMode = true
			m.inputValue = ""
			m.statusMessage = "Enter PR number and press Enter"
			return m, nil
		case "d":
			// Delete selected build
			if m.state.GetSelectedBuild() != nil {
				m.state.RemoveBuild(m.state.SelectedIndex)
				m.statusMessage = "Build deleted"
				// Save state after deletion
				_ = m.saveState()
			}
			return m, nil
		case "p":
			// Open PR in browser
			if build := m.state.GetSelectedBuild(); build != nil {
				return m, openURLCmd(build.PRURL)
			}
			return m, nil
		}

	case tea.KeyEnter:
		// Open build in browser
		if build := m.state.GetSelectedBuild(); build != nil && build.BuildURL != "" {
			return m, openURLCmd(build.BuildURL)
		}
		return m, nil

	case tea.KeyUp:
		m.state.MoveSelection("up")
		return m, nil

	case tea.KeyDown:
		m.state.MoveSelection("down")
		return m, nil

	case tea.KeyLeft:
		m.state.MoveSelection("left")
		return m, nil

	case tea.KeyRight:
		m.state.MoveSelection("right")
		return m, nil
	}

	return m, nil
}

// handleInputMode processes keyboard input when in input mode
func (m Model) handleInputMode(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyEnter:
		// Submit the PR number
		if m.inputValue != "" {
			// Create a loading build and add it to state
			build := models.Build{
				PRNumber: m.inputValue,
				Status:   models.StatusPending,
				Stage:    "Loading...",
				JobName:  "Fetching data...",
			}
			m.state.AddBuild(build)
			newIndex := len(m.state.Builds) - 1
			m.statusMessage = "✓ Added PR-" + m.inputValue + " - Fetching build & branch data..."
			
			// Save state after adding
			_ = m.saveState()
			
			// Reset input state
			m.inputMode = false
			m.inputValue = ""
			
			// Fetch Jenkins build data AND GitHub branch name
			if m.jenkinsClient != nil {
				return m, fetchBuildAndBranchCmd(m.jenkinsClient, build.PRNumber, newIndex)
			}
			return m, nil
		}
		return m, nil

	case tea.KeyEsc:
		// Cancel input
		m.inputMode = false
		m.inputValue = ""
		m.statusMessage = "Press 'a' to add a PR build, arrow keys to navigate"
		return m, nil

	case tea.KeyBackspace:
		if len(m.inputValue) > 0 {
			m.inputValue = m.inputValue[:len(m.inputValue)-1]
		}
		return m, nil
	
	case tea.KeyRunes:
		m.inputValue += string(msg.Runes)
		return m, nil
	}

	return m, nil
}

// Message types for async operations
type tickMsg time.Time
type blinkMsg time.Time
type timeTickMsg time.Time

// Commands
func tickCmd() tea.Cmd {
	return tea.Tick(10*time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func blinkCmd() tea.Cmd {
	return tea.Tick(800*time.Millisecond, func(t time.Time) tea.Msg {
		return blinkMsg(t)
	})
}

func timeTickCmd() tea.Cmd {
	return tea.Tick(1*time.Second, func(t time.Time) tea.Msg {
		return timeTickMsg(t)
	})
}

// View renders the UI
func (m Model) View() string {
	var sections []string

	// Title/Header
	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FFFFFF")).
		Padding(0, 1)
	header := headerStyle.Render("🔨 Jenkins Build Dashboard")
	sections = append(sections, header)
	sections = append(sections, "")

	// Main grid
	grid := RenderGrid(m.state.Builds, m.state.SelectedIndex, m.state.GridColumns, m.blinkState)
	sections = append(sections, grid)
	sections = append(sections, "")

	// Input field (if in input mode)
	if m.inputMode {
		inputStyle := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#00FFFF")).
			Padding(0, 1)
		
		prompt := "PR number: "
		inputDisplay := inputStyle.Render(prompt + m.inputValue + "█")
		sections = append(sections, inputDisplay)
		sections = append(sections, "")
	}

	// Status bar
	statusStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#888888")).
		Padding(0, 1)
	statusBar := statusStyle.Render(m.statusMessage)
	sections = append(sections, statusBar)

	// Footer with key bindings
	footerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#666666")).
		Padding(0, 1)
	footer := footerStyle.Render("a: Add PR | d: Delete | ↑↓←→: Navigate | enter: Open Build | p: Open PR | q: Quit")
	sections = append(sections, footer)

	return lipgloss.JoinVertical(lipgloss.Left, sections...)
}
