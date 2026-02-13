import { FileInfo, StorageSource } from "../types";
import { adaptFileNode } from "../lib/adapter";

const API_BASE = "/@api/v1";

/**
 * 基础请求封装
 */
async function request<T>(path: string, options: RequestInit = {}): Promise<T> {
  const token = localStorage.getItem("token");
  const headers = new Headers(options.headers);
  
  if (token) {
    headers.set("Authorization", `Bearer ${token}`);
  }
  
  // 根据后端 ClientCheck 中间件要求，可能需要 Client ID
  headers.set("X-Client-Id", "web-browser");

  const response = await fetch(path, { ...options, headers });

  if (!response.ok) {
    if (response.status === 401) {
      // 处理登录过期
      console.error("Unauthorized, redirecting to login...");
    }
    const errorData = await response.json().catch(() => ({}));
    throw new Error(errorData.message || "Network response was not ok");
  }

  return response.json();
}

/**
 * 真实 API 调用
 */
export const FileService = {
  /**
   * 获取指定路径下的文件列表
   * 后端路由: GET /:source_key/*path
   */
  async fetchFiles(path: string) {
    // 确保路径以 / 开头且不重复
    const url = `/${path.replace(/^\/+/, "")}`;
    const data = await request<FileInfo[]>(url);
    return data.map(f => adaptFileNode(f, path));
  },

  /**
   * 获取所有存储源
   * 后端路由: GET /@api/v1/sources
   */
  async fetchSources(): Promise<StorageSource[]> {
    return request<StorageSource[]>(`${API_BASE}/sources`);
  },

  /**
   * 文件操作 (移动/复制)
   * 后端路由: POST /@cp/:source_key/*path?dest=...
   */
  async performAction(action: 'cp' | 'mv', sourcePath: string, destPath: string) {
    const url = `/@${action}/${sourcePath.replace(/^\/+/, "")}?dest=${encodeURIComponent(destPath)}`;
    return request(url, { method: "POST" });
  },

  /**
   * 删除文件
   * 后端路由: DELETE /:source_key/*path
   */
  async deleteFile(path: string) {
    const url = `/${path.replace(/^\/+/, "")}`;
    return request(url, { method: "DELETE" });
  }
};
