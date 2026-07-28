// src/projection/StudentSummaryProjection.js
export function studentSummaryProjection(state, event) {
  if (event.type !== 'ATTENDANCE_SUBMITTED') return state
  // Asumsikan payload berisi detail per siswa? Untuk ringkas, kita bisa menggunakan summary kelas.
  // Tapi lebih baik jika attendance submitted juga membawa array student status?
  // Untuk saat ini, kita hanya update count global.
  // Nanti di Sprint Assessment, kita akan punya data lebih detail.
  return {
    ...state,
    globalStats: {
      ...state.globalStats,
      totalAttendanceSubmitted: state.globalStats.totalAttendanceSubmitted + 1,
    },
  }
}
