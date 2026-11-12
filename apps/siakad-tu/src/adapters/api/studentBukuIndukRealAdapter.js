// apps/siakad-tu/src/adapters/api/studentBukuIndukRealAdapter.js

import { api } from '@/boot/axios.js'

const handleApiError = (error) => ({
  success: false,
  data: null,
  error: error.response?.data?.error || error.message || 'Terjadi kesalahan.',
})

const handleSuccess = (data) => ({ success: true, data, error: null })

export const studentBukuIndukRealAdapter = {
  async list() {
    try {
      const response = await api.get('/kesiswaan/buku-induk-siswa')
      return handleSuccess(response.data.data || [])
    } catch (error) {
      return handleApiError(error)
    }
  },

  async getById(studentId) {
    try {
      const response = await api.get(`/kesiswaan/buku-induk-siswa/${studentId}`)
      return handleSuccess(response.data)
    } catch (error) {
      return handleApiError(error)
    }
  },

  async downloadPDF(studentId) {
    try {
      const response = await api.get(`/kesiswaan/buku-induk-siswa/${studentId}/pdf`, {
        responseType: 'blob',
      })
      return handleSuccess(response.data)
    } catch (error) {
      return handleApiError(error)
    }
  },
}
