import sys
sys.path.insert(0, '.')
import os
os.environ['JENKINS_USER'] = 'mpetters'
os.environ['JENKINS_TOKEN'] = '113e76ce0e6775a1dca5a29d857e025888'

from src.models import Build, BuildStatus
from src.components.build_tile import BuildTile

# Create tiles
builds = [
    Build(pr_number="3859", status=BuildStatus.PENDING, stage="Loading...", job_name="Fetching data..."),
    Build(pr_number="3860", status=BuildStatus.SUCCESS, stage="Deploy", job_name="account-eks", build_number=100),
    Build(pr_number="3861", status=BuildStatus.RUNNING, stage="Test", job_name="account-eks", build_number=200),
]

print("=== TILES WOULD DISPLAY AS ===\n")
for i, build in enumerate(builds):
    tile = BuildTile(build, is_selected=(i==0))
    print(f"Tile {i+1}:")
    print(tile._render_content())
    print()

