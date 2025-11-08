#!/bin/bash
# Clean run script - ensures absolutely latest code
cd "$(dirname "$0")"

echo "🧹 Deep cleaning..."
# Remove all Python cache
find . -type d -name "__pycache__" -exec rm -rf {} + 2>/dev/null
find . -name "*.pyc" -delete 2>/dev/null
find . -name "*.pyo" -delete 2>/dev/null
rm -rf .pytest_cache 2>/dev/null

# Sync files from worktree if they exist
if [ -d "/Users/mpetters/.cursor/worktrees/jenkins-dash/CuzHm/src" ]; then
    echo "📋 Syncing latest files from worktree..."
    cp -r /Users/mpetters/.cursor/worktrees/jenkins-dash/CuzHm/src/* src/ 2>/dev/null
fi

# Activate venv
if [ -d "venv" ]; then
    source venv/bin/activate
else
    echo "❌ No venv found"
    exit 1
fi

echo "✅ Running with latest code..."
python -m src.main "$@"

