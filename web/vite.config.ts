import { fileURLToPath, URL } from 'node:url'

import tailwindcss from '@tailwindcss/vite'
import react from '@vitejs/plugin-react'
import { defineConfig } from 'vite'
import { VitePWA } from 'vite-plugin-pwa'

// https://vite.dev/config/
export default defineConfig({
  plugins: [
    react(),
    tailwindcss(),
    VitePWA({
      registerType: 'autoUpdate',
      // Never cache API responses: data must always come from the server.
      workbox: { navigateFallbackDenylist: [/^\/api\//] },
      manifest: {
        name: 'FamilyShop Hub',
        short_name: 'FamilyShop',
        description: 'Quản lý cửa hàng cho hộ kinh doanh gia đình',
        lang: 'vi',
        start_url: '/',
        display: 'standalone',
        background_color: '#ffffff',
        theme_color: '#171717',
        icons: [{ src: '/favicon.svg', sizes: 'any', type: 'image/svg+xml', purpose: 'any' }],
      },
    }),
  ],
  resolve: {
    alias: { '@': fileURLToPath(new URL('./src', import.meta.url)) },
  },
  server: {
    // `make dev-web`: Vite serves the UI, the API runs on :8081 (go run).
    proxy: { '/api': 'http://localhost:8081' },
  },
})
