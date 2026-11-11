import { defineStore, acceptHMRUpdate } from 'pinia'
import { api } from '@/boot/axios'

export const useSuratStore = defineStore('surat', {
  state: () => ({
    suratMasuk: [],
    suratKeluar: [],
    agendaList: [],
    isLoading: false,
    isLoadingAgenda: false,
    error: null,
  }),

  actions: {
    async fetchSuratMasuk(filter = {}) {
      this.isLoading = true
      this.error = null
      try {
        const params = { type: 'INCOMING', ...filter }
        Object.keys(params).forEach(key => (!params[key]) && delete params[key])
        const res = await api.get('/correspondences', { params })
        this.suratMasuk = res.data || []
      } catch (err) {
        this.error = err.response?.data?.error || err.message
        throw err
      } finally {
        this.isLoading = false
      }
    },

    async fetchSuratKeluar(filter = {}) {
      this.isLoading = true
      this.error = null
      try {
        const params = { type: 'OUTGOING', ...filter }
        Object.keys(params).forEach(key => (!params[key]) && delete params[key])
        const res = await api.get('/correspondences', { params })
        this.suratKeluar = res.data || []
      } catch (err) {
        this.error = err.response?.data?.error || err.message
        throw err
      } finally {
        this.isLoading = false
      }
    },

    async createSurat(payload) {
      const res = await api.post('/correspondences', payload)
      return res.data
    },

    async updateSurat(id, payload) {
      const res = await api.put(`/correspondences/${id}`, payload)
      return res.data
    },

    async deleteSurat(id) {
      await api.delete(`/correspondences/${id}`)
    },

    async signSurat(id, payload) {
      const res = await api.post(`/correspondences/${id}/sign`, payload)
      return res.data
    },

    async disposisiSurat(id, payload) {
      const res = await api.post(`/correspondences/${id}/dispositions`, payload)
      return res.data
    },

    // === AKSI EKSPEDISI ===
    async createEkspedisi(correspondenceId, payload) {
      const res = await api.post(`/correspondences/${correspondenceId}/expeditions`, payload)
      return res.data
    },

    async fetchEkspedisi(correspondenceId) {
      const res = await api.get(`/correspondences/${correspondenceId}/expeditions`)
      return res.data
    },

    // === AKSI BUKU AGENDA ===
    async fetchAgenda(filter = {}) {
      this.isLoadingAgenda = true
      this.error = null
      try {
        const res = await api.get('/agenda', { params: filter })
        this.agendaList = res.data || []
      } catch (err) {
        this.error = err.response?.data?.error || err.message
        throw err
      } finally {
        this.isLoadingAgenda = false
      }
    },

    async exportAgendaExcel(filter = {}) {
      const res = await api.get('/agenda/export', { params: filter })
      return res.data.file_url // Mengembalikan URL file
    }
  },
})

if (import.meta.hot) {
  import.meta.hot.accept(acceptHMRUpdate(useSuratStore, import.meta.hot))
}
