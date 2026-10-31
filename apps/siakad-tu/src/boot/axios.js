import axios from 'axios'
import { useContextStore } from '../stores/contextStore.js'
import { API_BASE_URL } from '../config.js'

export const api = axios.create({
  baseURL: API_BASE_URL,
})

api.interceptors.request.use((config) => {
  try {
    const contextStore = useContextStore()
    const userId = contextStore.userId
    const schoolId = contextStore.schoolId
    const semesterId = contextStore.selectedSemesterId

    if (schoolId) config.headers['x-school-id'] = schoolId
    if (userId) config.headers['x-user-id'] = userId
    // PERBAIKAN: Sesuaikan dengan header yang dibaca Backend
    if (semesterId) config.headers['x-academic-period'] = semesterId
  } catch (e) {
    // Abaikan jika store belum siap
    console.log(e)
  }

  return config
})

export default ({ app }) => {
  app.config.globalProperties.$axios = api
}
