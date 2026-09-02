import { defineStore, acceptHMRUpdate } from 'pinia'

export const useMutasiStore = defineStore('mutasi', {
  state: () => ({
    list: [], // array pengajuan mutasi
  }),
  actions: {
    loadData() {
      const stored = localStorage.getItem('data_mutasi')
      if (stored) this.list = JSON.parse(stored)
      else this.list = []
    },
    saveData() {
      localStorage.setItem('data_mutasi', JSON.stringify(this.list))
    },
    // Tambah pengajuan mutasi
    tambah(mutasi) {
      mutasi.id = Date.now()
      mutasi.tanggal_pengajuan = new Date().toISOString()
      mutasi.status = 'pending' // pending, disetujui, ditolak
      mutasi.catatan_verifikasi = null
      this.list.push(mutasi)
      this.saveData()
      return mutasi.id
    },
    update(id, updated) {
      const idx = this.list.findIndex((m) => m.id === id)
      if (idx !== -1) {
        this.list[idx] = { ...this.list[idx], ...updated }
        this.saveData()
      }
    },
    setujui(id, catatan = '') {
      const idx = this.list.findIndex((m) => m.id === id)
      if (idx !== -1 && this.list[idx].status === 'pending') {
        this.list[idx].status = 'disetujui'
        this.list[idx].catatan_verifikasi = catatan
        this.list[idx].tanggal_verifikasi = new Date().toISOString()
        this.saveData()
        return true
      }
      return false
    },
    tolak(id, catatan) {
      const idx = this.list.findIndex((m) => m.id === id)
      if (idx !== -1 && this.list[idx].status === 'pending') {
        this.list[idx].status = 'ditolak'
        this.list[idx].catatan_verifikasi = catatan
        this.saveData()
        return true
      }
      return false
    },
    hapus(id) {
      this.list = this.list.filter((m) => m.id !== id)
      this.saveData()
    },
    getByStatus(status) {
      return this.list.filter((m) => m.status === status)
    },
    getByJenis(jenis) {
      return this.list.filter((m) => m.jenis === jenis) // 'masuk' atau 'keluar'
    },
  },
})

if (import.meta.hot) {
  import.meta.hot.accept(acceptHMRUpdate(useMutasiStore, import.meta.hot))
}
