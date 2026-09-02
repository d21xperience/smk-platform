import { acceptHMRUpdate, defineStore } from 'pinia'

export const useAbsensiStore = defineStore('absensi', {
  state: () => ({
    list: [], // { id, siswa_id, tanggal, status (hadir, izin, sakit, alpha), keterangan, jam_masuk, jam_keluar }
  }),

  actions: {
    // Muat data dari localStorage
    loadData() {
      const stored = localStorage.getItem('data_absensi')
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
      localStorage.setItem('data_absensi', JSON.stringify(this.list))
    },

    // Generate data dummy untuk bulan berjalan
    generateDummyData() {
      const dummy = []
      const today = new Date()
      const year = today.getFullYear()
      const month = today.getMonth() + 1
      const daysInMonth = new Date(year, month, 0).getDate()

      // Asumsikan ada 5 siswa dengan id 1..5
      for (let siswaId = 1; siswaId <= 5; siswaId++) {
        for (let day = 1; day <= daysInMonth; day++) {
          const date = `${year}-${String(month).padStart(2, '0')}-${String(day).padStart(2, '0')}`
          // Lewati hari Minggu (libur)
          const dayOfWeek = new Date(year, month - 1, day).getDay()
          if (dayOfWeek === 0) continue

          // Random status, 80% hadir, sisanya izin/sakit/alpha
          const rand = Math.random()
          let status = 'hadir'
          if (rand > 0.8) status = 'izin'
          if (rand > 0.9) status = 'sakit'
          if (rand > 0.95) status = 'alpha'

          dummy.push({
            id: Date.now() + siswaId * 1000 + day,
            siswa_id: siswaId,
            tanggal: date,
            status: status,
            keterangan: status === 'izin' ? 'Izin keluarga' : status === 'sakit' ? 'Sakit' : '',
            jam_masuk: '07:30',
            jam_keluar: '15:00',
          })
        }
      }
      return dummy
    },

    // Tambah absensi baru
    tambah(absen) {
      absen.id = Date.now()
      this.list.push(absen)
      this.saveData()
    },

    // Update absensi
    update(id, updated) {
      const idx = this.list.findIndex((a) => a.id === id)
      if (idx !== -1) {
        this.list[idx] = { ...this.list[idx], ...updated }
        this.saveData()
      }
    },

    // Hapus absensi
    hapus(id) {
      this.list = this.list.filter((a) => a.id !== id)
      this.saveData()
    },

    // Dapatkan absensi berdasarkan siswa dan tanggal
    getBySiswaDanTanggal(siswa_id, tanggal) {
      return this.list.find((a) => a.siswa_id === siswa_id && a.tanggal === tanggal)
    },

    // Dapatkan semua absensi siswa dalam rentang bulan
    getBySiswaBulan(siswa_id, tahun, bulan) {
      const monthStr = String(bulan).padStart(2, '0')
      return this.list.filter(
        (a) => a.siswa_id === siswa_id && a.tanggal.startsWith(`${tahun}-${monthStr}`),
      )
    },

    // Rekap absensi per kelas untuk bulan tertentu
    getRekapKelas(kelas_id, tahun, bulan) {
      // Memerlukan data siswa dari store lain, jadi panggil di komponen
      // Di sini hanya return list raw, filter manual di komponen
      const monthStr = String(bulan).padStart(2, '0')
      return this.list.filter((a) => a.tanggal.startsWith(`${tahun}-${monthStr}`))
    },

    // Hitung rata-rata kehadiran untuk semua siswa pada bulan ini
    getRataKehadiranBulanIni() {
      const now = new Date()
      const tahun = now.getFullYear()
      const bulan = now.getMonth() + 1
      const monthStr = String(bulan).padStart(2, '0')
      const absenBulanIni = this.list.filter((a) => a.tanggal.startsWith(`${tahun}-${monthStr}`))
      if (absenBulanIni.length === 0) return 0
      const totalHadir = absenBulanIni.filter((a) => a.status === 'hadir').length
      return Math.round((totalHadir / absenBulanIni.length) * 100)
    },

    // Hapus semua data (opsional)
    clearAll() {
      this.list = []
      this.saveData()
    },
  },
})

if (import.meta.hot) {
  import.meta.hot.accept(acceptHMRUpdate(useAbsensiStore, import.meta.hot))
}
