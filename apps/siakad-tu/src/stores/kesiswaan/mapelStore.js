import { defineStore, acceptHMRUpdate } from 'pinia'
import { api } from '@/boot/axios'

export const useMapelStore = defineStore('mapel', {
  state: () => ({
    list: [],
    isLoading: false,
  }),
  actions: {
    async fetchSubjects() {
      this.isLoading = true
      try {
        const res = await api.get('/subjects')
        this.list = res.data.data || []
      } finally {
        this.isLoading = false
      }
    },
    async createSubject(payload) {
      await api.post('/subjects', payload)
      await this.fetchSubjects()
    },
    async deleteSubject(id) {
      await api.delete(`/subjects/${id}`)
      await this.fetchSubjects()
    }
  },
})
if (import.meta.hot) {
  import.meta.hot.accept(acceptHMRUpdate(useMapelStore, import.meta.hot))
}
