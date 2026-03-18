import { type Component, createSignal } from "solid-js";
import { store } from "../store";
import { Monitor, Loader2, Lock, User } from "lucide-solid";

const Login: Component = () => {
  const [username, setUsername] = createSignal("");
  const [password, setPassword] = createSignal("");
  const [loading, setLoading] = createSignal(false);
  const [error, setError] = createSignal("");

  const handleSubmit = async (e: Event) => {
    e.preventDefault();
    setLoading(true);
    setError("");
    
    try {
      await store.login(username(), password());
    } catch (err) {
      setError((err as Error).message || "Login failed. Please check your credentials.");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div class="min-h-screen flex items-center justify-center bg-slate-100 p-4">
      <div class="max-w-md w-full bg-white rounded-xl shadow-lg p-8 border border-slate-200">
        <div class="flex flex-col items-center mb-8">
          <div class="w-16 h-16 bg-blue-600 rounded-2xl flex items-center justify-center text-white mb-4 shadow-blue-200 shadow-lg">
            <Monitor size={32} />
          </div>
          <h1 class="text-2xl font-bold text-slate-900">MyGoFileHub</h1>
          <p class="text-slate-500 text-sm mt-1">Sign in to manage your files</p>
        </div>

        <form onSubmit={handleSubmit} class="space-y-6">
          <div class="space-y-1">
            <label class="text-xs font-semibold text-slate-600 uppercase tracking-wider ml-1">Username</label>
            <div class="relative">
              <span class="absolute left-3 top-1/2 -translate-y-1/2 text-slate-400">
                <User size={18} />
              </span>
              <input
                type="text"
                required
                value={username()}
                onInput={(e) => setUsername(e.currentTarget.value)}
                class="w-full pl-10 pr-4 py-2 bg-slate-50 border border-slate-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500/20 focus:border-blue-500 transition-all text-slate-900"
                placeholder="Enter your username"
              />
            </div>
          </div>

          <div class="space-y-1">
            <label class="text-xs font-semibold text-slate-600 uppercase tracking-wider ml-1">Password</label>
            <div class="relative">
              <span class="absolute left-3 top-1/2 -translate-y-1/2 text-slate-400">
                <Lock size={18} />
              </span>
              <input
                type="password"
                required
                value={password()}
                onInput={(e) => setPassword(e.currentTarget.value)}
                class="w-full pl-10 pr-4 py-2 bg-slate-50 border border-slate-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500/20 focus:border-blue-500 transition-all text-slate-900"
                placeholder="••••••••"
              />
            </div>
          </div>

          {error() && (
            <div class="bg-red-50 border border-red-100 text-red-600 text-sm py-2 px-3 rounded-lg flex items-center gap-2">
              <div class="w-1 h-1 bg-red-600 rounded-full animate-pulse"></div>
              {error()}
            </div>
          )}

          <button
            type="submit"
            disabled={loading()}
            class="w-full py-2.5 bg-blue-600 hover:bg-blue-700 text-white font-semibold rounded-lg shadow-md shadow-blue-200 transition-all disabled:opacity-70 disabled:cursor-not-allowed flex items-center justify-center gap-2 mt-2"
          >
            {loading() ? (
              <Loader2 size={20} class="animate-spin" />
            ) : (
              "Sign In"
            )}
          </button>
        </form>

        <div class="mt-8 text-center">
          <p class="text-slate-400 text-[10px] uppercase tracking-widest font-bold">
            Private File Hub Service
          </p>
        </div>
      </div>
    </div>
  );
};

export default Login;
