# MyGoFileHub 项目级配置 (Project Configuration)

**Project**: MyGoFileHub (Personal Cloud Storage Service)
**Default Admin**: admin / admin123 (Auto-created on first run)
**Note**: 本配置是对全局 AI 助手规则的补充，专门定义本项目的技术栈与架构边界。

## 1. 🏗️ 后端架构红线 (Go 1.25+)

本项目严格遵循 Clean Architecture（整洁架构）。依赖关系必须单向向内：`interface & infrastructure` -> `application` -> `domain`。

* **`internal/domain/` (领域层)**：核心模型与接口声明（如 VFS 虚拟文件系统）。**绝对禁止**引入外部框架（如 Gin、GORM 等）或底层系统实现细节。
* **`internal/application/` (应用层)**：纯粹的业务逻辑用例。协调领域模型和基础设施。
* **`internal/infrastructure/` (基础设施层)**：数据库具体实现、存储驱动（local/smb）、加密解密组件。
* **`internal/interface/` (接口层)**：REST API (Gin 路由与中间件)、WebDAV 实现。**🚨 严禁在 Gin 的 API Handler 中直接编写业务逻辑。**
* **性能底线**：处理文件 I/O 和高并发传输时，优先使用 Go 原生标准库（如 `io.Reader`/`io.Writer`），必须严格控制内存占用，避免大文件一次性加载到内存。

---

## 2. 🎨 前端架构红线 (Solid.js + TS + TailwindCSS)

* **状态与响应式**：优先使用 Solid 的响应式原语（Signals, Stores）。处理异步请求和数据加载时，必须使用 `createResource`。
* **🚨 框架禁区**：本项目使用的是 Solid.js。**绝对禁止**幻觉输出 React 的 API（如 `useState`、`useEffect`）。必须且只能使用 `createSignal`、`createEffect` 等 Solid 原生方法。
* **目录规范**：
    * 组件统一存放在 `src/components/`，按业务/功能模块化拆分。
    * 后端 API 请求统一封装收敛至 `src/services/api.ts`，禁止在组件内散落零碎的 `fetch` 请求。

---

## 3. 🚀 运行、测试与调试 (Run & Debug)

### 3.1 脚本使用 (Script-First)

**⚠️ 重要**：所有环境的启停、构建和测试操作，**必须且只能**使用项目 `scripts/` 目录下的脚本。永远不要直接在全局执行 `npm run ...`、`pnpm ...` 或 `go run ...`。

| 脚本 | Bash | PowerShell | 说明 |
|------|------|------------|------|
| 运行后端 | `./scripts/run-backend.sh [port]` | `.\scripts\run-backend.ps1 [port]` | 开发模式运行后端 |
| 运行前端 | `./scripts/run-frontend.sh` | `.\scripts\run-frontend.ps1` | 开发模式运行前端 |
| 构建前端 | `./scripts/build-frontend.sh` | `.\scripts\build-frontend.ps1` | 构建前端生产包 |
| 构建全项目 | `./scripts/build.sh` | `.\scripts\build.ps1` | 前端 + 后端全量构建 |
| 清理 | `./scripts/clean.sh` | `.\scripts\clean.ps1` | 清理构建产物 |
| 测试 | `./scripts/test.sh [pattern]` | `.\scripts\test.ps1 [pattern]` | 运行测试 |

**Windows PowerShell 首次使用**：
```powershell
Set-ExecutionPolicy -Scope Process -ExecutionPolicy Bypass
```

**Linux/macOS 首次使用**：
```bash
chmod +x scripts/*.sh
```

### 3.2 日志收敛

运行或调试后端服务前，必须确保 Logger 配置正确，统一将日志输出到项目根目录的 `logs/` 文件夹中。

### 3.3 开发工作流

```bash
# 1. 启动前端开发服务器（新窗口/终端）
./scripts/run-frontend.sh

# 2. 启动后端开发服务器（新窗口/终端）
./scripts/run-backend.sh

# 3. 访问 http://localhost:5173 (Vite dev server) 或 http://localhost:3939 (后端 API)
```

---

## 4. 🏛️ 核心架构模式 (Core Architecture Patterns)

### 4.1 存储驱动抽象层

**接口定义**：[`internal/domain/vfs/driver.go`](internal/domain/vfs/driver.go)

