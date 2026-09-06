import tailwindcss from '@tailwindcss/vite'
import vue from '@vitejs/plugin-vue'
import { defineConfig } from 'vite'

// https://vite.dev/config/
export default defineConfig({
  plugins: [vue(), tailwindcss()],
  server: {
    host: '0.0.0.0',
    port: 5173,
    watch: {
      usePolling: true,
    },
    proxy: {
      '/auth': {
        target: process.env.VITE_API_URL ?? 'http://backend:8080',
        changeOrigin: true,
      },
      '/profils': {
        target: process.env.VITE_API_URL ?? 'http://backend:8080',
        changeOrigin: true,
      },
      '/quiz': {
        target: process.env.VITE_API_URL ?? 'http://backend:8080',
        changeOrigin: true,
        bypass(req) {
          if (req.headers.accept?.includes('text/html')) return '/index.html'
        },
      },
    },
  },
})
