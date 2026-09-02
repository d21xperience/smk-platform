import { acceptHMRUpdate, defineStore } from 'pinia'

export const useKonfigurasiStore = defineStore('konfigurasi', {
  state: () => ({
    // Pengaturan umum
    semester_aktif: 'Ganjil',
    tahun_ajaran_aktif: '2025/2026',
    nama_sekolah: 'SMK Pasundan',
    npsn: '12345678',

    // Pengaturan akademik
    kkm_default: 75,
    bobot_harian: 20,
    bobot_uts: 30,
    bobot_uas: 30,
    bobot_praktik: 15,
    bobot_portofolio: 5,

    // Pengaturan absensi
    toleransi_keterlambatan_menit: 15,
    batas_alpha_bulanan: 3,

    // Pengaturan ujian
    auto_lock_ujian: false,
    waktu_minimal_pengerjaan_menit: 30,
    maksimal_percobaan_ujian: 3,

    // Pengaturan sistem
    max_upload_size_mb: 2,
    enable_whatsapp_gateway: false,
    enable_email_gateway: true,
    waktu_backup_otomatis: '00:00',

    // Pengaturan tampilan
    logo_url: '',
    favicon_url: '',
    footer_text: '© SMK Pasundan - Sistem Informasi Akademik Terpadu',
  }),

  actions: {
    // Muat konfigurasi dari localStorage
    loadData() {
      const stored = localStorage.getItem('konfigurasi_sistem')
      if (stored) {
        const parsed = JSON.parse(stored)
        Object.assign(this, parsed)
      }
    },

    // Simpan konfigurasi ke localStorage
    saveData() {
      const toStore = {
        semester_aktif: this.semester_aktif,
        tahun_ajaran_aktif: this.tahun_ajaran_aktif,
        nama_sekolah: this.nama_sekolah,
        npsn: this.npsn,
        kkm_default: this.kkm_default,
        bobot_harian: this.bobot_harian,
        bobot_uts: this.bobot_uts,
        bobot_uas: this.bobot_uas,
        bobot_praktik: this.bobot_praktik,
        bobot_portofolio: this.bobot_portofolio,
        toleransi_keterlambatan_menit: this.toleransi_keterlambatan_menit,
        batas_alpha_bulanan: this.batas_alpha_bulanan,
        auto_lock_ujian: this.auto_lock_ujian,
        waktu_minimal_pengerjaan_menit: this.waktu_minimal_pengerjaan_menit,
        maksimal_percobaan_ujian: this.maksimal_percobaan_ujian,
        max_upload_size_mb: this.max_upload_size_mb,
        enable_whatsapp_gateway: this.enable_whatsapp_gateway,
        enable_email_gateway: this.enable_email_gateway,
        waktu_backup_otomatis: this.waktu_backup_otomatis,
        logo_url: this.logo_url,
        favicon_url: this.favicon_url,
        footer_text: this.footer_text,
      }
      localStorage.setItem('konfigurasi_sistem', JSON.stringify(toStore))
    },

    // Update konfigurasi (partial)
    update(updates) {
      Object.assign(this, updates)
      this.saveData()
    },

    // Reset ke default
    reset() {
      // Ambil state awal dari definisi store (gunakan factory)
      const defaultState = {
        semester_aktif: 'Ganjil',
        tahun_ajaran_aktif: '2025/2026',
        nama_sekolah: 'SMK Pasundan',
        npsn: '12345678',
        kkm_default: 75,
        bobot_harian: 20,
        bobot_uts: 30,
        bobot_uas: 30,
        bobot_praktik: 15,
        bobot_portofolio: 5,
        toleransi_keterlambatan_menit: 15,
        batas_alpha_bulanan: 3,
        auto_lock_ujian: false,
        waktu_minimal_pengerjaan_menit: 30,
        maksimal_percobaan_ujian: 3,
        max_upload_size_mb: 2,
        enable_whatsapp_gateway: false,
        enable_email_gateway: true,
        waktu_backup_otomatis: '00:00',
        logo_url: '',
        favicon_url: '',
        footer_text: '© SMK Pasundan - Sistem Informasi Akademik Terpadu',
      }
      Object.assign(this, defaultState)
      this.saveData()
    },

    // Validasi bobot total 100%
    isBobotValid() {
      const total =
        this.bobot_harian +
        this.bobot_uts +
        this.bobot_uas +
        this.bobot_praktik +
        this.bobot_portofolio
      return total === 100
    },

    // Dapatkan pesan error bobot jika tidak valid
    getBobotErrorMessage() {
      const total =
        this.bobot_harian +
        this.bobot_uts +
        this.bobot_uas +
        this.bobot_praktik +
        this.bobot_portofolio
      if (total !== 100) {
        return `Total bobot saat ini ${total}%, harus 100%`
      }
      return null
    },
  },
})

if (import.meta.hot) {
  import.meta.hot.accept(acceptHMRUpdate(useKonfigurasiStore, import.meta.hot))
}
