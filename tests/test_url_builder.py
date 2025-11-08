"""Unit tests for URL builder utilities."""
import os
import pytest
from unittest.mock import patch
from src.utils.url_builder import (
    parse_pr_number,
    build_jenkins_build_url,
    build_pr_url,
    infer_job_path_from_pr,
    get_default_repo_owner,
    get_default_repo_name
)


class TestParsePRNumber:
    """Test PR number parsing."""
    
    def test_parse_pr_with_prefix(self):
        """Test parsing PR number with PR- prefix."""
        assert parse_pr_number("PR-3859") == "3859"
        assert parse_pr_number("pr-3859") == "3859"
    
    def test_parse_pr_without_prefix(self):
        """Test parsing PR number without prefix."""
        assert parse_pr_number("3859") == "3859"
    
    def test_parse_pr_with_whitespace(self):
        """Test parsing PR number with whitespace."""
        assert parse_pr_number("  PR-3859  ") == "3859"
        assert parse_pr_number("  3859  ") == "3859"


class TestBuildJenkinsURL:
    """Test Jenkins URL building."""
    
    def test_build_jenkins_url_basic(self):
        """Test building basic Jenkins URL."""
        url = build_jenkins_build_url(
            "identity-manage/account/account-eks",
            "PR-3859",
            263
        )
        assert "identity-manage" in url
        assert "PR-3859" in url
        assert "263" in url
        assert "build.intuit.com" in url
    
    def test_build_jenkins_url_with_stage(self):
        """Test building Jenkins URL with stage ID."""
        url = build_jenkins_build_url(
            "identity-manage/account/account-eks",
            "PR-3859",
            263,
            stage_id=295
        )
        assert "263" in url
        assert "295" in url
        assert "/pipeline/" in url


class TestBuildPRURL:
    """Test PR URL building."""
    
    @patch.dict(os.environ, {
        "DEFAULT_REPO_OWNER": "identity-manage",
        "DEFAULT_REPO_NAME": "account"
    })
    def test_build_pr_url_defaults(self):
        """Test building PR URL with default repo."""
        url = build_pr_url("3859")
        assert "github.intuit.com" in url
        assert "identity-manage" in url
        assert "account" in url
        assert "3859" in url
    
    def test_build_pr_url_custom_repo(self):
        """Test building PR URL with custom repo."""
        url = build_pr_url("3859", repo_owner="custom-owner", repo_name="custom-repo")
        assert "custom-owner" in url
        assert "custom-repo" in url
        assert "3859" in url


class TestInferJobPath:
    """Test job path inference."""
    
    @patch.dict(os.environ, {
        "DEFAULT_REPO_OWNER": "identity-manage",
        "DEFAULT_REPO_NAME": "account"
    })
    def test_infer_job_path(self):
        """Test inferring job path from PR number."""
        job_path = infer_job_path_from_pr("3859")
        assert "identity-manage" in job_path
        assert "account" in job_path

