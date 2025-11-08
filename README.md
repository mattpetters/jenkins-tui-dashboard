# Jenkins TUI Dashboard

A terminal-based dashboard for monitoring Jenkins builds by PR number.

## Features

- Real-time build status monitoring
- Tiled layout with color-coded status indicators
- Quick add/edit/delete of PR builds
- Navigation with arrow keys
- Auto-refresh every 10 seconds
- Open builds and PRs in browser

## Installation

✅ **Already set up!** The virtual environment and dependencies are installed.

If you need to reinstall:
```bash
# Create virtual environment
python3 -m venv venv

# Activate virtual environment
source venv/bin/activate  # On macOS/Linux

# Install dependencies
pip install -r requirements.txt
```

## Configuration

✅ **Already configured!** Your `.env` file is set up with Jenkins credentials.

If you need to reconfigure:
```bash
cp .env.example .env
# Then edit .env with your credentials
```

## Usage

### Run the application:

```bash
# Option 1: Using the run script
./run.sh

# Option 2: Using Python module
python -m src.main

# Option 3: Direct Python
python src/main.py
```

### Keyboard Shortcuts

- `a` - Add a new PR build
- `e` - Edit selected tile
- `d` - Delete selected tile
- `Enter` - Open build in browser
- `p` - Open PR page in browser
- Arrow keys - Navigate between tiles
- `Esc` - Cancel input mode
- `q` - Quit application

## Example

Start with PR-3859:
1. Press `a` to add a new build
2. Enter `3859` when prompted (or `PR-3859`)
3. Press Enter to submit
4. The tile will appear showing:
   - PR number
   - Current stage name
   - Job name
   - Time spent in current job
   - Build number (bottom right)
5. The tile auto-refreshes every 10 seconds
6. Use arrow keys to navigate
7. Press `Enter` to open the build in your browser
8. Press `p` to open the PR page in your browser

## Features

- **Color-coded status**: 
  - 🔵 Blinking blue for running builds
  - 🟢 Green for successful builds
  - 🔴 Red for failed builds
- **Real-time updates**: Each tile refreshes every 10 seconds
- **Interactive navigation**: Arrow keys to move between tiles
- **Quick actions**: Open builds and PRs directly from the dashboard

## Testing

Run the test suite:

```bash
# Run all tests
./run_tests.sh

# Or manually
source venv/bin/activate
pytest tests/ -v

# With coverage report
pytest tests/ -v --cov=src --cov-report=html
```

The test suite includes:
- Unit tests for data models (`test_models.py`)
- URL builder tests (`test_url_builder.py`)
- Jenkins client tests (`test_jenkins_client.py`)
- Dashboard component tests (`test_dashboard.py`)
- Build tile component tests (`test_build_tile.py`)

