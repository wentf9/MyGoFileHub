import { type Component, createResource, createSignal, For, Show } from "solid-js";
import { AdminService } from "../../services/api";
import { Trash2, UserPlus, Shield, User as UserIcon } from "lucide-solid";
import type { User } from "../../types";

const UserManager: Component = () => {
  const [users, { refetch }] = createResource(AdminService.fetchUsers);
  const [isAdding, setIsAdding] = createSignal(false);
  const [loading, setLoading] = createSignal(false);

  // Form State
  const [newUser, setNewUser] = createSignal<Partial<User>>({
    username: "",
    password: "",
    role: "user"
  });

  const handleCreate = async (e: Event) => {
    e.preventDefault();
    setLoading(true);
    try {
      await AdminService.createUser(newUser());
      await refetch();
      setIsAdding(false);
      setNewUser({ username: "", password: "", role: "user" });
    } catch (err) {
      alert("Failed to create user: " + (err as Error).message);
    } finally {
      setLoading(false);
    }
  };

  const handleDelete = async (id: number) => {
    if (!confirm("Delete this user? This action cannot be undone.")) return;
    try {
      await AdminService.deleteUser(id);
      await refetch();
    } catch (err) {
      alert("Delete failed: " + (err as Error).message);
    }
  };

  return (
    <div class="space-y-6">
      <header class="flex items-center justify-between">
        <div>
          <h4 class="text-sm font-semibold text-slate-800">System Users</h4>
          <p class="text-xs text-slate-500">Manage account access and permissions</p>
        </div>
        <button 
          onClick={() => setIsAdding(true)}
          class="flex items-center gap-1.5 px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white text-xs font-bold rounded-lg transition-all shadow-md shadow-blue-100"
        >
          <UserPlus size={16} /> Create User
        </button>
      </header>

      {/* Add Form */}
      <Show when={isAdding()}>
        <div class="bg-slate-50 border border-slate-200 rounded-xl p-6 animate-in slide-in-from-top-2 duration-200">
          <form onSubmit={handleCreate} class="grid grid-cols-3 gap-4">
            <div class="space-y-1">
              <label class="text-[10px] font-bold text-slate-500 uppercase ml-1">Username</label>
              <input 
                required
                placeholder="Unique username"
                class="w-full px-3 py-2 bg-white border border-slate-200 rounded-lg text-sm outline-none focus:ring-2 focus:ring-blue-500/20"
                onInput={(e) => setNewUser({ ...newUser(), username: e.currentTarget.value })}
              />
            </div>
            <div class="space-y-1">
              <label class="text-[10px] font-bold text-slate-500 uppercase ml-1">Initial Password</label>
              <input 
                required
                type="password"
                placeholder="••••••••"
                class="w-full px-3 py-2 bg-white border border-slate-200 rounded-lg text-sm outline-none focus:ring-2 focus:ring-blue-500/20"
                onInput={(e) => setNewUser({ ...newUser(), password: e.currentTarget.value })}
              />
            </div>
            <div class="space-y-1">
              <label class="text-[10px] font-bold text-slate-500 uppercase ml-1">System Role</label>
              <select 
                class="w-full px-3 py-2 bg-white border border-slate-200 rounded-lg text-sm outline-none"
                onInput={(e) => setNewUser({ ...newUser(), role: e.currentTarget.value })}
              >
                <option value="user">Standard User</option>
                <option value="admin">Administrator</option>
              </select>
            </div>
            <div class="col-span-3 flex justify-end gap-2 mt-2">
              <button 
                type="button" 
                onClick={() => setIsAdding(false)}
                class="px-4 py-2 text-xs font-bold text-slate-600 hover:bg-slate-200 rounded-lg transition-colors"
              >
                Cancel
              </button>
              <button 
                type="submit"
                disabled={loading()}
                class="px-6 py-2 bg-blue-600 text-white text-xs font-bold rounded-lg hover:bg-blue-700 disabled:opacity-50 transition-all shadow-md shadow-blue-100"
              >
                {loading() ? "Creating..." : "Create User"}
              </button>
            </div>
          </form>
        </div>
      </Show>

      {/* List */}
      <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
        <For each={users()} fallback={
           <div class="col-span-full py-10 text-center text-slate-400">
             <Show when={users.loading} fallback={<span>No users found</span>}>
               <div class="flex items-center justify-center gap-2 font-medium">
                 <div class="w-4 h-4 border-2 border-blue-500 border-t-transparent rounded-full animate-spin"></div>
                 <span>Fetching users...</span>
               </div>
             </Show>
           </div>
        }>
          {(user) => (
            <div class="p-4 bg-white border border-slate-200 rounded-xl flex items-center justify-between hover:shadow-md transition-all group">
              <div class="flex items-center gap-3">
                <div class={`w-10 h-10 rounded-full flex items-center justify-center ${user.role === 'admin' ? 'bg-amber-50 text-amber-600' : 'bg-slate-50 text-slate-600'}`}>
                  {user.role === 'admin' ? <Shield size={20} /> : <UserIcon size={20} />}
                </div>
                <div>
                  <div class="font-bold text-slate-800 text-sm flex items-center gap-2">
                    {user.username}
                    <Show when={user.role === 'admin'}>
                      <span class="px-1.5 py-0.5 bg-amber-100 text-amber-700 text-[8px] font-black uppercase rounded tracking-tighter">Admin</span>
                    </Show>
                  </div>
                  <div class="text-[10px] text-slate-500 font-medium">ID: #{user.id} • Registered System Account</div>
                </div>
              </div>
              <button 
                onClick={() => handleDelete(user.id)}
                class="p-2 text-slate-300 hover:text-red-600 hover:bg-red-50 rounded-lg transition-all opacity-0 group-hover:opacity-100"
              >
                <Trash2 size={18} />
              </button>
            </div>
          )}
        </For>
      </div>
    </div>
  );
};

export default UserManager;
