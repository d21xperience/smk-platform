import { defineStore, acceptHMRUpdate } from 'pinia'

export const useCutiStore = defineStore('cuti', {
  state: () => ({
    list: [],
  }),
  actions: {
    loadData() {
      const stored = localStorage.getItem('data_cuti')
      if (stored) {
        this.list = JSON.parse(stored)
      } else {
        // Data dummy
        this.list = [
          {
            id: 1,
            guru_id: 1,
            guru_nama: 'Budi Santoso',
            tanggal_mulai: '2025-07-10',
            tanggal_selesai: '2025-07-15',
            jenis: 'Cuti Tahunan',
            alasan: 'Liburan keluarga',
            status: 'pending',
            catatan_kepsek: null,
            diajukan_pada: '2025-06-01T08:00:00',
          },
          {
            id: 2,
            guru_id: 2,
            guru_nama: 'Dewi Lestari',
            tanggal_mulai: '2025-07-20',
            tanggal_selesai: '2025-07-22',
            jenis: 'SPT (Surat Perintah Tugas)',
            alasan: 'Diklat di luar kota',
            status: 'disetujui',
            catatan_kepsek: 'Diklat penting, disetujui',
            diajukan_pada: '2025-06-02T10:30:00',
          },
        ]
        this.saveData()
      }
    },
    saveData() {
      localStorage.setItem('data_cuti', JSON.stringify(this.list))
    },
    tambah(cuti) {
      cuti.id = Date.now()
      cuti.status = 'pending'
      cuti.catatan_kepsek = null
      cuti.diajukan_pada = new Date().toISOString()
      this.list.push(cuti)
      this.saveData()
    },
    updateStatus(id, status, catatan = null) {
      const idx = this.list.findIndex((c) => c.id === id)
      if (idx !== -1) {
        this.list[idx].status = status
        if (catatan) this.list[idx].catatan_kepsek = catatan
        this.saveData()
      }
    },
    hapus(id) {
      this.list = this.list.filter((c) => c.id !== id)
      this.saveData()
    },
  },
})

if (import.meta.hot) {
  import.meta.hot.accept(acceptHMRUpdate(useCutiStore, import.meta.hot))
}
