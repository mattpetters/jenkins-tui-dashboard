"""Helper to add test builds for debugging."""
from ..models import Build, BuildStatus
from datetime import timedelta


def create_test_builds():
    """Create test builds for debugging the UI."""
    return [
        Build(
            pr_number="3859",
            build_number=263,
            status=BuildStatus.RUNNING,
            stage="Test",
            job_name="account-eks",
            duration=timedelta(minutes=5, seconds=30),
            build_url="https://build.intuit.com/identity/blue/organizations/jenkins/identity-manage%2Faccount%2Faccount-eks/detail/PR-3859/263/pipeline/295",
            pr_url="https://github.intuit.com/identity-manage/account/pull/3859",
            job_path="identity-manage/account/account-eks"
        ),
        Build(
            pr_number="3860",
            build_number=100,
            status=BuildStatus.SUCCESS,
            stage="Deploy",
            job_name="account-eks",
            duration=timedelta(minutes=10, seconds=15),
            build_url="https://build.intuit.com/identity/blue/organizations/jenkins/identity-manage%2Faccount%2Faccount-eks/detail/PR-3860/100",
            pr_url="https://github.intuit.com/identity-manage/account/pull/3860",
            job_path="identity-manage/account/account-eks"
        ),
        Build(
            pr_number="3861",
            build_number=200,
            status=BuildStatus.FAILURE,
            stage="Test",
            job_name="account-eks",
            duration=timedelta(minutes=2, seconds=45),
            build_url="https://build.intuit.com/identity/blue/organizations/jenkins/identity-manage%2Faccount%2Faccount-eks/detail/PR-3861/200",
            pr_url="https://github.intuit.com/identity-manage/account/pull/3861",
            job_path="identity-manage/account/account-eks"
        ),
    ]

