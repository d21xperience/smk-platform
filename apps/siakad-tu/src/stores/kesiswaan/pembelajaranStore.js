import { defineStore, acceptHMRUpdate } from 'pinia'
import { api } from '@/boot/axios'

export const usePembelajaranStore = defineStore('pembelajaran', {
  state: () => ({
    list: [],
    isLoading: false,
  }),
  actions: {
    async fetchAll(classId = '') {
      this.isLoading = true
      try {
        const res = await api.get('/pembelajaran', { params: { classId } })
        this.list = res.data.data || []
      } finally {
        this.isLoading = false
      }
    },
    async assignTeacher(payload) {
      await api.post('/pembelajaran', payload)
      await this.fetchAll(payload.classId)
    },
    async removeAssignment(id, classId) {
      await api.delete(`/pembelajaran/${id}`)
      await this.fetchAll(classId)
    }
  },
})
if (import.meta.hot) {
  import.meta.hot.accept(acceptHMRUpdate(usePembelajaranStore, import.meta.hot))
}
