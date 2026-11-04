import { defineStore, acceptHMRUpdate } from 'pinia'
import { api } from '@/boot/axios'

export const useKelasStore = defineStore('kelas', {
  state: () => ({
    list: [],
    isLoading: false,
    total: 0,
  }),
  actions: {
    async fetchKelas(params = { limit: 100 }) {
      this.isLoading = true
      try {
        // CUKUP TULIS '/classes' SAJA (baseURL '/api/v1' otomatis ditambahkan)
        const response = await api.get('/classes', { params })
        this.list = response.data.data.data || []
        this.total = response.data.data.total || 0
      } catch (error) {
        console.error('Gagal mengambil data kelas:', error)
        throw error
      } finally {
        this.isLoading = false
      }
    },

    async createKelas(payload) {
      this.isLoading = true
      try {
        const response = await api.post('/classes', payload)
        await this.fetchKelas()
        return response.data
      } finally {
        this.isLoading = false
      }
    },

    async updateKelas(id, payload) {
      this.isLoading = true
      try {
        const response = await api.put(`/classes/${id}`, payload)
        await this.fetchKelas()
        return response.data
      } finally {
        this.isLoading = false
      }
    },

    async deleteKelas(id) {
      this.isLoading = true
      try {
        const response = await api.delete(`/classes/${id}`)
        await this.fetchKelas()
        return response.data
      } finally {
        this.isLoading = false
      }
    }
  },
})

if (import.meta.hot) {
  import.meta.hot.accept(acceptHMRUpdate(useKelasStore, import.meta.hot))
}
