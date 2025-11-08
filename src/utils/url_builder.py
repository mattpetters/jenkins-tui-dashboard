"""URL building utilities for Jenkins and GitHub."""
import os
from typing import Optional
from urllib.parse import quote


def get_jenkins_base_url() -> str:
    """Get Jenkins base URL from environment."""
    return os.getenv("JENKINS_URL", "https://build.intuit.com")


def get_github_base_url() -> str:
    """Get GitHub base URL from environment."""
    return os.getenv("GITHUB_BASE_URL", "https://github.intuit.com")


def get_default_repo_owner() -> str:
    """Get default repository owner from environment."""
    return os.getenv("DEFAULT_REPO_OWNER", "identity-manage")


def get_default_repo_name() -> str:
    """Get default repository name from environment."""
    return os.getenv("DEFAULT_REPO_NAME", "account")


def parse_pr_number(pr_input: str) -> str:
    """Parse PR number from input (handles PR-3859 or 3859)."""
    pr_input = pr_input.strip()
    if pr_input.upper().startswith("PR-"):
        return pr_input[3:]
    return pr_input


def build_jenkins_build_url(
    job_path: str,
    pr_branch: str,
    build_number: int,
    stage_id: Optional[int] = None
) -> str:
    """Build Jenkins build URL.
    
    Args:
        job_path: Jenkins job path (e.g., 'identity-manage/account/account-eks')
        pr_branch: PR branch name (e.g., 'PR-3859')
        build_number: Build number
        stage_id: Optional stage ID for pipeline view
    
    Returns:
        Full Jenkins build URL
    """
    base_url = get_jenkins_base_url()
    # URL encode the job path
    encoded_job_path = quote(job_path, safe="")
    
    url = f"{base_url}/identity/blue/organizations/jenkins/{encoded_job_path}/detail/{pr_branch}/{build_number}"
    
    if stage_id:
        url += f"/pipeline/{stage_id}"
    
    return url


def build_pr_url(pr_number: str, repo_owner: Optional[str] = None, repo_name: Optional[str] = None) -> str:
    """Build GitHub PR URL.
    
    Args:
        pr_number: PR number (e.g., '3859')
        repo_owner: Repository owner (defaults to env var)
        repo_name: Repository name (defaults to env var)
    
    Returns:
        Full GitHub PR URL
    """
    base_url = get_github_base_url()
    owner = repo_owner or get_default_repo_owner()
    repo = repo_name or get_default_repo_name()
    
    return f"{base_url}/{owner}/{repo}/pull/{pr_number}"


def infer_job_path_from_pr(pr_number: str) -> str:
    """Infer Jenkins job path from PR number.
    
    For now, uses default repo structure. Can be extended later.
    
    Args:
        pr_number: PR number
    
    Returns:
        Inferred job path
    """
    # Default structure based on example: identity-manage/account/account-eks
    owner = get_default_repo_owner()
    repo = get_default_repo_name()
    # Assume job name is repo-eks for now
    return f"{owner}/{repo}/{repo}-eks"

