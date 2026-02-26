# MyGoFileHub

**MyGoFileHub** 是一个轻量级个人私有云存储服务，支持多种存储后端（本地文件系统、SMB 网络共享），提供 Web 界面和 WebDAV 访问功能。

## 特性

- **多存储后端支持**
  - 本地文件系统（Local）
  - SMB/CIFS 网络共享（SMB）
  - 可扩展的驱动架构，支持自定义存储后端

- **用户权限管理**
  - 基于路径前缀的细粒度权限控制
  - 读/写权限分离
  - 管理员/普通用户角色区分

- **多种访问方式**
  - 响应式 Web 界面（Solid.js + TailwindCSS）
  - WebDAV 协议支持（可挂载为网络驱动器）
  - RESTful API

- **安全可靠**
  - JWT 认证
  - 敏感配置 AES-256 加密存储
  - IP 白名单支持
  - 路径遍历防护

## 技术栈

**后端**
- Go 1.25+
- Gin Web 框架
- GORM + SQLite
- Clean Architecture 架构

**前端**
- Solid.js
- TypeScript
- TailwindCSS
- Vite 构建工具

## 快速开始

### 环境变量配置

| 变量名 | 默认值 | 说明 |
|--------|--------|------|
| `MY_GO_FILE_HUB_SERVER_PORT` | `3939` | HTTP 服务端口 |
| `MY_GO_FILE_HUB_LISTEN` | `localhost` | 绑定地址 |
| `MY_GO_FILE_HUB_DATA_DIR` | `./data` | 数据库和数据存储目录 |
| `MY_GO_FILE_HUB_WHITE_LIST` | `127.0.0.1` | IP 白名单（逗号分隔或 `*`） |
| `MY_GO_FILE_HUB_SECRET_KEY` | 自动生成 | 32 字符 AES 密钥 |
| `MY_GO_FILE_HUB_MODE` | `prod` | 运行模式：`dev` / `prod` |

### 开发模式运行

**1. 启动前端开发服务器**
```bash
# Linux/macOS
./scripts/run-frontend.sh

# Windows PowerShell
.\scripts\run-frontend.ps1
```

**2. 启动后端服务**
```bash
# Linux/macOS
./scripts/run-backend.sh

# Windows PowerShell
.\scripts\run-backend.ps1
```

**3. 访问应用**
- 前端开发服务器：http://localhost:5173
- 后端 API：http://localhost:3939

### 生产模式构建

```bash
# Linux/macOS
./scripts/build.sh

# Windows PowerShell
.\scripts\build.ps1
```

### Docker 部署（待实现）

```bash
docker run -d \
  -p 3939:3939 \
  -v /path/to/data:/app/data \
  -e MY_GO_FILE_HUB_SECRET_KEY="your-32-character-secret-key" \
  wentf9/mygofilehub:latest
```

## 默认账户

首次启动时会自动创建默认管理员账户：

- **用户名**: `admin`
- **密码**: `admin123`

**请首次登录后立即修改密码！**

## API 文档

### 认证接口

#### 登录
```bash
POST /@api/v1/login
Content-Type: application/json

{
  "username": "admin",
  "password": "admin123"
}
```

响应：
```json
{
  "code": 0,
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  },
  "msg": "success"
}
```

### 文件操作接口

| 方法 | 路由 | 说明 |
|------|------|------|
| GET | `/:source_key/*path` | 获取文件列表或文件信息 |
| POST | `/:source_key/*path` | 上传文件或创建目录 |
| PUT | `/:source_key/*path` | 重命名文件 |
| DELETE | `/:source_key/*path` | 删除文件 |
| POST | `/@cp/:source_key/*path` | 复制文件 |
| POST | `/@mv/:source_key/*path` | 移动文件 |

### WebDAV 访问

```
http://localhost:3939/webdav/:source_key/
```

支持标准 WebDAV 协议，可被以下客户端使用：
- Windows 资源管理器
- macOS Finder
- Cyberduck
- RaiDrive
- 等

## 项目结构

```
MyGoFileHub/
├── cmd/                        # 命令行工具
├── config/                     # 全局配置模块
├── frontend/                   # 前端项目
│   └── my-go-file-hub-ui/      # Solid.js + Vite
├── internal/                   # 核心业务代码
│   ├── domain/                 # 领域层：实体、接口定义
│   ├── application/            # 应用层：业务逻辑
│   ├── infrastructure/         # 基础设施层：具体实现
│   └── interface/              # 接口层：API Handlers、Middleware
├── scripts/                    # 构建/运行脚本
├── main.go                     # 应用入口
└── README.md                   # 项目文档
```

## 架构设计

MyGoFileHub 采用 **Clean Architecture（整洁架构）**，依赖关系单向向内：

```
interface & infrastructure → application → domain
                                    ↓
                              (依赖向内单向)
```

### 核心组件

- **StorageDriver** - 存储驱动抽象接口
- **SecureDriver** - 装饰器模式，包裹底层驱动进行权限检查
- **PermissionService** - 最长前缀匹配算法进行细粒度权限控制

## 开发指南

### 添加新的存储驱动

1. 在 `internal/infrastructure/drivers/` 下创建新驱动目录
2. 实现 `vfs.StorageDriver` 接口
3. 在 `init()` 函数中注册驱动

```go
func init() {
    schema := model.StorageDriverSchema{
        Type: "your-driver",
        Name: "Your Driver Name",
        Config: []model.ConfigItem{
            // 配置项定义
        },
    }
    drivers.Register("your-driver", NewYourDriver, schema)
}
```

### 运行测试

```bash
# Linux/macOS
./scripts/test.sh

# Windows PowerShell
.\scripts\test.ps1
```

## 路线图

详见 [ROADMAP.md](ROADMAP.md)

## 许可证

本项目采用 MIT 开源许可证。详见 [LICENSE](LICENSE) 文件。

## 贡献

欢迎提交 Issue 和 Pull Request！

## 致谢

- [Gin Web Framework](https://gin-gonic.com/)
- [Solid.js](https://www.solidjs.com/)
- [GORM](https://gorm.io/)
- [TailwindCSS](https://tailwindcss.com/)
