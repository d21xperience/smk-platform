// src/boot/axios.js
// ✅ PERBAIKAN: Impor menggunakan boot dari #q-app/wrappers
import { boot } from 'quasar/wrappers'
import axios from 'axios'
import { useAuthStore } from '@/stores/authStore'
import { useContextStore } from '@/stores/contextStore' // <-- Import context store

import {
  ApplicationError,
  NotFoundError,
  ConflictError,
  ValidationError,
  UnauthorizedError,
  BusinessRuleError,
  NetworkError,
  TimeoutError,
} from '../domain/errors/index.js'
const api = axios.create({
  baseURL: import.meta.env.QCLI_API_BASE_URL || '/api/v1',
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json',
    Accept: 'application/json',
  },
})

export default boot(async ({ app }) => {
  // Request interceptor: sisipkan token JWT jika ada
  api.interceptors.request.use(
    (config) => {
      const auth = useAuthStore()
      const contextStore = useContextStore()

      if (auth.token) {
        config.headers.Authorization = `Bearer ${auth.token}`
      }
      // 2. Inject Operational Context Headers (WAJIB untuk TU-Core Backend)
      const ctx = contextStore.current

      // Fallback untuk debugging jika context belum di-set (misal belum login)
      const finalCtx = ctx || {
        schoolId: 'SCHOOL_001',
        academicYear: '2024/2025',
        semester: 1,
        userId: 'user-debug-001',
        userRole: 'admin',
      }

      if (finalCtx) {
        config.headers['x-school-id'] = finalCtx.schoolId
        config.headers['x-academic-year'] = finalCtx.academicYear
        config.headers['x-semester'] = String(finalCtx.semester)
        config.headers['x-user-id'] = finalCtx.userId || 'system'
        config.headers['x-user-role'] = finalCtx.userRole || 'admin'

        console.log('[Axios Interceptor] Injected context headers:', {
          'x-school-id': config.headers['x-school-id'],
          'x-user-id': config.headers['x-user-id'],
        })
      }

      return config
    },
    (error) => Promise.reject(error),
  )
  // Response interceptor: tangani 401 Unauthorized
  api.interceptors.response.use(
    (response) => response,
    (error) => {
      // A. Timeout Error
      if (error.code === 'ECONNABORTED' || error.message?.includes('timeout')) {
        return Promise.reject(
          new TimeoutError('Request exceeded time limit', 'TIMEOUT', error.config),
        )
      }

      // B. Network Error (No response received from server)
      if (!error.response) {
        return Promise.reject(
          new NetworkError('Network error or server unreachable', 'NETWORK_ERROR', error.message),
        )
      }

      // C. HTTP Status Code Mapping
      const { status, data } = error.response

      // Extract error details flexibly based on backend DTO structure
      const errorMessage = data?.message || data?.error || 'An unexpected error occurred'
      const errorCode = data?.code || `HTTP_${status}`
      const errorDetails = data?.details || data?.errors || null

      switch (status) {
        case 400:
          return Promise.reject(new ValidationError(errorMessage, errorCode, errorDetails))
        case 401:
          // Optional: Trigger global logout or token refresh flow here
          return Promise.reject(new UnauthorizedError(errorMessage, errorCode, errorDetails))
        case 403:
          return Promise.reject(
            new UnauthorizedError('Forbidden: Insufficient permissions', 'FORBIDDEN', errorDetails),
          )
        case 404:
          return Promise.reject(new NotFoundError(errorMessage, errorCode, errorDetails))
        case 409:
          return Promise.reject(new ConflictError(errorMessage, errorCode, errorDetails))
        case 422:
          return Promise.reject(new BusinessRuleError(errorMessage, errorCode, errorDetails))
        case 500:
        case 502:
        case 503:
        case 504:
          return Promise.reject(
            new NetworkError(`Server error (${status})`, `SERVER_${status}`, errorDetails),
          )
        default:
          return Promise.reject(new ApplicationError(errorMessage, errorCode, status, errorDetails))
      }
    },
  )
  // Pasang instance axios ke global properties agar bisa diakses adapter/komponen
  app.config.globalProperties.$axios = api
  app.config.globalProperties.$api = api // Tambahan opsional jika ada komponen yang memanggil $api
})
export { api }
