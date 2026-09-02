import { defineStore, acceptHMRUpdate } from 'pinia'

export const useArsipStore = defineStore('arsip', {
  state: () => ({
    list: [],
  }),
  actions: {
    loadData() {
      const stored = localStorage.getItem('data_arsip')
      if (stored) {
        this.list = JSON.parse(stored)
      } else {
        // Data dummy awal
        this.list = [
          {
            id: 1,
            guru_id: 1,
            guru_nama: 'Budi Santoso',
            jenis_dokumen: 'SK Pengangkatan',
            nomor_dokumen: 'SK/001/2020',
            tanggal_terbit: '2020-01-15',
            tanggal_kadaluarsa: null,
            file_name: 'sk_budi.pdf',
            file_base64: '',
            catatan: 'SK tetap',
          },
          {
            id: 2,
            guru_id: 1,
            guru_nama: 'Budi Santoso',
            jenis_dokumen: 'Sertifikat Serdik',
            nomor_dokumen: 'SERD/1234',
            tanggal_terbit: '2015-06-10',
            tanggal_kadaluarsa: '2025-06-10',
            file_name: 'serdik_budi.pdf',
            file_base64: '',
            catatan: 'Akan kadaluarsa',
          },
        ]
        this.saveData()
      }
    },
    saveData() {
      localStorage.setItem('data_arsip', JSON.stringify(this.list))
    },
    tambah(dokumen) {
      dokumen.id = Date.now()
      this.list.push(dokumen)
      this.saveData()
    },
    update(id, updated) {
      const idx = this.list.findIndex((d) => d.id === id)
      if (idx !== -1) {
        this.list[idx] = { ...this.list[idx], ...updated }
        this.saveData()
      }
    },
    hapus(id) {
      this.list = this.list.filter((d) => d.id !== id)
      this.saveData()
    },
    getExpiringDocs(daysThreshold = 30) {
      const today = new Date()
      const future = new Date()
      future.setDate(today.getDate() + daysThreshold)
      return this.list.filter((doc) => {
        if (!doc.tanggal_kadaluarsa) return false
        const expiry = new Date(doc.tanggal_kadaluarsa)
        return expiry >= today && expiry <= future
      })
    },
  },
})

if (import.meta.hot) {
  import.meta.hot.accept(acceptHMRUpdate(useArsipStore, import.meta.hot))
}
