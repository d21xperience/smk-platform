// FILE: src/composables/useHomeroomPrintReport.js
// STATUS: NEW
// STATUS IMPLEMENTASI: COMPLETE

import { computed } from 'vue'
import { useHomeroomStudent } from './useHomeroomStudent.js'
import { useHomeroomProgress } from './useHomeroomProgress.js'
import { useHomeroomAttendance } from './useHomeroomAttendance.js'
import { useHomeroomJournal } from './useHomeroomJournal.js'
import { useHomeroomCommunication } from './useHomeroomCommunication.js'
import { useContext } from './useContext.js'
import { useAuth } from './useAuth.js'

/**
 * Composable Facade untuk Laporan Cetak Wali Kelas.
 * Mengumpulkan data dari seluruh store homeroom yang sudah ada.
 * Tidak membuat store/service/adapter baru — data sudah tersedia.
 * Page yang memutuskan kapan memanggil loadAllData() dan printReport().
 */
export function useHomeroomPrintReport() {
  const { students, className, loading: studentLoading, loadStudents } = useHomeroomStudent()

  const {
    summary: progressSummary,
    students: progressStudents,
    loading: progressLoading,
    loadProgressSummary,
  } = useHomeroomProgress()

  const {
    summary: attendanceSummary,
    students: attendanceStudents,
    months,
    loading: attendanceLoading,
    loadAttendanceSummary,
  } = useHomeroomAttendance()

  const {
    notes,
    summary: journalSummary,
    loading: journalLoading,
    loadNotes,
  } = useHomeroomJournal()

  const {
    communications,
    summary: communicationSummary,
    loading: communicationLoading,
    loadCommunications,
  } = useHomeroomCommunication()

  const { currentContext } = useContext()
  const { currentUser } = useAuth()

  const isAllLoading = computed(() => {
    return (
      studentLoading.value ||
      progressLoading.value ||
      attendanceLoading.value ||
      journalLoading.value ||
      communicationLoading.value
    )
  })

  const reportHeader = computed(() => {
    return {
      schoolName: 'SMP Negeri 1 Digital',
      academicYear: currentContext.value?.academicYearId || '-',
      semester: currentContext.value?.semesterId || '-',
      className: className.value || '-',
      teacherName: currentUser.value?.name || '-',
    }
  })

  /**
   * Memuat seluruh data yang dibutuhkan untuk laporan cetak.
   * Memanggil load function dari semua store homeroom secara paralel.
   */
  async function loadAllData() {
    if (!currentContext.value || !currentUser.value) return

    await Promise.all([
      loadStudents(),
      loadProgressSummary(),
      loadAttendanceSummary(),
      loadNotes(),
      loadCommunications(),
    ])
  }

  /**
   * Memicu dialog cetak browser.
   * CSS @media print di Page akan mengatur tampilan cetak.
   */
  function printReport() {
    window.print()
  }

  return {
    students,
    className,
    progressSummary,
    progressStudents,
    attendanceSummary,
    attendanceStudents,
    months,
    notes,
    journalSummary,
    communications,
    communicationSummary,
    currentContext,
    currentUser,
    reportHeader,
    isAllLoading,
    loadAllData,
    printReport,
  }
}
