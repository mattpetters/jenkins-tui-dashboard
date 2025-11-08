"""Unit tests for Dashboard component."""
import pytest
from unittest.mock import Mock, patch, MagicMock
from textual.app import App
from textual.widgets import Input
from src.components.dashboard import Dashboard
from src.models import Build, BuildStatus, DashboardState
from src.jenkins_client import JenkinsClient


class TestDashboard:
    """Test Dashboard component."""
    
    @pytest.fixture
    def dashboard(self):
        """Create a dashboard instance for testing."""
        # Create a minimal app context
        app = Mock(spec=App)
        dashboard = Dashboard()
        # Use __dict__ to bypass property setter restriction
        dashboard.__dict__['app'] = app
        
        # Mock query_one for status bar and other queries
        def mock_query_one(selector, *args, **kwargs):
            if selector == "#status-bar":
                mock_status = Mock()
                mock_status.update = Mock()
                return mock_status
            elif selector == "#dashboard-container":
                return dashboard
            elif selector == "#build-grid":
                mock_grid = Mock()
                mock_grid.mount = Mock()
                return mock_grid
            return Mock()
        
        dashboard.query_one = mock_query_one
        dashboard.call_after_refresh = Mock()  # Mock call_after_refresh
        return dashboard
    
    def test_dashboard_initialization(self, dashboard):
        """Test dashboard initialization."""
        assert dashboard.state is not None
        assert isinstance(dashboard.state, DashboardState)
        assert dashboard.jenkins_client is not None
        assert len(dashboard.tiles) == 0
        assert dashboard.input_mode is False
    
    def test_action_add_build(self, dashboard):
        """Test adding a build."""
        # Mock the input widget creation
        with patch.object(dashboard, 'query_one') as mock_query:
            mock_container = Mock()
            mock_query.return_value = mock_container
            
            dashboard.action_add_build()
            
            assert dashboard.input_mode is True
            assert dashboard.input_widget is not None
    
    def test_action_delete_build_empty(self, dashboard):
        """Test deleting build when none exist."""
        dashboard.action_delete_build()
        assert len(dashboard.state.builds) == 0
    
    def test_action_delete_build_with_builds(self, dashboard):
        """Test deleting a build."""
        build = Build(pr_number="3859", build_number=263)
        dashboard.state.add_build(build)
        
        with patch.object(dashboard, 'refresh_grid') as mock_refresh, \
             patch.object(dashboard, 'update_status') as mock_update:
            dashboard.action_delete_build()
            mock_refresh.assert_called_once()
            assert len(dashboard.state.builds) == 0
    
    def test_action_open_build_no_selection(self, dashboard):
        """Test opening build with no selection."""
        with patch('webbrowser.open') as mock_open:
            dashboard.action_open_build()
            mock_open.assert_not_called()
    
    def test_action_open_build_with_selection(self, dashboard):
        """Test opening build in browser."""
        build = Build(
            pr_number="3859",
            build_number=263,
            build_url="https://build.intuit.com/test/263"
        )
        dashboard.state.add_build(build)
        
        with patch('webbrowser.open') as mock_open, \
             patch.object(dashboard, 'update_status') as mock_update:
            dashboard.action_open_build()
            mock_open.assert_called_once_with("https://build.intuit.com/test/263")
    
    def test_action_open_pr_no_selection(self, dashboard):
        """Test opening PR with no selection."""
        with patch('webbrowser.open') as mock_open:
            dashboard.action_open_pr()
            mock_open.assert_not_called()
    
    def test_action_open_pr_with_selection(self, dashboard):
        """Test opening PR in browser."""
        build = Build(
            pr_number="3859",
            pr_url="https://github.intuit.com/identity-manage/account/pull/3859"
        )
        dashboard.state.add_build(build)
        
        with patch('webbrowser.open') as mock_open, \
             patch.object(dashboard, 'update_status') as mock_update:
            dashboard.action_open_pr()
            mock_open.assert_called_once_with(
                "https://github.intuit.com/identity-manage/account/pull/3859"
            )
    
    def test_action_move_up(self, dashboard):
        """Test moving selection up."""
        build1 = Build(pr_number="3859")
        build2 = Build(pr_number="3860")
        dashboard.state.add_build(build1)
        dashboard.state.add_build(build2)
        dashboard.state.selected_index = 1
        
        with patch.object(dashboard, '_update_selection') as mock_update:
            dashboard.action_move_up()
            assert dashboard.state.selected_index == 0
            mock_update.assert_called_once()
    
    def test_action_move_down(self, dashboard):
        """Test moving selection down."""
        build1 = Build(pr_number="3859")
        build2 = Build(pr_number="3860")
        dashboard.state.add_build(build1)
        dashboard.state.add_build(build2)
        dashboard.state.selected_index = 0
        
        with patch.object(dashboard, '_update_selection') as mock_update:
            dashboard.action_move_down()
            assert dashboard.state.selected_index == 1
            mock_update.assert_called_once()
    
    def test_action_move_left(self, dashboard):
        """Test moving selection left."""
        build1 = Build(pr_number="3859")
        build2 = Build(pr_number="3860")
        dashboard.state.add_build(build1)
        dashboard.state.add_build(build2)
        dashboard.state.selected_index = 1
        
        with patch.object(dashboard, '_update_selection') as mock_update:
            dashboard.action_move_left()
            assert dashboard.state.selected_index == 0
            mock_update.assert_called_once()
    
    def test_action_move_right(self, dashboard):
        """Test moving selection right."""
        build1 = Build(pr_number="3859")
        build2 = Build(pr_number="3860")
        dashboard.state.add_build(build1)
        dashboard.state.add_build(build2)
        dashboard.state.selected_index = 0
        
        with patch.object(dashboard, '_update_selection') as mock_update:
            dashboard.action_move_right()
            assert dashboard.state.selected_index == 1
            mock_update.assert_called_once()
    
    def test_add_build_from_input(self, dashboard):
        """Test adding build from input."""
        with patch.object(dashboard.jenkins_client, 'create_build_from_pr') as mock_create:
            mock_build = Build(pr_number="3859", build_number=263)
            mock_create.return_value = mock_build
            
            with patch.object(dashboard, 'refresh_grid') as mock_refresh, \
                 patch.object(dashboard, 'update_status') as mock_update:
                dashboard._add_build_from_input("3859")
                mock_create.assert_called_once_with("3859")
                mock_refresh.assert_called_once()
                assert len(dashboard.state.builds) == 1
    
    def test_on_input_submitted(self, dashboard):
        """Test handling input submission."""
        dashboard.input_mode = True
        mock_input = Mock()
        mock_input.id = "pr-input"
        mock_input.value = "3859"
        mock_input.remove = Mock()
        
        event = Mock()
        event.input = mock_input
        event.stop = Mock()
        
        with patch.object(dashboard, '_add_build_from_input') as mock_add, \
             patch.object(dashboard, 'update_status') as mock_update:
            dashboard.on_input_submitted(event)
            mock_add.assert_called_once_with("3859")
            assert dashboard.input_mode is False
            event.stop.assert_called_once()
    
    def test_on_key_escape(self, dashboard):
        """Test handling escape key."""
        dashboard.input_mode = True
        mock_input = Mock()
        mock_input.remove = Mock()
        dashboard.input_widget = mock_input
        
        event = Mock()
        event.key = "escape"
        event.stop = Mock()
        
        with patch.object(dashboard, 'update_status') as mock_update:
            dashboard.on_key(event)
        
        assert dashboard.input_mode is False
        assert dashboard.input_widget is None
        event.stop.assert_called_once()
    
    def test_on_key_shortcuts(self, dashboard):
        """Test keyboard shortcuts."""
        event = Mock()
        event.stop = Mock()
        
        # Test 'a' key
        event.key = "a"
        with patch.object(dashboard, 'action_add_build') as mock_add:
            dashboard.on_key(event)
            mock_add.assert_called_once()
            event.stop.assert_called_once()
        
        # Test 'd' key
        event.key = "d"
        event.stop.reset_mock()
        with patch.object(dashboard, 'action_delete_build') as mock_delete:
            dashboard.on_key(event)
            mock_delete.assert_called_once()
            event.stop.assert_called_once()
        
        # Test 'e' key
        event.key = "e"
        event.stop.reset_mock()
        with patch.object(dashboard, 'action_edit_build') as mock_edit:
            dashboard.on_key(event)
            mock_edit.assert_called_once()
            event.stop.assert_called_once()
        
        # Test 'p' key
        event.key = "p"
        event.stop.reset_mock()
        with patch.object(dashboard, 'action_open_pr') as mock_open_pr:
            dashboard.on_key(event)
            mock_open_pr.assert_called_once()
            event.stop.assert_called_once()
        
        # Test arrow keys
        for arrow_key in ["up", "down", "left", "right"]:
            event.key = arrow_key
            event.stop.reset_mock()
            with patch.object(dashboard, f'action_move_{arrow_key.replace("left", "left").replace("right", "right")}') as mock_move:
                dashboard.on_key(event)
                mock_move.assert_called_once()
                event.stop.assert_called_once()

