# Testing Guide

## Overview

This document describes the test suite for the Jenkins TUI Dashboard. The tests cover all major components and features.

## Running Tests

```bash
# Run all tests
./run_tests.sh

# Run specific test file
pytest tests/test_models.py -v

# Run with coverage
pytest tests/ -v --cov=src --cov-report=html

# Run specific test
pytest tests/test_models.py::TestBuild::test_build_creation -v
```

## Test Coverage

### 1. `test_models.py` - Data Models
Tests for core data structures:
- **Build model**: Creation, status checks, duration formatting
- **DashboardState**: Adding/removing builds, selection management, navigation

**Key Tests:**
- `test_build_creation` - Verify build object creation
- `test_build_is_running` - Test running status detection
- `test_build_is_success` - Test success status detection
- `test_build_is_failure` - Test failure status detection
- `test_format_duration` - Test duration formatting (seconds, minutes, hours)
- `test_add_build` - Test adding builds to dashboard
- `test_remove_build` - Test removing builds
- `test_move_selection` - Test arrow key navigation

### 2. `test_url_builder.py` - URL Utilities
Tests for URL construction:
- PR number parsing (with/without PR- prefix)
- Jenkins build URL construction
- GitHub PR URL construction
- Job path inference

**Key Tests:**
- `test_parse_pr_number` - Parse PR numbers correctly
- `test_build_jenkins_url_basic` - Build Jenkins URLs
- `test_build_pr_url_defaults` - Build PR URLs with defaults

### 3. `test_jenkins_client.py` - Jenkins API Client
Tests for Jenkins API integration:
- Client initialization
- Build status fetching
- Error handling
- Data parsing

**Key Tests:**
- `test_client_initialization` - Verify client setup
- `test_get_build_status_success` - Fetch build status
- `test_get_build_status_failure` - Handle API errors
- `test_create_build_from_pr` - Create build from PR number
- `test_parse_build_data_*` - Parse different build states

### 4. `test_dashboard.py` - Dashboard Component
Tests for dashboard functionality:
- Keyboard shortcuts (a, e, d, p, enter, arrows)
- Build management (add, edit, delete)
- Navigation
- Browser opening
- Input handling

**Key Tests:**
- `test_action_add_build` - Test 'a' key to add build
- `test_action_edit_build` - Test 'e' key to edit build
- `test_action_delete_build` - Test 'd' key to delete build
- `test_action_open_build` - Test Enter key to open build
- `test_action_open_pr` - Test 'p' key to open PR
- `test_action_move_*` - Test arrow key navigation
- `test_on_key_shortcuts` - Test all keyboard shortcuts

### 5. `test_build_tile.py` - Build Tile Component
Tests for tile rendering and behavior:
- Tile initialization
- Status rendering (running, success, failure)
- Blinking animation
- Selection highlighting

**Key Tests:**
- `test_tile_initialization` - Verify tile creation
- `test_render_content_running` - Render running builds (blue)
- `test_render_content_success` - Render successful builds (green)
- `test_render_content_failure` - Render failed builds (red)
- `test_start_blinking` - Test blink animation start
- `test_stop_blinking` - Test blink animation stop

## Keyboard Shortcut Tests

All keyboard shortcuts are tested in `test_dashboard.py`:

- **'a'** - Add build: `test_action_add_build`
- **'e'** - Edit build: `test_action_edit_build`
- **'d'** - Delete build: `test_action_delete_build`
- **Enter** - Open build: `test_action_open_build`
- **'p'** - Open PR: `test_action_open_pr`
- **Arrow keys** - Navigation: `test_action_move_up/down/left/right`
- **Esc** - Cancel input: `test_on_key_escape`

## Fixes Applied

### Keyboard Binding Fix
The original implementation used Textual's BINDINGS which don't work well with Container widgets. Fixed by:
1. Making Dashboard focusable (`can_focus = True`)
2. Adding `self.focus()` in `on_mount()`
3. Handling key events directly in `on_key()` method

### Test Fixtures
- Fixed BuildTile tests to initialize `_blink_timer`
- Fixed Dashboard tests to properly mock the `app` property
- Fixed Jenkins client tests to handle API response parsing

## Coverage Report

Current coverage: ~48% (284/544 statements)

Areas with good coverage:
- Models (93%)
- URL Builder (100%)
- Colors (92%)

Areas needing more coverage:
- Dashboard component (22%) - UI interactions are hard to test
- Build tile component (34%) - Rendering logic needs more tests
- Main app (0%) - Entry point typically not unit tested

## Future Test Improvements

1. Integration tests for full user workflows
2. Mock Jenkins API responses more comprehensively
3. Test refresh timer functionality
4. Test grid layout and tile positioning
5. Test error handling and edge cases

