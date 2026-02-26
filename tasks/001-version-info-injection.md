
# 任务 001: 添加版本号信息注入 (Version Info Injection)

## 进度摘要

**状态**: 已完成 ✅

- [x] 步骤 1: 创建 version 包定义版本信息结构
- [x] 步骤 2: 创建 build version 脚本支持 ldflags 注入
- [x] 步骤 3: 添加 version API 端点暴露版本信息
- [x] 步骤 4: 更新构建脚本集成 ldflags
- [x] 步骤 5: 添加 --version CLI 参数支持
- [x] 步骤 6: 验证构建和版本输出

### 验证结果

```bash
# CLI 版本输出
$ ./bin/MyGoFileHub --version
MyGoFileHub a9a2d96-dirty
  Git Commit: a9a2d96
  Build Time: 2026-02-26T06:18:52Z
  Go Version: go1.25.5
  Platform:   windows/amd64
```

**构建脚本输出**：
- 前端构建成功
- 后端成功注入版本信息
- CGO_ENABLED=0 解决 Windows cgo 编译问题

**版本发布**：
- Git Commit: `06a645b`
- Git Tag: `v0.1.0`

## 概述

为 MyGoFileHub 项目添加完整的版本信息注入机制，在编译时通过 Go 链接器的 `-ldflags` 参数将版本号、Git commit hash、构建时间等信息注入到二进制文件中。这是 Go 社区的最佳实践，可以让发布的二进制文件包含完整的构建元数据。

## 当前状态分析

当前项目的 `config/config.go` 已经通过环境变量实现了配置管理，并使用了 `github.com/caarlos0/env/v11` 库。项目中没有专门的 version 包或相关文件。

**相关文件:**
- `config/config.go` - 当前配置管理
- `main.go` - 程序入口
- `scripts/build.sh` / `scripts/build.ps1` - 构建脚本（需要更新）
- `scripts/README.md` - 脚本使用说明（需要更新）

## 目标状态

1. **version 包**: 创建 `version/version.go` 定义版本信息结构和打印函数
2. **注入变量**: 使用标准 Go 模式定义可注入的变量
   ```go
   var (
       Version   = "v0.1.0-dev"
       GitCommit = "unknown"
       BuildTime = "unknown"
       GoVersion = "unknown"
   )
   ```
3. **构建脚本**: 更新 `scripts/build.sh` 和 `scripts/build.ps1` 支持自动获取 Git 信息并通过 `-ldflags` 注入
4. **版本查询 API**: 添加 `GET /@api/v1/version` 端点
5. **CLI 参数**: 支持 `--version` 或 `-v` 参数打印版本信息

## 实现步骤

### 步骤 1: 创建 version 包

创建 `internal/application/version.go`（应用层）或 `version/version.go`（独立包）：

**涉及文件:**
- `internal/application/version.go` (新建)

```go
package application

import (
    "fmt"
    "runtime"
    "runtime/debug"
)

// 这些变量将在编译时通过 -ldflags 注入
var (
    version   = "v0.1.0-dev"
    gitCommit = "unknown"
    buildTime = "unknown"
)

// BuildInfo 包含完整的构建信息
type BuildInfo struct {
    Version   string `json:"version"`
    GitCommit string `json:"git_commit"`
    BuildTime string `json:"build_time"`
    GoVersion string `json:"go_version"`
    Platform  string `json:"platform"`
}

// GetBuildInfo 返回当前构建信息
func GetBuildInfo() BuildInfo {
    goVer := version
    if goVer == "(devel)" {
        if info, ok := debug.ReadBuildInfo(); ok {
            goVer = info.GoVersion
        }
    }

    return BuildInfo{
        Version:   version,
        GitCommit: gitCommit,
        BuildTime: buildTime,
        GoVersion: goVer,
        Platform:  fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
    }
}

// PrintVersion 打印版本信息到 stdout
func PrintVersion() {
    info := GetBuildInfo()
    fmt.Printf("MyGoFileHub %s\n", info.Version)
    fmt.Printf("  Git Commit: %s\n", info.GitCommit)
    fmt.Printf("  Build Time: %s\n", info.BuildTime)
    fmt.Printf("  Go Version: %s\n", info.GoVersion)
    fmt.Printf("  Platform:   %s\n", info.Platform)
}
```

### 步骤 2: 创建构建辅助脚本

创建 `scripts/get_version.sh` 和 `scripts/get_version.ps1` 用于获取版本信息：

**涉及文件:**
- `scripts/get_version.sh` (新建)
- `scripts/get_version.ps1` (新建)

