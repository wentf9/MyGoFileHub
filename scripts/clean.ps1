# Clean build artifacts and temporary files
# Usage: .\scripts\clean.ps1

$ErrorActionPreference = "Stop"

# Get script and root directory
$ScriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path
$RootDir = Split-Path -Parent $ScriptDir
$FrontendDir = Join-Path $RootDir "frontend" "my-go-file-hub-ui"

Write-Host "Cleaning build artifacts..."

# Clean frontend build
$FrontendDist = Join-Path $FrontendDir "dist"
if (Test-Path $FrontendDist) {
    Write-Host "  Removing frontend/dist..."
    Remove-Item -Recurse -Force $FrontendDist
}

# Clean frontend node_modules (optional, uncomment if needed)
# $NodeModules = Join-Path $FrontendDir "node_modules"
# if (Test-Path $NodeModules) {
#     Write-Host "  Removing frontend/node_modules..."
#     Remove-Item -Recurse -Force $NodeModules
# }

# Clean backend binary
$BinDir = Join-Path $RootDir "bin"
if (Test-Path $BinDir) {
    Write-Host "  Removing bin/..."
    Remove-Item -Recurse -Force $BinDir
}

# Clean Go cache
Write-Host "  Cleaning Go cache..."
go clean -cache

# Clean logs (optional, uncomment if needed)
# $LogDir = Join-Path $RootDir "logs"
# if (Test-Path $LogDir) {
#     Write-Host "  Removing logs/..."
#     Remove-Item -Recurse -Force $LogDir
# }

Write-Host "Clean completed!"
