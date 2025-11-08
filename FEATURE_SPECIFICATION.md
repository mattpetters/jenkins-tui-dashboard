# Jenkins Build Dashboard - Complete Feature Specification

## Overview
A terminal-based dashboard for monitoring Jenkins build statuses for multiple PRs simultaneously. Users can add PR numbers, view real-time build status, and navigate to Jenkins/GitHub in their browser.

## Core Features

### 1. Build Tile Display

#### Tile Appearance
Each build is displayed as a tile with the following information:
```
┌──────────────────────────────┐
│          PR-3859             │  <- PR number (centered)
├──────────────────────────────┤
│ Stage: Build                 │  <- Current pipeline stage
│ Job: maven-build             │  <- Job name
│ Time: 5m 23s                 │  <- Duration
│                        #142  │  <- Build number (right-aligned)
└──────────────────────────────┘
```

#### Tile Colors (Background)
- **Green**: Build succeeded
- **Red**: Build failed
- **Blue**: Build is running (should blink/pulse between blue and bright_blue)
- **Yellow**: Build pending/loading/unknown status

#### Tile Dimensions
- Width: 36 characters
- Height: 11 lines
- Fixed size to maintain grid alignment

#### Selected Tile
- Selected tile has **bold** text
- Only one tile can be selected at a time
- Selection persists when adding/removing tiles

### 2. Grid Layout

#### Layout Behavior
- Tiles arranged in a grid (default: 3 columns)
- Grid should adapt to terminal width:
  - Small terminal: 1-2 columns
  - Medium terminal: 3 columns
  - Large terminal: 4 columns
- Tiles maintain consistent size regardless of content
- Grid should scroll if more tiles than fit on screen

