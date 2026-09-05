// apps/siakad-tu/src/composables/landing/useLandingKurikulum.js

import { ref, onMounted } from 'vue'
import { useLandingLayout } from './useLandingLayout'

/**
 * Composable untuk halaman LandingKurikulum.
 * Siap untuk integrasi backend di masa depan.
 */
export function useLandingKurikulum() {
  const { requestAccess } = useLandingLayout()

  const isLoading = ref(false)
  const error = ref(null)

  const features = ref([
    {
      title: 'Kurikulum Merdeka',
      description: 'Implementasi kurikulum merdeka',
      targetRoute: 'kurikulum-merdeka',
    },
    {
      title: 'Jadwal Pelajaran',
      description: 'Kelola jadwal pelajaran',
      targetRoute: 'jadwal-pelajaran',
    },
    {
      title: 'Penilaian',
      description: 'Sistem penilaian siswa',
      targetRoute: 'penilaian-nilai',
    },
  ])

  function handleViewDetail(targetRoute) {
    requestAccess(targetRoute)
  }

  async function loadFeatures() {
    isLoading.value = true
    try {
      // TODO: Ganti dengan API call ke backend
      await new Promise(resolve => setTimeout(resolve, 300))
      console.log('[useLandingKurikulum] Features loaded:', features.value)
    } catch (err) {
      error.value = { code: 'LOAD_ERROR', message: err.message }
    } finally {
      isLoading.value = false
    }
  }

  onMounted(() => {
    loadFeatures()
  })

  return {
    isLoading,
    error,
    features,
    handleViewDetail,
  }
}
