// apps/siakad-tu/src/composables/landing/useLandingKesiswaan.js

import { ref, computed, onMounted } from 'vue'
// import { useRouter } from 'vue-router'
import { useLandingLayout } from './useLandingLayout'

/**
 * Composable untuk halaman LandingKesiswaan.
 *
 * Tanggung jawab:
 * - Mengelola data statistik kesiswaan (siap untuk integrasi backend)
 * - Menangani navigasi ke halaman detail dengan cek authentication
 * - Menyediakan data untuk tampilan (tidak mengubah UI)
 *
 * Catatan: Saat ini menggunakan data statis, ke depannya bisa diganti
 * dengan data dari store/backend tanpa mengubah template.
 */
export function useLandingKesiswaan() {
  // const router = useRouter()
  const { requestAccess } = useLandingLayout()

  // === STATE (Siap untuk integrasi backend) ===
  const isLoading = ref(false)
  const error = ref(null)

  // Data statistik - ke depannya akan diambil dari store/backend
  const statistics = ref({
    totalSiswaAktif: 1234,
    mutasiBulanIni: 45,
    pendaftarBaru: 89,
    lulusTahunIni: 156,
  })

  // Data fitur cards - ke depannya bisa diambil dari backend
  const features = ref([
    {
      title: 'Data Siswa',
      description: 'Kelola data siswa',
      targetRoute: 'student-list',
    },
    {
      title: 'Mutasi Siswa',
      description: 'Proses mutasi siswa',
      targetRoute: 'manajemen-mutasi-siswa',
    },
    {
      title: 'PSB',
      description: 'Penerimaan siswa baru',
      targetRoute: 'tambah-siswa',
    },
  ])

  // === COMPUTED ===
  const hasError = computed(() => error.value !== null)

  // === ACTIONS ===

  /**
   * Handle klik tombol "Lihat Detail" pada card fitur.
   * Menggunakan requestAccess dari useLandingLayout untuk cek authentication.
   *
   * @param {string} targetRoute - Nama route tujuan
   */
  function handleViewDetail(targetRoute) {
    requestAccess(targetRoute)
  }

  /**
   * Load data statistik dari backend.
   * Saat ini menggunakan data statis, ke depannya akan memanggil API.
   *
   * Contoh implementasi ke depannya:
   * ```
   * const result = await kesiswaanService.getStatistics()
   * if (result.success) {
   *   statistics.value = result.data
   * }
   * ```
   */
  async function loadStatistics() {
    isLoading.value = true
    error.value = null

    try {
      // TODO: Ganti dengan API call ke backend
      // const result = await kesiswaanService.getStatistics()
      // if (result.success) {
      //   statistics.value = result.data
      // } else {
      //   error.value = result.error
      // }

      // Simulasi delay network
      await new Promise((resolve) => setTimeout(resolve, 300))

      // Data statis (sudah di-set di state)
      console.log('[useLandingKesiswaan] Statistics loaded:', statistics.value)
    } catch (err) {
      error.value = { code: 'LOAD_ERROR', message: err.message }
      console.error('[useLandingKesiswaan] Error loading statistics:', err)
    } finally {
      isLoading.value = false
    }
  }

  /**
   * Load data fitur dari backend.
   * Saat ini menggunakan data statis.
   */
  async function loadFeatures() {
    // TODO: Ganti dengan API call ke backend jika needed
    console.log('[useLandingKesiswaan] Features loaded:', features.value)
  }

  // === LIFECYCLE ===
  onMounted(async () => {
    await Promise.all([loadStatistics(), loadFeatures()])
  })

  // === RETURN ===
  return {
    // State
    isLoading,
    error,
    statistics,
    features,

    // Computed
    hasError,

    // Actions
    handleViewDetail,
    loadStatistics,
    loadFeatures,
  }
}
