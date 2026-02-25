import { type Component, Show, createSignal, createEffect } from "solid-js";
import { store } from "../store";
import { FileService } from "../services/api";
import type { TabSession } from "../types";
import { X, Edit2, Loader2 } from "lucide-solid";

const RenameModal: Component = () => {
    const [loading, setLoading] = createSignal(false);
    const [error, setError] = createSignal("");
    const [newName, setNewName] = createSignal("");

    const isOpen = () => store.state.renameModal.isOpen;
    const file = () => store.state.renameModal.file;

    // Reset state when opening
    createEffect(() => {
        if (isOpen() && file()) {
            setError("");
            setLoading(false);
            setNewName(file()!.name);
        }
    });

    const handleRename = async (e: Event) => {
        e.preventDefault();
        const currentFile = file();
        if (!currentFile || !newName().trim() || newName().trim() === currentFile.name) {
            store.closeRenameModal();
            return;
        }

        setLoading(true);
        setError("");

        try {
            // Calculate new path based on current path's directory
            const segments = currentFile.fullPath.split("/").filter(Boolean);
            segments.pop(); // Remove old name
            const basePath = segments.length > 0 ? `/${segments.join("/")}/` : "/";
            const targetPath = basePath + newName().trim();

            await FileService.renameFile(currentFile.fullPath, targetPath);

            // Refresh current tab
            if (store.state.activeTabId) {
                const tab = store.state.tabs.find((t: TabSession) => t.id === store.state.activeTabId);
                if (tab) {
                    store.navigate(tab.id, tab.currentPath);
                }
            }
            store.closeRenameModal();
        } catch (err) {
            setError((err as Error).message || "Failed to rename item");
        } finally {
            setLoading(false);
        }
    };

    return (
        <Show when={isOpen()}>
            <div class="fixed inset-0 z-50 flex items-center justify-center bg-slate-900/50 backdrop-blur-sm">
                <div class="bg-white rounded-xl shadow-xl w-full max-w-sm overflow-hidden animate-in fade-in zoom-in duration-200">
                    <div class="flex items-center justify-between px-4 py-3 border-b border-slate-100">
                        <h2 class="text-sm font-semibold flex items-center gap-2 text-slate-800">
                            <Edit2 size={16} class="text-blue-500" />
                            Rename Item
                        </h2>
                        <button
                            onClick={() => store.closeRenameModal()}
                            class="p-1 text-slate-400 hover:text-slate-600 hover:bg-slate-100 rounded-md transition-colors"
                        >
                            <X size={16} />
                        </button>
                    </div>

                    <form onSubmit={handleRename}>
                        <div class="p-4">
                            <label class="block text-xs font-medium text-slate-700 mb-1">
                                New Name
                            </label>
                            <input
                                type="text"
                                value={newName()}
                                onInput={(e) => setNewName(e.currentTarget.value)}
                                autofocus
                                class="w-full px-3 py-2 border border-slate-300 rounded-md text-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500 transition-shadow"
                                placeholder="Enter new name"
                            />

                            <Show when={error()}>
                                <div class="mt-4 p-2 bg-red-50 text-red-600 text-xs rounded-md border border-red-100">
                                    {error()}
                                </div>
                            </Show>
                        </div>

                        <div class="flex items-center justify-end gap-2 px-4 py-3 bg-slate-50 border-t border-slate-100">
                            <button
                                type="button"
                                disabled={loading()}
                                onClick={() => store.closeRenameModal()}
                                class="px-3 py-1.5 text-xs font-medium text-slate-600 hover:bg-slate-200 bg-slate-100 rounded-md transition-colors"
                            >
                                Cancel
                            </button>
                            <button
                                type="submit"
                                disabled={loading() || !newName().trim()}
                                class="flex items-center gap-1 px-3 py-1.5 text-xs font-medium text-white bg-blue-600 hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed rounded-md transition-colors"
                            >
                                <Show when={loading()} fallback={<span>Rename</span>}>
                                    <Loader2 size={12} class="animate-spin" />
                                    <span>Saving...</span>
                                </Show>
                            </button>
                        </div>
                    </form>
                </div>
            </div>
        </Show>
    );
};

export default RenameModal;
