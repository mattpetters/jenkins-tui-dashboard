"""Jenkins API client for fetching build status."""
import os
import requests
from typing import Optional, Dict, Any
from datetime import datetime, timedelta
from requests.auth import HTTPBasicAuth

from .models import Build, BuildStatus
from .utils.url_builder import build_jenkins_build_url, build_pr_url, infer_job_path_from_pr, parse_pr_number


class JenkinsClient:
    """Client for interacting with Jenkins API."""
    
    def __init__(self):
        """Initialize Jenkins client with credentials from environment."""
        self.base_url = os.getenv("JENKINS_URL", "https://build.intuit.com")
        self.username = os.getenv("JENKINS_USER", "mpetters")
        self.token = os.getenv("JENKINS_TOKEN", "")
        self.auth = HTTPBasicAuth(self.username, self.token) if self.token else None
    
    def get_build_status(
        self,
        job_path: str,
        pr_branch: str,
        build_number: Optional[int] = None
    ) -> Optional[Build]:
        """Get build status for a PR branch.
        
        Args:
            job_path: Jenkins job path (e.g., 'identity-manage/account/account-eks')
            pr_branch: PR branch name (e.g., 'PR-3859')
            build_number: Optional build number (if None, gets latest)
        
        Returns:
            Build object with status information, or None if error
        """
        try:
            # Try using the Jenkins MCP server first
            # If that's not available, fall back to direct API calls
            return self._fetch_build_status_direct(job_path, pr_branch, build_number)
        except Exception as e:
            # Return a build with error status
            pr_number = parse_pr_number(pr_branch.replace("PR-", ""))
            build = Build(
                pr_number=pr_number,
                status=BuildStatus.ERROR,
                error_message=str(e),
                job_path=job_path,
                build_url=build_jenkins_build_url(job_path, pr_branch, build_number or 0) if build_number else None,
                pr_url=build_pr_url(pr_number)
            )
            return build
    
    def _fetch_build_status_direct(
        self,
        job_path: str,
        pr_branch: str,
        build_number: Optional[int] = None
    ) -> Optional[Build]:
        """Fetch build status using direct Jenkins API calls."""
        # Construct API URL
        # Jenkins Blue Ocean API endpoint
        encoded_job_path = job_path.replace("/", "%2F")
        api_url = f"{self.base_url}/blue/rest/organizations/jenkins/pipelines/{encoded_job_path}/branches/{pr_branch}/runs"
        
        if build_number:
            api_url += f"/{build_number}"
        else:
            # Get latest build
            api_url += "?latest=true"
        
        try:
            response = requests.get(api_url, auth=self.auth, timeout=10)
            response.raise_for_status()
            data = response.json()
            
            return self._parse_build_data(data, job_path, pr_branch)
        except requests.exceptions.RequestException as e:
            # Try alternative API endpoint
            return self._try_alternative_api(job_path, pr_branch, build_number)
    
    def _try_alternative_api(
        self,
        job_path: str,
        pr_branch: str,
        build_number: Optional[int] = None
    ) -> Optional[Build]:
        """Try alternative Jenkins API endpoint."""
        pr_number = parse_pr_number(pr_branch.replace("PR-", ""))
        
        # Try to get build info from the job API
        encoded_job_path = job_path.replace("/", "/job/")
        api_url = f"{self.base_url}/job/{encoded_job_path}/api/json"
        
        try:
            response = requests.get(api_url, auth=self.auth, timeout=10)
            if response.status_code == 200:
                job_data = response.json()
                # Try to find the build for this PR branch
                builds = job_data.get("builds", [])
                for build_info in builds:
                    build_num = build_info.get("number")
                    if build_number and build_num != build_number:
                        continue
                    
                    # Get detailed build info
                    build_detail_url = build_info.get("url", "")
                    if build_detail_url:
                        return self._fetch_build_details(build_detail_url, job_path, pr_branch, pr_number)
        except Exception:
            pass
        
        # Return a build with pending status if we can't fetch
        return Build(
            pr_number=pr_number,
            status=BuildStatus.PENDING,
            job_path=job_path,
            build_url=build_jenkins_build_url(job_path, pr_branch, build_number or 0) if build_number else None,
            pr_url=build_pr_url(pr_number)
        )
    
    def _fetch_build_details(
        self,
        build_url: str,
        job_path: str,
        pr_branch: str,
        pr_number: str
    ) -> Build:
        """Fetch detailed build information."""
        try:
            api_url = f"{build_url}api/json"
            response = requests.get(api_url, auth=self.auth, timeout=10)
            response.raise_for_status()
            data = response.json()
            
            return self._parse_build_data(data, job_path, pr_branch)
        except Exception as e:
            return Build(
                pr_number=pr_number,
                status=BuildStatus.ERROR,
                error_message=str(e),
                job_path=job_path,
                build_url=build_url,
                pr_url=build_pr_url(pr_number)
            )
    
    def _parse_build_data(
        self,
        data: Dict[str, Any],
        job_path: str,
        pr_branch: str
    ) -> Build:
        """Parse build data from Jenkins API response."""
        pr_number = parse_pr_number(pr_branch.replace("PR-", ""))
        
        # Extract build information
        build_number = data.get("number") or data.get("id")
        state_val = data.get("state")
        result_val = data.get("result")
        # Handle None values safely
        state = str(state_val).lower() if state_val is not None else ""
        result = str(result_val).lower() if result_val is not None else ""
        
        # Determine status
        if state == "running" or state == "in_progress":
            status = BuildStatus.RUNNING
        elif result == "success":
            status = BuildStatus.SUCCESS
        elif result == "failure":
            status = BuildStatus.FAILURE
        elif result == "unstable":
            status = BuildStatus.UNSTABLE
        elif result == "aborted":
            status = BuildStatus.ABORTED
        else:
            status = BuildStatus.PENDING
        
        # Extract duration
        duration_ms = data.get("durationInMillis", 0)
        duration = timedelta(milliseconds=duration_ms) if duration_ms else None
        
        # Extract stage and job information
        stage = None
        job_name = None
        
        # Try to get current stage from pipeline steps
        stages = data.get("stages", [])
        if stages:
            # Find the currently running or last stage
            for stage_info in stages:
                stage_state_val = stage_info.get("state") or ""
                stage_state = str(stage_state_val).lower() if stage_state_val is not None else ""
                if stage_state == "running" or stage_state == "in_progress":
                    stage = stage_info.get("name", "Unknown")
                    break
            
            # If no running stage, get the last one
            if not stage and stages:
                last_stage = stages[-1]
                stage = last_stage.get("name", "Unknown")
        
        # Get job name from pipeline or job path
        job_name = data.get("pipeline", "") or job_path.split("/")[-1]
        
        # Build URLs
        build_url = build_jenkins_build_url(job_path, pr_branch, build_number)
        pr_url = build_pr_url(pr_number)
        
        return Build(
            pr_number=pr_number,
            build_number=build_number,
            status=status,
            stage=stage,
            job_name=job_name,
            duration=duration,
            build_url=build_url,
            pr_url=pr_url,
            job_path=job_path,
            last_updated=datetime.now()
        )
    
    def create_build_from_pr(self, pr_number: str, job_path: Optional[str] = None) -> Build:
        """Create a Build object from a PR number.
        
        Args:
            pr_number: PR number (e.g., '3859')
            job_path: Optional job path (if None, will infer)
        
        Returns:
            Build object
        """
        pr_number = parse_pr_number(pr_number)
        pr_branch = f"PR-{pr_number}"
        
        if not job_path:
            job_path = infer_job_path_from_pr(pr_number)
        
        # Fetch build status
        build = self.get_build_status(job_path, pr_branch)
        
        if not build:
            # Return a pending build if we can't fetch
            build = Build(
                pr_number=pr_number,
                status=BuildStatus.PENDING,
                job_path=job_path,
                build_url=None,
                pr_url=build_pr_url(pr_number)
            )
        
        return build

