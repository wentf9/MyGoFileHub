
# 任务 003: 静态资源缓存优化 (Static Asset Caching)

## 进度摘要

**状态**: 已完成 ✅

- [x] 步骤 1: 创建静态资源缓存中间件
- [x] 步骤 2: 修改 route.go 应用缓存中间件
- [x] 步骤 3: 构建并验证缓存头
- [x] 步骤 4: 更新任务文档并提交

### 验证结果

**JS 文件缓存头**:
```
Cache-Control: public, max-age=31536000, immutable
```

**HTML 文件缓存头**:
```
Cache-Control: no-cache
```

**CSS 文件缓存头**:
```
Cache-Control: public, max-age=31536000, immutable
```

## 概述

为生产模式下的静态资源（JS/CSS/图片等）添加 HTTP 缓存头，优化前端加载性能，减少重复请求。

**背景**：
- 当前 `internal/interface/api/route.go` 使用 `r.StaticFS()` 服务前端静态文件
- 没有设置任何缓存头，浏览器每次都会重新请求资源
- Vite 构建产物已使用哈希命名（如 `index-DUJWmTRd.js`），可安全设置长期缓存

**目标**：
1. 为静态资源添加 `Cache-Control` 头
2. 为 HTML 文件设置不缓存或短期缓存（确保用户获取最新版本）
3. 保留开发模式不缓存的行为

## 实现方案

### 方案 A：自定义中间件（推荐）

在 `r.StaticFS()` 之前添加中间件，根据文件类型设置不同的缓存策略：

| 资源类型 | Cache-Control | 说明 |
|----------|---------------|------|
| JS/CSS/图片 | `public, max-age=31536000, immutable` | 一年缓存，内容哈希文件名确保更新 |
| HTML | `no-cache` 或 `max-age=0` | 确保获取最新版本 |
| 其他 | `public, max-age=86400` | 默认 1 天缓存 |

**优点**：
- 细粒度控制，可为不同文件类型设置不同策略
- 不影响现有路由结构

**缺点**：
- 需要中间件判断文件类型

### 方案 B：使用 Gin 的 StaticFS 配置

通过自定义 `http.FileSystem` 包装器添加缓存头。

**优点**：
- 更符合 Go 标准库设计

**缺点**：
- 实现复杂度较高

## 相关文件

### 需要修改的文件

- `internal/interface/api/route.go` - 添加缓存中间件

### 可能需要的新文件

- `internal/interface/api/middleware/static_cache.go` - 静态资源缓存中间件

## 验收标准

### 功能要求

- [x] JS/CSS 文件返回 `Cache-Control: public, max-age=31536000, immutable`
- [x] HTML 文件返回 `Cache-Control: no-cache`
- [x] 图片等静态资源返回合理的缓存时间
- [x] 开发模式不设置长期缓存

### 验证方法

```bash
# 检查 JS 文件缓存头
curl -I http://localhost:3939/ui/assets/index-DUJWmTRd.js

# 检查 HTML 文件缓存头
curl -I http://localhost:3939/ui/index.html
```

预期输出：
- JS 文件：`Cache-Control: public, max-age=31536000, immutable`
- HTML 文件：`Cache-Control: no-cache`

## 注意事项

1. **Vite 哈希文件名**：Vite 构建时使用内容哈希，文件名变化意味着内容变化，可安全设置长期缓存
2. **HTML 不缓存**：HTML 文件不包含哈希，需要确保用户获取最新版本以引用新的 JS/CSS
3. **开发模式**：开发时应禁用缓存以便调试

## 后续优化建议

1. 添加 ETag 支持，允许条件请求
2. 为 API 响应添加适当的缓存策略
