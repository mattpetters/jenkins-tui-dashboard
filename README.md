# Jenkins Build Dashboard 🔨

A beautiful terminal-based dashboard for monitoring Jenkins builds across multiple PRs, built with Bubbletea and strict TDD.

![Jenkins Dashboard](https://img.shields.io/badge/Go-1.24+-00ADD8?style=flat&logo=go)
![Tests](https://img.shields.io/badge/tests-30%2B%20passing-success)
![Coverage](https://img.shields.io/badge/coverage-66--89%25-green)


<img width="1470" height="540" alt="Screenshot 2025-11-08 at 8 29 42 AM" src="https://github.com/user-attachments/assets/d878b35e-f386-49da-a1aa-44d8204fdb35" />

## Features

### Core Functionality
- 🎨 **Beautiful pastel colors** - Soft green/red/blue/yellow, easy on the eyes
- 🔄 **Auto-refresh** - Updates every 10 seconds automatically
- ⚡ **Manual refresh** - Press 'r' to refresh immediately
- 🧹 **Clear cache** - Press 'c' to clear and refetch everything
- ⏱️ **Live time** - Running builds show elapsed time updating every second
- 💾 **Persistent** - Auto-saves to `~/.jenkins-dash-builds.json`
- 🎯 **Clear selection** - Bright green border on selected tile
- 🌐 **Browser integration** - Blue Ocean and GitHub integration

### Jenkins Integration
- ✅ Dual API calls (standard + wfapi for pipeline stages)
- ✅ Basic Auth with username:token
- ✅ Real-time pipeline stage tracking
- ✅ Parallel stage detection
- ✅ Completion timestamps in Pacific Time

### GitHub Integration
- ✅ Auto-fetches Git branch names (e.g., "IDLMP-2038-aggregate")
- ✅ Shows PR check status (e.g., "5/8 checks", "all passed")
- ✅ Direct links to PRs and commits

## Requirements

- **Go 1.24+** (automatically managed by go.mod)
- **Terminal with Unicode support** (for box drawing characters)
- **Jenkins credentials** (username + API token)
- **GitHub token** (optional, for branch names and check status)

## Installation

```bash
# Clone the repository
git clone https://github.intuit.com/mpetters/jenkins-tui-dashboard.git
cd jenkins-tui-dashboard

# Build
go build -o jenkins-dash ./cmd/jenkins-dash

# Or use the build script
./run.sh
```

## Configuration

Create a `.env` file in the project root (see `env.example`):

### Required
```bash
JENKINS_USER=your_username
JENKINS_TOKEN=your_jenkins_api_token
```

### Optional (defaults to identity-manage/account project)
```bash
# GitHub integration (recommended for branch names and check status)
GITHUB_TOKEN=your_github_token

# Customize for your project
JENKINS_JOB_PATH=your-org/job/your-project/job/your-job
GITHUB_REPO=your-org/your-repo

# Advanced (usually not needed)
JENKINS_BASE_URL=https://build.intuit.com
GITHUB_BASE_URL=https://github.intuit.com
```

### Getting Tokens

**Jenkins API Token:**
1. Log into Jenkins
2. Click your name → Configure
3. API Token → Add new Token
4. Copy the generated token

**GitHub Token:**
1. Go to https://github.intuit.com/settings/tokens
2. Generate new token (classic)
3. Select scopes: `repo`, `read:org`
4. Copy the token

### Defaults

If you don't set environment variables, the app uses:
- **Jenkins Job**: `identity/job/identity-manage/job/account/job/account-eks`
- **GitHub Repo**: `identity-manage/account`
- **Jenkins URL**: `https://build.intuit.com`
- **GitHub URL**: `https://github.intuit.com`

## Keyboard Controls

| Key | Action |
|-----|--------|
| `a` | Add new PR build |
| `c` | Clear cache & refetch all data |
| `d` | Delete selected build |
| `r` | Refresh all builds now |
| `↑↓←→` | Navigate between builds |
| `Enter` | Open build in Blue Ocean pipeline view |
| `p` | Open PR in GitHub |
| `q` | Quit |

## Display Format

### Tile Layout
```
┌──────────────────────────────┐
│          PR-3934             │  ← PR number
│    IDLMP-2038-aggregate      │  ← Git branch from GitHub
├──────────────────────────────┤
│ Stage: BUILD:                │  ← Pipeline phase
│ Job: Run Unit Tests          │  ← Actual Jenkins task
│ Time: 32m 15s                │  ← Duration (live for running)
│ 11/7 10:45pm         #263    │  ← Completion time (PT) + Build #
│ PR: 5/8 checks               │  ← GitHub check status
└──────────────────────────────┘
```

### Status Colors
- 🟢 **Green (Passed)**: Build succeeded
- 🔴 **Red (Failed)**: Build failed
- 🔵 **Blue (Running)**: Build in progress (blinks)
- 🟡 **Yellow (Pending)**: Loading data

### Stage & Job Logic
- **Completed builds**: Simple "Passed" or "Failed"
- **Running builds**: Actual pipeline stages from Jenkins
  - Stage: Outer phase (e.g., "BUILD:", "QAL:", "E2E EAST:")
  - Job: Nested task (e.g., "Podman Multi-Stage Build", "Run Unit Tests")
  - Parallel stages: Multiple tasks shown with commas

## Architecture

```
jenkins-dash/
├── cmd/jenkins-dash/     # Main entry point
├── internal/
│   ├── browser/         # URL opening
│   ├── jenkins/         # API client & parsers
│   ├── models/          # Data structures
│   ├── persistence/     # Save/load builds
│   ├── testdata/        # Test fixtures
│   └── ui/              # Bubbletea UI components
└── go.mod
```

## Development

### Run Tests
```bash
go test ./...                    # All tests
go test ./... -v                 # Verbose
go test ./... -cover             # With coverage
```

### Test Coverage
```
browser:      67% coverage
jenkins:      66% coverage  
models:       89% coverage
persistence:  75% coverage
ui:           56% coverage
```

## API Integration

### Jenkins
- Fetches from standard `/api/json` endpoint (basic build info)
- Fetches from `/wfapi/describe` endpoint (pipeline stages)
- Merges data for complete picture
- Uses Basic Auth (username:token)

### GitHub
- Fetches branch names from PR API
- Fetches check run status
- Uses Bearer token authentication
- Caches results in persistence file

## Development

### Run Tests
```bash
go test ./...                    # All tests
go test ./... -v                 # Verbose
go test ./... -cover             # With coverage
```

### Test Coverage
```
✅ 30+ tests, all passing
✅ 66-89% code coverage
✅ All features TDD'd
```

## License

MIT

## Author

Built with strict Test-Driven Development using Bubbletea, Lipgloss, and Go.
