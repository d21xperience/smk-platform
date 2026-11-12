// apps/siakad-tu/src/adapters/api/pklPlottingRealAdapter.js

import { api } from '@/boot/axios.js'

const handleApiError = (error) => ({
  success: false,
  data: null,
  error: error.response?.data?.error || error.message || 'Terjadi kesalahan.',
})

const handleSuccess = (data) => ({ success: true, data, error: null })

export const pklPlottingRealAdapter = {
  async list() {
    try {
      const response = await api.get('/kesiswaan/pkl-plotting')
      return handleSuccess(response.data.data || [])
    } catch (error) {
      return handleApiError(error)
    }
  },

  async create(payload) {
    try {
      const response = await api.post('/kesiswaan/pkl-plotting', payload)
      return handleSuccess(response.data)
    } catch (error) {
      return handleApiError(error)
    }
  },

  async update(id, payload) {
    try {
      const response = await api.put(`/kesiswaan/pkl-plotting/${id}`, payload)
      return handleSuccess(response.data)
    } catch (error) {
      return handleApiError(error)
    }
  },

  async delete(id) {
    try {
      const response = await api.delete(`/kesiswaan/pkl-plotting/${id}`)
      return handleSuccess(response.data)
    } catch (error) {
      return handleApiError(error)
    }
  },
}
