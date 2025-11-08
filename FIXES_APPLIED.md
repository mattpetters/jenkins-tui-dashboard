# Fixes Applied for TUI Issues

## Issues Reported

### ✅ Issue 1: 'a' key adds tile with text entry
**Status:** WORKING
- Input field appears when you press 'a'
- Keyboard bindings moved to App level for reliability

### ✅ Issue 2: Entering valid PR number shows nothing
**Fixed:**
- Added white text on colored backgrounds (was black text on dark colors)
- Added loading feedback: "Fetching build for PR-3859..."
- Added success feedback: "✓ Added PR-3859 - Use arrow keys to navigate..."
- Changed "Unknown Stage/Job" to "Loading..." for better UX
- Tiles now show visible text in all color states

**Text Colors:**
- White text on blue (running), green (success), red (failure)
- Black text on yellow (pending)

### ✅ Issue 3: Adding more than one tile makes screen go black
**Fixed:**
- Set fixed tile dimensions (width: 34, height: 10)
- Added proper grid gutters and padding
- Ensured tiles have consistent sizing
- Background colors are now properly contained to each tile

## Additional Improvements

1. **Tile Visibility**
   - All text now has proper contrast with backgrounds
   - Fixed-width tiles (32 chars) for consistency
   - Centered PR number for better aesthetics

2. **Feedback Messages**
   - Status bar shows feedback for all actions
   - Loading indicator when fetching builds
   - Success/error messages with checkmarks (✓/✗)
   - Helpful hints in status messages

3. **Grid Layout**
   - Proper spacing between tiles (grid-gutter: 1 2)
   - Grid padding for better layout
   - Fixed tile dimensions prevent layout collapse

4. **Build Status Indicators**
   - Blue background (blinking) for running
   - Green background for success
   - Red background for failure
   - Yellow background for pending
   - All with white/black text for visibility

## Testing

Run the dashboard:
```bash
./run.sh
```

Then:
1. Press `a` - Input appears
2. Type `3859` - You can see what you're typing
3. Press Enter - Status bar shows "Fetching..."
4. Tile appears with visible text on colored background
5. Press `a` again to add another tile
6. Multiple tiles display correctly side-by-side

All 57 tests still pass after these changes!

