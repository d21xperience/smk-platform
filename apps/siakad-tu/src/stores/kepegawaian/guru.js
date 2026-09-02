import { defineStore, acceptHMRUpdate } from 'pinia'

export const useGuruStore = defineStore('guru', {
  state: () => ({
    list: [],
  }),
  actions: {
    loadData() {
      const stored = localStorage.getItem('data_guru')
      if (stored) {
        this.list = JSON.parse(stored)
      } else {
        // Data dummy awal
        this.list = [
          {
            id: 1,
            nama: 'Budi Santoso',
            nuptk: '123456789012345',
            jabatan: 'Guru Matematika',
            masa_tugas: 10,
            status: 'aktif',
          },
          {
            id: 2,
            nama: 'Dewi Lestari',
            nuptk: '987654321098765',
            jabatan: 'Guru Bahasa Inggris',
            masa_tugas: 5,
            status: 'aktif',
          },
        ]
        this.saveData()
      }
    },
    saveData() {
      localStorage.setItem('data_guru', JSON.stringify(this.list))
    },
    tambah(guru) {
      guru.id = Date.now()
      this.list.push(guru)
      this.saveData()
    },
    update(id, updated) {
      const index = this.list.findIndex((g) => g.id === id)
      if (index !== -1) {
        this.list[index] = { ...this.list[index], ...updated }
        this.saveData()
      }
    },
    hapus(id) {
      this.list = this.list.filter((g) => g.id !== id)
      this.saveData()
    },
  },
})

if (import.meta.hot) {
  import.meta.hot.accept(acceptHMRUpdate(useGuruStore, import.meta.hot))
}
