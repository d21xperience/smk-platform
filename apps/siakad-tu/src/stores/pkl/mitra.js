import { defineStore, acceptHMRUpdate } from 'pinia'

export const useMitraStore = defineStore('mitra', {
  state: () => ({
    list: [],
  }),
  actions: {
    loadData() {
      const stored = localStorage.getItem('data_mitra')
      if (stored) {
        this.list = JSON.parse(stored)
      } else {
        // Data dummy
        this.list = [
          {
            id: 1,
            nama: 'PT. Teknologi Nusantara',
            bidang: 'IT',
            alamat: 'Jl. Merdeka No.10',
            kuota: 5,
            rating: 4.8,
            kontak: '021-123456',
          },
          {
            id: 2,
            nama: 'CV. Karya Mandiri',
            bidang: 'Akuntansi',
            alamat: 'Jl. Sudirman No.20',
            kuota: 3,
            rating: 4.5,
            kontak: '022-987654',
          },
        ]
        this.saveData()
      }
    },
    saveData() {
      localStorage.setItem('data_mitra', JSON.stringify(this.list))
    },
    tambah(mitra) {
      mitra.id = Date.now()
      this.list.push(mitra)
      this.saveData()
    },
    update(id, updated) {
      const idx = this.list.findIndex((m) => m.id === id)
      if (idx !== -1) {
        this.list[idx] = { ...this.list[idx], ...updated }
        this.saveData()
      }
    },
    hapus(id) {
      this.list = this.list.filter((m) => m.id !== id)
      this.saveData()
    },
  },
})

if (import.meta.hot) {
  import.meta.hot.accept(acceptHMRUpdate(useMitraStore, import.meta.hot))
}
