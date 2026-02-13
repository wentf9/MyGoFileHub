import { type Component, createEffect, For, Show, onMount } from "solid-js";
import { store } from "./store";
import { Plus, X, Folder, HardDrive, Settings, Monitor, ChevronRight, File, Image, Video, FileText, Loader2, LogOut, User as UserIcon, ArrowUp, RotateCw } from "lucide-solid";
import { type FileNode } from "./types";
import Login from "./components/Login";
import SettingsModal from "./components/Settings/SettingsModal";

const FileIcon: Component<{ file: FileNode }> = (props) => {
  switch (props.file.mimeType) {
    case 'folder': return <Folder size={48} class="text-yellow-500" />;
    case 'image': return <Image size={48} class="text-blue-500" />;
    case 'video': return <Video size={48} class="text-purple-500" />;
    case 'text': return <FileText size={48} class="text-slate-500" />;
    default: return <File size={48} class="text-slate-400" />;
  }
};

const App: Component = () => {
  onMount(() => {
    if (store.state.isAuthenticated) {
      store.loadDrives();
    }
  });

  // Initialize with a default tab if authenticated and empty
  createEffect(() => {
    if (store.state.isAuthenticated && store.state.tabs.length === 0) {
      store.addTab("/", "Home");
    }
  });

  const activeTab = () => store.state.tabs.find(t => t.id === store.state.activeTabId);

  const handleUp = () => {
    const tab = activeTab();
    if (!tab || tab.currentPath === "/") return;
    const segments = tab.currentPath.split("/").filter(p => p !== "");
    segments.pop();
    const parentPath = segments.length === 0 ? "/" : "/" + segments.join("/");
    store.navigate(tab.id, parentPath);
  };

  const handleRefresh = () => {
    const tab = activeTab();
    if (!tab) return;
    store.navigate(tab.id, tab.currentPath);
  };

  return (
    <Show when={store.state.isAuthenticated} fallback={<Login />}>
      <div class="flex flex-col h-screen w-full bg-slate-50 text-slate-900 select-none">
        {/* Settings Modal */}
        <SettingsModal />

        {/* TopBar */}
        <header class="h-10 flex items-center px-4 bg-white border-b border-slate-200 gap-4">
          <div class="flex items-center gap-2 font-bold text-blue-600">
            <Monitor size={18} />
            <span>MyGoFileHub</span>
          </div>
          <div class="flex-1"></div>
          
          <div class="flex items-center gap-3">
            <div class="flex items-center gap-2 px-2 py-1 bg-slate-100 rounded-lg text-xs font-medium text-slate-600">
              <UserIcon size={14} />
              <span>{store.state.user?.username || "Admin"}</span>
            </div>
            <button 
              onClick={() => store.logout()}
              class="p-1.5 hover:bg-red-50 hover:text-red-600 rounded-lg text-slate-500 transition-colors"
              title="Logout"
            >
              <LogOut size={18} />
            </button>
            <button 
              onClick={() => store.toggleSettings(true)}
              class="p-1.5 hover:bg-slate-100 rounded-lg text-slate-500 transition-colors"
              title="System Settings"
            >
              <Settings size={18} />
            </button>
          </div>
        </header>

        {/* Ribbon / Toolbar */}
        <div class="h-12 flex items-center px-4 bg-white border-b border-slate-200 gap-2">
          <button class="px-3 py-1 hover:bg-slate-100 rounded text-sm flex items-center gap-1 border border-transparent hover:border-slate-200">
            <Plus size={16} /> New
          </button>
          <div class="w-px h-6 bg-slate-200 mx-2"></div>
          <button class="px-3 py-1 hover:bg-slate-100 rounded text-sm disabled:opacity-50" disabled={!store.state.clipboard.items.length}>Paste</button>
        </div>

        {/* TabBar */}
        <div class="flex items-center bg-slate-100 px-2 pt-1 gap-0.5 overflow-x-auto border-b border-slate-200">
          <For each={store.state.tabs}>
            {(tab) => (
              <div
                onClick={() => store.setActiveTab(tab.id)}
                classList={{
                  "group flex items-center gap-2 px-3 py-1.5 min-w-[120px] max-w-[200px] text-xs rounded-t-md border-t border-x cursor-default transition-colors": true,
                  "bg-white border-slate-200 font-medium": store.state.activeTabId === tab.id,
                  "bg-transparent border-transparent text-slate-500 hover:bg-slate-200": store.state.activeTabId !== tab.id
                }}
              >
                <Folder size={14} class="text-yellow-500 flex-shrink-0" />
                <span class="truncate flex-1">{tab.title}</span>
                <button 
                  onClick={(e) => { e.stopPropagation(); store.closeTab(tab.id); }}
                  class="opacity-0 group-hover:opacity-100 p-0.5 hover:bg-slate-300 rounded-sm transition-opacity"
                >
                  <X size={12} />
                </button>
              </div>
            )}
          </For>
          <button 
            onClick={() => store.addTab("/", "New Tab")}
            class="p-1.5 text-slate-500 hover:bg-slate-200 rounded-md ml-1"
          >
            <Plus size={16} />
          </button>
        </div>

        {/* Main Layout Area */}
        <div class="flex-1 flex overflow-hidden">
          {/* SideBar */}
          <aside class="w-64 bg-white border-r border-slate-200 flex flex-col">
            <div class="p-2">
              <div class="text-[10px] font-bold text-slate-400 px-2 py-1 uppercase tracking-wider">Quick Access</div>
              <nav class="space-y-0.5">
                <div 
                  onClick={() => store.addTab("/", "This PC")}
                  class="flex items-center gap-2 px-2 py-1.5 text-xs hover:bg-slate-100 text-slate-700 rounded-md cursor-default"
                >
                  <Monitor size={14} /> This PC
                </div>
              </nav>
              
              <div class="text-[10px] font-bold text-slate-400 px-2 py-1 mt-4 uppercase tracking-wider">Storage Sources</div>
              <nav class="space-y-0.5">
                <For each={store.state.drives}>
                  {(drive) => (
                    <div 
                      onClick={() => store.addTab(drive.key, drive.name)}
                      class="flex items-center gap-2 px-2 py-1.5 text-xs text-slate-600 hover:bg-slate-100 rounded-md cursor-default"
                    >
                      <HardDrive size={14} /> {drive.name}
                    </div>
                  )}
                </For>
              </nav>
            </div>
          </aside>

          {/* Content Area */}
          <main class="flex-1 flex flex-col bg-white">
            {/* Breadcrumbs / Address Bar */}
            <div class="h-9 flex items-center px-2 bg-white border-b border-slate-100 text-xs gap-1">
              <div class="flex items-center gap-0.5 mr-1">
                <button 
                  onClick={handleUp}
                  disabled={activeTab()?.currentPath === "/"}
                  class="p-1 hover:bg-slate-100 rounded disabled:opacity-30 transition-colors text-slate-600"
                  title="Up one level"
                >
                  <ArrowUp size={14} />
                </button>
                <button 
                  onClick={handleRefresh}
                  class="p-1 hover:bg-slate-100 rounded transition-colors text-slate-600"
                  title="Refresh"
                >
                  <RotateCw size={14} class={activeTab()?.loading ? "animate-spin" : ""} />
                </button>
              </div>

              <div class="w-px h-4 bg-slate-200 mx-1"></div>

              <div 
                onClick={() => store.navigate(store.state.activeTabId!, "/")}
                class="flex items-center gap-1 text-slate-500 hover:bg-slate-100 px-1.5 py-1 rounded cursor-pointer transition-colors ml-1"
              >
                <Monitor size={14} />
                <span>This PC</span>
              </div>
              
              <For each={activeTab()?.currentPath.split("/").filter(p => p !== "")}>
                {(segment, index) => {
                  const pathUntilNow = () => "/" + activeTab()?.currentPath.split("/").filter(p => p !== "").slice(0, index() + 1).join("/");
                  return (
                    <div class="flex items-center gap-1">
                      <ChevronRight size={14} class="text-slate-300 flex-shrink-0" />
                      <button 
                        onClick={() => store.navigate(store.state.activeTabId!, pathUntilNow())}
                        class="px-1.5 py-1 hover:bg-slate-100 rounded text-slate-700 hover:text-blue-600 transition-colors truncate max-w-[120px]"
                      >
                        {segment}
                      </button>
                    </div>
                  );
                }}
              </For>
  
              <div class="flex-1 min-w-4"></div>
              
              <Show when={activeTab()?.loading}>
                <Loader2 size={14} class="animate-spin text-blue-500 mr-2" />
              </Show>
            </div>
            
            {/* Files Grid/List */}
            <div class="flex-1 overflow-y-auto p-4">
              <Show when={activeTab()} fallback={<div class="flex items-center justify-center h-full text-slate-400 font-light">Open a tab to get started</div>}>
                <div class="grid grid-cols-[repeat(auto-fill,minmax(100px,1fr))] gap-2">
                  <For each={activeTab()?.files}>
                    {(file) => (
                      <div 
                        onDblClick={() => file.isDir && store.navigate(store.state.activeTabId!, file.fullPath)}
                        class="flex flex-col items-center gap-1 p-2 hover:bg-blue-50 rounded group cursor-default border border-transparent hover:border-blue-100 transition-colors"
                      >
                        <div class="w-16 h-16 flex items-center justify-center group-hover:bg-blue-100/50 rounded transition-colors">
                          <FileIcon file={file} />
                        </div>
                        <span class="text-xs text-center line-clamp-2 break-all px-1 w-full">{file.name}</span>
                      </div>
                    )}
                  </For>
                  <Show when={!activeTab()?.loading && activeTab()?.files.length === 0}>
                    <div class="col-span-full flex flex-col items-center justify-center py-20 text-slate-400">
                      <Folder size={48} class="opacity-20 mb-2" />
                      <p class="text-sm">This folder is empty</p>
                    </div>
                  </Show>
                </div>
              </Show>
            </div>
          </main>
        </div>

        {/* StatusBar */}
        <footer class="h-6 flex items-center px-3 bg-white border-t border-slate-200 text-[10px] text-slate-500 gap-4">
          <Show when={activeTab()}>
            <div>{activeTab()?.files.length || 0} items</div>
            <div class="w-px h-3 bg-slate-200"></div>
            <div class="truncate max-w-md">{activeTab()?.currentPath}</div>
          </Show>
          <div class="flex-1"></div>
          <div>{activeTab()?.loading ? "Loading..." : "Ready"}</div>
        </footer>
      </div>
    </Show>
  );
};

export default App;
