#!/bin/bash
# Run backend development server
# Usage: ./scripts/run-backend.sh [port]

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(dirname "$SCRIPT_DIR")"
cd "$ROOT_DIR"

# Default port
PORT=${1:-3939}

# Ensure logs directory exists
mkdir -p "$ROOT_DIR/logs"

# Ensure data directory exists
mkdir -p "$ROOT_DIR/data"

echo "Starting backend server on port $PORT..."
echo "Logs will be written to $ROOT_DIR/logs/"

# Set environment variables and run
MY_GO_FILE_HUB_SERVER_PORT="$PORT" \
MY_GO_FILE_HUB_DATA_DIR="$ROOT_DIR/data" \
MY_GO_FILE_HUB_LISTEN="localhost" \
go run main.go 2>&1 | tee "$ROOT_DIR/logs/backend.log"
