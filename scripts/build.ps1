# Build entire project (frontend + backend)
# Usage: .\scripts\build.ps1

$ErrorActionPreference = "Stop"

# Get script and root directory
$ScriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path
$RootDir = Split-Path -Parent $ScriptDir

Write-Host "========================================" -ForegroundColor Cyan
Write-Host "Building MyGoFileHub" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan

# Step 1: Build frontend
Write-Host ""
Write-Host "[Step 1/2] Building frontend..." -ForegroundColor Green
Write-Host "----------------------------------------" -ForegroundColor Gray
& "$ScriptDir\build-frontend.ps1"

# Step 2: Build backend
Write-Host ""
Write-Host "[Step 2/2] Building backend..." -ForegroundColor Green
Write-Host "----------------------------------------" -ForegroundColor Gray

# Ensure directories exist
$logDir = Join-Path $RootDir "logs"
$dataDir = Join-Path $RootDir "data"
$binDir = Join-Path $RootDir "bin"

if (!(Test-Path $logDir)) {
    New-Item -ItemType Directory -Path $logDir | Out-Null
}
if (!(Test-Path $dataDir)) {
    New-Item -ItemType Directory -Path $dataDir | Out-Null
}
if (!(Test-Path $binDir)) {
    New-Item -ItemType Directory -Path $binDir | Out-Null
}

# Build backend binary
Write-Host ""
Write-Host "Getting version info..." -ForegroundColor Gray

# Get version info
try {
    $VersionOutput = & "$ScriptDir\get_version.ps1"
    $VersionOutput | ForEach-Object {
        if ($_ -match "^VERSION=(.*)$") { $VERSION = $matches[1] }
        elseif ($_ -match "^GIT_COMMIT=(.*)$") { $GIT_COMMIT = $matches[1] }
        elseif ($_ -match "^BUILD_TIME=(.*)$") { $BUILD_TIME = $matches[1] }
    }
} catch {
    Write-Host "Warning: Failed to get version info, using defaults" -ForegroundColor Yellow
    $VERSION = "v0.1.0-dev"
    $GIT_COMMIT = "unknown"
    $BUILD_TIME = (Get-Date -AsUTC -Format "yyyy-MM-ddTHH:mm:ssZ")
}

Write-Host "  Version: $VERSION"
Write-Host "  Git Commit: $GIT_COMMIT"
Write-Host "  Build Time: $BUILD_TIME"

# Build ldflags
$LDFLAGS = "-s -w"
$LDFLAGS += " -X 'github.com/wentf9/MyGoFileHub/internal/application.version=$VERSION'"
$LDFLAGS += " -X 'github.com/wentf9/MyGoFileHub/internal/application.gitCommit=$GIT_COMMIT'"
$LDFLAGS += " -X 'github.com/wentf9/MyGoFileHub/internal/application.buildTime=$BUILD_TIME'"

Write-Host "Compiling Go binary with ldflags..."
Set-Location $RootDir
$env:CGO_ENABLED = "0"
go build -ldflags "$LDFLAGS" -o "$binDir\MyGoFileHub.exe" main.go

Write-Host ""
Write-Host "========================================" -ForegroundColor Cyan
Write-Host "Build completed successfully!" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host "Binary: $binDir\MyGoFileHub.exe"
Write-Host "Frontend: $RootDir\frontend\my-go-file-hub-ui\dist"
