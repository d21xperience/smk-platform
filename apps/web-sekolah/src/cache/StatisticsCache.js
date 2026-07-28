// src/cache/StatisticsCache.js
export const createInitialCache = () => ({
  classAttendance: {}, // key: className, value: summary
  teacherDashboard: {}, // key: teacherId, value: dashboard data
  studentSummary: {}, // key: studentId, value: summary
  homeroomSummary: {}, // key: className, value: homeroom view
  globalStats: {
    totalSessions: 0,
    totalAttendanceSubmitted: 0,
    totalJournals: 0,
  },
})