```bash
#!/bin/bash
# 获取版本信息用于 -ldflags

# 获取 git commit hash (短版本)
GIT_COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")

# 获取构建时间 (RFC3339 格式)
BUILD_TIME=$(date -u +"%Y-%m-%dT%H:%M:%SZ")

# 获取当前 tag 作为版本号 (如果没有 tag 则使用 dev 版本)
VERSION=$(git describe --tags --always --dirty 2>/dev/null || echo "v0.1.0-dev")

echo "VERSION=$VERSION"
echo "GIT_COMMIT=$GIT_COMMIT"
echo "BUILD_TIME=$BUILD_TIME"
```

### 步骤 3: 更新构建脚本集成 ldflags

**涉及文件:**
- `scripts/build.sh` (修改)
- `scripts/build.ps1` (修改)

在构建脚本中添加：
```bash
# 获取版本信息
eval $(./scripts/get_version.sh)

# 构建 ldflags
LDFLAGS="-s -w"
LDFLAGS="$LDFLAGS -X 'github.com/wentf9/MyGoFileHub/internal/application.version=$VERSION'"
LDFLAGS="$LDFLAGS -X 'github.com/wentf9/MyGoFileHub/internal/application.gitCommit=$GIT_COMMIT'"
LDFLAGS="$LDFLAGS -X 'github.com/wentf9/MyGoFileHub/internal/application.buildTime=$BUILD_TIME'"

# 编译
go build -ldflags "$LDFLAGS" -o "$ROOT_DIR/bin/MyGoFileHub" main.go
```

### 步骤 4: 添加版本查询 API 端点

**涉及文件:**
- `internal/interface/api/handlers/version_handler.go` (新建)
- `internal/interface/api/route.go` (修改)

```go
// version_handler.go
package handlers

import (
    "net/http"
    "github.com/wentf9/MyGoFileHub/internal/application"
    "github.com/gin-gonic/gin"
)

func VersionHandler(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{
        "data": application.GetBuildInfo(),
    })
}
```

### 步骤 5: 添加 CLI --version 参数支持

**涉及文件:**
- `main.go` (修改)

在 main 函数开头添加：
```go
func main() {
    // 检查 --version 或 -v 参数
    for _, arg := range os.Args[1:] {
        if arg == "--version" || arg == "-v" {
            application.PrintVersion()
            return
        }
    }

    // ... 其余逻辑
}
```

## 验收标准

### 功能要求

- [x] 编译后的二进制文件包含版本信息
- [x] 执行 `MyGoFileHub --version` 显示版本信息
- [x] 访问 `GET /@api/v1/version` 返回 JSON 格式版本信息
- [x] 构建脚本自动获取 Git commit hash 和构建时间
- [x] 如果不在 Git 仓库中，使用默认值 "unknown"

### 技术要求

- [x] 遵循 Go 社区 ldflags 最佳实践
- [x] 版本变量使用小写（遵循 Go 命名约定）
- [x] 构建脚本兼容 Bash 和 PowerShell
- [x] 不影响现有构建流程
- [x] 代码通过 `go fmt` 和 `go vet` 检查

## 相关文件

### 新建文件

- `internal/application/version.go` - 版本信息定义
- `internal/interface/api/handlers/version_handler.go` - 版本 API handler
- `scripts/get_version.sh` - Bash 版本获取脚本
- `scripts/get_version.ps1` - PowerShell 版本获取脚本

### 修改文件

- `scripts/build.sh` - 集成 ldflags 注入
- `scripts/build.ps1` - 集成 ldflags 注入
- `scripts/README.md` - 更新使用说明
- `main.go` - 添加 --version 参数支持
- `internal/interface/api/route.go` - 添加版本路由
- `docs/API.md` - 添加版本信息 API 文档

## 注意事项

1. **ldflags 路径**: `-X` 参数中的包路径必须完全匹配，包括变量名（必须是全小写的非导出变量才能被注入）
2. **跨平台兼容**: 确保 Bash 和 PowerShell 脚本功能一致
3. **Git 缺失处理**: 在非 Git 环境（如 CI/CD 的某些场景）中构建时，使用合理的默认值
4. **版本号格式**: 推荐使用 SemVer 格式 (vMAJOR.MINOR.PATCH)
5. **构建时间格式**: 使用 RFC3339/ISO8601 格式，便于机器解析

## 参考资料

- [Go Build Flags](https://pkg.go.dev/cmd/go#hdr-Compile_packages_and_dependencies)
- [Go linker -X option](https://pkg.go.dev/cmd/link)
- [Semantic Versioning](https://semver.org/)
