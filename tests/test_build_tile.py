"""Unit tests for BuildTile component."""
import pytest
from unittest.mock import Mock, patch
from src.components.build_tile import BuildTile
from src.models import Build, BuildStatus


class TestBuildTile:
    """Test BuildTile component."""
    
    @pytest.fixture
    def build(self):
        """Create a test build."""
        return Build(
            pr_number="3859",
            build_number=263,
            status=BuildStatus.SUCCESS,
            stage="Test",
            job_name="account-eks",
            duration=None
        )
    
    @pytest.fixture
    def tile(self, build):
        """Create a test tile."""
        tile = BuildTile(build)
        # Use __dict__ to bypass property setter
        tile.__dict__['app'] = Mock()
        # Mock set_interval for blinking tests
        tile.set_interval = Mock(return_value=Mock())
        return tile
    
    def test_tile_initialization(self, build):
        """Test tile initialization."""
        tile = BuildTile(build, is_selected=False)
        # Ensure _blink_timer is initialized
        assert hasattr(tile, '_blink_timer')
        assert tile.build.pr_number == "3859"
        assert tile.is_selected is False
        assert tile.blink_state is False
    
    def test_tile_selected(self, build):
        """Test tile with selected state."""
        tile = BuildTile(build, is_selected=True)
        # Ensure _blink_timer is initialized
        assert hasattr(tile, '_blink_timer')
        assert tile.is_selected is True
    
    def test_render_content_running(self):
        """Test rendering running build."""
        build = Build(pr_number="3859", status=BuildStatus.RUNNING)
        tile = BuildTile(build)
        # Ensure _blink_timer is initialized
        assert hasattr(tile, '_blink_timer')
        content = tile._render_content()
        
        assert "PR-3859" in content
        assert "blue" in content or "bright_blue" in content
    
    def test_render_content_success(self, build):
        """Test rendering successful build."""
        tile = BuildTile(build)
        # Ensure _blink_timer is initialized
        assert hasattr(tile, '_blink_timer')
        content = tile._render_content()
        
        assert "PR-3859" in content
        assert "green" in content
        assert "#263" in content
    
    def test_render_content_failure(self):
        """Test rendering failed build."""
        build = Build(pr_number="3859", status=BuildStatus.FAILURE)
        tile = BuildTile(build)
        # Ensure _blink_timer is initialized
        assert hasattr(tile, '_blink_timer')
        content = tile._render_content()
        
        assert "PR-3859" in content
        assert "red" in content
    
    def test_render_content_with_stage_and_job(self, build):
        """Test rendering with stage and job info."""
        tile = BuildTile(build)
        # Ensure _blink_timer is initialized
        assert hasattr(tile, '_blink_timer')
        content = tile._render_content()
        
        assert "Stage:" in content
        assert "Job:" in content
        assert "Test" in content
        assert "account-eks" in content
    
    def test_start_blinking(self, tile):
        """Test starting blink animation."""
        tile.build.status = BuildStatus.RUNNING
        tile._start_blinking()
        assert tile._blink_timer is not None
    
    def test_stop_blinking(self, tile):
        """Test stopping blink animation."""
        tile.build.status = BuildStatus.RUNNING
        tile._start_blinking()
        tile._stop_blinking()
        assert tile._blink_timer is None
        assert tile.blink_state is False
    
    def test_toggle_blink(self, tile):
        """Test toggling blink state."""
        initial_state = tile.blink_state
        tile._toggle_blink()
        assert tile.blink_state != initial_state
    
    def test_watch_build_running(self, tile):
        """Test watching build change to running."""
        new_build = Build(pr_number="3859", status=BuildStatus.RUNNING)
        with patch.object(tile, '_start_blinking') as mock_start:
            tile.watch_build(new_build)
            mock_start.assert_called_once()
    
    def test_watch_build_not_running(self, tile):
        """Test watching build change to not running."""
        new_build = Build(pr_number="3859", status=BuildStatus.SUCCESS)
        with patch.object(tile, '_stop_blinking') as mock_stop:
            tile.watch_build(new_build)
            mock_stop.assert_called_once()
    
    def test_watch_is_selected(self, tile):
        """Test watching selection changes."""
        with patch.object(tile, 'update') as mock_update:
            tile.watch_is_selected(True)
            mock_update.assert_called_once()
    
    def test_watch_blink_state(self, tile):
        """Test watching blink state changes."""
        tile.build.status = BuildStatus.RUNNING
        tile.update = Mock()  # Mock update method
        tile.watch_blink_state(True)
        # Assert update was called
        assert tile.update.called

