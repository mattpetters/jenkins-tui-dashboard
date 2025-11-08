# Jenkins Build Dashboard - Technical Reference

## Current Python Implementation Structure

### Project Layout
```
jenkins-dash/
├── src/
│   ├── __init__.py
│   ├── __main__.py              # Entry point
│   ├── main.py                  # Textual app class
│   ├── models.py                # Data models
│   ├── jenkins_client.py        # Jenkins API client
│   ├── github_helper.py         # GitHub API helper (if used)
│   ├── components/
│   │   ├── dashboard.py         # Main dashboard widget
│   │   ├── build_tile.py        # Individual tile widget
│   │   └── test_build_helper.py # Test data generator
│   └── utils/
│       ├── colors.py            # Color definitions
│       └── url_builder.py       # URL construction utilities
├── tests/                       # Test files
├── requirements.txt             # Python dependencies
└── .env                         # Environment variables (JENKINS_TOKEN)
```

### Data Models (models.py)

#### BuildStatus Enum
```python
from enum import Enum

class BuildStatus(Enum):
    PENDING = "pending"
    RUNNING = "running"
    SUCCESS = "success"
    FAILURE = "failure"
    ERROR = "error"
```

#### Build Class
```python
from dataclasses import dataclass
from typing import Optional

@dataclass
class Build:
    pr_number: str
    status: BuildStatus = BuildStatus.PENDING
    stage: Optional[str] = None
    job_name: Optional[str] = None
    job_path: Optional[str] = None
    build_number: Optional[int] = None
    build_url: Optional[str] = None
    pr_url: Optional[str] = None
    duration_seconds: int = 0
    timestamp: float = 0.0
    error_message: Optional[str] = None
    
    def is_running(self) -> bool:
        return self.status == BuildStatus.RUNNING
    
    def is_success(self) -> bool:
        return self.status == BuildStatus.SUCCESS
    
    def is_failure(self) -> bool:
        return self.status == BuildStatus.FAILURE
    
    def format_duration(self) -> str:
        """Format duration as 'Xh Ym Zs' or 'Ym Zs' or 'Zs'"""
        if self.duration_seconds == 0:
            return "0s"
        
        hours = self.duration_seconds // 3600
        minutes = (self.duration_seconds % 3600) // 60
        seconds = self.duration_seconds % 60
        
        if hours > 0:
            return f"{hours}h {minutes}m {seconds}s"
        elif minutes > 0:
            return f"{minutes}m {seconds}s"
        else:
            return f"{seconds}s"
```

#### DashboardState Class
```python
@dataclass
class DashboardState:
    builds: list[Build] = field(default_factory=list)
    selected_index: int = 0
    
    def add_build(self, build: Build) -> None:
        """Add a build to the list"""
        self.builds.append(build)
    
    def remove_build(self, index: int) -> bool:
        """Remove build at index, return True if successful"""
        if 0 <= index < len(self.builds):
            self.builds.pop(index)
            if self.selected_index >= len(self.builds) and len(self.builds) > 0:
                self.selected_index = len(self.builds) - 1
            return True
        return False
    
    def get_selected_build(self) -> Optional[Build]:
        """Get currently selected build"""
        if 0 <= self.selected_index < len(self.builds):
            return self.builds[self.selected_index]
        return None
    
    def move_selection(self, direction: str) -> None:
        """Move selection up/down/left/right in grid"""
        # Implementation calculates grid position and moves accordingly
        pass
```

### Jenkins API Client (jenkins_client.py)

