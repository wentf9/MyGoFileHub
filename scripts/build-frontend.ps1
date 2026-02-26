# Build frontend production bundle
# Usage: .\scripts\build-frontend.ps1

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

Write-Host "Building frontend..."
Write-Host "Working directory: $FrontendDir"

npm run build

# Check output
$DistDir = Join-Path $FrontendDir "dist"

if (Test-Path $DistDir) {
    Write-Host "Frontend built successfully!"
    Write-Host "Output directory: $DistDir"
} else {
    Write-Host "Error: Build failed, dist directory not found" -ForegroundColor Red
    exit 1
}
