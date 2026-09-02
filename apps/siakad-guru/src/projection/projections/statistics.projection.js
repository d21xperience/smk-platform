import { projectionManager } from '../ProjectionManager'

const INITIAL_STATE = {
  totalAttendanceToday: 0,
  totalPresentToday: 0,
  totalSickToday: 0,
  totalPermitToday: 0,
  totalAbsentToday: 0,
  totalLateToday: 0,
  attendanceRate: 0,
  totalSessionsToday: 0,
  totalSessionsThisWeek: 0,
  lastUpdated: null,
}

function updateStatisticsProjection(state, event) {
  const { type, payload } = event

  switch (type) {
    case 'ATTENDANCE_SUBMITTED': {
      const { data } = payload
      const { date, records } = data

      // Hanya hitung untuk hari ini
      const today = new Date().toISOString().split('T')[0]
      if (date !== today) return state

      const summary = records.reduce(
        (acc, rec) => {
          acc[rec.status] = (acc[rec.status] || 0) + 1
          acc.total++
          return acc
        },
        { present: 0, sick: 0, permit: 0, absent: 0, late: 0, total: 0 },
      )

      // Update total
      const newState = {
        ...state,
        totalAttendanceToday: state.totalAttendanceToday + summary.total,
        totalPresentToday: state.totalPresentToday + (summary.present || 0),
        totalSickToday: state.totalSickToday + (summary.sick || 0),
        totalPermitToday: state.totalPermitToday + (summary.permit || 0),
        totalAbsentToday: state.totalAbsentToday + (summary.absent || 0),
        totalLateToday: state.totalLateToday + (summary.late || 0),
        totalSessionsToday: state.totalSessionsToday + 1,
        lastUpdated: new Date().toISOString(),
      }
      // Update attendance rate
      const total = newState.totalAttendanceToday
      if (total > 0) {
        newState.attendanceRate =
          ((newState.totalPresentToday + newState.totalLateToday) / total) * 100
      }
      return newState
    }
    case 'SESSION_ENDED': {
      // Bisa digunakan untuk total sesi hari ini/minggu ini
      // Untuk sederhana kita abaikan atau bisa tambahakan hitungan
      return state
    }
    default:
      return state
  }
}

// Register statistics projection
export const statisticsProjection = projectionManager.register({
  name: 'statistics',
  initialState: INITIAL_STATE,
  updateFn: updateStatisticsProjection,
})
