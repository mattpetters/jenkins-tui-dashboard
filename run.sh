#!/bin/bash
# Run the Jenkins TUI Dashboard
cd "$(dirname "$0")"

# Clean Python cache to ensure latest code runs
echo "🧹 Cleaning Python cache..."
find . -type d -name "__pycache__" -exec rm -rf {} + 2>/dev/null
find . -name "*.pyc" -delete 2>/dev/null
rm -rf .pytest_cache 2>/dev/null

# Activate virtual environment if it exists
if [ -d "venv" ]; then
    source venv/bin/activate
else
    echo "❌ Virtual environment not found. Run: python3 -m venv venv && source venv/bin/activate && pip install -r requirements.txt"
    exit 1
fi

echo "✅ Starting Jenkins Build Dashboard..."
python -m src.main "$@"

