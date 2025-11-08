"""Unit tests for data models."""
import pytest
from datetime import timedelta
from src.models import Build, BuildStatus, DashboardState


class TestBuild:
    """Test Build model."""
    
    def test_build_creation(self):
        """Test creating a build."""
        build = Build(pr_number="3859")
        assert build.pr_number == "3859"
        assert build.status == BuildStatus.PENDING
        assert build.build_number is None
    
    def test_build_is_running(self):
        """Test is_running method."""
        build = Build(pr_number="3859", status=BuildStatus.RUNNING)
        assert build.is_running() is True
        
        build.status = BuildStatus.SUCCESS
        assert build.is_running() is False
    
    def test_build_is_success(self):
        """Test is_success method."""
        build = Build(pr_number="3859", status=BuildStatus.SUCCESS)
        assert build.is_success() is True
        
        build.status = BuildStatus.FAILURE
        assert build.is_success() is False
    
    def test_build_is_failure(self):
        """Test is_failure method."""
        build = Build(pr_number="3859", status=BuildStatus.FAILURE)
        assert build.is_failure() is True
        
        build.status = BuildStatus.ABORTED
        assert build.is_failure() is True
        
        build.status = BuildStatus.ERROR
        assert build.is_failure() is True
        
        build.status = BuildStatus.SUCCESS
        assert build.is_failure() is False
    
    def test_format_duration(self):
        """Test duration formatting."""
        build = Build(pr_number="3859")
        assert build.format_duration() == "N/A"
        
        build.duration = timedelta(seconds=45)
        assert build.format_duration() == "45s"
        
        build.duration = timedelta(minutes=5, seconds=30)
        assert build.format_duration() == "5m 30s"
        
        build.duration = timedelta(hours=2, minutes=15, seconds=30)
        assert build.format_duration() == "2h 15m 30s"


class TestDashboardState:
    """Test DashboardState model."""
    
    def test_initial_state(self):
        """Test initial dashboard state."""
        state = DashboardState()
        assert len(state.builds) == 0
        assert state.selected_index == 0
    
    def test_add_build(self):
        """Test adding a build."""
        state = DashboardState()
        build = Build(pr_number="3859")
        state.add_build(build)
        assert len(state.builds) == 1
        assert state.builds[0].pr_number == "3859"
    
    def test_remove_build(self):
        """Test removing a build."""
        state = DashboardState()
        build1 = Build(pr_number="3859")
        build2 = Build(pr_number="3860")
        state.add_build(build1)
        state.add_build(build2)
        
        assert len(state.builds) == 2
        assert state.remove_build(0) is True
        assert len(state.builds) == 1
        assert state.builds[0].pr_number == "3860"
    
    def test_remove_build_invalid_index(self):
        """Test removing build with invalid index."""
        state = DashboardState()
        assert state.remove_build(0) is False
        assert state.remove_build(-1) is False
    
    def test_remove_build_adjusts_selection(self):
        """Test that removing build adjusts selected index."""
        state = DashboardState()
        build1 = Build(pr_number="3859")
        build2 = Build(pr_number="3860")
        state.add_build(build1)
        state.add_build(build2)
        state.selected_index = 1
        
        state.remove_build(1)
        assert state.selected_index == 0
    
    def test_get_selected_build(self):
        """Test getting selected build."""
        state = DashboardState()
        assert state.get_selected_build() is None
        
        build = Build(pr_number="3859")
        state.add_build(build)
        selected = state.get_selected_build()
        assert selected is not None
        assert selected.pr_number == "3859"
    
    def test_move_selection(self):
        """Test moving selection."""
        state = DashboardState()
        build1 = Build(pr_number="3859")
        build2 = Build(pr_number="3860")
        build3 = Build(pr_number="3861")
        state.add_build(build1)
        state.add_build(build2)
        state.add_build(build3)
        
        state.selected_index = 1
        state.move_selection("up")
        assert state.selected_index == 0
        
        state.move_selection("down")
        assert state.selected_index == 1
        
        state.move_selection("down")
        assert state.selected_index == 2
        
        # Should not go beyond bounds
        state.move_selection("down")
        assert state.selected_index == 2
        
        state.move_selection("up")
        assert state.selected_index == 1
        
        state.move_selection("up")
        assert state.selected_index == 0
        
        # Should not go below 0
        state.move_selection("up")
        assert state.selected_index == 0
    
    def test_move_selection_empty(self):
        """Test moving selection with no builds."""
        state = DashboardState()
        state.move_selection("up")
        assert state.selected_index == 0

