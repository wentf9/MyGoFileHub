import { type Component, Show, createSignal, createEffect } from "solid-js";
import { store } from "../store";
import { FileService } from "../services/api";
import type { TabSession } from "../types";
import { X, AlertTriangle, Loader2 } from "lucide-solid";

const DeleteModal: Component = () => {
    const [loading, setLoading] = createSignal(false);
    const [error, setError] = createSignal("");

    const isOpen = () => store.state.deleteModal.isOpen;
    const file = () => store.state.deleteModal.file;

    // Reset state when opening
    createEffect(() => {
        if (isOpen()) {
            setError("");
            setLoading(false);
        }
    });

    const handleDelete = async () => {
        const currentFile = file();
        if (!currentFile) return;

        setLoading(true);
        setError("");

        try {
            await FileService.deleteFile(currentFile.fullPath);
            // Refresh current tab
            if (store.state.activeTabId) {
                const tab = store.state.tabs.find((t: TabSession) => t.id === store.state.activeTabId);
                if (tab) {
                    store.navigate(tab.id, tab.currentPath);
                }
            }
            store.closeDeleteModal();
        } catch (err) {
            setError((err as Error).message || "Failed to delete item");
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
                            <AlertTriangle size={16} class="text-red-500" />
                            Delete Confirmation
                        </h2>
                        <button
                            onClick={() => store.closeDeleteModal()}
                            class="p-1 text-slate-400 hover:text-slate-600 hover:bg-slate-100 rounded-md transition-colors"
                        >
                            <X size={16} />
                        </button>
                    </div>

                    <div class="p-4">
                        <p class="text-sm text-slate-600">
                            Are you sure you want to delete <span class="font-semibold text-slate-800 break-all">{file()?.name}</span>?
                        </p>
                        <p class="text-xs text-red-500 mt-2">This action cannot be undone.</p>

                        <Show when={error()}>
                            <div class="mt-4 p-2 bg-red-50 text-red-600 text-xs rounded-md border border-red-100">
                                {error()}
                            </div>
                        </Show>
                    </div>

                    <div class="flex items-center justify-end gap-2 px-4 py-3 bg-slate-50 border-t border-slate-100">
                        <button
                            disabled={loading()}
                            onClick={() => store.closeDeleteModal()}
                            class="px-3 py-1.5 text-xs font-medium text-slate-600 hover:bg-slate-200 bg-slate-100 rounded-md transition-colors"
                        >
                            Cancel
                        </button>
                        <button
                            disabled={loading()}
                            onClick={handleDelete}
                            class="flex items-center gap-1 px-3 py-1.5 text-xs font-medium text-white bg-red-500 hover:bg-red-600 disabled:opacity-50 disabled:cursor-not-allowed rounded-md transition-colors"
                        >
                            <Show when={loading()} fallback={<span>Delete</span>}>
                                <Loader2 size={12} class="animate-spin" />
                                <span>Deleting...</span>
                            </Show>
                        </button>
                    </div>
                </div>
            </div>
        </Show>
    );
};

export default DeleteModal;
