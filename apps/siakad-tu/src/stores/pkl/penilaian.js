import { defineStore, acceptHMRUpdate } from 'pinia'

export const usePenilaianStore = defineStore('penilaianPkl', {
  state: () => ({
    list: [],
  }),
  actions: {
    loadData() {
      const stored = localStorage.getItem('data_penilaian_pkl')
      if (stored) {
        this.list = JSON.parse(stored)
      } else {
        this.list = []
        this.saveData()
      }
    },
    saveData() {
      localStorage.setItem('data_penilaian_pkl', JSON.stringify(this.list))
    },
    tambah(nilai) {
      nilai.id = Date.now()
      this.list.push(nilai)
      this.saveData()
    },
    update(id, updated) {
      const idx = this.list.findIndex((n) => n.id === id)
      if (idx !== -1) {
        this.list[idx] = { ...this.list[idx], ...updated }
        this.saveData()
      }
    },
    getRataRata(siswa_id, penempatan_id) {
      const nilaiSiswa = this.list.filter(
        (n) => n.siswa_id === siswa_id && n.penempatan_id === penempatan_id,
      )
      if (nilaiSiswa.length === 0) return null
      const total = nilaiSiswa.reduce((sum, n) => sum + (n.nilai_hard + n.nilai_soft) / 2, 0)
      return total / nilaiSiswa.length
    },
  },
})

if (import.meta.hot) {
  import.meta.hot.accept(acceptHMRUpdate(usePenilaianStore, import.meta.hot))
}
