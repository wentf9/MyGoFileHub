import { defineConfig } from 'vite'
import solid from 'vite-plugin-solid'

export default defineConfig({
  plugins: [solid()],
  server: {
    host: '0.0.0.0',
    port: 5173,
    strictPort: true,
    proxy: {
      // 1. 明确代理 API 路径
      '/@api': {
        target: 'http://localhost:3939',
        changeOrigin: true,
      },
      // 2. 代理文件操作路径
      // 使用正则表达式匹配，确保根路径 / 也能进入代理逻辑
      '^/(?!(src|node_modules|@vite|@id|@solid|@fs|assets)).*': {
        target: 'http://localhost:3939',
        changeOrigin: true,
        bypass: (req) => {
          const url = req.url || "";
          const pathname = url.split('?')[0];

          // 排除已知的前端静态资源后缀
          const isStaticAsset = /\.(js|ts|jsx|tsx|css|svg|png|jpg|jpeg|gif|webp|ico|woff2?|json|html)$/i.test(pathname);
          
          // 关键点：只有当浏览器明确请求 HTML（页面导航）时，才返回 index.html
          // 这样可以确保普通的 fetch('/') 请求能够通过代理到达后端
          const isPageRequest = req.headers.accept?.includes('text/html');

          if (isStaticAsset || isPageRequest) {
            return req.url; 
          }
          
          // 其他请求（如 fetch('/')）继续执行代理转发
          return;
        }
      }
    }
  },
})
