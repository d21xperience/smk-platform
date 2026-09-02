import { api } from '@/boot/axios'

export const ReportingApi = {
  getDashboardStats: (params) => api.get('/reporting/dashboard', { params }),
  getClassProgress: (classId, semesterId) =>
    api.get(`/reporting/class/${classId}/progress`, { params: { semesterId } }),
  generateReport: (params) => api.post('/reporting/generate', params),
}
