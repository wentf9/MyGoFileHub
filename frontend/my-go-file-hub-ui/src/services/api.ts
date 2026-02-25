import type { StorageSource, StorageDriverSchema, AuthResponse, FileNode, User } from "../types";
import { adaptFileNode } from "../lib/adapter";

const API_BASE = "/@api/v1";

interface ApiResponse<T> {
  code: number;
  data: T;
  msg: string;
}

/**
 * 基础请求封装
 */
async function request<T>(path: string, options: RequestInit = {}): Promise<T> {
  const token = localStorage.getItem("token");
  const headers = new Headers(options.headers);

  if (token) {
    headers.set("Authorization", `Bearer ${token}`);
  }

  headers.set("X-Client-Id", "web-browser");

  const response = await fetch(path, { ...options, headers });

  if (!response.ok) {
    if (response.status === 401) {
      localStorage.removeItem("token");
    }
    const errorData = await response.json().catch(() => ({}));
    throw new Error(errorData.error || errorData.msg || "Network response was not ok");
  }

  const json = await response.json();

  // 适配后端统一响应格式：解包 data
  if (json && typeof json === 'object' && 'code' in json && 'data' in json) {
    return (json as ApiResponse<T>).data;
  }

  return json as T;
}

/**
 * 认证服务
 */
export const AuthService = {
  async login(username: string, password: string): Promise<AuthResponse> {
    return request<AuthResponse>(`${API_BASE}/login`, {
      method: "POST",
      body: JSON.stringify({ username, password }),
      headers: { "Content-Type": "application/json" }
    });
  }
};

/**
 * 后台管理服务
 */
export const AdminService = {
  // 存储源管理
  async createSource(source: Partial<StorageSource>) {
    return request(`${API_BASE}/sources`, {
      method: "POST",
      body: JSON.stringify(source),
      headers: { "Content-Type": "application/json" }
    });
  },
  async updateSource(id: number, source: Partial<StorageSource>) {
    return request(`${API_BASE}/sources/${id}`, {
      method: "PUT",
      body: JSON.stringify(source),
      headers: { "Content-Type": "application/json" }
    });
  },
  async deleteSource(id: number) {
    return request(`${API_BASE}/sources/${id}`, { method: "DELETE" });
  },
  async fetchSourceSchemas(): Promise<StorageDriverSchema[]> {
    return request<StorageDriverSchema[]>(`${API_BASE}/sources/schema`);
  },

  // 用户管理
  async fetchUsers(): Promise<User[]> {
    return request<User[]>(`${API_BASE}/users`);
  },
  async createUser(user: Partial<User>) {
    return request(`${API_BASE}/users`, {
      method: "POST",
      body: JSON.stringify(user),
      headers: { "Content-Type": "application/json" }
    });
  },
  async updateUser(id: number, user: Partial<User>) {
    return request(`${API_BASE}/users/${id}`, {
      method: "PUT",
      body: JSON.stringify(user),
      headers: { "Content-Type": "application/json" }
    });
  },
  async deleteUser(id: number) {
    return request(`${API_BASE}/users/${id}`, { method: "DELETE" });
  }
};

/**
 * 真实 API 调用
 */
export const FileService = {
  /**
   * 获取指定路径下的文件列表
   * 后端路由: GET /:source_key/*path
   */
  async fetchFiles(path: string) {
    const cleanPath = path.replace(/\/+$/, "") || "/";

    if (cleanPath === "/") {
      const sources = await this.fetchSources();
      return (Array.isArray(sources) ? sources : []).map(s => ({
        id: s.key,
        name: s.name,
        size: 0,
        isDir: true,
        modTime: s.updatedAt || new Date().toISOString(),
        fullPath: s.key,
        extension: "",
        mimeType: "folder"
      } as FileNode));
    }

    const url = `/${cleanPath.replace(/^\/+/, "")}`;
    const responseData = await request<any>(url);

    // 兼容处理：responseData 可能是 { files: [...] } 或直接是数组
    const files = Array.isArray(responseData) ? responseData : (responseData?.files || []);

    return (Array.isArray(files) ? files : []).map(f => adaptFileNode(f, cleanPath));
  },

  /**
   * 获取所有存储源
   */
  async fetchSources(): Promise<StorageSource[]> {
    return request<StorageSource[]>(`${API_BASE}/sources`);
  },

  /**
   * 文件操作 (移动/复制)
   */
  async performAction(action: 'cp' | 'mv', sourcePath: string, destPath: string) {
    const url = `/@${action}/${sourcePath.replace(/^\/+/, "")}?dest=${encodeURIComponent(destPath)}`;
    return request(url, { method: "POST" });
  },

  /**
   * 删除文件
   */
  async deleteFile(path: string) {
    const url = `/${path.replace(/^\/+/, "")}`;
    return request(url, { method: "DELETE" });
  },

  /**
   * 重命名文件或移动
   */
  async renameFile(sourcePath: string, newPath: string) {
    const url = `/${sourcePath.replace(/^\/+/, "")}`;
    return request(url, {
      method: "PUT",
      body: JSON.stringify({ new_path: newPath }),
      headers: { "Content-Type": "application/json" }
    });
  },

  /**
   * 新建文件夹
   */
  async createFolder(path: string) {
    const url = `/${path.replace(/^\/+/, "")}?type=dir`;
    return request(url, { method: "POST" });
  }
};