#### JenkinsClient Class
```python
import os
import requests
from typing import Optional

class JenkinsClient:
    def __init__(self):
        self.base_url = "https://build.intuit.com"
        self.token = os.getenv("JENKINS_TOKEN")
        self.headers = {
            "Authorization": f"Bearer {self.token}",
            "Accept": "application/json"
        }
    
    def get_build_status(
        self, 
        job_path: str, 
        branch: str, 
        build_number: Optional[int] = None
    ) -> Optional[Build]:
        """
        Fetch build status from Jenkins API.
        
        Args:
            job_path: Full job path (e.g., "intuit-authentication/job/...")
            branch: Branch name (e.g., "PR-3859")
            build_number: Specific build number, or None for lastBuild
        
        Returns:
            Build object with status information, or None if error
        """
        # Construct URL
        build_ref = build_number if build_number else "lastBuild"
        url = f"{self.base_url}/{job_path}/job/{branch}/{build_ref}/api/json"
        
        try:
            response = requests.get(url, headers=self.headers, timeout=10)
            response.raise_for_status()
            data = response.json()
            
            # Parse response into Build object
            return self._parse_build_response(data, branch, job_path)
            
        except requests.exceptions.RequestException as e:
            # Return error build
            return Build(
                pr_number=branch.replace("PR-", ""),
                status=BuildStatus.ERROR,
                error_message=str(e)
            )
    
    def _parse_build_response(self, data: dict, branch: str, job_path: str) -> Build:
        """Parse Jenkins API JSON response into Build object"""
        # Extract PR number from branch
        pr_number = branch.replace("PR-", "")
        
        # Determine status
        if data.get("building", False):
            status = BuildStatus.RUNNING
        elif data.get("result") == "SUCCESS":
            status = BuildStatus.SUCCESS
        elif data.get("result") == "FAILURE":
            status = BuildStatus.FAILURE
        elif data.get("result") is None:
            status = BuildStatus.PENDING
        else:
            status = BuildStatus.ERROR
        
        # Get current stage (from stages if available)
        stage = None
        if "stages" in data and data["stages"]:
            # Find first non-SUCCESS stage or last stage
            for s in data["stages"]:
                if s.get("status") != "SUCCESS":
                    stage = s.get("name")
                    break
            if not stage:
                stage = data["stages"][-1].get("name")
        
        # Calculate duration
        duration_ms = data.get("duration", 0)
        duration_seconds = duration_ms // 1000
        
        # If still running, calculate from timestamp
        if status == BuildStatus.RUNNING:
            timestamp_ms = data.get("timestamp", 0)
            current_time_ms = time.time() * 1000
            duration_seconds = int((current_time_ms - timestamp_ms) / 1000)
        
        return Build(
            pr_number=pr_number,
            status=status,
            stage=stage,
            job_name=self._extract_job_name(data.get("fullDisplayName", "")),
            job_path=job_path,
            build_number=data.get("number"),
            build_url=data.get("url"),
            pr_url=build_pr_url(pr_number),
            duration_seconds=duration_seconds,
            timestamp=data.get("timestamp", 0) / 1000
        )
    
    def _extract_job_name(self, full_name: str) -> str:
        """Extract job name from full display name"""
        # Example: "auth-service » PR-3859 » #142" -> "auth-service"
        parts = full_name.split("»")
        return parts[0].strip() if parts else "Unknown"
```

### URL Builder Utilities (utils/url_builder.py)

```python
def infer_job_path_from_pr(pr_number: str) -> str:
    """
    Infer Jenkins job path from PR number.
    Currently hardcoded - could be made configurable.
    """
    return "intuit-authentication/job/auth-service-security/job/authentication-service-pr-ci"

def build_pr_url(pr_number: str) -> str:
    """Build GitHub PR URL from PR number"""
    return f"https://github.com/IntuitDeveloper/authentication-service/pull/{pr_number}"

def build_jenkins_build_url(job_path: str, branch: str, build_number: int) -> str:
    """Build full Jenkins build URL"""
    return f"https://build.intuit.com/{job_path}/job/{branch}/{build_number}"

def parse_pr_number(input_str: str) -> str:
    """Parse and validate PR number from user input"""
    # Remove "PR-" prefix if present
    cleaned = input_str.strip().upper().replace("PR-", "")
    
    # Validate it's numeric
    if not cleaned.isdigit():
        raise ValueError(f"Invalid PR number: {input_str}")
    
    return cleaned
```

### Test Data Generator (components/test_build_helper.py)

