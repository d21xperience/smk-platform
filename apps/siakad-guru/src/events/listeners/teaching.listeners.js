import { EVENTS } from '../EventCatalog'
import { useTeachingStore } from '@/stores/teaching.store'

/**
 * Teaching Event Listeners
 * Bereaksi terhadap event dari domain lain
 */
export function registerTeachingListeners(eventBus) {
  // 1. Ketika absensi di-submit untuk sebuah sesi,
  //    tandai journal terkait bahwa absensi sudah di-submit.
  eventBus.on(
    EVENTS.ATTENDANCE_SUBMITTED,
    async (payload) => {
      const teachingStore = useTeachingStore()
      const { sessionId } = payload

      // Jika ada journal aktif untuk sessionId, tandai attendanceSubmitted = true
      if (teachingStore.currentJournal?.sessionId === sessionId) {
        try {
          // Panggil service untuk update attendanceSubmitted flag
          // (Ini akan menggunakan TeachingService.markAttendanceSubmitted)
          await teachingStore.markAttendanceSubmitted(sessionId)
          console.log(
            `[TeachingListener] Journal for session ${sessionId} marked as attendance submitted.`,
          )
        } catch (error) {
          console.error('[TeachingListener] Error marking attendance submitted:', error)
        }
      }
    },
    { priority: 10 }, // Priority tinggi karena ini adalah side-effect penting
  )

  // 2. Ketika sesi dimulai, log aktivitas (bisa juga trigger notifikasi)
  eventBus.on(
    EVENTS.SESSION_STARTED,
    (payload) => {
      // Bisa digunakan untuk UI feedback atau analytics
      console.log(`[TeachingListener] Session started: ${payload.id}`)
    },
    { priority: 0 },
  )

  // 3. Ketika jurnal di-finalisasi, trigger analitik atau update dashboard
  eventBus.on(
    EVENTS.JOURNAL_FINALIZED,
    (payload) => {
      console.log(`[TeachingListener] Journal finalized for session: ${payload.sessionId}`)
      // Bisa panggil API untuk update statistik
    },
    { priority: 0 },
  )
}
