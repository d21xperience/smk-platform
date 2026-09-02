import { acceptHMRUpdate, defineStore } from 'pinia'
import { tagihanService } from '../services/tagihanService.js'

export const useTagihanStore = defineStore('tagihan', {
  state: () => ({ items: [], loading: false }),
  actions: {
    async fetchBySiswa(siswaId) {
      this.loading = true
      this.items = await tagihanService.getBySiswa(siswaId)
      this.loading = false
    },
    async buatTagihan(data) {
      const baru = await tagihanService.create(data)
      this.items.push(baru)
    },
    async lunasiTagihan(id) {
      await tagihanService.updateStatus(id, 'lunas')
      const idx = this.items.findIndex((i) => i.id === id)
      if (idx !== -1) this.items[idx].status = 'lunas'
    },
  },
})
if (import.meta.hot) {
  import.meta.hot.accept(acceptHMRUpdate(useTagihanStore, import.meta.hot))
}
