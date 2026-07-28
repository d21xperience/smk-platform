// src/stores/auth-store.js
import { acceptHMRUpdate, defineStore } from 'pinia'
import { useCounterStore } from './example-store'

export const useUserStore = defineStore('user', {
  state: () => ({
    list: [],
  }),
  actions: {
    loadData() {
      const stored = localStorage.getItem('data_users')
      if (stored) {
        this.list = JSON.parse(stored)
      } else {
        // Data default
        this.list = [
          {
            id: 1,
            username: 'admin',
            password: 'admin123',
            nama: 'Administrator',
            role: 'admin',
            is_active: true,
          },
          {
            id: 2,
            username: 'operator',
            password: 'op123',
            nama: 'Operator',
            role: 'operator',
            is_active: true,
          },
          {
            id: 3,
            username: 'kepsek',
            password: 'kepsek123',
            nama: 'Kepala Sekolah',
            role: 'kepala_sekolah',
            is_active: true,
          },
          {
            id: 4,
            username: 'guru1',
            password: 'guru123',
            nama: 'Budi Santoso',
            role: 'guru',
            is_active: true,
          },
          {
            id: 5,
            username: 'siswa1',
            password: 'siswa123',
            nama: 'Ahmad Faizal',
            role: 'siswa',
            is_active: true,
          },
        ]
        this.saveData()
      }
    },
    saveData() {
      localStorage.setItem('data_users', JSON.stringify(this.list))
    },
    getUserByUsername(username) {
      return this.list.find((u) => u.username === username)
    },
    // ... method lain
  },
})
if (import.meta.hot) {
  import.meta.hot.accept(acceptHMRUpdate(useCounterStore, import.meta.hot))
}
