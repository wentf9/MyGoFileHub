import { type Component, createSignal, Show } from "solid-js";
import { store } from "../../store";
import { X, HardDrive, Users, Settings as SettingsIcon } from "lucide-solid";
import SourceManager from "./SourceManager";
import UserManager from "./UserManager";

const SettingsModal: Component = () => {
  const [activeTab, setActiveTab] = createSignal<"sources" | "users">("sources");

  return (
    <Show when={store.state.isSettingsOpen}>
      <div class="fixed inset-0 z-50 flex items-center justify-center bg-black/50 backdrop-blur-sm p-4 animate-in fade-in duration-200">
        <div 
          class="bg-white w-full max-w-5xl h-[80vh] rounded-xl shadow-2xl flex overflow-hidden border border-slate-200 animate-in zoom-in-95 duration-200"
          onClick={(e) => e.stopPropagation()}
        >
          {/* Sidebar */}
          <aside class="w-64 bg-slate-50 border-r border-slate-200 flex flex-col">
            <div class="p-6">
              <div class="flex items-center gap-2 text-blue-600 mb-1">
                <SettingsIcon size={20} />
                <h2 class="font-bold text-lg">System Settings</h2>
              </div>
              <p class="text-[10px] text-slate-400 font-bold uppercase tracking-widest">Administration</p>
            </div>
            
            <nav class="flex-1 px-3 space-y-1">
              <button
                onClick={() => setActiveTab("sources")}
                classList={{
                  "w-full flex items-center gap-3 px-3 py-2 rounded-lg text-sm font-medium transition-all": true,
                  "bg-white text-blue-600 shadow-sm border border-slate-200": activeTab() === "sources",
                  "text-slate-600 hover:bg-slate-100": activeTab() !== "sources"
                }}
              >
                <HardDrive size={18} />
                Storage Sources
              </button>
              <button
                onClick={() => setActiveTab("users")}
                classList={{
                  "w-full flex items-center gap-3 px-3 py-2 rounded-lg text-sm font-medium transition-all": true,
                  "bg-white text-blue-600 shadow-sm border border-slate-200": activeTab() === "users",
                  "text-slate-600 hover:bg-slate-100": activeTab() !== "users"
                }}
              >
                <Users size={18} />
                User Management
              </button>
            </nav>
            
            <div class="p-4 border-t border-slate-200">
              <div class="text-[10px] text-slate-400 text-center uppercase tracking-widest font-bold">
                MyGoFileHub v1.0.0
              </div>
            </div>
          </aside>

          {/* Main Content */}
          <main class="flex-1 flex flex-col min-w-0 bg-white">
            <header class="h-16 border-b border-slate-100 flex items-center justify-between px-8">
              <h3 class="font-bold text-slate-800 capitalize">
                {activeTab() === "sources" ? "Manage Storage Sources" : "Manage System Users"}
              </h3>
              <button 
                onClick={() => store.toggleSettings(false)}
                class="p-2 hover:bg-slate-100 rounded-full text-slate-400 transition-colors"
              >
                <X size={20} />
              </button>
            </header>
            
            <div class="flex-1 overflow-y-auto p-8">
              <Show when={activeTab() === "sources"}>
                <SourceManager />
              </Show>
              <Show when={activeTab() === "users"}>
                <UserManager />
              </Show>
            </div>
          </main>
        </div>
      </div>
    </Show>
  );
};

export default SettingsModal;