```go
type StorageDriver interface {
    DriverName() string
    Init(ctx context.Context, config map[string]any) error
    List(ctx context.Context, path string) ([]FileInfo, error)
    Open(ctx context.Context, path string) (io.ReadCloser, error)
    OpenFile(ctx context.Context, path string, flag int, perm fs.FileMode) (File, error)
    Create(ctx context.Context, path string, reader io.Reader, size int64) error
    Mkdir(ctx context.Context, path string, perm fs.FileMode) error
    Stat(ctx context.Context, path string) (FileInfo, error)
    Delete(ctx context.Context, path string) error
    Rename(ctx context.Context, srcPath, dstPath string) error
    Copy(ctx context.Context, srcPath, dstPath string) error
    Close() error
}
```

**驱动注册工厂**：[`internal/infrastructure/drivers/factory.go`](internal/infrastructure/drivers/factory.go)

- 各驱动在 `init()` 中调用 `drivers.Register(name, factory, schema)` 注册
- 通过 `drivers.CreateInstance(type)` 创建驱动实例
- `drivers.GetRegisteredSchemas()` 返回所有驱动的config schema（用于前端动态表单）

**已实现驱动**：
- `local` - 本地文件系统
- `smb` - SMB/CIFS 网络共享

### 4.2 SecureDriver 装饰器模式

**文件**：[`internal/domain/vfs/secure_driver.go`](internal/domain/vfs/secure_driver.go)

`SecureDriver` 包装底层驱动，在每次操作前进行权限检查：

```go
// 读操作 → 检查 "read" 权限
func (d *SecureDriver) List(ctx context.Context, path string) ([]FileInfo, error) {
    if ok, err := d.checker(ctx, path, "read"); !ok {
        return nil, os.ErrPermission
    }
    return d.base.List(ctx, path)
}

// 写操作 → 检查 "write" 权限
func (d *SecureDriver) Create(ctx context.Context, path string, reader io.Reader, size int64) error {
    if ok, err := d.checker(ctx, path, "write"); !ok {
        return os.ErrPermission
    }
    return d.base.Create(ctx, path, reader, size)
}

// Rename/Copy → 检查源路径 + 目标路径双重权限
```

### 4.3 权限检查 - 最长前缀匹配

**文件**：[`internal/application/permission_service.go`](internal/application/permission_service.go)

```go
// CheckPermission 检查用户在某源某路径下的权限
// 1. admin 角色拥有所有权限
// 2. 获取用户在该源下的所有权限规则
// 3. 最长前缀匹配：/work/photos 比 /work 更精确
// 4. 根据 action 返回 AllowRead 或 AllowWrite
```

**权限模型**：[`internal/domain/model/permission.go`](internal/domain/model/permission.go)
- `PathPrefix`: 路径前缀（`/` 代表整个源）
- `AllowRead`: 读/列出权限
- `AllowWrite`: 写/删/改权限

### 4.4 请求处理链路

```
HTTP Request
    → Middleware (CORS, ClientCheck)
    → JWTAuth / BasicAuth (注入 username 到 Context)
    → Handler (解析参数)
    → Service (FileService/AuthService/etc.)
    → SecureDriver (权限检查)
    → StorageDriver (local/smb 实际操作)
```

### 4.5 配置解密流程

**文件**：[`internal/application/file_service.go`](internal/application/file_service.go) `GetDriver()` 方法

1. 查询 `StorageSource` 获取配置
2. 从驱动 schema 获取需要解密的字段名（`type: "password"` 的 config items）
3. 对 `ENC:` 前缀的字段进行 AES 解密
4. 将解密后的 config 传入 `driver.Init()`

### 4.6 缓存策略

| 缓存对象 | 位置 | 失效时机 |
|----------|------|----------|
| `driverCache` | `FileService` 中 | `ClearDriverCache(sourceKey)` - 存储源配置更新后 |
| `userCache` | `PermissionService` 中 | 当前未实现失效，重启后失效 |
| `permissionCache` | `PermissionService` 中 | 当前未实现失效，重启后失效 |

---

## 5. 🌐 API 路由设计

### 5.1 文件操作 (`/`)

| 方法 | 路由 | 说明 |
|------|------|------|
| GET | `/:source_key/*path` | 获取文件列表或文件信息 |
| POST | `/:source_key/*path` | 上传文件或创建目录 (`?type=dir`) |
| PUT | `/:source_key/*path` | 重命名文件 (`body: {new_path: "..."}`) |
| DELETE | `/:source_key/*path` | 删除文件/目录 |
| POST | `/@cp/:source_key/*path` | 复制文件 (`?dest=...`) |
| POST | `/@mv/:source_key/*path` | 移动文件 (`?dest=...`) |

**特殊路由**：根路径 `/` 返回所有存储源列表（作为"虚拟根目录"）

### 5.2 WebDAV (`/webdav/:source_key/*path`)

