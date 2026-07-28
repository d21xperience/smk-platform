// // ============================================================
// // composables/useSchool.js
// // Shared logic, helpers, dan aksi yang dipakai banyak komponen
// // ============================================================
// import { computed } from 'vue'
// import { useQuasar } from 'quasar'
// import { SCHOOL_INFO } from '../data/schoolData'

// export function useSchool() {
//   const $q = useQuasar()

//   // ── Notifikasi helper ──────────────────────────────────────
//   function notify(message, type = 'positive', icon = 'check_circle') {
//     $q.notify({ message, type, icon, position: 'bottom-right', timeout: 3000 })
//   }

//   // ── Dark mode ──────────────────────────────────────────────
//   const isDark = computed(() => $q.dark.isActive)
//   function toggleDark() { $q.dark.toggle() }

//   // ── Buka URL eksternal ─────────────────────────────────────
//   function openUrl(url) {
//     if (url) window.open(url, '_blank', 'noopener,noreferrer')
//   }

//   // ── WhatsApp BKK ───────────────────────────────────────────
//   function hubungiBkk() {
//     openUrl(SCHOOL_INFO.whatsapp)
//   }

//   // ── Download brosur ────────────────────────────────────────
//   function downloadBrosur() {
//     openUrl(SCHOOL_INFO.brosurUrl)
//     notify('Memulai unduhan brosur digital…', 'info', 'file_download')
//   }

//   // ── Sosmed ────────────────────────────────────────────────
//   function bukaMedsos(platform) {
//     const urls = {
//       instagram: SCHOOL_INFO.instagram,
//       facebook:  SCHOOL_INFO.facebook,
//       youtube:   SCHOOL_INFO.youtube,
//       whatsapp:  SCHOOL_INFO.whatsapp,
//     }
//     openUrl(urls[platform])
//   }

//   // ── Smooth scroll ke section ──────────────────────────────
//   function scrollToSection(id) {
//     const el = document.getElementById(id)
//     if (el) el.scrollIntoView({ behavior: 'smooth', block: 'start' })
//   }

//   return {
//     isDark, toggleDark,
//     notify, openUrl, hubungiBkk, downloadBrosur, bukaMedsos, scrollToSection,
//   }
// }
import { reactive } from 'vue'
import { useApi } from './useApi'

export function useSchool() {
  const { api, request, loading } = useApi()

  // State reaktif untuk data sekolah
  const schoolInfo = reactive({
    name: '',
    address: '',
    phone: '',
    email: '',
    logoUrl: '',
  })

  // State reaktif untuk konfigurasi PPDB
  const ppdbConfig = reactive({
    tahunAjaran: '',
    startDate: null,
    endDate: null,
    isActive: false,
    targetUrl: '',
  })

  // --- ACTIONS ---

  /**
   * Mengambil informasi dasar sekolah
   * Endpoint Golang nanti: GET /api/v1/settings/school
   */
  const fetchSchoolInfo = async () => {
    const data = await request(() => api.get('/settings/school'), { silent: true })
    if (data) {
      Object.assign(schoolInfo, data)
    }
  }

  /**
   * Mengambil konfigurasi jadwal PPDB
   * Endpoint Golang nanti: GET /api/v1/settings/ppdb
   */
  const fetchPpdbConfig = async () => {
    const data = await request(() => api.get('/settings/ppdb'), { silent: true })
    if (data) {
      // Konversi string tanggal dari backend menjadi Date object untuk Vue
      ppdbConfig.tahunAjaran = data.tahunAjaran
      ppdbConfig.startDate = new Date(data.startDate)
      ppdbConfig.endDate = new Date(data.endDate)
      ppdbConfig.isActive = data.isActive
      ppdbConfig.targetUrl = data.targetUrl
    }
  }

  return {
    schoolInfo,
    ppdbConfig,
    loading,
    fetchSchoolInfo,
    fetchPpdbConfig,
  }
}
