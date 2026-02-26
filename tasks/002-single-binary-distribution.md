
# 任务 002: 单文件分发支持 (Single Binary Distribution)

## 进度摘要

**状态**: 已完成 ✅

- [x] 步骤 1: 创建 embedded.go 使用 go:embed 嵌入前端静态文件
- [x] 步骤 2: 添加 Gin 路由服务嵌入式静态文件
- [x] 步骤 3: 更新构建脚本支持嵌入前端资源
- [x] 步骤 4: 添加 SPA 路由回退支持 (HTML5 History API)
- [x] 步骤 5: 添加开发模式配置 (区分嵌入模式和开发模式)
- [x] 步骤 6: 验证单文件分发功能
- [x] 步骤 7: 更新任务文档并提交

### 验证结果

**生产模式启动输出**:
```
[INFO] Running in production mode (embedded frontend)
Server starting on 127.0.0.1:3939

Production mode:
  - Access: http://127.0.0.1:3939
```

**路由注册**:
- `/ui/*` - 前端静态文件服务
- `/` - 文件操作 API
- `/@api/v1/*` - 管理 API
- `/webdav/*` - WebDAV 服务

**访问方式**:
- 生产模式：`http://localhost:3939/ui` 访问前端
- 开发模式：`http://localhost:5173` 访问 Vite 开发服务器

## 概述

实现前后端一体化部署，使用 Go 的 `embed` 功能将前端构建产物嵌入到 Go 可执行文件中，实现单文件分发。同时保留开发模式下前后端分离调试的能力。

**目标**：
1. 生产模式：单个二进制文件包含前后端所有资源
2. 开发模式：前后端独立运行，前端使用 Vite 开发服务器，后端提供 API
3. 统一访问入口：所有流量通过后端一个端口

## 实现方案

### 架构设计

**生产模式**:
```
MyGoFileHub (单个二进制文件)
├── 后端 API
│   ├── /ui/*         → 前端静态文件 (embedded)
│   ├── /*            → 文件操作 API
│   ├── /@api/v1/*    → 管理 API
│   └── /webdav/*     → WebDAV 服务
└── 数据库 (运行时创建)
```

**开发模式**:
```
前端开发服务器 (Vite, :5173)
└── 代理 API 请求到后端 (:3939)

后端服务器 (Gin, :3939)
└── 仅提供 API，不服务静态文件
```

### 配置项

| 配置 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| `MY_GO_FILE_HUB_MODE` | `dev` / `prod` | `prod` | 运行模式 |
| `MY_GO_FILE_HUB_FRONTEND_DIR` | 路径 | `./frontend` | 开发模式前端目录路径 |

## 相关文件

### 新建文件

- `internal/interface/static/embedded.go` - 嵌入式文件系统定义

### 修改文件

- `config/config.go` - 添加 `Mode` 和 `FrontendDir` 配置项
- `main.go` - 添加 `go:embed` 和模式切换逻辑
- `internal/interface/api/route.go` - 添加 `/ui` 静态文件服务和 SPA 回退
- `scripts/build.sh` / `scripts/build.ps1` - 集成前端构建和复制到嵌入目录
- `.gitignore` - 排除 `frontend/dist/` 目录

## 验收标准

### 功能要求

- [x] 生产模式下，单个二进制文件可独立运行
- [x] 访问 `http://localhost:3939/ui` 返回前端页面
- [x] 访问 `http://localhost:3939/@api/v1/version` 返回 API 数据
- [x] SPA 路由刷新后不 404 (HTML5 History API 支持)
- [x] 开发模式下前后端可分离调试
- [x] 静态资源 (JS/CSS) 正确加载

### 技术要求

- [x] 使用 `go:embed all:frontend/dist` 指令嵌入静态文件
- [x] 嵌入目录在 `.gitignore` 中排除
- [x] 构建脚本自动处理前端构建和复制
- [x] 开发模式通过环境变量切换
- [x] 代码通过 `go fmt` 和 `go vet` 检查

## 注意事项

1. **前端访问路径**: 生产模式下前端通过 `/ui` 路径访问，需要在 Vite 配置中设置 `base: '/ui/'`
2. **SPA 路由**: 使用 `NoRoute` 处理 SPA 路由回退，返回 `index.html`
3. **CGO 依赖**: `go-sqlite3` 需要 CGO 支持，编译时不能禁用 CGO
4. **构建顺序**: 必须先构建前端，复制 dist 目录，然后编译 Go 二进制

## 后续优化建议

1. ~~**根路径访问**: 可在根路径 `/` 添加重定向到 `/ui`，提升用户体验~~ - 当前实现已支持通过 `/` 访问
2. ~~**Vite base 配置**: 设置 `base: '/ui/'` 使前端资源使用相对路径~~ ✅ 已完成
3. **缓存策略**: 为静态资源添加缓存头，优化性能
