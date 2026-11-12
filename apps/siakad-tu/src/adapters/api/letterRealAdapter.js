import { api } from '@/boot/axios.js'

export const letterRealAdapter = {
  async getLetters() {
    try {
      const response = await api.get('/letters')
      return { success: true, data: response.data.data || [], error: null }
    } catch (error) {
      return { success: false, data: [], error: error.response?.data?.error || error.message }
    }
  },
  async createLetter(payload) {
    try {
      const response = await api.post('/letters', payload)
      return { success: true, data: response.data, error: null }
    } catch (error) {
      return { success: false, data: null, error: error.response?.data?.error || error.message }
    }
  },
  async completeLetter(id) {
    try {
      const response = await api.post(`/letters/${id}/complete`)
      return { success: true, data: response.data, error: null }
    } catch (error) {
      return { success: false, data: null, error: error.response?.data?.error || error.message }
    }
  }
}
