import { defineStore, acceptHMRUpdate } from 'pinia'
import { api } from '@/boot/axios'

export const useJurusanStore = defineStore('jurusan', {
  state: () => ({
    list: [],
    isLoading: false,
  }),
  actions: {
    async fetchJurusan() {
      this.isLoading = true
      try {
        const res = await api.get('/jurusan')
        this.list = res.data.data || []
      } catch (e) {
        console.error('Gagal mengambil data jurusan:', e)
        throw e
      } finally {
        this.isLoading = false
      }
    }
  },
})

if (import.meta.hot) {
  import.meta.hot.accept(acceptHMRUpdate(useJurusanStore, import.meta.hot))
}
