# Run tests
# Usage: .\scripts\test.ps1 [pattern]

param(
    [string]$Pattern = ""
)

$ErrorActionPreference = "Stop"

# Get script and root directory
$ScriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path
$RootDir = Split-Path -Parent $ScriptDir

Write-Host "Running Go tests..."

Set-Location $RootDir

# Run all tests or specific pattern
if ($Pattern) {
    go test -v ./... -run $Pattern
} else {
    go test -v ./...
}

Write-Host ""
Write-Host "Running frontend tests (if any)..."

$FrontendDir = Join-Path $RootDir -ChildPath "frontend\my-go-file-hub-ui"
Set-Location $FrontendDir

# Check if there's a test script in package.json
$PackageJson = Get-Content "package.json" | ConvertFrom-Json
if ($PackageJson.scripts.PSObject.Properties.Name -contains "test") {
    npm test
} else {
    Write-Host "No test script found in frontend package.json"
}
