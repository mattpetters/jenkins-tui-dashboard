#!/bin/bash
# Quick build and run script for Jenkins Dashboard

set -e

echo "🔨 Building Jenkins Dashboard..."
go build -o jenkins-dash ./cmd/jenkins-dash

echo "✅ Build successful!"
echo "🚀 Starting dashboard..."
echo ""

./jenkins-dash

