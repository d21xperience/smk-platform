/* eslint-disable no-unused-vars */
import { EVENTS } from '../EventCatalog'
import { studentProgressService } from '@/services/StudentProgressService'

export function registerStudentProgressListeners(eventBus) {
  // 1. Mendengarkan finalisasi assessment untuk update progress
  eventBus.on(EVENTS.ASSESSMENT_FINALIZED, async (payload) => {
    console.log('[StudentProgressListener] Assessment finalized, updating progress...')

    // Asumsikan payload mengandung data per siswa
    // Untuk batch, kita loop semua results di dalam assessment
    if (payload.results && Array.isArray(payload.results)) {
      for (const result of payload.results) {
        try {
          // Ambil data core subjects grades? Bisa dari parameter atau dari konteks assessment.
          // Untuk mock, kita sederhanakan.
          const syncData = {
            studentId: result.studentId,
            finalScore: result.finalScore,
            coreSubjectsGrades: result.coreSubjectsGrades || {},
          }
          await studentProgressService.syncProgressFromAssessment(syncData)
        } catch (error) {
          console.error(
            '[StudentProgressListener] Error updating progress for student:',
            result.studentId,
            error,
          )
        }
      }
    } else if (payload.studentId) {
      // Single student update
      await studentProgressService.syncProgressFromAssessment({
        studentId: payload.studentId,
        finalScore: payload.finalScore,
        coreSubjectsGrades: payload.coreSubjectsGrades || {},
      })
    }
  })

  // 2. (Opsional) Mendengarkan attendance submit untuk update stats kehadiran
  // Nanti bisa ditambahkan di Sprint C/D setelah API attendance siap
  eventBus.on(EVENTS.ATTENDANCE_SUBMITTED, (payload) => {
    console.log('[StudentProgressListener] Attendance submitted, updating attendance stats...')
    // Akan dipanggil ke service untuk update total present/sick/absent dll.
    // studentProgressService.syncAttendanceStats(payload);
  })

  // 3. Mendengarkan achievement/violation untuk update otomatis (sudah di Service)
}
