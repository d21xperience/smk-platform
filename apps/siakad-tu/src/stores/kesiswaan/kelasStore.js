import { defineStore, acceptHMRUpdate } from 'pinia'
import { api } from '@/boot/axios'

export const useKelasStore = defineStore('kelas', {
  state: () => ({
    list: [],
    isLoading: false,
    total: 0,
    // State untuk Anggota Kelas
    members: [],
    availableStudents: [],
    isLoadingMembers: false,
    isLoadingAvailable: false,
  }),
  actions: {
    async fetchKelas(params = { limit: 100 }) {
      this.isLoading = true
      try {
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
    },
    // === ENROLLMENT ACTIONS ===
    async fetchClassMembers(classId) {
      this.isLoadingMembers = true
      try {
        const response = await api.get(`/classes/${classId}/members`)
        this.members = response.data.data || []
      } finally {
        this.isLoadingMembers = false
      }
    },
    async fetchAvailableStudents(params = { limit: 100 }) {
      this.isLoadingAvailable = true
      try {
        const response = await api.get('/students/available', { params })
        this.availableStudents = response.data.data || []
      } finally {
        this.isLoadingAvailable = false
      }
    },
    async addMember(classId, studentId) {
      await api.post(`/classes/${classId}/members`, { studentId })
    },
    async removeMember(classId, studentId) {
      await api.delete(`/classes/${classId}/members/${studentId}`)
    },
    // Tambahkan ini di dalam actions kelasStore.js
    async promoteClasses() {
      this.isLoading = true
      try {
        await api.post('/classes/promote')
        await this.fetchKelas() // Refresh data setelah kenaikan
      } finally {
        this.isLoading = false
      }
    },
  },
})

if (import.meta.hot) {
  import.meta.hot.accept(acceptHMRUpdate(useKelasStore, import.meta.hot))
}
