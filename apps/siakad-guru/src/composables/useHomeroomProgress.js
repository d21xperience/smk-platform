// FILE: src/composables/useHomeroomProgress.js
// STATUS: NEW
// STATUS IMPLEMENTASI: COMPLETE

import { computed } from 'vue'
import { useHomeroomProgressStore } from '../stores/homeroomProgressStore.js'
import { useContext } from './useContext.js'
import { useAuth } from './useAuth.js'

/**
 * Composable Facade untuk UI Homeroom Progress.
 * Composable bertanggung jawab:
 * - Mengambil context dari useContext()
 * - Mengambil auth dari useAuth()
 * - Meneruskan parameter eksplisit ke Store
 * Page yang memutuskan kapan memanggil loadProgressSummary().
 */
export function useHomeroomProgress() {
  const store = useHomeroomProgressStore()
  const { currentContext } = useContext()
  const { currentUser } = useAuth()

  const classId = computed(() => store.classId)
  const className = computed(() => store.className)
  const summary = computed(() => store.summary)
  const students = computed(() => store.students)
  const loading = computed(() => store.loading)
  const error = computed(() => store.error)
  const hasData = computed(() => store.hasData)
  const hasSummary = computed(() => store.hasSummary)

  /**
   * Memuat rekap nilai siswa.
   * Composable mengambil context dan auth, lalu meneruskan parameter eksplisit ke store.
   */
  async function loadProgressSummary() {
    if (!currentContext.value || !currentUser.value) return
    await store.loadProgressSummary({
      userId: currentUser.value.id,
      schoolId: currentContext.value.schoolId,
      academicYearId: currentContext.value.academicYearId,
      semesterId: currentContext.value.semesterId,
    })
  }

  /**
   * Mengosongkan data store.
   */
  function clearData() {
    store.clearData()
  }

  return {
    classId,
    className,
    summary,
    students,
    loading,
    error,
    hasData,
    hasSummary,
    loadProgressSummary,
    clearData,
  }
}
