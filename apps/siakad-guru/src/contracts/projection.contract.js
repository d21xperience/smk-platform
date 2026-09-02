/**
 * FRONTEND DRIVEN CONTRACT
 * Projection Contract untuk SQLite Read Model (Backend)
 * Backend harus menyediakan data dalam format ini.
 */
export const ProjectionContract = {
  attendance: {
    byClass: {
      // classId -> { date: { summary: { present, sick, permit, absent, late, total } } }
    },
    byTeacher: {
      // teacherId -> { date: { summary: ... } }
    },
    byStudent: {
      // studentId -> { date: { sessionId, status } }
    },
    lastUpdated: 'ISO datetime',
  },
  teaching: {
    byTeacher: {
      // teacherId -> { date: [{ sessionId, classId, subjectId, startTime, endTime, status }] }
    },
    byClass: {
      // classId -> { date: [{ ... }] }
    },
    journals: {
      // sessionId -> TeachingJournal data
    },
    lastUpdated: 'ISO datetime',
  },
  statistics: {
    totalAttendanceToday: 'number',
    totalPresentToday: 'number',
    totalSickToday: 'number',
    totalPermitToday: 'number',
    totalAbsentToday: 'number',
    totalLateToday: 'number',
    attendanceRate: 'number (0-100)',
    totalSessionsToday: 'number',
    totalSessionsThisWeek: 'number',
    lastUpdated: 'ISO datetime',
  },
}
