#!/bin/bash
# Clean build artifacts and temporary files
# Usage: ./scripts/clean.sh

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(dirname "$SCRIPT_DIR")"
FRONTEND_DIR="$ROOT_DIR/frontend/my-go-file-hub-ui"

echo "Cleaning build artifacts..."

# Clean frontend build
if [ -d "$FRONTEND_DIR/dist" ]; then
    echo "  Removing frontend/dist..."
    rm -rf "$FRONTEND_DIR/dist"
fi

# Clean frontend node_modules (optional, uncomment if needed)
# if [ -d "$FRONTEND_DIR/node_modules" ]; then
#     echo "  Removing frontend/node_modules..."
#     rm -rf "$FRONTEND_DIR/node_modules"
# fi

# Clean backend binary
if [ -d "$ROOT_DIR/bin" ]; then
    echo "  Removing bin/..."
    rm -rf "$ROOT_DIR/bin"
fi

# Clean Go cache
echo "  Cleaning Go cache..."
go clean -cache

# Clean logs (optional, uncomment if needed)
# if [ -d "$ROOT_DIR/logs" ]; then
#     echo "  Removing logs/..."
#     rm -rf "$ROOT_DIR/logs"
# fi

echo "Clean completed!"
