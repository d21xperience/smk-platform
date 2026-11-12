// apps/siakad-tu/src/adapters/api/kesiswaanReportRealAdapter.js

import { api } from '@/boot/axios.js'

const handleApiError = (error) => ({
  success: false,
  data: null,
  error: error.response?.data?.error || error.message || 'Terjadi kesalahan.',
})

const handleSuccess = (data) => ({ success: true, data, error: null })

export const kesiswaanReportRealAdapter = {
  async getStatistik() {
    try {
      const response = await api.get('/kesiswaan/pelaporan/statistik')
      return handleSuccess(response.data)
    } catch (error) {
      return handleApiError(error)
    }
  },

  async downloadExcel() {
    try {
      const response = await api.get('/kesiswaan/pelaporan/export/excel', {
        responseType: 'blob',
      })
      return handleSuccess(response.data)
    } catch (error) {
      return handleApiError(error)
    }
  },

  async downloadPDF() {
    try {
      const response = await api.get('/kesiswaan/pelaporan/export/pdf', {
        responseType: 'blob',
      })
      return handleSuccess(response.data)
    } catch (error) {
      return handleApiError(error)
    }
  },
}
