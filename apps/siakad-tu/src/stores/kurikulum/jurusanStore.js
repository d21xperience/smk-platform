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
        const rawData = res.data.data || []
        
        // ACL (Anti-Corruption Layer): Map API response ke format UI Quasar
        this.list = rawData.map(item => ({
          id: item.jurusan_id,
          name: item.nama_jurusan,
          level: item.level_bidang_id
        }))
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
