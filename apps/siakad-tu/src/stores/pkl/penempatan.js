import { defineStore, acceptHMRUpdate } from 'pinia'

export const usePenempatanStore = defineStore('penempatan', {
  state: () => ({
    list: [],
  }),
  actions: {
    loadData() {
      const stored = localStorage.getItem('data_penempatan')
      if (stored) {
        this.list = JSON.parse(stored)
      } else {
        this.list = []
        this.saveData()
      }
    },
    saveData() {
      localStorage.setItem('data_penempatan', JSON.stringify(this.list))
    },
    tambah(penempatan) {
      penempatan.id = Date.now()
      penempatan.status = 'ditempatkan' // atau 'pending', etc
      this.list.push(penempatan)
      this.saveData()
    },
    update(id, updated) {
      const idx = this.list.findIndex((p) => p.id === id)
      if (idx !== -1) {
        this.list[idx] = { ...this.list[idx], ...updated }
        this.saveData()
      }
    },
    hapus(id) {
      this.list = this.list.filter((p) => p.id !== id)
      this.saveData()
    },
  },
})

if (import.meta.hot) {
  import.meta.hot.accept(acceptHMRUpdate(usePenempatanStore, import.meta.hot))
}
