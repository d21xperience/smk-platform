import { acceptHMRUpdate, defineStore } from 'pinia'

export const usePendaftaranStore = defineStore('pendaftaran', {
  state: () => ({
    list: []
  }),
  actions: {
    loadData() {
      const stored = localStorage.getItem('data_pendaftaran')
      if (stored) this.list = JSON.parse(stored)
      else this.list = []
    },
    saveData() {
      localStorage.setItem('data_pendaftaran', JSON.stringify(this.list))
    },
    tambah(calon) {
      calon.id = Date.now()
      calon.tanggal_daftar = new Date().toISOString()  // otomatis
      calon.pembayaran = {
        status: 'belum',
        total_tagihan: 500000,
        jumlah_bayar: 0,
        sisa: 500000,
        tanggal_bayar_terakhir: null,
        riwayat: []
      }
      calon.berkas = calon.berkas || { ijazah: null, skl: null, kk: null, akte: null, foto: null }
      this.list.push(calon)
      this.saveData()
    },
    update(id, updated) {
      const idx = this.list.findIndex(c => c.id === id)
      if (idx !== -1) {
        this.list[idx] = { ...this.list[idx], ...updated }
        this.saveData()
      }
    },
    updateBerkas(id, jenisBerkas, fileBase64) {
      const idx = this.list.findIndex(c => c.id === id)
      if (idx !== -1) {
        if (!this.list[idx].berkas) this.list[idx].berkas = {}
        this.list[idx].berkas[jenisBerkas] = fileBase64
        this.saveData()
      }
    },
    // Pembayaran: tambah bayar (bisa sebagian)
    tambahPembayaran(id, jumlah, catatan = '') {
      const idx = this.list.findIndex(c => c.id === id)
      if (idx !== -1 && jumlah > 0) {
        const calon = this.list[idx]
        const bayar = calon.pembayaran
        const newJumlahBayar = bayar.jumlah_bayar + jumlah
        let sisa = bayar.total_tagihan - newJumlahBayar
        let status = 'belum'
        if (sisa <= 0) {
          status = 'lunas'
          sisa = 0
        } else if (newJumlahBayar > 0) {
          status = 'sebagian'
        }
        const riwayatBaru = {
          tanggal: new Date().toISOString(),
          jumlah: jumlah,
          sisa_setelah: sisa,
          catatan: catatan
        }
        calon.pembayaran = {
          ...bayar,
          jumlah_bayar: newJumlahBayar,
          sisa: sisa,
          status: status,
          tanggal_bayar_terakhir: new Date().toISOString(),
          riwayat: [...(bayar.riwayat || []), riwayatBaru]
        }
        this.saveData()
        return true
      }
      return false
    },
    // Set total tagihan (jika perlu diubah)
    setTotalTagihan(id, totalBaru) {
      const idx = this.list.findIndex(c => c.id === id)
      if (idx !== -1 && totalBaru > 0) {
        const calon = this.list[idx]
        const bayar = calon.pembayaran
        let sisa = totalBaru - bayar.jumlah_bayar
        let status = bayar.status
        if (sisa <= 0) {
          status = 'lunas'
          sisa = 0
        } else if (bayar.jumlah_bayar > 0) {
          status = 'sebagian'
        } else {
          status = 'belum'
        }
        calon.pembayaran.total_tagihan = totalBaru
        calon.pembayaran.sisa = sisa
        calon.pembayaran.status = status
        this.saveData()
      }
    },
    hapus(id) {
      this.list = this.list.filter(c => c.id !== id)
      this.saveData()
    },
    getKelengkapanBerkas(calon) {
      const required = ['ijazah', 'skl', 'kk', 'akte', 'foto']
      const uploaded = Object.keys(calon.berkas || {}).filter(key => calon.berkas[key])
      const missing = required.filter(r => !uploaded.includes(r))
      return { total: required.length, uploaded: uploaded.length, missing }
    }
  }
})

if (import.meta.hot) {
  import.meta.hot.accept(acceptHMRUpdate(usePendaftaranStore, import.meta.hot))
}
