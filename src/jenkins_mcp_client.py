"""Jenkins MCP client integration for fetching build status."""
import os
from typing import Optional, Dict, Any
from datetime import datetime, timedelta

from .models import Build, BuildStatus
from .utils.url_builder import build_jenkins_build_url, build_pr_url, parse_pr_number


class JenkinsMCPClient:
    """Client that uses Jenkins MCP server for fetching build data."""
    
    def __init__(self):
        """Initialize Jenkins MCP client."""
        self.base_url = os.getenv("JENKINS_URL", "https://build.intuit.com")
    
    def get_pr_build_using_mcp(
        self,
        pr_number: str,
        repo_owner: str = "identity-manage",
        repo_name: str = "account"
    ) -> Optional[Build]:
        """Get build for a PR using Jenkins MCP and GitHub MCP.
        
        Args:
            pr_number: PR number
            repo_owner: Repository owner
            repo_name: Repository name
        
        Returns:
            Build object or None
        """
        try:
            pr_num = parse_pr_number(pr_number)
            
            # Construct Jenkins job path
            # Format: identity/job/{owner}/job/{repo}/job/{repo}-eks/job/PR-{number}
            job_path = f"identity/job/{repo_owner}/job/{repo_name}/job/{repo_name}-eks/job/PR-{pr_num}"
            
            # This would call the Jenkins MCP server
            # For now, return None to use the HTTP client
            # In a future enhancement, we'd integrate with the actual MCP server here
            
            return None
        except Exception as e:
            return None
    
    @staticmethod
    def parse_job_info_to_build(job_info: Dict[str, Any], pr_number: str) -> Build:
        """Parse Jenkins job info from MCP into Build object.
        
        Args:
            job_info: Job info dict from Jenkins MCP
            pr_number: PR number
        
        Returns:
            Build object
        """
        pr_num = parse_pr_number(pr_number)
        
        # Extract last build info
        last_build = job_info.get("lastBuild", {})
        build_number = last_build.get("number")
        build_url = last_build.get("url", "")
        
        # Get pipeline info
        pipeline = job_info.get("pipeline", {})
        stages = pipeline.get("stages", [])
        
        # Find the last completed or running stage
        current_stage = None
        status = BuildStatus.PENDING
        duration_ms = 0
        
        for stage in stages:
            stage_status = stage.get("status", "")
            if stage_status == "IN_PROGRESS":
                current_stage = stage.get("name")
                status = BuildStatus.RUNNING
                duration_ms += stage.get("durationMillis", 0)
            elif stage_status == "SUCCESS":
                current_stage = stage.get("name")
                duration_ms += stage.get("durationMillis", 0)
        
        # Determine overall status from health report or last build
        if "lastSuccessfulBuild" in job_info and job_info.get("lastSuccessfulBuild", {}).get("number") == build_number:
            status = BuildStatus.SUCCESS
        elif "lastFailedBuild" in job_info and job_info.get("lastFailedBuild", {}).get("number") == build_number:
            status = BuildStatus.FAILURE
        
        duration = timedelta(milliseconds=duration_ms) if duration_ms > 0 else None
        
        return Build(
            pr_number=pr_num,
            build_number=build_number,
            status=status,
            stage=current_stage,
            job_name=job_info.get("name", f"PR-{pr_num}"),
            duration=duration,
            build_url=build_url,
            pr_url=build_pr_url(pr_num),
            job_path=job_info.get("url", ""),
            last_updated=datetime.now()
        )

