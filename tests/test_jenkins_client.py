"""Unit tests for Jenkins client."""
import pytest
from unittest.mock import Mock, patch, MagicMock
from datetime import timedelta
from src.jenkins_client import JenkinsClient
from src.models import Build, BuildStatus


class TestJenkinsClient:
    """Test JenkinsClient."""
    
    @patch.dict('os.environ', {
        'JENKINS_USER': 'test_user',
        'JENKINS_TOKEN': 'test_token',
        'JENKINS_URL': 'https://test.jenkins.com'
    })
    def test_client_initialization(self):
        """Test client initialization."""
        client = JenkinsClient()
        assert client.username == "test_user"
        assert client.token == "test_token"
        assert client.base_url == "https://test.jenkins.com"
    
    @patch('src.jenkins_client.requests.get')
    def test_get_build_status_success(self, mock_get):
        """Test getting build status successfully."""
        # Mock successful API response
        mock_response = Mock()
        mock_response.status_code = 200
        mock_response.json.return_value = {
            "number": 263,
            "state": "running",
            "result": None,
            "durationInMillis": 120000,
            "stages": [
                {"name": "Build", "state": "running"}
            ],
            "pipeline": "account-eks"
        }
        mock_response.raise_for_status = Mock()
        mock_get.return_value = mock_response
        
        client = JenkinsClient()
        build = client.get_build_status(
            "identity-manage/account/account-eks",
            "PR-3859",
            263
        )
        
        assert build is not None
        assert build.build_number == 263
        assert build.status == BuildStatus.RUNNING
    
    @patch('src.jenkins_client.requests.get')
    def test_get_build_status_failure(self, mock_get):
        """Test handling build status fetch failure."""
        # Mock failed API response
        mock_get.side_effect = Exception("Connection error")
        
        client = JenkinsClient()
        build = client.get_build_status(
            "identity-manage/account/account-eks",
            "PR-3859",
            263
        )
        
        assert build is not None
        assert build.status == BuildStatus.ERROR
        assert build.error_message is not None
    
    @patch('src.jenkins_client.requests.get')
    def test_create_build_from_pr(self, mock_get):
        """Test creating build from PR number."""
        # Mock API response
        mock_response = Mock()
        mock_response.status_code = 200
        mock_response.json.return_value = {
            "number": 263,
            "state": "success",
            "result": "success",
            "durationInMillis": 300000,
            "stages": [
                {"name": "Test", "state": "success"}
            ],
            "pipeline": "account-eks"
        }
        mock_response.raise_for_status = Mock()
        mock_get.return_value = mock_response
        
        client = JenkinsClient()
        build = client.create_build_from_pr("3859")
        
        assert build is not None
        assert build.pr_number == "3859"
        assert build.build_number == 263
        assert build.status == BuildStatus.SUCCESS
    
    def test_parse_build_data_running(self):
        """Test parsing build data for running build."""
        client = JenkinsClient()
        data = {
            "number": 263,
            "state": "running",
            "result": None,
            "durationInMillis": 120000,
            "stages": [
                {"name": "Build", "state": "running"}
            ],
            "pipeline": "account-eks"
        }
        
        build = client._parse_build_data(data, "identity-manage/account/account-eks", "PR-3859")
        assert build.status == BuildStatus.RUNNING
        assert build.stage == "Build"
        assert build.build_number == 263
    
    def test_parse_build_data_success(self):
        """Test parsing build data for successful build."""
        client = JenkinsClient()
        data = {
            "number": 263,
            "state": "finished",
            "result": "success",
            "durationInMillis": 300000,
            "stages": [
                {"name": "Test", "state": "success"}
            ],
            "pipeline": "account-eks"
        }
        
        build = client._parse_build_data(data, "identity-manage/account/account-eks", "PR-3859")
        assert build.status == BuildStatus.SUCCESS
        assert build.duration == timedelta(milliseconds=300000)
    
    def test_parse_build_data_failure(self):
        """Test parsing build data for failed build."""
        client = JenkinsClient()
        data = {
            "number": 263,
            "state": "finished",
            "result": "failure",
            "durationInMillis": 150000,
            "stages": [
                {"name": "Test", "state": "failure"}
            ],
            "pipeline": "account-eks"
        }
        
        build = client._parse_build_data(data, "identity-manage/account/account-eks", "PR-3859")
        assert build.status == BuildStatus.FAILURE

