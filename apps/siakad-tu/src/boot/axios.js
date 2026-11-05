import axios from 'axios'
import { useContextStore } from '../stores/contextStore.js'
import { useAuthStore } from '../stores/authStore.js'
import { API_BASE_URL } from '../config.js'

export const api = axios.create({
  baseURL: API_BASE_URL,
})

// Endpoint yang TIDAK butuh operational context headers.
const PUBLIC_URL_PATTERNS = ['/auth/login', '/public/']

api.interceptors.request.use((config) => {
  const url = config.url || ''

  // Skip header injection untuk endpoint publik (login, landing stats, dsb).
  if (PUBLIC_URL_PATTERNS.some((pattern) => url.includes(pattern))) {
    return config
  }

  const contextStore = useContextStore()
  const authStore = useAuthStore()

  // Tanpa login, jangan tempelkan header palsu.
  if (!authStore.isAuthenticated) {
    return config
  }

  const userId = contextStore.userId
  const schoolId = contextStore.schoolId
  
  // FIX: Fallback ke authStore jika contextStore kosong (terjadi saat refresh halaman)
  const semesterId = contextStore.selectedSemesterId || String(authStore.user?.academic_period_id || '')
  const academicYear = contextStore.academicYear || authStore.user?.academic_year || ''
  const semester = contextStore.semester || authStore.user?.semester || 1
  const role = authStore.userRole

  if (userId) config.headers['x-user-id'] = userId
  if (schoolId) config.headers['x-school-id'] = schoolId
  if (role) config.headers['x-user-role'] = role
  if (semesterId) config.headers['x-academic-period-id'] = semesterId
  if (academicYear) config.headers['x-academic-year'] = academicYear
  if (semester) config.headers['x-semester'] = String(semester)

  return config
})

export default ({ app }) => {
  app.config.globalProperties.$axios = api
}
