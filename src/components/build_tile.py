"""Build tile component for displaying individual build status."""
from textual.widgets import Static
from textual.message import Message
from textual.reactive import reactive
from textual import events
from datetime import datetime

from ..models import Build, BuildStatus
from ..utils.colors import Colors


class BuildTile(Static):
    """A tile widget displaying build status information."""
    
    # Reactive attributes
    build: reactive[Build] = reactive(Build(pr_number=""))
    is_selected: reactive[bool] = reactive(False)
    blink_state: reactive[bool] = reactive(False)
    
    # Class-level initialization of _blink_timer to avoid watcher issues
    _blink_timer = None
    
    def __init__(self, build: Build, is_selected: bool = False, *args, **kwargs):
        """Initialize build tile.
        
        Args:
            build: Build object to display
            is_selected: Whether this tile is currently selected
        """
        # Initialize _blink_timer as instance attribute before super()
        # This ensures it exists when watchers fire during reactive init
        object.__setattr__(self, '_blink_timer', None)
        super().__init__(*args, **kwargs)
        self.build = build
        self.is_selected = is_selected
        self.update(self._render_content())
    
    def on_mount(self) -> None:
        """Set up blinking timer for running builds."""
        if self.build.is_running():
            self._start_blinking()
    
    def watch_build(self, build: Build) -> None:
        """Watch for build changes and update blinking."""
        # Ensure _blink_timer is initialized (defensive check)
        if not hasattr(self, '_blink_timer'):
            object.__setattr__(self, '_blink_timer', None)
        
        if build.is_running():
            self._start_blinking()
        else:
            self._stop_blinking()
        self.update(self._render_content())
    
    def watch_is_selected(self, is_selected: bool) -> None:
        """Watch for selection changes."""
        self.update(self._render_content())
    
    def watch_blink_state(self, blink_state: bool) -> None:
        """Watch for blink state changes."""
        if self.build.is_running():
            self.update(self._render_content())
    
    def _start_blinking(self) -> None:
        """Start blinking animation."""
        # Ensure _blink_timer is initialized (defensive check)
        if not hasattr(self, '_blink_timer'):
            object.__setattr__(self, '_blink_timer', None)
        
        if self._blink_timer is None:
            try:
                self._blink_timer = self.set_interval(0.8, self._toggle_blink)
            except:
                pass
    
    def _stop_blinking(self) -> None:
        """Stop blinking animation."""
        # Ensure _blink_timer is initialized (defensive check)
        if not hasattr(self, '_blink_timer'):
            object.__setattr__(self, '_blink_timer', None)
            return
        
        if self._blink_timer:
            try:
                self._blink_timer.stop()
            except:
                pass
            self._blink_timer = None
            self.blink_state = False
    
    def _toggle_blink(self) -> None:
        """Toggle blink state."""
        self.blink_state = not self.blink_state
    
    def _render_content(self) -> str:
        """Render the tile content as text."""
        # Determine background color based on status
        if self.build.is_running():
            bg_color_name = "blue" if not self.blink_state else "bright_blue"
            text_color = "white"
        elif self.build.is_success():
            bg_color_name = "green"
            text_color = "white"
        elif self.build.is_failure():
            bg_color_name = "red"
            text_color = "white"
        else:
            bg_color_name = "yellow"
            text_color = "black"
        
        # Build the tile content
        pr_text = f"PR-{self.build.pr_number}"
        
        # Show status-appropriate text
        if self.build.status == BuildStatus.PENDING:
            stage_text = "Loading..."[:18]
            job_text = "Fetching data..."[:18]
        elif self.build.error_message:
            stage_text = "Error"[:18]
            job_text = self.build.error_message[:18]
        else:
            stage_text = (self.build.stage or "Unknown")[:18]
            job_text = (self.build.job_name or "Unknown")[:18]
        
        duration_text = self.build.format_duration()
        build_num_text = f"#{self.build.build_number}" if self.build.build_number else "..."
        
        # Fixed width for consistent tiles
        width = 32
        
        # Build the tile with borders and styling
        lines = []
        
        # Top border
        top_border = "┌" + "─" * (width - 2) + "┐"
        lines.append(top_border)
        
        # PR number (centered)
        pr_padded = pr_text.center(width - 4)
        pr_line = f"│ {pr_padded} │"
        lines.append(pr_line)
        
        # Separator
        separator = "├" + "─" * (width - 2) + "┤"
        lines.append(separator)
        
        # Stage
        stage_label = "Stage:"
        stage_line = f"│ {stage_label} {stage_text:<{width-len(stage_label)-5}} │"
        lines.append(stage_line)
        
        # Job
        job_label = "Job:"
        job_line = f"│ {job_label} {job_text:<{width-len(job_label)-3}} │"
        lines.append(job_line)
        
        # Duration
        time_label = "Time:"
        duration_line = f"│ {time_label} {duration_text:<{width-len(time_label)-3}} │"
        lines.append(duration_line)
        
        # Build number (bottom right)
        build_num_line = f"│ {'':<{width-len(build_num_text)-5}}{build_num_text} │"
        lines.append(build_num_line)
        
        # Bottom border
        bottom_border = "└" + "─" * (width - 2) + "┘"
        lines.append(bottom_border)
        
        # Combine all lines
        content = "\n".join(lines)
        
        # Apply styling with text color and background color
        if self.is_selected:
            styled_content = f"[bold {text_color} on {bg_color_name}]{content}[/]"
        else:
            styled_content = f"[{text_color} on {bg_color_name}]{content}[/]"
        
        return styled_content
    
    def on_click(self, event: events.Click) -> None:
        """Handle click events."""
        self.post_message(self.Selected(self))
    
    class Selected(Message):
        """Message sent when tile is selected."""
        def __init__(self, tile: "BuildTile"):
            super().__init__()
            self.tile = tile

