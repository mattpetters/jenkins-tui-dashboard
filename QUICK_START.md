# Quick Start Guide

## Setup (One-time)

```bash
cd /Users/mpetters/code/jenkins-dash

# 1. Install dependencies (already done!)
source venv/bin/activate

# 2. Verify credentials are set
cat .env  # Should show your Jenkins credentials
```

## Running the Dashboard

```bash
# From the jenkins-dash directory
./run.sh
```

## Using the Dashboard

### Adding Your First Build (PR-3859)

1. Press `a` - An input field appears at the top
2. Type `3859` (or `PR-3859`)
3. Press `Enter`
4. The tile appears and starts auto-refreshing every 10 seconds

### Keyboard Shortcuts

| Key | Action |
|-----|--------|
| `a` | Add a new PR build |
| `e` | Edit selected build |
| `d` | Delete selected build |
| `Enter` | Open build in browser |
| `p` | Open PR page in browser |
| `↑↓←→` | Navigate between tiles |
| `Esc` | Cancel input mode |
| `q` | Quit application |

### Tile Status Colors

- 🔵 **Blinking Blue** - Build is running
- 🟢 **Green** - Build succeeded
- 🔴 **Red** - Build failed
- 🟡 **Yellow** - Build pending/unknown

### Tile Information

Each tile shows:
```
┌──────────────────────────────┐
│ PR-3859                      │  ← PR number
├──────────────────────────────┤
│ Stage: Test                  │  ← Current stage
│ Job:   account-eks           │  ← Job name
│ Time:  5m 30s                │  ← Time spent
│                        #263 │  ← Build number
└──────────────────────────────┘
```

## Troubleshooting

### Keys not working?
- Make sure the terminal window is focused
- The dashboard automatically focuses on startup
- You should see the footer showing available commands

### Build not showing?
- Check your `.env` file has correct credentials
- The dashboard will show "Error" status if it can't fetch the build
- Try pressing `d` to delete and `a` to re-add

### Tests

Run the test suite:
```bash
./run_tests.sh
# or
source venv/bin/activate
pytest tests/ -v
```

All 57 tests should pass!

