// FILE: src/composables/useHomeroomAttendance.js
// STATUS: NEW
// STATUS IMPLEMENTASI: COMPLETE

import { computed } from 'vue'
import { useHomeroomAttendanceStore } from '../stores/homeroomAttendanceStore.js'
import { useContext } from './useContext.js'
import { useAuth } from './useAuth.js'

/**
 * Composable Facade untuk UI Homeroom Attendance.
 * Composable bertanggung jawab:
 * - Mengambil context dari useContext()
 * - Mengambil auth dari useAuth()
 * - Meneruskan parameter eksplisit ke Store
 * Page yang memutuskan kapan memanggil loadAttendanceSummary().
 */
export function useHomeroomAttendance() {
  const store = useHomeroomAttendanceStore()
  const { currentContext } = useContext()
  const { currentUser } = useAuth()

  const classId = computed(() => store.classId)
  const className = computed(() => store.className)
  const months = computed(() => store.months)
  const summary = computed(() => store.summary)
  const students = computed(() => store.students)
  const loading = computed(() => store.loading)
  const error = computed(() => store.error)
  const hasData = computed(() => store.hasData)
  const hasSummary = computed(() => store.hasSummary)
  const hasMonths = computed(() => store.hasMonths)

  /**
   * Memuat rekap kehadiran siswa.
   * Composable mengambil context dan auth, lalu meneruskan parameter eksplisit ke store.
   */
  async function loadAttendanceSummary() {
    if (!currentContext.value || !currentUser.value) return
    await store.loadAttendanceSummary({
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
    months,
    summary,
    students,
    loading,
    error,
    hasData,
    hasSummary,
    hasMonths,
    loadAttendanceSummary,
    clearData,
  }
}