#### Auto-refresh
- Each tile refreshes its build status every 10 seconds
- Refresh is independent per tile (don't refresh all at once)
- Running builds should show animation (blinking)

### 3. Adding Builds

#### User Flow
1. User presses 'a' key
2. Input field appears at bottom of screen (or appropriate location)
3. Input prompt: "Enter PR number (e.g., 3859)"
4. User types PR number (just the number, not "PR-")
5. User presses Enter:
   - Input closes immediately
   - Loading tile appears with "PR-XXXX" and "Loading..." status
   - API call happens in background (non-blocking)
   - Tile updates when data arrives
6. If user presses Escape: input closes, no build added

#### PR Number Resolution
Given PR number (e.g., 3859), the system infers:
- Job path: `intuit-authentication/job/auth-service-security/job/authentication-service-pr-ci`
- Branch: `PR-3859`
- PR URL: `https://github.com/IntuitDeveloper/authentication-service/pull/3859`
- Jenkins URL: `https://build.intuit.com/intuit-authentication/job/auth-service-security/job/authentication-service-pr-ci/job/PR-3859/lastBuild`

### 4. Keyboard Navigation

#### Global Keys (always active)
- `a` - Add new PR build
- `e` - Edit selected build's PR number
- `d` - Delete selected build
- `q` - Quit application
- `Enter` - Open selected build in browser (Jenkins)
- `p` - Open selected PR in browser (GitHub)

#### Arrow Keys (for selection)
- `Up` - Move selection up one row
- `Down` - Move selection down one row  
- `Left` - Move selection left one tile
- `Right` - Move selection right one tile

#### Grid Navigation Logic
- Selection moves within the grid
- If at top edge and press Up: stay at current position (or wrap to bottom)
- If at bottom edge and press Down: stay at current position (or wrap to top)
- If at left edge and press Left: stay at current position
- If at right edge and press Right: stay at current position
- When tiles are added/removed, try to maintain selection on same index

### 5. Editing Builds

#### Edit Flow
1. User selects a tile and presses 'e'
2. Input field appears with current PR number pre-filled
3. User modifies the number
4. On Enter: tile updates to new PR number (fetches new data)
5. On Escape: input closes, no changes

### 6. Deleting Builds

#### Delete Flow
1. User selects a tile and presses 'd'
2. Tile is immediately removed
3. Grid reflows to fill the gap
4. Selection moves to next tile (or previous if last tile deleted)
5. Status message: "✓ Deleted PR-XXXX. N build(s) remaining"

### 7. Opening Links

#### Open Build (Enter)
- Opens Jenkins build URL in default browser
- URL format: `https://build.intuit.com/{job_path}/job/{branch}/{buildNumber}`
- If build number is 0 or unknown: opens lastBuild
- Status message: "Opened build #142 in browser"

#### Open PR (p)
- Opens GitHub PR URL in default browser
- URL format: `https://github.com/IntuitDeveloper/authentication-service/pull/{pr_number}`
- Status message: "Opened PR-3859 in browser"

### 8. Status Bar

#### Location
- Bottom of screen
- Always visible
- Height: 3 lines
- Docked to bottom

#### Content
- Shows current action/result messages
- Examples:
  - Default: "Press 'a' to add a PR build, arrow keys to navigate"
  - After adding: "✓ Added PR-3859 - Fetching build data..."
  - After deleting: "✓ Deleted PR-3859. 5 build(s) remaining"
  - Error: "✗ Error fetching PR-3859: Connection timeout"

### 9. Header

#### Content
- Application title: "Jenkins Build Dashboard"
- Clock showing current time
- Always visible at top

### 10. Footer

#### Content
- Shows available keyboard shortcuts
- Format: `a Add PR  e Edit  d Delete  ↵ Open Build  p Open PR  ↑↓←→ Navigate  q Quit`
- Always visible at bottom (above status bar)

## Data Model

### Build Object
```python
class Build:
    pr_number: str           # e.g., "3859"
    status: BuildStatus      # PENDING, RUNNING, SUCCESS, FAILURE, ERROR
    stage: str              # e.g., "Build", "Test", "Deploy"
    job_name: str           # e.g., "maven-build"
    job_path: str           # Full Jenkins job path
    build_number: int       # e.g., 142
    build_url: str          # Full Jenkins build URL
    pr_url: str             # Full GitHub PR URL
    duration_seconds: int   # How long build has been running/ran
    timestamp: float        # When build started
    error_message: str      # If status is ERROR
```

### BuildStatus Enum
```python
PENDING = "pending"      # Waiting to start / loading data
RUNNING = "running"      # Currently executing
SUCCESS = "success"      # Passed
FAILURE = "failure"      # Failed
ERROR = "error"          # Error fetching data
```

## API Integration

### Jenkins Client

#### Get Build Status
```python
def get_build_status(job_path: str, branch: str, build_number: Optional[int] = None) -> Build:
    """
    Fetches build status from Jenkins.
    
    Args:
        job_path: e.g., "intuit-authentication/job/auth-service-security/job/authentication-service-pr-ci"
        branch: e.g., "PR-3859"
        build_number: Specific build number, or None for lastBuild
    
    Returns:
        Build object with status, stage, duration, etc.
    
    API Endpoint:
        GET https://build.intuit.com/{job_path}/job/{branch}/{build_number}/api/json
    
    Authentication:
        Bearer token from JENKINS_TOKEN environment variable
    """
```

#### Jenkins API Response Structure
```json
{
  "number": 142,
  "result": "SUCCESS",  // null if running, "SUCCESS", "FAILURE", "ABORTED"
  "duration": 323000,   // milliseconds
  "timestamp": 1699564800000,
  "url": "https://build.intuit.com/...",
  "building": false,
  "stages": [
    {
      "name": "Checkout",
      "status": "SUCCESS"
    },
    {
      "name": "Build",
      "status": "SUCCESS"
    }
  ]
}
```

### URL Building

#### Job Path Inference
```python
def infer_job_path_from_pr(pr_number: str) -> str:
    """
    Default: intuit-authentication/job/auth-service-security/job/authentication-service-pr-ci
    
    This is hardcoded for now but could be configurable
    """
```

#### PR URL Building
```python
def build_pr_url(pr_number: str) -> str:
    """
    Returns: https://github.com/IntuitDeveloper/authentication-service/pull/{pr_number}
    
    This is hardcoded for now but could be configurable
    """
```

## Error Handling

### Network Errors
- If Jenkins API is unreachable: Show tile with ERROR status and message "Connection failed"
- If Jenkins returns 404: Show "Build not found"
- If authentication fails: Show "Authentication failed - check JENKINS_TOKEN"

### Invalid Input
- PR number must be numeric (no letters)
- Empty input is ignored
- Leading/trailing whitespace is trimmed

### Missing Configuration
- If JENKINS_TOKEN not set: Show warning notification on startup
- App should still run but API calls will fail

## Testing Features

### Test Builds
For development/testing, the app includes a helper to create test builds:
```python
def create_test_builds() -> list[Build]:
    """Returns 3-5 sample builds with different statuses for testing UI"""
```

These should be automatically added on startup (can be disabled with a flag).

## Visual Design

### Color Scheme
- Dark theme (dark background)
- Status colors:
  - Success: Green (#00FF00 or similar)
  - Failure: Red (#FF0000 or similar)
  - Running: Blue (#0000FF or similar) with blinking effect
  - Pending: Yellow (#FFFF00 or similar)
- Selected tile: Bold text
- Status bar: Panel background color

### Box Drawing Characters
Tiles use Unicode box-drawing characters:
- ┌ ─ ┐ (top)
- │ (sides)
- ├ ─ ┤ (separator)
- └ ─ ┘ (bottom)

### Text Formatting
- PR numbers: Centered
- Labels: Left-aligned (e.g., "Stage:", "Job:")
- Values: Left-aligned after label
- Build number: Right-aligned at bottom

## Performance Requirements

### Responsiveness
- UI should remain responsive during API calls
- API calls must be non-blocking (async/background)
- Grid refresh should be smooth (no flashing)

### Refresh Rate
- Build status: Poll every 10 seconds per tile
- UI updates: Immediate on user action
- Animation: Blink running builds every 0.8 seconds

### Resource Usage
- Should support 10-20 tiles without noticeable lag
- HTTP connections should be reused/pooled
- Old tiles should clean up timers when removed

## Configuration

### Environment Variables
```bash
JENKINS_TOKEN=your_token_here  # Required for API access
```

### Future Configuration (Optional)
```bash
JENKINS_BASE_URL=https://build.intuit.com
GITHUB_BASE_URL=https://github.com
DEFAULT_JOB_PATH=intuit-authentication/job/...
DEFAULT_REPO=IntuitDeveloper/authentication-service
```

## Edge Cases

### Empty State
- No builds added yet
- Show grid with no tiles
- Status: "Press 'a' to add a PR build, arrow keys to navigate"

### Single Tile
- Selection should work
- Deletion should leave empty grid
- Navigation keys should not crash

### Many Tiles
- Grid should scroll
- All tiles should be accessible via scrolling
- Performance should remain acceptable

### Long Running Builds
- Duration should format properly (e.g., "1h 23m 45s")
- Blinking should continue indefinitely while running
- Status should update when build completes

### Network Interruption
- If API call fails mid-operation: Show error in tile
- Retry logic: User can press 'e' to re-fetch
- Don't crash on network errors

## Implementation Notes for Bubbletea

### Recommended Architecture
```
main.go
├── model.go          # Tea model, state management
├── view.go           # Rendering logic
├── update.go         # Event handling, key bindings
├── api.go            # Jenkins API client
├── tile.go           # Build tile component
└── grid.go           # Grid layout logic
```

### Key Bubbletea Patterns
1. **State Management**: Use Tea model pattern
2. **Async API Calls**: Use Tea commands (Cmd) for non-blocking
3. **Grid Layout**: Use lipgloss for styling and layout
4. **Input Handling**: Use bubbles/textinput for PR input
5. **Status Bar**: Use lipgloss styling at bottom
6. **Colors**: Use lipgloss color definitions
7. **Timers**: Use tea.Tick for refresh and blinking

### Suggested Libraries
- `github.com/charmbracelet/bubbletea` - TUI framework
- `github.com/charmbracelet/lipgloss` - Styling
- `github.com/charmbracelet/bubbles` - Pre-built components (textinput)
- Standard library `net/http` - HTTP client

## Success Criteria

The rewrite is successful when:
1. ✅ All test builds display with correct colors on startup
2. ✅ Can add new PR builds via 'a' key
3. ✅ Tiles show correct information (PR, stage, job, time, build #)
4. ✅ Arrow keys navigate between tiles
5. ✅ Selected tile is visually distinct (bold)
6. ✅ 'd' key deletes selected tile
7. ✅ Enter opens Jenkins URL in browser
8. ✅ 'p' opens GitHub PR in browser
9. ✅ Running builds blink/pulse
10. ✅ Build status refreshes every 10 seconds
11. ✅ Status bar shows helpful messages
12. ✅ No crashes on network errors or invalid input
13. ✅ UI remains responsive during API calls

