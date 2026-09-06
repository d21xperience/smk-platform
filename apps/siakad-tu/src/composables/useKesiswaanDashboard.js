import { computed } from 'vue'
import { useKesiswaanStore } from '@/stores/kesiswaanStore'
import { useContextStore } from '@/stores/contextStore'
import { useLandingLayout } from '@/composables/landing/useLandingLayout'

export function useKesiswaanDashboard() {
  const store = useKesiswaanStore()
  const contextStore = useContextStore()
  const { requestAccess } = useLandingLayout()

  // ---- UI HELPERS (Business Logic UI) ----
  /**
   * Menghitung tinggi batang tren berdasarkan nilai maksimum
   */
  function trendHeight(jumlah) {
    const maxTren = Math.max(...(store.dashboardData?.trenPelanggaran?.map((b) => b.jumlah) || [1]))
    return Math.max(8, Math.round((jumlah / maxTren) * 120))
  }

  /**
   * Menghitung lebar bar distribusi poin berdasarkan nilai maksimum
   */
  function barWidth(total) {
    const maxPoin = Math.max(...(store.dashboardData?.poinPerTingkat?.map((p) => p.total) || [1]))
    return Math.round((total / maxPoin) * 100)
  }

  /**
   * Menentukan warna berdasarkan kategori kasus
   */
  function kategoriWarna(kategori) {
    if (kategori === 'Berat') return 'red-8'
    if (kategori === 'Sedang') return 'orange-8'
    return 'blue-grey-6'
  }

  // ---- INTERAKSI & ROUTING ----
  /**
   * Mapping tipe KPI ke target route
   * Sesuaikan dengan nama route yang terdaftar di router
   */
  const kpiRouteMap = {
    totalSiswaAktif: { name: 'student-list' },
    totalPelanggaranBulanIni: { name: 'pelanggaran-list' },
    totalBerprestasi: { name: 'prestasi-list' },
    totalEkskul: { name: 'ekskul-list' },
  }

  /**
   * Menangani klik pada KPI card
   * @param {string} kpiType - Tipe KPI (totalSiswaAktif, dll)
   */
  function handleKpiClick(kpiType) {
    const target = kpiRouteMap[kpiType]
    if (target) {
      // Gunakan requestAccess dari useLandingLayout untuk proteksi akses
      requestAccess(target)
    } else {
      console.warn(`Route untuk KPI "${kpiType}" belum didefinisikan.`)
    }
  }

  // ---- LOAD DATA ----
  /**
   * Memuat data dashboard menggunakan context dari contextStore
   * @returns {Promise<void>}
   */
  async function load() {
    const context = contextStore.current
    if (!context) {
      console.warn('Context belum tersedia, tidak bisa load dashboard.')
      return
    }
    await store.loadDashboard(context)
  }

  // ---- EXPOSE STATE (reaktif via computed) ----
  const data = computed(() => store.dashboardData)
  const loading = computed(() => store.isLoading)
  const error = computed(() => store.error)

  // ---- RETURN ----
  return {
    // State
    data,
    loading,
    error,
    // Actions
    load,
    // UI Helpers
    trendHeight,
    barWidth,
    kategoriWarna,
    // Interactions
    handleKpiClick,
  }
}
