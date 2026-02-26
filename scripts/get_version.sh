#!/bin/bash
# 获取版本信息用于 -ldflags 注入
# 用法：eval $(./scripts/get_version.sh)

# 获取 git commit hash (短版本)
GIT_COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")

# 获取构建时间 (RFC3339 格式，UTC 时间)
BUILD_TIME=$(date -u +"%Y-%m-%dT%H:%M:%SZ")

# 获取当前 tag 作为版本号 (如果没有 tag 则使用 dev 版本)
VERSION=$(git describe --tags --always --dirty 2>/dev/null || echo "v0.1.0-dev")

# 输出为 shell 变量格式
echo "VERSION=\"$VERSION\""
echo "GIT_COMMIT=\"$GIT_COMMIT\""
echo "BUILD_TIME=\"$BUILD_TIME\""
