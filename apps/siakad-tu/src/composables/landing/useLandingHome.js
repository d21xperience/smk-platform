// apps/siakad-tu/src/composables/landing/useLandingHome.js

import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'

/**
 * Composable untuk halaman landing home.
 *
 * Tanggung jawab:
 * - Mengelola data statis (features, menu items)
 * - Menyediakan helper functions untuk navigasi
 * - Mengkoordinasi interaksi antara Page dan Store (jika ada)
 */
export function useLandingHome() {
  const router = useRouter()

  // === STATIC DATA ===
  const features = ref([
    {
      icon: 'school',
      title: 'Manajemen Kesiswaan',
      description: 'Kelola data siswa, mutasi, dan kelulusan dengan mudah',
      color: 'primary',
    },
    {
      icon: 'account_balance',
      title: 'Keuangan Sekolah',
      description: 'Monitoring tagihan, pembayaran, dan laporan keuangan',
      color: 'secondary',
    },
    {
      icon: 'event_available',
      title: 'Absensi Digital',
      description: 'Catat kehadiran siswa secara real-time',
      color: 'accent',
    },
    {
      icon: 'assessment',
      title: 'Penilaian & Rapor',
      description: 'Input nilai dan generate rapor otomatis',
      color: 'positive',
    },
    {
      icon: 'mail',
      title: 'Persuratan',
      description: 'Kelola surat masuk, keluar, dan disposisi',
      color: 'info',
    },
    {
      icon: 'inventory',
      title: 'Sarpras',
      description: 'Manajemen sarana dan prasarana sekolah',
      color: 'warning',
    },
  ])

  const modules = ref([
    { name: 'Kesiswaan', path: '/kesiswaan', icon: 'school' },
    { name: 'Keuangan', path: '/keuangan', icon: 'account_balance' },
    { name: 'Absensi', path: '/absensi', icon: 'event_available' },
    { name: 'Penilaian', path: '/penilaian', icon: 'assessment' },
    { name: 'Persuratan', path: '/surat', icon: 'mail' },
    { name: 'PSB', path: '/psb', icon: 'person_add' },
  ])

  // === COMPUTED ===
  const isMaintenanceMode = computed(() => {
    // Bisa diambil dari environment variable atau API
    return import.meta.env?.VITE_MAINTENANCE_MODE === 'true'
  })

  // === ACTIONS ===
  function navigateTo(path) {
    router.push(path)
  }

  function goToLogin() {
    router.push('/login')
  }

  function goToRegister() {
    router.push('/register')
  }

  function getModuleIcon(moduleName) {
    const module = modules.value.find(m => m.name === moduleName)
    return module?.icon || 'dashboard'
  }

  // === RETURN ===
  return {
    // Data
    features,
    modules,

    // Computed
    isMaintenanceMode,

    // Actions
    navigateTo,
    goToLogin,
    goToRegister,
    getModuleIcon,
  }
}
