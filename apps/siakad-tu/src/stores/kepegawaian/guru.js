import { defineStore, acceptHMRUpdate } from 'pinia'
import { api } from '@/boot/axios'

export const useGuruStore = defineStore('guru', {
  state: () => ({
    list: [],
    isLoading: false,
  }),
  actions: {
    async fetchTeachers() {
      this.isLoading = true
      try {
        const response = await api.get('/teachers', { params: { limit: 100 } })
        this.list = response.data.data || []
      } catch (error) {
        console.error('Gagal mengambil data guru:', error)
        throw error
      } finally {
        this.isLoading = false
      }
    },

    async createTeacher(payload) {
      this.isLoading = true
      try {
        const response = await api.post('/teachers', payload)
        // Asumsi backend dimodifikasi untuk mengembalikan ID guru baru
        return response.data
      } finally {
        this.isLoading = false
      }
    },

    async updateTeacher(id, payload) {
      this.isLoading = true
      try {
        await api.put(`/teachers/${id}`, payload)
      } finally {
        this.isLoading = false
      }
    },

    async deleteTeacher(id) {
      this.isLoading = true
      try {
        await api.delete(`/teachers/${id}`)
      } finally {
        this.isLoading = false
      }
    }
  },
})

if (import.meta.hot) {
  import.meta.hot.accept(acceptHMRUpdate(useGuruStore, import.meta.hot))
}
