// stores/kategoriTagihanStore.js
import { acceptHMRUpdate, defineStore } from 'pinia'
import { kategoriTagihanService } from '@/services/kategoriTagihanService'

export const useKategoriTagihanStore = defineStore('kategoriTagihan', {
  state: () => ({
    items: [],
    loading: false,
    error: null,
  }),
  getters: {
    getKategoriById: (state) => (id) => state.items.find((k) => k.id === id),
    getKategoriOptions: (state) => state.items.map((k) => ({ label: k.nama, value: k.id })),
  },
  actions: {
    async fetchAll() {
      this.loading = true
      this.error = null
      try {
        this.items = await kategoriTagihanService.getAll()
      } catch (err) {
        this.error = err.message
        console.error('Failed to fetch kategori tagihan:', err)
      } finally {
        this.loading = false
      }
    },
    async tambahKategori(data) {
      this.loading = true
      try {
        const newItem = await kategoriTagihanService.create(data)
        this.items.push(newItem)
        return newItem
      } catch (err) {
        this.error = err.message
        throw err
      } finally {
        this.loading = false
      }
    },
    async updateKategori(id, data) {
      this.loading = true
      try {
        const updated = await kategoriTagihanService.update(id, data)
        const index = this.items.findIndex((i) => i.id === id)
        if (index !== -1) this.items[index] = updated
        return updated
      } catch (err) {
        this.error = err.message
        throw err
      } finally {
        this.loading = false
      }
    },
    async hapusKategori(id) {
      this.loading = true
      try {
        await kategoriTagihanService.delete(id)
        this.items = this.items.filter((i) => i.id !== id)
      } catch (err) {
        this.error = err.message
        throw err
      } finally {
        this.loading = false
      }
    },
  },
})

if (import.meta.hot) {
  import.meta.hot.accept(acceptHMRUpdate(useKategoriTagihanStore, import.meta.hot))
}
