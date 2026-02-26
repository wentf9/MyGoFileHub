# Scripts

本目录包含项目的所有运行脚本，提供 Bash 和 PowerShell 两个版本。

## 快速开始

### Windows (PowerShell)

```powershell
# 首次使用，需要允许执行脚本
Set-ExecutionPolicy -Scope Process -ExecutionPolicy Bypass

# 开发模式运行
.\scripts\run-frontend.ps1   # 运行前端开发服务器
.\scripts\run-backend.ps1    # 运行后端服务器（新窗口）

# 或直接使用全项目构建
.\scripts\build.ps1
```

### Linux / macOS (Bash)

```bash
# 添加执行权限
chmod +x scripts/*.sh

# 开发模式运行
./scripts/run-frontend.sh &  # 运行前端开发服务器
./scripts/run-backend.sh     # 运行后端服务器

# 全项目构建
./scripts/build.sh
```

## 脚本说明

| 脚本 | 说明 | 使用方法 |
|------|------|----------|
| `run-backend` | 运行后端开发服务器 | `./scripts/run-backend.sh [port]` |
| `run-frontend` | 运行前端开发服务器 | `./scripts/run-frontend.sh` |
| `build-frontend` | 构建前端生产包 | `./scripts/build-frontend.sh` |
| `build` | 构建全项目（前端 + 后端） | `./scripts/build.sh` |
| `clean` | 清理构建产物 | `./scripts/clean.sh` |
| `test` | 运行测试 | `./scripts/test.sh [pattern]` |

## 环境要求

- Go 1.25+
- Node.js 18+
- Git

## 目录结构

```
scripts/
├── run-backend.sh        # Bash 版后端运行脚本
├── run-backend.ps1       # PowerShell 版后端运行脚本
├── run-frontend.sh       # Bash 版前端运行脚本
├── run-frontend.ps1      # PowerShell 版前端运行脚本
├── build-frontend.sh     # Bash 版前端构建脚本
├── build-frontend.ps1    # PowerShell 版前端构建脚本
├── build.sh              # Bash 版全项目构建脚本
├── build.ps1             # PowerShell 版全项目构建脚本
├── clean.sh              # Bash 版清理脚本
├── clean.ps1             # PowerShell 版清理脚本
├── test.sh               # Bash 版测试脚本
├── test.ps1              # PowerShell 版测试脚本
└── README.md             # 本说明文档
```
