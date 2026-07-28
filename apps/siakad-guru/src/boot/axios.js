// src/boot/axios.js
// ✅ PERBAIKAN: Impor menggunakan boot dari #q-app/wrappers
import { boot } from 'quasar/wrappers'
import axios from 'axios'
import { useAuthStore } from '@/stores/auth.store'

const api = axios.create({
  baseURL: import.meta.env.QCLI_API_BASE_URL || 'http://localhost:8000/api',
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json',
    Accept: 'application/json',
  },
})

export default boot(async ({ app, router }) => {
  // Request interceptor: sisipkan token JWT jika ada
  api.interceptors.request.use(
    (config) => {
      const auth = useAuthStore()
      if (auth.token) {
        config.headers.Authorization = `Bearer ${auth.token}`
      }
      return config
    },
    (error) => Promise.reject(error),
  )

  // Response interceptor: tangani 401 Unauthorized
  api.interceptors.response.use(
    (response) => response,
    (error) => {
      if (error.response?.status === 401) {
        const auth = useAuthStore()
        auth.logout()
        // router.push('/login')
        router.push('/unauthorized')
      }
      return Promise.reject(error)
    },
  )
  // Pasang instance axios ke global properties agar bisa diakses adapter/komponen
  app.config.globalProperties.$axios = api
  app.config.globalProperties.$api = api // Tambahan opsional jika ada komponen yang memanggil $api
})
export { api }
