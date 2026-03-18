# MyGoFileHub API Documentation (v1)

本文档描述了 MyGoFileHub 后端提供的 RESTful API。

## 1. 通用响应说明
系统采用统一的响应格式：
- **成功状态码**: `200 OK`
- **业务状态码**: `code: 0` 代表成功
- **响应结构**:
```json
{
    "code": 0,
    "msg": "success", 
    "data": { ... }
}
```

## 2. 认证接口 (Auth)

### 2.1 用户登录
- **方法**: `POST`
- **URL**: `/@api/v1/login`
- **Headers**:
  - `X-Client-Id`: `web-browser` (用于 ClientCheck 中间件)
- **Request Body**:
  ```json
  {
      "username": "admin",
      "password": "password"
  }
  ```
- **Response Body**:
  ```json
  {
      "token": "JWT_TOKEN_STRING",
      "msg": "Login successful"
  }
  ```

---

## 3. 文件系统接口 (Files)
*注意：所有请求需携带 Header `Authorization: Bearer <token>`。*

### 3.1 获取文件列表/详情
- **方法**: `GET`
- **URL**: `/:source_key/*path`
- **路径参数**:
  - `source_key`: 存储源唯一标识
  - `path`: 文件相对路径 (e.g., `folder/file.txt`)
- **Query 参数**:
  - `download`: `true` 时触发文件下载。
- **Response (目录内容)**:
  ```json
  {
      "code": 0,
      "data": {
          "files": [
              { "name": "docs", "isDir": true, "size": 0, "modTime": "..." },
              { "name": "readme.md", "isDir": false, "size": 1024, "modTime": "..." }
          ]
      }
  }
  ```

### 3.2 上传文件 / 创建目录
- **方法**: `POST`
- **URL**: `/:source_key/*path`
- **Query 参数**:
  - `type`: `dir` 代表创建目录，缺省为上传文件。
- **Request Body (上传文件)**: `multipart/form-data`
  - `file`: 待上传的文件字段。
- **Response**: `{"code": 0, "msg": "success"}`

### 3.3 重命名
- **方法**: `PUT`
- **URL**: `/:source_key/*path`
- **Request Body**:
  ```json
  { "new_path": "new_name_or_path" }
  ```

### 3.4 删除文件或目录
- **方法**: `DELETE`
- **URL**: `/:source_key/*path`

### 3.5 复制与移动 (Actions)
- **方法**: `POST`
- **URL**: `/@cp/:source_key/*path` (复制) 或 `/@mv/:source_key/*path` (移动)
- **Query 参数**:
  - `dest`: 目标完整路径 (e.g., `/local_disk/backup/file.txt`)

---

## 4. 存储源管理 (Sources)
*仅限管理员权限。*

### 4.1 获取存储源列表
- **方法**: `GET`
- **URL**: `/@api/v1/sources`

### 4.2 获取所有存储驱动配置结构 (Schemas)
- **方法**: `GET`
- **URL**: `/@api/v1/sources/schema`
- **Response**:
  ```json
  {
      "code": 0,
      "data": [
          {
              "type": "local",
              "name": "Local Folder",
              "config": [
                  {
                      "name": "root_path",
                      "label": "Root Path",
                      "type": "string",
                      "required": true,
                      "description": "...",
                      "default": ""
                  }
              ]
          }
      ],
      "msg": "success"
  }
  ```

### 4.3 创建存储源
- **方法**: `POST`
- **URL**: `/@api/v1/sources`
- **Request Body**:
  ```json
  {
      "key": "nas",
      "name": "My NAS",
      "type": "smb",
      "config": {
          "address": "192.168.1.100",
          "username": "admin"
      }
  }
  ```

---

## 5. 用户管理 (Users)
*仅限管理员权限。*

### 5.1 获取用户列表
- **方法**: `GET`
- **URL**: `/@api/v1/users`

### 5.2 创建用户
- **方法**: `POST`
- **URL**: `/@api/v1/users`
- **Request Body**:
  ```json
  {
      "username": "newuser",
      "password": "securepassword",
      "role": "user"
  }
  ```

---

## 6. 版本信息 (Version)

### 6.1 获取系统版本信息
- **方法**: `GET`
- **URL**: `/@api/v1/version`
- **权限**: 公开（无需认证）
- **Response**:
  ```json
  {
      "data": {
          "version": "v0.1.0-dev",
          "git_commit": "a9a2d96",
          "build_time": "2026-02-26T06:18:52Z",
          "go_version": "go1.25.5",
          "platform": "windows/amd64"
      }
  }
  ```
