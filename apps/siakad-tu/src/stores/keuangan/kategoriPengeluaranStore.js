import { acceptHMRUpdate, defineStore } from 'pinia'
import { kategoriPengeluaranService } from '../../services/kategoriPengeluaranService.js'

export const useKategoriPengeluaranStore = defineStore('kategoriPengeluaran', {
  state: () => ({
    items: [],
    loading: false,
  }),
  actions: {
    async fetchAll() {
      this.loading = true
      try {
        this.items = await kategoriPengeluaranService.getAll()
      } finally {
        this.loading = false
      }
    },
    async tambah(kategori) {
      const baru = await kategoriPengeluaranService.create(kategori)
      this.items.push(baru)
    },
    // update, hapus bisa ditambahkan
  },
})

if (import.meta.hot) {
  import.meta.hot.accept(acceptHMRUpdate(useKategoriPengeluaranStore, import.meta.hot))
}
