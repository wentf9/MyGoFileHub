import { defineConfig } from 'vite'
import solid from 'vite-plugin-solid'

export default defineConfig({
  plugins: [solid()],
  server: {
    host: '0.0.0.0',
    port: 5173,
    strictPort: true,
    proxy: {
      '/@api': {
        target: 'http://localhost:8080',
        changeOrigin: true,
      },
      // 代理文件操作路径
      // 排除前端必要的静态路径
      '^/(?!(src|node_modules|@vite|@id|index.html|favicon.ico|@solid|assets)).*': {
        target: 'http://localhost:8080',
        changeOrigin: true,
      }
    }
  },
})
