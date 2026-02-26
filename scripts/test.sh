#!/bin/bash
# Run tests
# Usage: ./scripts/test.sh [pattern]

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(dirname "$SCRIPT_DIR")"

cd "$ROOT_DIR"

echo "Running Go tests..."

# Run all tests or specific pattern
if [ -n "$1" ]; then
    go test -v ./... -run "$1"
else
    go test -v ./...
fi

echo ""
echo "Running frontend tests (if any)..."
cd "$ROOT_DIR/frontend/my-go-file-hub-ui"

# Check if there's a test script in package.json
if npm run | grep -q " test "; then
    npm test
else
    echo "No test script found in frontend package.json"
fi