支持的方法：`OPTIONS`, `HEAD`, `GET`, `PUT`, `POST`, `DELETE`, `PROPFIND`, `PROPPATCH`, `MKCOL`, `COPY`, `MOVE`, `LOCK`, `UNLOCK`

认证方式：Basic Auth（用户名/密码）

### 5.3 管理 API (`/@api/v1/`)

| 路由 | 方法 | 权限 | 说明 |
|------|------|------|------|
| `/login` | POST | 公开 | JWT 登录 |
| `/users` | GET/POST/PUT/DELETE | Admin | 用户管理 |
| `/sources` | GET/POST/PUT/DELETE | Admin | 存储源管理 |
| `/sources/schema` | GET | Admin | 获取驱动 schema |

---

## 6. 🖥️ 前端关键设计

### 6.1 数据适配器

**文件**：[`frontend/my-go-file-hub-ui/src/lib/adapter.ts`](frontend/my-go-file-hub-ui/src/lib/adapter.ts)

```typescript
// adaptFileNode: 将后端 FileInfo 转换为前端 FileNode
// - 生成唯一 id (fullPath)
// - 推断 mimeType (folder/image/video/file)
// - 提取 extension
```

### 6.2 API 服务封装

**文件**：[`frontend/my-go-file-hub-ui/src/services/api.ts`](frontend/my-go-file-hub-ui/src/services/api.ts)

- `AuthService.login()` - 登录
- `AdminService` - 用户和存储源管理
- `FileService` - 文件操作（自动处理根路径特殊逻辑）

### 6.3 状态管理

**文件**：[`frontend/my-go-file-hub-ui/src/store/index.ts`](frontend/my-go-file-hub-ui/src/store/index.ts)

- 多标签页支持，每个标签页独立维护：`currentPath`, `files`, `history`, `scrollTop`
- 剪贴板支持：`clipboard` (sourcePath, action: 'copy'|'move')
- 选中状态：`selection[tabId]` 存储选中的文件路径数组

---

## 7. 🗄️ 数据库模型 (SQLite + GORM)

### 7.1 User

```go
type User struct {
    ID       uint   `gorm:"primaryKey"`
    Username string `gorm:"size:32;uniqueIndex;not null"`
    Password string // bcrypt hashed
    Role     string // "admin" | "user"
}
```

### 7.2 StorageSource

```go
type StorageSource struct {
    ID     uint    `gorm:"primaryKey"`
    Key    string  `gorm:"uniqueIndex"` // 唯一标识，也是访问路径前缀
    Name   string  // 显示名称
    Type   string  // "local" | "smb"
    Config JSONMap // 驱动配置，敏感字段以 "ENC:" 前缀加密存储
}
```

### 7.3 UserPermission

```go
type UserPermission struct {
    ID         uint `gorm:"primaryKey"`
    UserID     uint `gorm:"index:idx_user_source"`
    SourceID   uint `gorm:"index:idx_user_source"`
    PathPrefix string  // 路径前缀
    AllowRead  bool
    AllowWrite bool
}
```

---

## 8. 🔐 安全机制

| 机制 | 实现位置 | 说明 |
|------|----------|------|
| JWT 认证 | `middleware/jwt_auth.go` | Bearer Token，注入 username 到 Context |
| Basic Auth | `middleware/basic_auth.go` | WebDAV 使用 |
| 密码加密 | `application/auth_service.go` | bcrypt |
| 配置加密 | `infrastructure/crypto/` | AES-256, `ENC:` 前缀标识 |
| 路径遍历防护 | 各驱动的 `safePath()` | `filepath.Clean()` + 前缀检查 |
| IP 白名单 | `config/config.go` | 支持 IP 和 CIDR 格式 |
| 角色检查 | `middleware/role_check.go` | `AdminOnly()`, `UserOnly()` |

---

## 9. 📦 配置 (Environment Variables)

| 变量 | 默认值 | 说明 |
|------|--------|------|
| `MY_GO_FILE_HUB_SERVER_PORT` | `3939` | HTTP 服务端口 |
| `MY_GO_FILE_HUB_LISTEN` | `localhost` | 绑定地址 |
| `MY_GO_FILE_HUB_DATA_DIR` | `./data` | SQLite 数据库和数据存储目录 |
| `MY_GO_FILE_HUB_WHITE_LIST` | `127.0.0.1` | IP 白名单（逗号分隔或 `*`） |
| `MY_GO_FILE_HUB_SECRET_KEY` | 自动生成 | 32 字符 AES 密钥，持久化到 `.secret_key` |
