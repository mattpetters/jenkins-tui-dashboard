#!/bin/bash
# Run tests for Jenkins TUI Dashboard
cd "$(dirname "$0")"

# Activate virtual environment if it exists
if [ -d "venv" ]; then
    source venv/bin/activate
fi

# Run pytest with coverage
python -m pytest tests/ -v --cov=src --cov-report=term-missing "$@"

