# 任务 004: 前后端输入校验与安全防护 (Input Validation and Security)

## 进度摘要

**状态**: 已完成 ✅

- [x] 步骤 1: 创建通用验证函数 validation.go
- [x] 步骤 2: 修改 source_service.go 添加 Key 校验
- [x] 步骤 3: 修改 file_handler.go 添加文件名校验
- [x] 步骤 4: 编写单元测试
- [x] 步骤 5: 运行测试验证
- [x] 步骤 6: 更新任务文档并提交

### 验证结果

```
=== RUN   TestValidateStorageKey
--- PASS: TestValidateStorageKey (0.00s)
=== RUN   TestValidateFileName
--- PASS: TestValidateFileName (0.00s)
=== RUN   TestValidatePath
--- PASS: TestValidatePath (0.00s)
PASS
```

## 概述

为 MyGoFileHub 添加全面的输入校验和安全防护措施，重点解决存储源 key 的路由冲突问题，并防止所有形式的注入攻击和路径遍历攻击。

## 安全风险分析

### 已识别的风险点

#### 1. 存储源 Key 路由冲突 🔴 **高风险**

**问题描述**：
- 存储源 key 直接用于 URL 路径（如 `/:source_key/*path`）
- 没有校验规则，用户可能创建与系统路由冲突的 key
- 潜在冲突：`@api`、`@cp`、`@mv`、`webdav`、`ui` 等系统路径

**攻击场景**：
```
攻击者创建 source_key = "@api" 的存储源
→ 访问 /@api/v1/users 可能被路由到文件处理器而非管理 API
→ 导致路由混乱或权限绕过
```

**当前代码位置**：
- `internal/interface/api/route.go` - 路由定义
- `internal/interface/api/handlers/source_handler.go` - Create/Update 方法
- `internal/application/source_service.go` - 没有 key 校验

#### 2. 路径遍历攻击 🟡 **中风险**

**当前防护**：
- ✅ `LocalDriver.safePath()` 使用 `filepath.Clean()` 和前缀检查
- ✅ `SecureDriver` 装饰器进行权限检查

**潜在问题**：
- 需要确认所有驱动实现都正确实现了 safePath
- URL 解码后的路径可能绕过检查

#### 3. SQL 注入 🟢 **低风险**

**当前防护**：
- ✅ 使用 GORM ORM，参数化查询
- ✅ `c.ShouldBindJSON()` 自动绑定和转义

#### 4. XSS 攻击 🟡 **中风险**

**问题**：
- 响应中的用户输入（文件名、存储源名称）未进行 HTML 转义
- Gin 默认使用 `gin.JSON()`，需要确认是否自动转义

#### 5. 文件名/路径注入 🟡 **中风险**

**问题**：
- 文件上传时，文件名来自用户输入
- 没有校验文件名中的特殊字符（如 `\0`、`..`、`/`）
- 可能创建非法路径或隐藏文件

## 实现方案

### 方案 A：综合校验中间件 + 服务层验证（推荐）

#### 存储源 Key 校验规则

```go
// 定义合法 key 的模式
// 1. 只能包含：小写字母、数字、下划线、中划线
// 2. 不能以 @ 开头（避免与系统路由冲突）
// 3. 长度限制：1-32 字符
// 4. 保留字黑名单：@api, @cp, @mv, webdav, ui, static 等
var keyRegex = regexp.MustCompile(`^[a-z][a-z0-9_-]{0,31}$`)

var reservedKeys = map[string]bool{
    "api": true, "cp": true, "mv": true,
    "webdav": true, "ui": true, "static": true,
}
```

#### 文件路径校验

```go
// 验证上传文件名
func ValidateFileName(name string) error {
    // 1. 不能包含路径分隔符
    // 2. 不能以 . 开头（隐藏文件）
    // 3. 不能包含 \0 等控制字符
    // 4. 长度限制
}
```

### 方案 B：使用验证库 (go-playground/validator)

使用成熟的验证库，通过 struct tag 定义规则。

**优点**：
- 功能强大，支持复杂验证规则
- 社区活跃，文档完善

**缺点**：
- 增加依赖
- 对于简单的 key 校验可能过度设计

## 相关文件

### 需要修改的文件

| 文件 | 修改内容 |
|------|----------|
| `internal/application/source_service.go` | 添加 Key 校验逻辑 |
| `internal/interface/api/handlers/file_handler.go` | 添加路径/文件名校验 |
| `internal/interface/api/middleware/input_validation.go` | 新增通用验证中间件 |

### 可能需要的新文件

- `internal/domain/model/validation.go` - 通用验证函数
- `internal/interface/api/middleware/input_validation.go` - 输入验证中间件

## 验收标准

### 存储源 Key 校验

- [x] 只能包含小写字母、数字、下划线、中划线
- [x] 必须以字母开头
- [x] 长度 1-32 字符
- [x] 不能使用保留字（@api, @cp, @mv, webdav, ui 等）
- [x] 创建/更新时返回明确的错误信息

### 文件路径校验

- [x] 拒绝包含 `..` 的路径
- [x] 拒绝以 `/` 开头的绝对路径
- [x] 拒绝包含 `\0` 的路径
- [x] 文件名长度限制

### 通用安全

- [x] 所有用户输入经过校验
- [x] SQL 参数化查询（已满足）
- [x] 响应内容适当转义

## 注意事项

1. **向后兼容**：已存在的存储源 key 如果不符合新规则，应允许继续使用但提示用户
2. **错误信息**：错误提示不应泄露系统内部信息
3. **前端配合**：前端表单应同步添加验证，提供即时反馈

## 后续优化建议

1. 实现速率限制（Rate Limiting）防止暴力破解
2. 添加审计日志记录敏感操作
3. 实现 CSRF 防护
