#!/bin/bash
# Build entire project (frontend + backend)
# Usage: ./scripts/build.sh

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(dirname "$SCRIPT_DIR")"

cd "$ROOT_DIR"

echo "========================================"
echo "Building MyGoFileHub"
echo "========================================"

# Step 1: Build frontend
echo ""
echo "[Step 1/2] Building frontend..."
echo "----------------------------------------"
"$SCRIPT_DIR/build-frontend.sh"

# Step 2: Build backend
echo ""
echo "[Step 2/2] Building backend..."
echo "----------------------------------------"

# Ensure logs and data directories exist
mkdir -p "$ROOT_DIR/logs"
mkdir -p "$ROOT_DIR/data"

# Get version info
echo "Getting version info..."
if command -v bash &> /dev/null; then
    eval $("bash" "$SCRIPT_DIR/get_version.sh")
else
    echo "Warning: bash not found, using default version info"
    VERSION="v0.1.0-dev"
    GIT_COMMIT="unknown"
    BUILD_TIME=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
fi

echo "  Version: $VERSION"
echo "  Git Commit: $GIT_COMMIT"
echo "  Build Time: $BUILD_TIME"

# Build ldflags
LDFLAGS="-s -w"
LDFLAGS="$LDFLAGS -X 'github.com/wentf9/MyGoFileHub/internal/application.version=$VERSION'"
LDFLAGS="$LDFLAGS -X 'github.com/wentf9/MyGoFileHub/internal/application.gitCommit=$GIT_COMMIT'"
LDFLAGS="$LDFLAGS -X 'github.com/wentf9/MyGoFileHub/internal/application.buildTime=$BUILD_TIME'"

# Build backend binary
echo "Compiling Go binary with ldflags..."
CGO_ENABLED=0 go build -ldflags "$LDFLAGS" -o "$ROOT_DIR/bin/MyGoFileHub" main.go

echo ""
echo "========================================"
echo "Build completed successfully!"
echo "========================================"
echo "Binary: $ROOT_DIR/bin/MyGoFileHub"
echo "Frontend: $ROOT_DIR/frontend/my-go-file-hub-ui/dist"
