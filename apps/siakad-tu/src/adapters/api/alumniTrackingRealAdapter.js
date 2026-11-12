// apps/siakad-tu/src/adapters/api/alumniTrackingRealAdapter.js

import { api } from '@/boot/axios.js'

const handleApiError = (error) => ({
  success: false,
  data: null,
  error: error.response?.data?.error || error.message || 'Terjadi kesalahan.',
})

const handleSuccess = (data) => ({ success: true, data, error: null })

export const alumniTrackingRealAdapter = {
  async list() {
    try {
      const response = await api.get('/kesiswaan/tracer-study')
      return handleSuccess(response.data.data || [])
    } catch (error) {
      return handleApiError(error)
    }
  },

  async update(id, payload) {
    try {
      const response = await api.put(`/kesiswaan/tracer-study/${id}`, payload)
      return handleSuccess(response.data)
    } catch (error) {
      return handleApiError(error)
    }
  },
}
