import { createStore } from "solid-js/store";
import type { AppState, TabSession, FileNode } from "../types";
import { FileService } from "../services/api";

const initialState: AppState = {
  tabs: [],
  activeTabId: null,
  clipboard: {
    items: [],
    action: null,
  },
  drives: [],
};

const [state, setState] = createStore<AppState>(initialState);

export const store = {
  state,
  
  // 初始化加载存储源
  loadDrives: async () => {
    const drives = await FileService.fetchSources();
    setState("drives", drives);
  },

  // Tab Actions
  addTab: async (path: string, title: string = "New Tab") => {
    const id = crypto.randomUUID();
    const newTab: TabSession = {
      id,
      title,
      currentPath: path,
      files: [],
      loading: true,
      historyStack: [path],
      historyIndex: 0,
    };
    setState("tabs", (tabs) => [...tabs, newTab]);
    setState("activeTabId", id);

    // 初始加载文件
    const files = await FileService.fetchFiles(path);
    setState("tabs", (t) => t.id === id, { files, loading: false });
  },

  closeTab: (id: string) => {
    const tabIndex = state.tabs.findIndex((t) => t.id === id);
    if (tabIndex === -1) return;

    setState("tabs", (tabs) => tabs.filter((t) => t.id !== id));
    
    if (state.activeTabId === id) {
      if (state.tabs.length > 0) {
        const nextActiveTab = state.tabs[Math.max(0, tabIndex - 1)];
        setState("activeTabId", nextActiveTab.id);
      } else {
        setState("activeTabId", null);
      }
    }
  },

  setActiveTab: (id: string) => {
    setState("activeTabId", id);
  },

  /**
   * 跳转到新路径 (Maps/Navigate Action)
   */
  navigate: async (tabId: string, path: string, title?: string) => {
    const tab = state.tabs.find(t => t.id === tabId);
    if (!tab) return;

    // 更新状态为加载中
    setState("tabs", t => t.id === tabId, { loading: true, currentPath: path });

    try {
      const files = await FileService.fetchFiles(path);
      
      setState("tabs", t => t.id === tabId, (t) => {
        const newStack = t.historyStack.slice(0, t.historyIndex + 1);
        newStack.push(path);
        
        return {
          files,
          loading: false,
          title: title ?? (path === "/" ? "Home" : path.split("/").pop() || path),
          historyStack: newStack,
          historyIndex: newStack.length - 1
        };
      });
    } catch (error) {
      console.error("Failed to navigate:", error);
      setState("tabs", t => t.id === tabId, { loading: false });
    }
  },

  // Clipboard Actions
  setClipboard: (items: FileNode[], action: 'copy' | 'cut') => {
    setState("clipboard", { items, action });
  },

  clearClipboard: () => {
    setState("clipboard", { items: [], action: null });
  },
};
