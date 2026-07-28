// src/composables/useDashboard.js
import { ref } from 'vue'
import { useApi } from '../useApi'

export function useDashboard() {
  const { api, request, loading } = useApi()

  const dashboardData = ref({
    stats: [],
    weeklyAttendance: [],
    todaySchedule: [],
    criticalAlerts: []
  })

  /**
   * Mengambil semua data ringkasan dashboard dalam satu request
   * GET /api/v1/siakad/dashboard/overview
   * (Parameter 'periode' otomatis di-inject oleh axios interceptor)
   */
  const fetchOverview = async () => {
    const data = await request(() => api.get('/siakad/dashboard/overview'))
    if (data) {
      dashboardData.value = {
        stats: data.stats || [],
        weeklyAttendance: data.weeklyAttendance || [],
        todaySchedule: data.todaySchedule || [],
        criticalAlerts: data.criticalAlerts || []
      }
    }
  }

  return {
    loading,
    dashboardData,
    fetchOverview
  }
}
