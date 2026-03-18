import { type Component, createResource, createSignal, createEffect, For, Show } from "solid-js";
import { AdminService, FileService } from "../../services/api";
import { store } from "../../store";
import { Plus, Trash2, HardDrive, Edit2 } from "lucide-solid";
import type { StorageSource, StorageDriverSchema } from "../../types";

const SourceManager: Component = () => {
  const [sources, { refetch }] = createResource(FileService.fetchSources);
  const [schemas] = createResource(AdminService.fetchSourceSchemas);

  const [isAdding, setIsAdding] = createSignal(false);
  const [editingId, setEditingId] = createSignal<number | null>(null);
  const [loading, setLoading] = createSignal(false);

  const activeSchema = () => {
    if (!schemas()) return null;
    return schemas()?.find((s: StorageDriverSchema) => s.type === sourceForm().type) || null;
  };

  // Form State
  const [sourceForm, setSourceForm] = createSignal<Partial<StorageSource>>({
    key: "",
    name: "",
    type: "local",
    config: {}
  });

  // 当选择的驱动类型改变时，或者编辑状态重置时，如果缺少字段则注入默认值
  createEffect(() => {
    const schema = activeSchema();
    if (schema) {
      const currentConfig = Object.assign({}, sourceForm().config) || {};
      let changed = false;
      schema.config.forEach(item => {
        if (currentConfig[item.name] === undefined && item.default !== undefined) {
          currentConfig[item.name] = item.default;
          changed = true;
        }
      });
      if (changed) {
        setSourceForm(prev => ({ ...prev, config: currentConfig }));
      }
    }
  });

  const resetForm = () => {
    setSourceForm({ key: "", name: "", type: "local", config: {} });
    setIsAdding(false);
    setEditingId(null);
  };

  const handleEdit = (source: StorageSource) => {
    // 确保深拷贝 config，防止引用修改
    setSourceForm({ ...source, config: { ...source.config } });
    setEditingId(source.id);
    setIsAdding(false);
  };

  const handleSubmit = async (e: Event) => {
    e.preventDefault();
    setLoading(true);
    try {
      const id = editingId();
      if (id !== null) {
        await AdminService.updateSource(id, sourceForm());
      } else {
        await AdminService.createSource(sourceForm());
      }
      await refetch();
      await store.loadDrives(); // 同步更新全局侧边栏
      resetForm();
    } catch (err) {
      alert("Action failed: " + (err as Error).message);
    } finally {
      setLoading(false);
    }
  };

  const handleDelete = async (id: number) => {
    if (!confirm("Are you sure you want to delete this storage source?")) return;
    try {
      await AdminService.deleteSource(id);
      await refetch();
      await store.loadDrives(); // 同步更新全局侧边栏
    } catch (err) {
      alert("Delete failed: " + (err as Error).message);
    }
  };

  return (
    <div class="space-y-6">
      <header class="flex items-center justify-between">
        <div>
          <h4 class="text-sm font-semibold text-slate-800">Storage Sources</h4>
          <p class="text-xs text-slate-500">Mount external drives or local folders</p>
        </div>
        <button
          onClick={() => { resetForm(); setIsAdding(true); }}
          class="flex items-center gap-1.5 px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white text-xs font-bold rounded-lg transition-all shadow-md shadow-blue-100"
        >
          <Plus size={16} /> Add Source
        </button>
      </header>

      {/* Form (Add or Edit) */}
      <Show when={isAdding() || editingId() !== null}>
        <div class="bg-slate-50 border border-slate-200 rounded-xl p-6 animate-in slide-in-from-top-2 duration-200">
          <div class="flex items-center gap-2 mb-4 text-slate-800">
            <Show when={editingId() !== null} fallback={<Plus size={18} class="text-blue-600" />}>
              <Edit2 size={18} class="text-blue-600" />
            </Show>
            <h5 class="font-bold text-sm">{editingId() !== null ? "Edit Storage Source" : "Add New Source"}</h5>
          </div>
          <form onSubmit={handleSubmit} class="grid grid-cols-2 gap-4">
            <div class="space-y-1">
              <label class="text-[10px] font-bold text-slate-500 uppercase ml-1">Unique Key</label>
              <input
                required
                value={sourceForm().key || ""}
                placeholder="e.g. nas_works"
                class="w-full px-3 py-2 bg-white border border-slate-200 rounded-lg text-sm focus:ring-2 focus:ring-blue-500/20 outline-none disabled:bg-slate-100"
                disabled={editingId() !== null}
                onInput={(e) => setSourceForm({ ...sourceForm(), key: e.currentTarget.value })}
              />
            </div>
            <div class="space-y-1">
              <label class="text-[10px] font-bold text-slate-500 uppercase ml-1">Display Name</label>
              <input
                required
                value={sourceForm().name || ""}
                placeholder="e.g. My NAS"
                class="w-full px-3 py-2 bg-white border border-slate-200 rounded-lg text-sm focus:ring-2 focus:ring-blue-500/20 outline-none"
                onInput={(e) => setSourceForm({ ...sourceForm(), name: e.currentTarget.value })}
              />
            </div>
            <div class="space-y-1">
              <label class="text-[10px] font-bold text-slate-500 uppercase ml-1">Driver Type</label>
              <select
                class="w-full px-3 py-2 bg-white border border-slate-200 rounded-lg text-sm outline-none disabled:bg-slate-100"
                value={sourceForm().type || "local"}
                disabled={editingId() !== null} /* 不允许编辑和修改驱动类型 */
                onChange={(e) => {
                  // 切换类型时，清空之前类型的配置
                  setSourceForm({ ...sourceForm(), type: e.currentTarget.value, config: {} });
                }}
              >
                <For each={schemas()}>
                  {(schema) => (
                    <option value={schema.type}>{schema.name}</option>
                  )}
                </For>
              </select>
            </div>

            {/* Dynamic Configuration Fields */}
            <Show when={activeSchema()}>
              <div class="col-span-2 grid grid-cols-2 gap-4 mt-2 p-4 bg-white border border-slate-100 rounded-lg shadow-sm">
                <div class="col-span-2">
                  <h6 class="text-[11px] font-bold text-slate-400 uppercase tracking-wider border-b border-slate-100 pb-2 mb-2">Driver Configuration</h6>
                </div>
                <For each={activeSchema()?.config}>
                  {(field) => {
                    const value = () => sourceForm().config?.[field.name] ?? field.default ?? "";

                    return (
                      <div class="space-y-1" classList={{ "col-span-2": field.type === "string" && field.name === "root_path" }}>
                        <label class="text-[10px] font-bold text-slate-600 uppercase ml-1 flex justify-between">
                          <span>{field.label} {field.required && <span class="text-red-500">*</span>}</span>
                        </label>
                        <Show when={field.description}>
                          <p class="text-[10px] text-slate-400 ml-1 mb-1">{field.description}</p>
                        </Show>
                        <input
                          required={field.required}
                          type={field.type === "password" ? "password" : field.type === "number" ? "number" : "text"}
                          value={value()}
                          class="w-full px-3 py-2 bg-slate-50 border border-slate-200 rounded-lg text-sm focus:bg-white focus:ring-2 focus:ring-blue-500/20 outline-none transition-colors"
                          onInput={(e) => {
                            let val: string | number = e.currentTarget.value;
                            if (field.type === "number") {
                              val = val === "" ? "" : Number(val);
                            }
                            setSourceForm({ ...sourceForm(), config: { ...sourceForm().config, [field.name]: val } })
                          }}
                        />
                      </div>
                    );
                  }}
                </For>
              </div>
            </Show>
            <div class="col-span-2 flex justify-end gap-2 mt-2">
              <button
                type="button"
                onClick={resetForm}
                class="px-4 py-2 text-xs font-bold text-slate-600 hover:bg-slate-200 rounded-lg transition-colors"
              >
                Cancel
              </button>
              <button
                type="submit"
                disabled={loading()}
                class="px-6 py-2 bg-blue-600 text-white text-xs font-bold rounded-lg hover:bg-blue-700 disabled:opacity-50 transition-all shadow-md shadow-blue-100"
              >
                {loading() ? "Saving..." : (editingId() !== null ? "Update Source" : "Save Source")}
              </button>
            </div>
          </form>
        </div>
      </Show>

      {/* List */}
      <div class="border border-slate-200 rounded-xl overflow-hidden">
        <table class="w-full text-sm text-left">
          <thead class="bg-slate-50 border-b border-slate-200">
            <tr>
              <th class="px-6 py-3 font-bold text-slate-500 text-[10px] uppercase">Name</th>
              <th class="px-6 py-3 font-bold text-slate-500 text-[10px] uppercase">Key</th>
              <th class="px-6 py-3 font-bold text-slate-500 text-[10px] uppercase">Type</th>
              <th class="px-6 py-3 font-bold text-slate-500 text-[10px] uppercase text-right">Actions</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-slate-100">
            <For each={sources()} fallback={
              <tr>
                <td colspan="4" class="px-6 py-10 text-center text-slate-400">
                  <Show when={sources.loading} fallback={
                    <div class="flex flex-col items-center gap-2">
                      <HardDrive size={32} class="opacity-20" />
                      <span>No storage sources configured</span>
                    </div>
                  }>
                    <div class="flex items-center justify-center gap-2">
                      <div class="w-4 h-4 border-2 border-blue-500 border-t-transparent rounded-full animate-spin"></div>
                      <span>Loading...</span>
                    </div>
                  </Show>
                </td>
              </tr>
            }>
              {(source) => (
                <tr class="hover:bg-slate-50 transition-colors group">
                  <td class="px-6 py-4 font-medium text-slate-700">{source.name}</td>
                  <td class="px-6 py-4 text-slate-500 font-mono text-xs">{source.key}</td>
                  <td class="px-6 py-4">
                    <span class="px-2 py-0.5 bg-slate-100 text-slate-600 rounded text-[10px] font-bold uppercase">
                      {source.type}
                    </span>
                  </td>
                  <td class="px-6 py-4 text-right">
                    <div class="flex justify-end gap-1 opacity-0 group-hover:opacity-100 transition-opacity">
                      <button
                        onClick={() => handleEdit(source)}
                        class="p-1.5 text-slate-400 hover:text-blue-600 hover:bg-blue-50 rounded transition-all"
                        title="Edit Source"
                      >
                        <Edit2 size={16} />
                      </button>
                      <button
                        onClick={() => handleDelete(source.id)}
                        class="p-1.5 text-slate-400 hover:text-red-600 hover:bg-red-50 rounded transition-all"
                        title="Delete Source"
                      >
                        <Trash2 size={16} />
                      </button>
                    </div>
                  </td>
                </tr>
              )}
            </For>
          </tbody>
        </table>
      </div>
    </div>
  );
};

export default SourceManager;
