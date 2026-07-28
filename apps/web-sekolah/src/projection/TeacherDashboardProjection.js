// src/projection/TeacherDashboardProjection.js
export function teacherDashboardProjection(state, event) {
  const { teacherId } = event.payload
  if (!teacherId) return state
  const teacherKey = teacherId.toString()
  const current = state.teacherDashboard[teacherKey] || {
    totalSessions: 0,
    onTime: 0,
    late: 0,
    totalMinutes: 0,
  }
  switch (event.type) {
    case 'TEACHING_SESSION_COMPLETED':
      return {
        ...state,
        teacherDashboard: {
          ...state.teacherDashboard,
          [teacherKey]: {
            ...current,
            totalSessions: current.totalSessions + 1,
            totalMinutes: current.totalMinutes + (event.payload.durationMinutes || 0),
          },
        },
      }
    case 'TEACHER_LATE':
      return {
        ...state,
        teacherDashboard: {
          ...state.teacherDashboard,
          [teacherKey]: {
            ...current,
            late: current.late + 1,
          },
        },
      }
    default:
      return state
  }
}
