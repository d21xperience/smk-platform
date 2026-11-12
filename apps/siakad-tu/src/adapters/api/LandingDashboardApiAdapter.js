import { api } from '@/boot/axios'

export const LandingDashboardApiAdapter = {
  async getDashboardData(context) {
    // PERBAIKAN: Hapus '/api/v1' karena kemungkinan sudah ada di baseURL axios.js
    const response = await api.get('/dashboard/summary', {
      headers: {
        'x-school-id': context.schoolId,
        'x-user-id': context.userId,
        'x-academic-period-id': context.semesterId,
      },
    })

    return response.data.data
  },
}
