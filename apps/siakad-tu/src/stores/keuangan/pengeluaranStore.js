import { defineStore } from 'pinia'
import { pengeluaranService } from '../../services/pengeluaranService.js'

export const usePengeluaranStore = defineStore('pengeluaran', {
  state: () => ({
    items: [],
    loading: false,
  }),
  actions: {
    async fetchAll() {
      this.loading = true
      try {
        this.items = await pengeluaranService.getAll()
      } finally {
        this.loading = false
      }
    },
    async simpan(pengeluaran) {
      if (pengeluaran.id) {
        const updated = await pengeluaranService.update(pengeluaran.id, pengeluaran)
        const index = this.items.findIndex((i) => i.id === pengeluaran.id)
        if (index !== -1) this.items[index] = updated
      } else {
        const baru = await pengeluaranService.create(pengeluaran)
        this.items.push(baru)
      }
    },
    async hapus(id) {
      await pengeluaranService.delete(id)
      this.items = this.items.filter((i) => i.id !== id)
    },
  },
})
