import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'
import tailwindcss from '@tailwindcss/vite'

const apiTarget = process.env.ROUTEBOARD_API_URL || 'http://localhost:8080'

export default defineConfig({
  plugins: [react(), tailwindcss()],
  server: {
    proxy: {
      '/api': apiTarget,
      '/health': apiTarget,
    },
  },
})
