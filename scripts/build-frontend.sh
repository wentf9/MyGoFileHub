#!/bin/bash
# Build frontend production bundle
# Usage: ./scripts/build-frontend.sh

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(dirname "$SCRIPT_DIR")"
FRONTEND_DIR="$ROOT_DIR/frontend/my-go-file-hub-ui"

cd "$FRONTEND_DIR"

# Check if node_modules exists
if [ ! -d "node_modules" ]; then
    echo "Installing dependencies..."
    npm install
fi

echo "Building frontend..."
echo "Working directory: $FRONTEND_DIR"

npm run build

# Copy dist to backend embed directory if needed
DIST_DIR="$FRONTEND_DIR/dist"
EMBED_DIR="$ROOT_DIR/frontend/dist"

if [ -d "$DIST_DIR" ]; then
    echo "Frontend built successfully!"
    echo "Output directory: $DIST_DIR"
else
    echo "Error: Build failed, dist directory not found"
    exit 1
fi
