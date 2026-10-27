import { defineConfig } from 'vitest/config'
import vue from '@vitejs/plugin-vue'
import { fileURLToPath } from 'node:url'

export default defineConfig({
  plugins: [vue()],
  test: {
    environment: 'jsdom',
    globals: true,
    // Gunakan threads agar jsdom tidak dibuat berulang kali (lebih cepat)
    pool: 'threads',
    exclude: [
      '**/node_modules/**',
      '**/dist/**',
      // Kecualikan tes lama yang menggunakan console.assert (dijalankan via test:pure)
      'src/adapters/mock/__tests__/**',
      'src/composables/student/__tests__/**',
      'src/domain/context/__tests__/context.test.js',
      'src/domain/student/__tests__/**',
      'src/engine/student/__test__/**',
      'src/services/student/__tests__/**',
      'src/stores/student/__tests__/**',
    ],
  },
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url)),
      src: fileURLToPath(new URL('./src', import.meta.url)),
    },
  },
})
