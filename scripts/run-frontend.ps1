# Run frontend development server
# Usage: .\scripts\run-frontend.ps1

$ErrorActionPreference = "Stop"

# Get script and root directory
$ScriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path
$RootDir = Split-Path -Parent $ScriptDir
$FrontendDir = Join-Path $RootDir "frontend" "my-go-file-hub-ui"

# Change to frontend directory
Set-Location $FrontendDir

# Check if node_modules exists
if (!(Test-Path "node_modules")) {
    Write-Host "Installing dependencies..."
    npm install
}

Write-Host "Starting frontend dev server..."
Write-Host "Working directory: $FrontendDir"

npm run dev
