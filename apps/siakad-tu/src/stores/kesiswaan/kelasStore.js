import { defineStore, acceptHMRUpdate } from 'pinia'

export const useKelasStore = defineStore('kelas', {
  state: () => ({
    list: [],
  }),
  actions: {
    loadData() {
      const stored = localStorage.getItem('data_kelas')
      if (stored) {
        this.list = JSON.parse(stored)
      } else {
        this.list = [
          { id: 1, nama: '12 RPL 1' },
          { id: 2, nama: '12 RPL 2' },
        ]
        this.saveData()
      }
    },
    saveData() {
      localStorage.setItem('data_kelas', JSON.stringify(this.list))
    },
  },
})

if (import.meta.hot) {
  import.meta.hot.accept(acceptHMRUpdate(useKelasStore, import.meta.hot))
}
