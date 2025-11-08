"""GitHub helper to get PR branch names and details."""
import os
from typing import Optional, Dict, Any


def get_pr_branch_from_github(pr_number: str, repo_owner: str = "identity-manage", repo_name: str = "account") -> Optional[str]:
    """Get the actual branch name for a PR from GitHub.
    
    Uses GitHub API to find the head ref (branch name) for a PR.
    This is needed because PR branches might not follow the PR-XXXX pattern.
    
    Args:
        pr_number: PR number (e.g., '3859')
        repo_owner: Repository owner
        repo_name: Repository name
    
    Returns:
        Branch name or None if not found
    """
    try:
        import requests
        from requests.auth import HTTPBasicAuth
        
        # Get GitHub token from environment
        github_token = os.getenv("GITHUB_TOKEN")
        if not github_token:
            return None
        
        # GitHub API endpoint
        api_url = f"https://github.intuit.com/api/v3/repos/{repo_owner}/{repo_name}/pulls/{pr_number}"
        
        headers = {
            "Authorization": f"token {github_token}",
            "Accept": "application/vnd.github.v3+json"
        }
        
        response = requests.get(api_url, headers=headers, timeout=10)
        if response.status_code == 200:
            data = response.json()
            # Get the head ref (branch name)
            return data.get("head", {}).get("ref")
        
        return None
    except Exception as e:
        print(f"Error getting PR branch from GitHub: {e}")
        return None


def get_pr_details(pr_number: str, repo_owner: str = "identity-manage", repo_name: str = "account") -> Optional[Dict[str, Any]]:
    """Get PR details from GitHub including branch name, status, etc.
    
    Args:
        pr_number: PR number
        repo_owner: Repository owner
        repo_name: Repository name
    
    Returns:
        Dictionary with PR details or None
    """
    try:
        import requests
        
        github_token = os.getenv("GITHUB_TOKEN")
        if not github_token:
            return None
        
        api_url = f"https://github.intuit.com/api/v3/repos/{repo_owner}/{repo_name}/pulls/{pr_number}"
        
        headers = {
            "Authorization": f"token {github_token}",
            "Accept": "application/vnd.github.v3+json"
        }
        
        response = requests.get(api_url, headers=headers, timeout=10)
        if response.status_code == 200:
            data = response.json()
            return {
                "branch": data.get("head", {}).get("ref"),
                "sha": data.get("head", {}).get("sha"),
                "title": data.get("title"),
                "state": data.get("state"),
                "draft": data.get("draft"),
                "html_url": data.get("html_url"),
            }
        
        return None
    except Exception:
        return None

