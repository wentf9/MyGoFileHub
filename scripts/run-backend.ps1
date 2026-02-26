# Run backend development server
# Usage: .\scripts\run-backend.ps1 [port]

param(
    [string]$Port = "3939"
)

$ErrorActionPreference = "Stop"

# Get script and root directory
$ScriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path
$RootDir = Split-Path -Parent $ScriptDir

# Change to root directory
Set-Location $RootDir

# Ensure directories exist
$logDir = Join-Path $RootDir "logs"
$dataDir = Join-Path $RootDir "data"

if (!(Test-Path $logDir)) {
    New-Item -ItemType Directory -Path $logDir | Out-Null
}
if (!(Test-Path $dataDir)) {
    New-Item -ItemType Directory -Path $dataDir | Out-Null
}

Write-Host "Starting backend server on port $Port..."
Write-Host "Logs will be written to $logDir"

# Set environment variables and run
$env:MY_GO_FILE_HUB_SERVER_PORT = $Port
$env:MY_GO_FILE_HUB_DATA_DIR = $dataDir
$env:MY_GO_FILE_HUB_LISTEN = "localhost"

# Run and capture output
go run main.go 2>&1 | Tee-Object -FilePath "$logDir\backend.log"
