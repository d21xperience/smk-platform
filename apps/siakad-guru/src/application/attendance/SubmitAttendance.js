// application/attendance/SubmitAttendance.js
import { attendanceMock } from '@/adapters/mock/handlers/attendance.mock'
import { isLate, getAttendanceStatus } from '@/domain/attendance/engine/AttendanceEngine.js'
import { useContextStore } from '@/stores/context.store.js'

export const SubmitAttendanceUseCase = {
  execute: async ({ studentId, classId, status, checkInTime, checkOutTime, notes }) => {
    const contextStore = useContextStore()
    const { semesterId, userId: teacherId, schoolId } = contextStore

    const isLateFlag = checkInTime ? isLate(checkInTime, '07:15') : false
    const finalStatus = status || getAttendanceStatus(status !== 'ALPHA', isLateFlag)

    const payload = {
      studentId,
      classId,
      date: new Date().toISOString().split('T')[0],
      status: finalStatus,
      checkInTime,
      checkOutTime,
      notes,
      teacherId,
      semesterId,
      schoolId,
      academicYearId: contextStore.academicYearId, // penting
    }

    const result = await attendanceMock.submit(payload)
    return result
  },
}
