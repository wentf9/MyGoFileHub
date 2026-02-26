package application

import (
	"fmt"
	"runtime"
)

// 这些变量将在编译时通过 -ldflags 注入
// 变量名必须小写，才能被外部包通过 -X 参数修改
var (
	// version 语义化版本号 (e.g., v1.0.0)
	version = "v0.1.0-dev"

	// gitCommit Git commit hash (短版本)
	gitCommit = "unknown"

	// buildTime 构建时间 (RFC3339 格式)
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
	return BuildInfo{
		Version:   version,
		GitCommit: gitCommit,
		BuildTime: buildTime,
		GoVersion: runtime.Version(),
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
