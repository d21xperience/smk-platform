// @/pages/siakad/teaching/useAttendancePage.js

import { ref, computed, watch } from 'vue'
import { useOperationalContext } from '@/composables/context/useOperationalContext'
import { usePermission } from '@/composables/permission/usePermission'
// Nanti akan ada composable domain
// import { useTeaching } from '@/composables/domain/useTeaching'
// import { useAttendance } from '@/composables/domain/useAttendance'

export function useAttendancePage() {
  // ------------------------------------------------------------
  // 1. REACTIVE STATE
  // ------------------------------------------------------------
  const pageState = ref('initializing') // 'initializing' | 'loading' | 'ready' | 'error'
  const errorMessage = ref('')
  const retryCount = ref(0)

  // State untuk business context (akan diisi nanti)
  const teachingSession = ref(null)
  const attendanceList = ref([])

  // ------------------------------------------------------------
  // 2. COMPUTED
  // ------------------------------------------------------------
  const isLoading = computed(
    () => pageState.value === 'loading' || pageState.value === 'initializing',
  )
  const isReady = computed(() => pageState.value === 'ready')
  const hasError = computed(() => pageState.value === 'error')

  // ------------------------------------------------------------
  // 3. COMPOSABLE - Operational Context & Permission
  // ------------------------------------------------------------
  const { context, loadContext, isContextReady } = useOperationalContext()
  const { checkPagePermission } = usePermission()

  // ------------------------------------------------------------
  // 4. LIFECYCLE - initialize() adalah SATU pintu masuk
  // ------------------------------------------------------------
  async function initialize() {
    if (pageState.value === 'loading' || pageState.value === 'initializing') return
    pageState.value = 'initializing'
    errorMessage.value = ''

    try {
      // STEP A: Load Operational Context
      await loadContext()
      if (!isContextReady.value) {
        throw new Error('Operational context gagal dimuat')
      }

      // STEP B: Permission Check
      const allowed = await checkPagePermission('attendance', context.value)
      if (!allowed) {
        throw new Error('Anda tidak memiliki akses ke halaman ini')
      }

      // STEP C: Load Business Context (akan diimplementasikan di Sprint B)
      await loadBusinessContext()

      // STEP D: Siap
      pageState.value = 'ready'
    } catch (err) {
      pageState.value = 'error'
      errorMessage.value = err.message || 'Terjadi kesalahan saat memuat halaman'
    }
  }

  // ------------------------------------------------------------
  // 5. METHODS - Business Context Loader (Sprint B nanti)
  // ------------------------------------------------------------
  async function loadBusinessContext() {
    // Di Sprint B, kita akan panggil useTeaching() dan useAttendance()
    // Contoh:
    // const { fetchTeachingSession } = useTeaching()
    // teachingSession.value = await fetchTeachingSession(context.value.currentScheduleId)
    // const { fetchAttendance } = useAttendance()
    // attendanceList.value = await fetchAttendance(teachingSession.value.id)

    // Untuk Sprint A, kita hanya mensimulasikan
    await new Promise((resolve) => setTimeout(resolve, 500))
    teachingSession.value = { id: 1, name: 'Matematika Kelas X-A' }
    attendanceList.value = []
  }

  // ------------------------------------------------------------
  // 6. WATCH (opsional)
  // ------------------------------------------------------------
  // Misal kita ingin reload ketika retryCount berubah
  watch(retryCount, () => {
    if (pageState.value === 'error') {
      initialize()
    }
  })

  // ------------------------------------------------------------
  // 7. EXPOSE - apa yang diberikan ke template
  // ------------------------------------------------------------
  return {
    // State
    pageState,
    errorMessage,
    teachingSession,
    attendanceList,
    // Computed
    isLoading,
    isReady,
    hasError,
    // Methods
    initialize,
    retry: () => {
      retryCount.value++
    },
  }
}
