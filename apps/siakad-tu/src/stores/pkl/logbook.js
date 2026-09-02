import { defineStore, acceptHMRUpdate } from 'pinia'

export const useLogbookStore = defineStore('logbook', {
  state: () => ({
    list: [],
  }),
  actions: {
    loadData() {
      const stored = localStorage.getItem('data_logbook')
      if (stored) {
        this.list = JSON.parse(stored)
      } else {
        this.list = []
        this.saveData()
      }
    },
    saveData() {
      localStorage.setItem('data_logbook', JSON.stringify(this.list))
    },
    tambah(log) {
      log.id = Date.now()
      log.tanggal = log.tanggal || new Date().toISOString().split('T')[0]
      log.status = 'pending' // pending, disetujui, ditolak
      this.list.push(log)
      this.saveData()
    },
    update(id, updated) {
      const idx = this.list.findIndex((l) => l.id === id)
      if (idx !== -1) {
        this.list[idx] = { ...this.list[idx], ...updated }
        this.saveData()
      }
    },
    hapus(id) {
      this.list = this.list.filter((l) => l.id !== id)
      this.saveData()
    },
    approve(id, catatan = null) {
      const idx = this.list.findIndex((l) => l.id === id)
      if (idx !== -1) {
        this.list[idx].status = 'disetujui'
        if (catatan) this.list[idx].catatan_pembimbing = catatan
        this.saveData()
      }
    },
    reject(id, alasan) {
      const idx = this.list.findIndex((l) => l.id === id)
      if (idx !== -1) {
        this.list[idx].status = 'ditolak'
        this.list[idx].alasan_penolakan = alasan
        this.saveData()
      }
    },
  },
})

if (import.meta.hot) {
  import.meta.hot.accept(acceptHMRUpdate(useLogbookStore, import.meta.hot))
}
