import { EVENTS } from '../EventCatalog'
import { useAttendanceStore } from '@/stores/attendanceStore'

export function registerAttendanceListeners(eventBus) {
  // Ketika sesi mengajar dimulai, otomatis load absensi untuk sesi tersebut
  // (Opsional: integrasi dengan Teaching)
  eventBus.on(
    EVENTS.SESSION_STARTED,
    async (payload) => {
      const attendanceStore = useAttendanceStore()
      const { classId, subjectId, date, id: sessionId } = payload

      // Jika ada classId dan date, load absensi
      if (classId && subjectId && date) {
        try {
          await attendanceStore.loadOrCreateSession({
            classId,
            subjectId,
            date,
            sessionId, // Bisa digunakan untuk mapping
          })
          console.log(`[AttendanceListener] Auto-loaded attendance for session ${sessionId}`)
        } catch (error) {
          console.warn(`[AttendanceListener] Failed to auto-load attendance: ${error.message}`)
        }
      }
    },
    { priority: 5 }, // Priority sedang
  )

  // Ketika context berubah, clear attendance session (agar tidak terjadi inkonsistensi)
  eventBus.on(
    EVENTS.CONTEXT_CHANGED,
    () => {
      const attendanceStore = useAttendanceStore()
      attendanceStore.clearSession()
      console.log('[AttendanceListener] Attendance session cleared due to context change.')
    },
    { priority: 10 },
  )
}
