// apps/siakad-tu/src/adapters/api/alumniLegalisirRealAdapter.js

import { api } from '@/boot/axios.js'

const handleApiError = (error) => ({
  success: false,
  data: null,
  error: error.response?.data?.error || error.message || 'Terjadi kesalahan.',
})

const handleSuccess = (data) => ({ success: true, data, error: null })

export const alumniLegalisirRealAdapter = {
  async list() {
    try {
      const response = await api.get('/kesiswaan/legalisir')
      return handleSuccess(response.data.data || [])
    } catch (error) {
      return handleApiError(error)
    }
  },

  async create(payload) {
    try {
      const response = await api.post('/kesiswaan/legalisir', payload)
      return handleSuccess(response.data)
    } catch (error) {
      return handleApiError(error)
    }
  },

  async updateStatus(id, status) {
    try {
      const response = await api.put(`/kesiswaan/legalisir/${id}/status`, { status })
      return handleSuccess(response.data)
    } catch (error) {
      return handleApiError(error)
    }
  },

  async delete(id) {
    try {
      const response = await api.delete(`/kesiswaan/legalisir/${id}`)
      return handleSuccess(response.data)
    } catch (error) {
      return handleApiError(error)
    }
  },
}
