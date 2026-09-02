/* eslint-disable no-unused-vars */
import { projectionManager } from '../ProjectionManager'

const INITIAL_STATE = {
  byClass: {}, // classId -> { date: { summary: { present, sick, permit, absent, late, total } } }
  byTeacher: {}, // teacherId -> { date: { ... } }
  byStudent: {}, // studentId -> { date: { sessionId, status } }
  lastUpdated: null,
}

function updateAttendanceProjection(state, event) {
  const { type, payload } = event

  switch (type) {
    case 'ATTENDANCE_SUBMITTED': {
      const { data } = payload
      const { classId, teacherId, date, records } = data

      // Update byClass
      if (!state.byClass[classId]) state.byClass[classId] = {}
      if (!state.byClass[classId][date])
        state.byClass[classId][date] = {
          summary: { present: 0, sick: 0, permit: 0, absent: 0, late: 0, total: 0 },
        }
      const classSummary = state.byClass[classId][date].summary
      // Reset atau tambah? Karena submit menggantikan, kita rebuild dari records.
      // Lebih baik rebuild total.
      const newSummary = records.reduce(
        (acc, rec) => {
          acc[rec.status] = (acc[rec.status] || 0) + 1
          acc.total++
          return acc
        },
        { present: 0, sick: 0, permit: 0, absent: 0, late: 0, total: 0 },
      )
      state.byClass[classId][date].summary = newSummary

      // Update byTeacher
      if (!state.byTeacher[teacherId]) state.byTeacher[teacherId] = {}
      if (!state.byTeacher[teacherId][date])
        state.byTeacher[teacherId][date] = {
          summary: { present: 0, sick: 0, permit: 0, absent: 0, late: 0, total: 0 },
        }
      state.byTeacher[teacherId][date].summary = newSummary

      // Update byStudent (untuk riwayat siswa)
      records.forEach((rec) => {
        if (!state.byStudent[rec.studentId]) state.byStudent[rec.studentId] = {}
        state.byStudent[rec.studentId][date] = { sessionId: data.id, status: rec.status }
      })

      state.lastUpdated = new Date().toISOString()
      return { ...state }
    }
    case 'ATTENDANCE_RECORD_UPDATED': {
      // Untuk update record sebelum submit, projection belum di-update (karena belum final)
      // Bisa diabaikan atau update sementara. Kita pilih diabaikan untuk menjaga konsistensi eventual.
      return state
    }
    default:
      return state
  }
}

// Register projection
export const attendanceProjection = projectionManager.register({
  name: 'attendance',
  initialState: INITIAL_STATE,
  updateFn: updateAttendanceProjection,
})
