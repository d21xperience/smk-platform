import { defineStore, acceptHMRUpdate } from 'pinia'

export const useSuratStore = defineStore('surat', {
  state: () => ({
    list: [], // Array of surat objects
  }),

  actions: {
    // Muat data dari localStorage
    loadData() {
      const stored = localStorage.getItem('data_surat')
      if (stored) {
        this.list = JSON.parse(stored)
      } else {
        // Data dummy untuk demo
        this.list = this.generateDummyData()
        this.saveData()
      }
    },

    // Simpan ke localStorage
    saveData() {
      localStorage.setItem('data_surat', JSON.stringify(this.list))
    },

    // Generate data dummy untuk demo
    generateDummyData() {
      // const now = new Date()
      return [
        {
          id: 1,
          nomor_surat: '001/Kel/SMK.Pasundan/2025',
          template_id: 1,
          judul: 'Surat Tugas Diklat',
          konten:
            '<p>Memberikan tugas kepada <strong>Budi Santoso</strong> untuk kegiatan Diklat.</p>',
          data_variabel: { nama_guru: 'Budi Santoso', kegiatan: 'Diklat' },
          jenis: 'keluar',
          status: 'selesai',
          prioritas: 'sedang',
          tanggal_dibuat: '2025-01-10T08:00:00',
          deadline: '2025-01-20',
          catatan_verifikasi: 'Disetujui',
          tanda_tangan: 'ttd_kepsek.png',
          created_by: 'admin',
          approved_by: 'kepsek',
        },
        {
          id: 2,
          nomor_surat: '002/Masuk/SMK.Pasundan/2025',
          template_id: null,
          judul: 'Surat Permohonan PKL',
          konten: '<p>Memohon tempat PKL untuk siswa kelas XII.</p>',
          data_variabel: {},
          jenis: 'masuk',
          status: 'menunggu_verifikasi',
          prioritas: 'tinggi',
          tanggal_dibuat: '2025-02-01T10:30:00',
          deadline: '2025-02-15',
          catatan_verifikasi: null,
          tanda_tangan: null,
          created_by: 'operator',
          approved_by: null,
        },
      ]
    },

    // Generate nomor surat otomatis
    generateNomor(jenis, tahun = new Date().getFullYear()) {
      const prefix = jenis === 'keluar' ? 'Kel' : 'Masuk'
      // Hitung jumlah surat dengan jenis dan tahun yang sama
      const count =
        this.list.filter(
          (s) =>
            s.nomor_surat &&
            s.nomor_surat.includes(`/${prefix}/`) &&
            s.nomor_surat.includes(tahun.toString()),
        ).length + 1
      const noUrut = String(count).padStart(3, '0')
      return `${noUrut}/${prefix}/SMK.Pasundan/${tahun}`
    },

    // Tambah surat baru (status default = draft)
    tambah(surat) {
      surat.id = Date.now()
      surat.status = surat.status || 'draft'
      surat.tanggal_dibuat = new Date().toISOString()
      this.list.push(surat)
      this.saveData()
    },

    // Update surat
    update(id, updated) {
      const idx = this.list.findIndex((s) => s.id === id)
      if (idx !== -1) {
        this.list[idx] = { ...this.list[idx], ...updated }
        this.saveData()
      }
    },

    // Submit ke verifikasi (guru → kaprodi/wakasek)
    submitVerifikasi(id) {
      const idx = this.list.findIndex((s) => s.id === id)
      if (idx !== -1 && this.list[idx].status === 'draft') {
        this.list[idx].status = 'menunggu_verifikasi'
        this.saveData()
        return true
      }
      return false
    },

    // Setujui surat
    approve(id, catatan = '', tanda_tangan = null, approvedBy = 'verifikator') {
      const idx = this.list.findIndex((s) => s.id === id)
      if (idx !== -1 && this.list[idx].status === 'menunggu_verifikasi') {
        this.list[idx].status = 'disetujui'
        this.list[idx].catatan_verifikasi = catatan
        this.list[idx].tanda_tangan = tanda_tangan
        this.list[idx].approved_by = approvedBy
        this.saveData()
        return true
      }
      return false
    },

    // Tolak surat
    reject(id, catatan) {
      const idx = this.list.findIndex((s) => s.id === id)
      if (idx !== -1 && this.list[idx].status === 'menunggu_verifikasi') {
        this.list[idx].status = 'ditolak'
        this.list[idx].catatan_verifikasi = catatan
        this.saveData()
        return true
      }
      return false
    },

    // Finalisasi (tanda tangan digital/selesai)
    finalize(id, tanda_tangan = null) {
      const idx = this.list.findIndex((s) => s.id === id)
      if (idx !== -1 && this.list[idx].status === 'disetujui') {
        this.list[idx].status = 'selesai'
        if (tanda_tangan) this.list[idx].tanda_tangan = tanda_tangan
        this.saveData()
        return true
      }
      return false
    },

    // Hapus surat
    hapus(id) {
      this.list = this.list.filter((s) => s.id !== id)
      this.saveData()
    },

    // Filter berdasarkan status
    getByStatus(status) {
      return this.list.filter((s) => s.status === status)
    },

    // Filter berdasarkan jenis (keluar/masuk)
    getByJenis(jenis) {
      return this.list.filter((s) => s.jenis === jenis)
    },

    // Dapatkan surat terakhir berdasarkan jenis (untuk nomor surat otomatis)
    getLastSuratByJenis(jenis, tahun) {
      const filtered = this.list.filter(
        (s) => s.jenis === jenis && s.nomor_surat && s.nomor_surat.includes(tahun.toString()),
      )
      filtered.sort((a, b) => {
        const noA = parseInt(a.nomor_surat.split('/')[0]) || 0
        const noB = parseInt(b.nomor_surat.split('/')[0]) || 0
        return noB - noA
      })
      return filtered[0] || null
    },

    // Hapus semua data (opsional)
    clearAll() {
      this.list = []
      this.saveData()
    },
  },
})

if (import.meta.hot) {
  import.meta.hot.accept(acceptHMRUpdate(useSuratStore, import.meta.hot))
}
