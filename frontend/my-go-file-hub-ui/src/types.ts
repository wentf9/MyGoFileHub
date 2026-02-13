// 对应后端: type FileInfo struct
export interface FileInfo {
  name: string;
  size: number;
  isDir: boolean;
  modTime: string;      // ISO string
}

// 前端使用的扩展对象 (UI Model)
export interface FileNode extends FileInfo {
  id: string;           // 唯一标识 (fullPath)
  fullPath: string;
  extension: string;
  mimeType: string;     // 'folder', 'image', 'video', 'file' 等
}

// 对应后端: type StorageSource struct
export interface StorageSource {
  id: number;
  key: string;          // 唯一标识符，也是根路径
  name: string;
  type: string;
  updatedAt: string;
}

export interface TabSession {
  id: string;
  title: string;
  currentPath: string;
  files: FileNode[];    // 当前路径下的文件列表
  loading: boolean;
  historyStack: string[];
  historyIndex: number;
}

export interface Clipboard {
  items: FileNode[];
  action: 'copy' | 'cut' | null;
}

export interface User {
  id: number;
  username: string;
  role: string;
}

export interface AuthResponse {
  token: string;
  msg: string;
  error?: string;
}

export interface AppState {
  tabs: TabSession[];
  activeTabId: string | null;
  clipboard: Clipboard;
  drives: StorageSource[];
  user: User | null;
  isAuthenticated: boolean;
}
