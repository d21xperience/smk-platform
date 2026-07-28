// src/projection/ClassAttendanceProjection.js
export function classAttendanceProjection(state, event) {
  if (event.type !== 'ATTENDANCE_SUBMITTED') return state
  const { className, summary } = event.payload
  return {
    ...state,
    classAttendance: {
      ...state.classAttendance,
      [className]: summary, // langsung replace dengan summary terbaru
    },
  }
}
