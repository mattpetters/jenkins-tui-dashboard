"""Color palette definitions for the dashboard."""
from textual.color import Color


class Colors:
    """Aesthetic color palette for the dashboard."""
    
    # Status colors
    RUNNING_PRIMARY = Color.parse("#4A90E2")
    RUNNING_SECONDARY = Color.parse("#6BA3E8")
    SUCCESS = Color.parse("#2ECC71")
    FAILURE = Color.parse("#E74C3C")
    UNSTABLE = Color.parse("#F39C12")
    ABORTED = Color.parse("#95A5A6")
    PENDING = Color.parse("#9B59B6")
    
    # Background colors
    BG_PRIMARY = Color.parse("#0D1117")
    BG_SECONDARY = Color.parse("#161B22")
    BG_TERTIARY = Color.parse("#1E1E1E")
    
    # Text colors
    TEXT_PRIMARY = Color.parse("#F0F0F0")
    TEXT_SECONDARY = Color.parse("#8B949E")
    TEXT_MUTED = Color.parse("#6E7681")
    
    # Border colors
    BORDER_DEFAULT = Color.parse("#30363D")
    BORDER_SELECTED = Color.parse("#58A6FF")
    BORDER_HOVER = Color.parse("#6BA3E8")
    
    # Status color mapping
    STATUS_COLORS = {
        "running": RUNNING_PRIMARY,
        "success": SUCCESS,
        "failure": FAILURE,
        "unstable": UNSTABLE,
        "aborted": ABORTED,
        "pending": PENDING,
    }
    
    @classmethod
    def get_status_color(cls, status: str) -> Color:
        """Get color for a build status."""
        return cls.STATUS_COLORS.get(status.lower(), cls.PENDING)
    
    @classmethod
    def get_running_color(cls, blink_state: bool) -> Color:
        """Get blinking color for running builds."""
        return cls.RUNNING_SECONDARY if blink_state else cls.RUNNING_PRIMARY