```python
import time
from ..models import Build, BuildStatus

def create_test_builds() -> list[Build]:
    """Create sample builds for testing UI"""
    return [
        Build(
            pr_number="3859",
            status=BuildStatus.SUCCESS,
            stage="Deploy",
            job_name="maven-build",
            job_path="intuit-authentication/job/auth-service-security/job/authentication-service-pr-ci",
            build_number=142,
            build_url="https://build.intuit.com/.../PR-3859/142",
            pr_url="https://github.com/IntuitDeveloper/authentication-service/pull/3859",
            duration_seconds=323,
            timestamp=time.time() - 600
        ),
        Build(
            pr_number="3860",
            status=BuildStatus.FAILURE,
            stage="Test",
            job_name="maven-test",
            job_path="intuit-authentication/job/auth-service-security/job/authentication-service-pr-ci",
            build_number=143,
            build_url="https://build.intuit.com/.../PR-3860/143",
            pr_url="https://github.com/IntuitDeveloper/authentication-service/pull/3860",
            duration_seconds=180,
            timestamp=time.time() - 300
        ),
        Build(
            pr_number="3861",
            status=BuildStatus.RUNNING,
            stage="Build",
            job_name="maven-build",
            job_path="intuit-authentication/job/auth-service-security/job/authentication-service-pr-ci",
            build_number=144,
            build_url="https://build.intuit.com/.../PR-3861/144",
            pr_url="https://github.com/IntuitDeveloper/authentication-service/pull/3861",
            duration_seconds=45,
            timestamp=time.time() - 45
        ),
    ]
```

## Jenkins API Details

### Authentication
- Uses Bearer token from `JENKINS_TOKEN` environment variable
- Token should be passed in `Authorization` header
- Format: `Authorization: Bearer <token>`

### API Endpoint Structure
```
GET https://build.intuit.com/{job_path}/job/{branch}/{build_number_or_lastBuild}/api/json
```

### Example API Response
```json
{
  "number": 142,
  "url": "https://build.intuit.com/intuit-authentication/job/auth-service-security/job/authentication-service-pr-ci/job/PR-3859/142/",
  "building": false,
  "result": "SUCCESS",
  "duration": 323000,
  "timestamp": 1699564800000,
  "fullDisplayName": "auth-service » PR-3859 » #142",
  "stages": [
    {
      "id": "1",
      "name": "Checkout",
      "status": "SUCCESS",
      "startTimeMillis": 1699564800000,
      "durationMillis": 5000
    },
    {
      "id": "2",
      "name": "Build",
      "status": "SUCCESS",
      "startTimeMillis": 1699564805000,
      "durationMillis": 180000
    },
    {
      "id": "3",
      "name": "Test",
      "status": "SUCCESS",
      "startTimeMillis": 1699564985000,
      "durationMillis": 120000
    },
    {
      "id": "4",
      "name": "Deploy",
      "status": "SUCCESS",
      "startTimeMillis": 1699565105000,
      "durationMillis": 18000
    }
  ]
}
```

### Result Values
- `null` - Build is running or queued
- `"SUCCESS"` - Build passed
- `"FAILURE"` - Build failed
- `"ABORTED"` - Build was manually stopped
- `"UNSTABLE"` - Build passed but with warnings
- `"NOT_BUILT"` - Build was skipped

### Stage Status Values
- `"SUCCESS"` - Stage completed successfully
- `"FAILED"` - Stage failed
- `"IN_PROGRESS"` - Stage currently running
- `"NOT_EXECUTED"` - Stage not reached yet
- `"PAUSED_PENDING_INPUT"` - Waiting for user input
- `"ABORTED"` - Stage was cancelled

## Color Definitions (utils/colors.py)

```python
class Colors:
    """Color definitions for build statuses"""
    SUCCESS = "green"
    FAILURE = "red"
    RUNNING = "blue"
    RUNNING_BRIGHT = "bright_blue"
    PENDING = "yellow"
    ERROR = "red"
    
    TEXT_LIGHT = "white"
    TEXT_DARK = "black"
```

## Grid Navigation Logic

### Converting Between Index and Grid Position
```python
def index_to_grid_pos(index: int, columns: int) -> tuple[int, int]:
    """Convert flat index to (row, col) position"""
    row = index // columns
    col = index % columns
    return (row, col)

def grid_pos_to_index(row: int, col: int, columns: int) -> int:
    """Convert (row, col) position to flat index"""
    return row * columns + col
```

