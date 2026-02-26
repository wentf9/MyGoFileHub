# 获取版本信息用于 -ldflags 注入
# 用法：& .\scripts\get_version.ps1

# 获取 git commit hash (短版本)
try {
    $GIT_COMMIT = (git rev-parse --short HEAD 2>$null)
    if ($null -eq $GIT_COMMIT) {
        $GIT_COMMIT = "unknown"
    }
} catch {
    $GIT_COMMIT = "unknown"
}

# 获取构建时间 (RFC3339 格式，UTC 时间)
$BUILD_TIME = (Get-Date -AsUTC -Format "yyyy-MM-ddTHH:mm:ssZ")

# 获取当前 tag 作为版本号 (如果没有 tag 则使用 dev 版本)
try {
    $VERSION = (git describe --tags --always --dirty 2>$null)
    if ($null -eq $VERSION) {
        $VERSION = "v0.1.0-dev"
    }
} catch {
    $VERSION = "v0.1.0-dev"
}

# 输出为 PowerShell 变量格式
Write-Output "VERSION=$VERSION"
Write-Output "GIT_COMMIT=$GIT_COMMIT"
Write-Output "BUILD_TIME=$BUILD_TIME"
