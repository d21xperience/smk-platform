import { api } from 'boot/axios'

export const landingService = {
  getKpiMetrics() {
    return api.get('/dashboard/kpi-metrics').then((res) => res.data.data)
  },

  getStudentDistribution() {
    return api.get('/dashboard/student-distribution').then((res) => res.data.data)
  },

  getRecentActivities() {
    return api.get('/dashboard/recent-activities').then((res) => res.data.data)
  },
}
