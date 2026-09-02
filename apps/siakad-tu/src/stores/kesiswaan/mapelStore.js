import { defineStore, acceptHMRUpdate } from 'pinia'

export const useMapelStore = defineStore('mapel', {
  state: () => ({
    list: [],
  }),
  actions: {
    loadData() {
      const stored = localStorage.getItem('data_mapel')
      if (stored) {
        this.list = JSON.parse(stored)
      } else {
        this.list = [
          { id: 1, nama: 'Matematika', jenis: 'umum' },
          { id: 2, nama: 'Pemrograman Web', jenis: 'kejuruan' },
          { id: 3, nama: 'Basis Data', jenis: 'kejuruan' },
        ]
        this.saveData()
      }
    },
    saveData() {
      localStorage.setItem('data_mapel', JSON.stringify(this.list))
    },
  },
})

if (import.meta.hot) {
  import.meta.hot.accept(acceptHMRUpdate(useMapelStore, import.meta.hot))
}
