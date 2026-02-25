import { type Component, Show, onCleanup } from "solid-js";
import { store } from "../store";
import { FolderOpen, Download, Copy, Scissors, Edit2, Trash2 } from "lucide-solid";

const ContextMenu: Component = () => {

    // Handle click outside to close
    const handleClickOutside = (_e: MouseEvent) => {
        if (store.state.contextMenu.isOpen) {
            store.closeContextMenu();
        }
    };

    // Attach global click listener
    window.addEventListener('click', handleClickOutside);
    onCleanup(() => {
        window.removeEventListener('click', handleClickOutside);
    });

    const handleOpen = () => {
        const file = store.state.contextMenu.file;
        if (file && file.isDir && store.state.activeTabId) {
            store.navigate(store.state.activeTabId, file.fullPath);
        }
        store.closeContextMenu();
    };

    const handleDownload = () => {
        const file = store.state.contextMenu.file;
        if (file && !file.isDir) {
            // Create a direct download link matching the API
            const token = localStorage.getItem("token");
            const url = `/${file.fullPath.replace(/^\/+/, "")}?download=true&token=${token || ''}`;

            const a = document.createElement('a');
            a.href = url;
            a.download = file.name;
            document.body.appendChild(a);
            a.click();
            document.body.removeChild(a);
        }
        store.closeContextMenu();
    };

    const handleCopy = () => {
        if (store.state.contextMenu.file) {
            store.setClipboard([store.state.contextMenu.file], 'copy');
        }
        store.closeContextMenu();
    };

    const handleCut = () => {
        if (store.state.contextMenu.file) {
            store.setClipboard([store.state.contextMenu.file], 'cut');
        }
        store.closeContextMenu();
    };

    const handleRename = () => {
        if (store.state.contextMenu.file) {
            store.openRenameModal(store.state.contextMenu.file);
        } else {
            store.closeContextMenu();
        }
    };

    const handleDelete = () => {
        if (store.state.contextMenu.file) {
            store.openDeleteModal(store.state.contextMenu.file);
        } else {
            store.closeContextMenu();
        }
    };

    // Calculate safe position to prevent menu from overflowing the screen
    const getStyle = () => {
        const menuWidth = 192; // equivalent to w-48 (12rem)
        const menuHeight = 220; // approximate height for 6 items

        let { x, y } = store.state.contextMenu;
        const { innerWidth, innerHeight } = window;

        if (x + menuWidth > innerWidth) {
            x = innerWidth - menuWidth - 8;
        }
        if (y + menuHeight > innerHeight) {
            y = innerHeight - menuHeight - 8;
        }

        return {
            left: `${x}px`,
            top: `${y}px`,
        };
    };

    return (
        <Show when={store.state.contextMenu.isOpen && store.state.contextMenu.file}>
            <div
                class="fixed z-50 w-48 bg-white/95 backdrop-blur shadow-lg border border-slate-200 rounded-lg overflow-hidden text-sm py-1 animate-in fade-in zoom-in-95 duration-100"
                style={getStyle()}
                onContextMenu={(e) => {
                    e.preventDefault(); // Prevent native right click on the context menu itself
                }}
            >
                <Show when={store.state.contextMenu.file?.isDir}>
                    <button
                        onClick={handleOpen}
                        class="w-full flex items-center gap-2 px-3 py-1.5 text-slate-700 hover:bg-blue-50 hover:text-blue-600 transition-colors"
                    >
                        <FolderOpen size={16} /> Open
                    </button>
                </Show>

                <Show when={!store.state.contextMenu.file?.isDir}>
                    <button
                        onClick={handleDownload}
                        class="w-full flex items-center gap-2 px-3 py-1.5 text-slate-700 hover:bg-blue-50 hover:text-blue-600 transition-colors"
                    >
                        <Download size={16} /> Download
                    </button>
                </Show>

                <div class="h-px bg-slate-150 my-1 mx-2"></div>

                <button
                    onClick={handleCopy}
                    class="w-full flex items-center gap-2 px-3 py-1.5 text-slate-700 hover:bg-blue-50 hover:text-blue-600 transition-colors"
                >
                    <Copy size={16} /> Copy
                </button>

                <button
                    onClick={handleCut}
                    class="w-full flex items-center gap-2 px-3 py-1.5 text-slate-700 hover:bg-blue-50 hover:text-blue-600 transition-colors"
                >
                    <Scissors size={16} /> Cut
                </button>

                <div class="h-px bg-slate-150 my-1 mx-2"></div>

                <button
                    onClick={handleRename}
                    class="w-full flex items-center gap-2 px-3 py-1.5 text-slate-700 hover:bg-blue-50 hover:text-blue-600 transition-colors"
                >
                    <Edit2 size={16} /> Rename
                </button>

                <button
                    onClick={handleDelete}
                    class="w-full flex items-center gap-2 px-3 py-1.5 text-red-600 hover:bg-red-50 transition-colors"
                >
                    <Trash2 size={16} /> Delete
                </button>
            </div>
        </Show>
    );
};

export default ContextMenu;
