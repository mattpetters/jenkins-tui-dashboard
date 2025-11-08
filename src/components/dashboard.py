"""Main dashboard component with grid layout and navigation."""
import webbrowser
import threading
from textual.app import ComposeResult
from textual.containers import Grid, Container
from textual.widgets import Input, Static, Button
from textual.message import Message
from textual import events
from textual.binding import Binding

from ..models import Build, DashboardState, BuildStatus
from ..jenkins_client import JenkinsClient
from .build_tile import BuildTile


class Dashboard(Container):
    """Main dashboard widget with grid layout."""
    
    def __init__(self, *args, **kwargs):
        """Initialize dashboard."""
        super().__init__(*args, **kwargs)
        self.state = DashboardState()
        self.jenkins_client = JenkinsClient()
        self.tiles: list[BuildTile] = []
        self.refresh_timers: dict[int, object] = {}
        self.input_mode = False
        self.input_widget: Input | None = None
    
    def compose(self) -> ComposeResult:
        """Compose the dashboard layout."""
        with Container(id="dashboard-container"):
            yield Grid(id="build-grid")
            yield Static("", id="status-bar")
    
    def on_mount(self) -> None:
        """Set up dashboard on mount."""
        self.update_status("Press 'a' to add a PR build, arrow keys to navigate")
        self.refresh_grid()
    
    def refresh_grid(self) -> None:
        """Refresh the grid layout with current builds."""
        grid = self.query_one("#build-grid", Grid)
        
        # Remove all existing tiles
        for tile in self.tiles:
            try:
                tile.remove()
            except:
                pass
        
        self.tiles.clear()
        
        # Stop all refresh timers
        for timer in self.refresh_timers.values():
            if timer:
                try:
                    timer.stop()
                except:
                    pass
        self.refresh_timers.clear()
        
        # Create tiles for each build
        for idx, build in enumerate(self.state.builds):
            is_selected = idx == self.state.selected_index
            tile = BuildTile(build, is_selected=is_selected, id=f"tile-{idx}")
            self.tiles.append(tile)
            grid.mount(tile)
            
            # Start refresh timer for this tile
            self._start_refresh_timer(idx)
        
        # Update grid columns based on terminal width
        self._update_grid_layout()
        
        # Update selection highlighting
        self._update_selection()
    
    def _update_grid_layout(self) -> None:
        """Update grid layout based on terminal size."""
        grid = self.query_one("#build-grid", Grid)
        num_tiles = len(self.tiles)
        
        if num_tiles == 0:
            return
        
        # Calculate optimal columns (aim for 3-4 columns)
        try:
            terminal_width = self.app.size.width if hasattr(self.app, 'size') else 120
        except:
            terminal_width = 120
        
        tile_width = 35  # Approximate tile width
        cols = max(1, min(4, terminal_width // tile_width))
        
        # Update grid columns
        try:
            grid.grid_size_columns = cols
        except:
            # Fallback if grid_size_columns doesn't exist
            pass
    
    def _start_refresh_timer(self, tile_index: int) -> None:
        """Start refresh timer for a tile."""
        def refresh_build():
            if tile_index < len(self.state.builds):
                build = self.state.builds[tile_index]
                self.refresh_build(tile_index)
        
        timer = self.set_interval(10.0, refresh_build)
        self.refresh_timers[tile_index] = timer
    
    def refresh_build(self, tile_index: int) -> None:
        """Refresh a specific build."""
        if tile_index >= len(self.state.builds):
            return
        
        build = self.state.builds[tile_index]
        
        # Fetch updated build status
        if build.job_path:
            updated_build = self.jenkins_client.get_build_status(
                build.job_path,
                f"PR-{build.pr_number}",
                build.build_number
            )
            
            if updated_build:
                self.state.builds[tile_index] = updated_build
                # Update the tile
                if tile_index < len(self.tiles):
                    self.tiles[tile_index].build = updated_build
                    self.tiles[tile_index].refresh()
    
    def _update_selection(self) -> None:
        """Update selection highlighting."""
        for idx, tile in enumerate(self.tiles):
            tile.is_selected = idx == self.state.selected_index
            tile.refresh()
    
    def action_add_build(self) -> None:
        """Add a new build."""
        if self.input_mode:
            return
        
        self.input_mode = True
        
        # Get the container
        try:
            container = self.query_one("#dashboard-container", Container)
        except:
            # Fallback: use self as container
            container = self
        
        # Create input widget
        self.input_widget = Input(
            placeholder="Enter PR number (e.g., 3859)",
            id="pr-input"
        )
        container.mount(self.input_widget)
        
        # Focus the input widget after it's mounted
        def focus_input():
            if self.input_widget:
                self.input_widget.focus()
        
        self.call_after_refresh(focus_input)
        
        self.update_status("Enter PR number and press Enter, or Esc to cancel")
    
    def action_edit_build(self) -> None:
        """Edit selected build."""
        if self.input_mode or len(self.state.builds) == 0:
            return
        
        selected = self.state.get_selected_build()
        if not selected:
            return
        
        self.input_mode = True
        
        # Get the container
        try:
            container = self.query_one("#dashboard-container", Container)
        except:
            # Fallback: use self as container
            container = self
        
        # Create input widget with current PR number
        self.input_widget = Input(
            value=selected.pr_number,
            placeholder="Enter PR number",
            id="pr-input"
        )
        container.mount(self.input_widget)
        
        # Focus the input widget after it's mounted
        def focus_input():
            if self.input_widget:
                self.input_widget.focus()
        
        self.call_after_refresh(focus_input)
        
        self.update_status("Edit PR number and press Enter, or Esc to cancel")
    
    def action_delete_build(self) -> None:
        """Delete selected build."""
        if self.input_mode or len(self.state.builds) == 0:
            return
        
        selected = self.state.get_selected_build()
        pr_num = selected.pr_number if selected else "unknown"
        
        if self.state.remove_build(self.state.selected_index):
            self.refresh_grid()
            self.update_status(f"✓ Deleted PR-{pr_num}. {len(self.state.builds)} build(s) remaining")
    
    def action_open_build(self) -> None:
        """Open selected build in browser."""
        if self.input_mode:
            return
        
        selected = self.state.get_selected_build()
        if selected and selected.build_url:
            webbrowser.open(selected.build_url)
            self.update_status(f"Opened build {selected.build_number} in browser")
    
    def action_open_pr(self) -> None:
        """Open selected PR in browser."""
        if self.input_mode:
            return
        
        selected = self.state.get_selected_build()
        if selected and selected.pr_url:
            webbrowser.open(selected.pr_url)
            self.update_status(f"Opened PR-{selected.pr_number} in browser")
    
    def action_move_up(self) -> None:
        """Move selection up."""
        if self.input_mode:
            return
        self.state.move_selection("up")
        self._update_selection()
    
    def action_move_down(self) -> None:
        """Move selection down."""
        if self.input_mode:
            return
        self.state.move_selection("down")
        self._update_selection()
    
    def action_move_left(self) -> None:
        """Move selection left."""
        if self.input_mode:
            return
        self.state.move_selection("left")
        self._update_selection()
    
    def action_move_right(self) -> None:
        """Move selection right."""
        if self.input_mode:
            return
        self.state.move_selection("right")
        self._update_selection()
    
    def on_input_submitted(self, event: Input.Submitted) -> None:
        """Handle input submission."""
        if event.input.id == "pr-input":
            pr_number = event.input.value.strip()
            if pr_number:
                self._add_build_from_input(pr_number)
            
            # Remove input widget
            try:
                event.input.remove()
            except:
                pass
            self.input_widget = None
            self.input_mode = False
            self.update_status("")
            event.stop()
    
    def on_input_changed(self, event: Input.Changed) -> None:
        """Handle input changes."""
        pass
    
    def on_key(self, event: events.Key) -> None:
        """Handle key events."""
        # Handle escape key for canceling input
        if event.key == "escape" and self.input_mode:
            # Cancel input
            if self.input_widget:
                try:
                    self.input_widget.remove()
                except:
                    pass
                self.input_widget = None
            self.input_mode = False
            self.update_status("Press 'a' to add a PR build, arrow keys to navigate")
            event.stop()
            return
        
        # Don't handle other keys if in input mode
        if self.input_mode:
            return
        
        # Handle keyboard shortcuts
        if event.key == "a":
            self.action_add_build()
            event.stop()
        elif event.key == "e":
            self.action_edit_build()
            event.stop()
        elif event.key == "d":
            self.action_delete_build()
            event.stop()
        elif event.key == "enter":
            self.action_open_build()
            event.stop()
        elif event.key == "p":
            self.action_open_pr()
            event.stop()
        elif event.key == "up":
            self.action_move_up()
            event.stop()
        elif event.key == "down":
            self.action_move_down()
            event.stop()
        elif event.key == "left":
            self.action_move_left()
            event.stop()
        elif event.key == "right":
            self.action_move_right()
            event.stop()
    
    def _add_build_from_input(self, pr_number: str) -> None:
        """Add build from PR number input."""
        from .build_tile import BuildTile
        from ..utils.url_builder import infer_job_path_from_pr, parse_pr_number, build_pr_url
        
        try:
            # Parse PR number
            pr_num = parse_pr_number(pr_number)
            
            # Create a loading build immediately (no blocking API call)
            job_path = infer_job_path_from_pr(pr_num)
            loading_build = Build(
                pr_number=pr_num,
                status=BuildStatus.PENDING,
                stage="Loading...",
                job_name="Loading...",
                job_path=job_path,
                pr_url=build_pr_url(pr_num)
            )
            
            # Add to state and display immediately
            self.state.add_build(loading_build)
            self.state.selected_index = len(self.state.builds) - 1
            self.refresh_grid()
            self.update_status(f"✓ Added PR-{pr_num} - Fetching build data...")
            
            # Fetch actual build data in background (non-blocking)
            build_index = len(self.state.builds) - 1
            self._fetch_build_data_async(pr_num, job_path, build_index)
            
        except Exception as e:
            self.update_status(f"✗ Error adding PR-{pr_number}: {str(e)}")
    
    def _fetch_build_data_async(self, pr_number: str, job_path: str, build_index: int) -> None:
        """Fetch build data asynchronously in the background."""
        import threading
        
        def fetch_and_update():
            try:
                # Fetch build status (this may take time)
                pr_branch = f"PR-{pr_number}"
                build = self.jenkins_client.get_build_status(job_path, pr_branch)
                
                # Update the build in the state
                if build and build_index < len(self.state.builds):
                    self.state.builds[build_index] = build
                    
                    # Update the tile display
                    if build_index < len(self.tiles):
                        self.tiles[build_index].build = build
                        self.tiles[build_index].refresh()
                    
                    # Update status bar
                    status_msg = f"✓ PR-{pr_number}: {build.status.value.upper()}"
                    if build.build_number:
                        status_msg += f" (Build #{build.build_number})"
                    self.call_from_thread(lambda: self.update_status(status_msg))
                    
            except Exception as e:
                # Update with error
                if build_index < len(self.state.builds):
                    self.state.builds[build_index].error_message = str(e)
                    self.state.builds[build_index].status = BuildStatus.ERROR
                    if build_index < len(self.tiles):
                        self.tiles[build_index].build = self.state.builds[build_index]
                        self.tiles[build_index].refresh()
                self.call_from_thread(lambda: self.update_status(f"✗ Error fetching PR-{pr_number}"))
        
        # Start background thread
        thread = threading.Thread(target=fetch_and_update, daemon=True)
        thread.start()
    
    def update_status(self, message: str) -> None:
        """Update status bar."""
        status_bar = self.query_one("#status-bar", Static)
        status_bar.update(message)
    
    def on_resize(self) -> None:
        """Handle resize events."""
        self._update_grid_layout()

