import axios from 'axios'
import { useContextStore } from '../stores/contextStore.js'
import { useAuthStore } from '../stores/authStore.js'
import { API_BASE_URL } from '../config.js'

export const api = axios.create({
  baseURL: API_BASE_URL,
})

const PUBLIC_URL_PATTERNS = ['/auth/login', '/public/']

api.interceptors.request.use((config) => {
  const url = config.url || ''

  if (PUBLIC_URL_PATTERNS.some((pattern) => url.includes(pattern))) {
    return config
  }

  const contextStore = useContextStore()
  const authStore = useAuthStore()

  if (!authStore.isAuthenticated) {
    return config
  }

  // AMBIL DENGAN FALLBACK BERLAPIS (Store -> AuthStore -> LocalStorage)
  const userId = contextStore.userId || authStore.user?.id || ''
  const schoolId = contextStore.schoolId || authStore.user?.school_id || ''
  const role = authStore.userRole || 'admin'
  
  const semesterId = contextStore.selectedSemesterId || authStore.user?.academic_period_id || ''
  let academicYear = contextStore.academicYear || authStore.user?.academic_year || ''
  const semester = contextStore.semester || authStore.user?.semester || 1

  // Jika academicYear masih kosong tapi semesterId ada (misal: "20261")
  if (!academicYear && semesterId && semesterId.length === 5) {
    const year = parseInt(semesterId.substring(0, 4), 10)
    if (!isNaN(year)) academicYear = `${year}/${year + 1}`
  }

  // Fallback terakhir jika semua kosong
  if (!academicYear) academicYear = '2026/2027'

  if (userId) config.headers['x-user-id'] = userId
  if (schoolId) config.headers['x-school-id'] = schoolId
  if (role) config.headers['x-user-role'] = role
  if (semesterId) config.headers['x-academic-period-id'] = semesterId
  config.headers['x-academic-year'] = academicYear
  config.headers['x-semester'] = String(semester)

  return config
})

export default ({ app }) => {
  app.config.globalProperties.$axios = api
}
