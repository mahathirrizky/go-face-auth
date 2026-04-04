import { defineConfig, loadEnv } from 'vite'
import vue from '@vitejs/plugin-vue'
import tailwindcss from '@tailwindcss/vite'

// https://vite.dev/config/
export default defineConfig(({ mode }) => {
  const env = loadEnv(mode, process.cwd(), '');
  return {
    plugins: [
      tailwindcss(),
      vue()
    ],
    define: {
      'process.env': JSON.stringify(env)
    }
    ,
    // Dev server proxy: forward API and WebSocket requests to backend
    // - Uses VITE_API_BASE_URL from env (fall back to http://localhost:8080)
    // - Proxies '/api' (HTTP) and '/ws' (WebSocket)
    server: {
      proxy: {
        // Proxy HTTP API calls prefixed with /api -> backend (strip /api prefix)
        '^/api': {
          target: (env.VITE_API_BASE_URL || 'http://localhost:8080').replace(/\/$/, ''),
          changeOrigin: true,
          secure: false,
          rewrite: (path) => path.replace(/^\/api/, '')
        },
        // Proxy websocket endpoints starting with /ws
        '^/ws': {
          target: (env.VITE_API_BASE_URL || 'http://localhost:8080').replace(/\/$/, ''),
          changeOrigin: true,
          secure: false,
          ws: true
        }
      }
    },
    build: {
      chunkSizeWarningLimit: 600,
      rollupOptions: {
        output: {
          manualChunks(id) {
            if (id.includes('node_modules')) {
              // Split PrimeVue into its own chunk
              if (id.includes('primevue')) {
                return 'vendor-primevue';
              }
              // Split Vue & Pinia related core libraries
              if (id.includes('vue') || id.includes('pinia')) {
                return 'vendor-core';
              }
              // Anything else in node_modules goes here
              return 'vendor-utils';
            }
          }
        }
      }
    }
  }
})