### Navigation Logic
```python
def move_selection(current_index: int, direction: str, total_items: int, columns: int) -> int:
    """
    Calculate new index after navigation.
    
    Args:
        current_index: Current selected index
        direction: "up", "down", "left", or "right"
        total_items: Total number of items in grid
        columns: Number of columns in grid
    
    Returns:
        New index after movement
    """
    row, col = index_to_grid_pos(current_index, columns)
    rows = (total_items + columns - 1) // columns  # Ceiling division
    
    if direction == "up":
        if row > 0:
            new_row = row - 1
            new_index = grid_pos_to_index(new_row, col, columns)
            # Check if new index is valid
            if new_index < total_items:
                return new_index
    
    elif direction == "down":
        new_row = row + 1
        new_index = grid_pos_to_index(new_row, col, columns)
        # Check if new index is valid
        if new_index < total_items:
            return new_index
    
    elif direction == "left":
        if col > 0:
            return current_index - 1
    
    elif direction == "right":
        if col < columns - 1:
            new_index = current_index + 1
            # Check if new index exists
            if new_index < total_items:
                return new_index
    
    # If we can't move, stay at current position
    return current_index
```

## Browser Opening

```python
import webbrowser

def open_url_in_browser(url: str) -> None:
    """Open URL in default browser"""
    webbrowser.open(url)
```

## Environment Setup

### Required Environment Variables
```bash
# .env file
JENKINS_TOKEN=your_jenkins_bearer_token_here
```

### Python Dependencies (requirements.txt)
```
textual==0.47.0  # TUI framework (note: being replaced with bubbletea)
requests==2.31.0  # HTTP client
python-dotenv==1.0.0  # Environment variable loading
rich==13.7.0  # Terminal formatting (used by textual)
```

## Known Issues with Current Textual Implementation

1. **Widget Visibility**: Tiles don't appear after mounting to Grid
2. **Refresh Timing**: `refresh(layout=True)` doesn't always trigger layout recalculation
3. **Async Mounting**: `tile.remove()` is async, causing ID conflicts on re-mount
4. **CSS Overrides**: App-level CSS overrides widget-level styling unpredictably
5. **Update Timing**: Calling `update()` in `__init__` before mount doesn't work
6. **Grid Height**: Grid collapses to 0 height with `height: auto`
7. **Focus Management**: Focus not restored after Input widget removed

## Recommended Bubbletea Implementation

### Main Model Structure
```go
type model struct {
    builds         []Build
    selectedIndex  int
    inputMode      bool
    inputValue     string
    statusMessage  string
    jenkinsClient  *JenkinsClient
    termWidth      int
    termHeight     int
    blinkState     bool  // For animating running builds
}
```

### Key Commands (for async operations)
```go
// Fetch build status in background
type buildFetchedMsg struct {
    build Build
    err   error
}

func fetchBuildCmd(client *JenkinsClient, prNumber string) tea.Cmd {
    return func() tea.Msg {
        build, err := client.GetBuildStatus(prNumber)
        return buildFetchedMsg{build: build, err: err}
    }
}

// Tick for refresh timer
func tickCmd() tea.Cmd {
    return tea.Tick(10*time.Second, func(t time.Time) tea.Msg {
        return tickMsg(t)
    })
}

// Blink for running builds
func blinkCmd() tea.Cmd {
    return tea.Tick(800*time.Millisecond, func(t time.Time) tea.Msg {
        return blinkMsg(t)
    })
}
```

### Rendering Approach
```go
func (m model) View() string {
    var s strings.Builder
    
    // Header
    s.WriteString(renderHeader())
    s.WriteString("\n")
    
    // Grid of tiles
    s.WriteString(renderGrid(m.builds, m.selectedIndex, m.blinkState))
    s.WriteString("\n")
    
    // Input field (if in input mode)
    if m.inputMode {
        s.WriteString(renderInput(m.inputValue))
        s.WriteString("\n")
    }
    
    // Status bar
    s.WriteString(renderStatusBar(m.statusMessage))
    s.WriteString("\n")
    
    // Footer
    s.WriteString(renderFooter())
    
    return s.String()
}
```

This technical reference should provide everything needed to reimplement the dashboard in Bubbletea (Go) or any other TUI framework!

