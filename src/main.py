"""Main application entry point for Jenkins TUI Dashboard."""
import os
import sys
from pathlib import Path
from dotenv import load_dotenv
from textual.app import App, ComposeResult
from textual.containers import Container
from textual.widgets import Header, Footer

from .components.dashboard import Dashboard
from .utils.colors import Colors


class JenkinsDashboardApp(App):
    """Main application class."""
    
    CSS = """
    Screen {
        background: $background;
    }
    
    #main-dashboard {
        width: 100%;
        height: 100%;
    }
    
    #dashboard-container {
        width: 100%;
        height: 100%;
        layout: vertical;
    }
    
    #build-grid {
        width: 100%;
        height: auto;
        grid-size-columns: 3;
        grid-gutter: 1 2;
        padding: 1;
    }
    
    #status-bar {
        width: 100%;
        height: 3;
        padding: 1;
        background: $panel;
        color: $text;
        text-align: center;
        dock: bottom;
    }
    
    BuildTile {
        width: 36;
        height: 11;
        min-width: 36;
        min-height: 11;
        border: none;
        background: transparent;
    }
    
    #pr-input {
        width: 60%;
        height: 3;
        margin: 1 2;
        border: tall $primary;
    }
    """
    
    TITLE = "Jenkins Build Dashboard"
    BINDINGS = [
        ("a", "add_build", "Add PR"),
        ("e", "edit_build", "Edit"),
        ("d", "delete_build", "Delete"),
        ("enter", "open_build", "Open Build"),
        ("p", "open_pr", "Open PR"),
        ("up", "move_up", "Up"),
        ("down", "move_down", "Down"),
        ("left", "move_left", "Left"),
        ("right", "move_right", "Right"),
        ("q", "quit", "Quit"),
    ]
    
    def __init__(self, *args, **kwargs):
        """Initialize the app."""
        super().__init__(*args, **kwargs)
        # Load environment variables
        env_path = Path(__file__).parent.parent / ".env"
        if env_path.exists():
            load_dotenv(env_path)
        else:
            # Try loading from current directory
            load_dotenv()
    
    def compose(self) -> ComposeResult:
        """Compose the application layout."""
        yield Header(show_clock=True)
        yield Dashboard(id="main-dashboard")
        yield Footer()
    
    def on_mount(self) -> None:
        """Handle app mount."""
        # Set theme colors
        self.dark = True
        
        # Check for credentials
        if not os.getenv("JENKINS_TOKEN"):
            self.notify(
                "Warning: JENKINS_TOKEN not set. Please configure .env file.",
                severity="warning"
            )
        
        # Get dashboard and focus it
        dashboard = self.query_one("#main-dashboard", Dashboard)
        # Focus the dashboard so it can receive key events
        self.call_after_refresh(dashboard.focus)
        
        # Pre-populate with PR-3859 if no builds exist (for testing)
        if len(dashboard.state.builds) == 0:
            # Optionally add PR-3859 as default for testing
            # Uncomment the next line to auto-add PR-3859 on startup
            # dashboard._add_build_from_input("3859")
            pass
    
    def action_add_build(self) -> None:
        """Add a new build - delegate to dashboard."""
        dashboard = self.query_one("#main-dashboard", Dashboard)
        dashboard.action_add_build()
    
    def action_edit_build(self) -> None:
        """Edit build - delegate to dashboard."""
        dashboard = self.query_one("#main-dashboard", Dashboard)
        dashboard.action_edit_build()
    
    def action_delete_build(self) -> None:
        """Delete build - delegate to dashboard."""
        dashboard = self.query_one("#main-dashboard", Dashboard)
        dashboard.action_delete_build()
    
    def action_open_build(self) -> None:
        """Open build - delegate to dashboard."""
        dashboard = self.query_one("#main-dashboard", Dashboard)
        dashboard.action_open_build()
    
    def action_open_pr(self) -> None:
        """Open PR - delegate to dashboard."""
        dashboard = self.query_one("#main-dashboard", Dashboard)
        dashboard.action_open_pr()
    
    def action_move_up(self) -> None:
        """Move up - delegate to dashboard."""
        dashboard = self.query_one("#main-dashboard", Dashboard)
        dashboard.action_move_up()
    
    def action_move_down(self) -> None:
        """Move down - delegate to dashboard."""
        dashboard = self.query_one("#main-dashboard", Dashboard)
        dashboard.action_move_down()
    
    def action_move_left(self) -> None:
        """Move left - delegate to dashboard."""
        dashboard = self.query_one("#main-dashboard", Dashboard)
        dashboard.action_move_left()
    
    def action_move_right(self) -> None:
        """Move right - delegate to dashboard."""
        dashboard = self.query_one("#main-dashboard", Dashboard)
        dashboard.action_move_right()
    
    def action_quit(self) -> None:
        """Quit the application."""
        self.exit()


def main():
    """Main entry point."""
    app = JenkinsDashboardApp()
    app.run()


if __name__ == "__main__":
    main()

