"""Data models for Jenkins builds and dashboard state."""
from dataclasses import dataclass, field
from enum import Enum
from typing import Optional
from datetime import datetime, timedelta


class BuildStatus(str, Enum):
    """Build status enumeration."""
    RUNNING = "running"
    SUCCESS = "success"
    FAILURE = "failure"
    UNSTABLE = "unstable"
    ABORTED = "aborted"
    PENDING = "pending"
    ERROR = "error"


@dataclass
class Build:
    """Represents a Jenkins build."""
    pr_number: str
    build_number: Optional[int] = None
    status: BuildStatus = BuildStatus.PENDING
    stage: Optional[str] = None
    job_name: Optional[str] = None
    duration: Optional[timedelta] = None
    build_url: Optional[str] = None
    pr_url: Optional[str] = None
    job_path: Optional[str] = None
    last_updated: datetime = field(default_factory=datetime.now)
    error_message: Optional[str] = None
    
    def is_running(self) -> bool:
        """Check if build is currently running."""
        return self.status == BuildStatus.RUNNING
    
    def is_success(self) -> bool:
        """Check if build succeeded."""
        return self.status == BuildStatus.SUCCESS
    
    def is_failure(self) -> bool:
        """Check if build failed."""
        return self.status in (BuildStatus.FAILURE, BuildStatus.ABORTED, BuildStatus.ERROR)
    
    def format_duration(self) -> str:
        """Format duration as human-readable string."""
        if not self.duration:
            return "N/A"
        
        total_seconds = int(self.duration.total_seconds())
        hours, remainder = divmod(total_seconds, 3600)
        minutes, seconds = divmod(remainder, 60)
        
        if hours > 0:
            return f"{hours}h {minutes}m {seconds}s"
        elif minutes > 0:
            return f"{minutes}m {seconds}s"
        else:
            return f"{seconds}s"


@dataclass
class DashboardState:
    """Manages dashboard state including tiles and selection."""
    builds: list[Build] = field(default_factory=list)
    selected_index: int = 0
    
    def add_build(self, build: Build) -> None:
        """Add a build to the dashboard."""
        self.builds.append(build)
    
    def remove_build(self, index: int) -> bool:
        """Remove a build at the given index."""
        if 0 <= index < len(self.builds):
            self.builds.pop(index)
            # Adjust selected index if needed
            if self.selected_index >= len(self.builds) and len(self.builds) > 0:
                self.selected_index = len(self.builds) - 1
            elif len(self.builds) == 0:
                self.selected_index = 0
            return True
        return False
    
    def get_selected_build(self) -> Optional[Build]:
        """Get the currently selected build."""
        if 0 <= self.selected_index < len(self.builds):
            return self.builds[self.selected_index]
        return None
    
    def move_selection(self, direction: str) -> None:
        """Move selection in the given direction."""
        if len(self.builds) == 0:
            return
        
        if direction == "up":
            self.selected_index = max(0, self.selected_index - 1)
        elif direction == "down":
            self.selected_index = min(len(self.builds) - 1, self.selected_index + 1)
        elif direction == "left":
            self.selected_index = max(0, self.selected_index - 1)
        elif direction == "right":
            self.selected_index = min(len(self.builds) - 1, self.selected_index + 1)

