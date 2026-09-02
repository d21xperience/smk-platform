import { defineStore, acceptHMRUpdate } from 'pinia'

export const useSiswaStore = defineStore('siswa', {
  state: () => ({
    list: [],
  }),
  actions: {
    loadData() {
      const stored = localStorage.getItem('data_siswa')
      if (stored) {
        this.list = JSON.parse(stored)
      } else {
        // Data dummy dengan orang tua dan riwayat mutasi
        this.list = [
          {
            id: 1,
            nis: '12345',
            nama: 'Ahmad Faizal',
            tempat_lahir: 'Bandung',
            tanggal_lahir: '2006-05-10',
            jenis_kelamin: 'Laki-laki',
            alamat: 'Jl. Merdeka No. 10',
            no_hp: '08123456789',
            kelas_id: 1,
            kelas_nama: '12 RPL 1',
            orang_tua: {
              ayah: 'Bapak Ahmad',
              ibu: 'Ibu Siti',
              pekerjaan_ayah: 'Swasta',
              pekerjaan_ibu: 'Ibu Rumah Tangga',
              no_hp_ortu: '08123456780',
            },
            riwayat_mutasi: [
              {
                id: 1,
                tanggal: '2023-07-01',
                asal_sekolah: 'SMP Negeri 1 Bandung',
                alasan: 'Pindah',
                dokumen: 'sk_mutasi_1.pdf',
              },
            ],
            status: 'aktif',
          },
          // ... siswa lain
        ]
        this.saveData()
      }
    },
    saveData() {
      localStorage.setItem('data_siswa', JSON.stringify(this.list))
    },
    tambah(siswa) {
      siswa.id = Date.now()
      siswa.riwayat_mutasi = siswa.riwayat_mutasi || []
      this.list.push(siswa)
      this.saveData()
    },
    update(id, updated) {
      const idx = this.list.findIndex((s) => s.id === id)
      if (idx !== -1) {
        this.list[idx] = { ...this.list[idx], ...updated }
        this.saveData()
      }
    },
    hapus(id) {
      this.list = this.list.filter((s) => s.id !== id)
      this.saveData()
    },
    tambahMutasi(id, mutasi) {
      const siswa = this.list.find((s) => s.id === id)
      if (siswa) {
        mutasi.id = Date.now()
        siswa.riwayat_mutasi.push(mutasi)
        this.saveData()
      }
    },
  },
})

if (import.meta.hot) {
  import.meta.hot.accept(acceptHMRUpdate(useSiswaStore, import.meta.hot))
}
