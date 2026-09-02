export const TeachingContract = {
  fetchSchedulesByContext: {
    request: {
      schoolId: 'string',
      academicYearId: 'string',
      semesterId: 'string',
      teacherId: 'string',
    },
    response: [
      {
        id: 'string',
        classId: 'string',
        className: 'string',
        subjectId: 'string',
        subjectName: 'string',
        dayOfWeek: 'number', // 1-7
        startTime: 'string', // HH:mm
        endTime: 'string', // HH:mm
        teacherId: 'string',
        teacherName: 'string',
        academicYearId: 'string',
        semesterId: 'string',
        schoolId: 'string',
      },
    ],
  },

  fetchSessionsByDate: {
    request: {
      schoolId: 'string',
      academicYearId: 'string',
      semesterId: 'string',
      teacherId: 'string',
      date: 'string', // YYYY-MM-DD
    },
    response: [
      {
        id: 'string',
        scheduleId: 'string',
        classId: 'string',
        className: 'string',
        subjectId: 'string',
        subjectName: 'string',
        date: 'string',
        startTime: 'string',
        endTime: 'string',
        status: 'string', // scheduled, active, completed, cancelled
        teacherId: 'string',
        teacherName: 'string',
        teacherPresence: 'string', // present, absent, late, delegated
        notes: 'string',
        academicYearId: 'string',
        semesterId: 'string',
        schoolId: 'string',
      },
    ],
  },

  createSession: {
    request: {
      scheduleId: 'string',
      date: 'string',
      schoolId: 'string',
      academicYearId: 'string',
      semesterId: 'string',
      teacherId: 'string',
    },
    response: {
      id: 'string',
      scheduleId: 'string',
      classId: 'string',
      className: 'string',
      subjectId: 'string',
      subjectName: 'string',
      date: 'string',
      startTime: 'string',
      endTime: 'string',
      status: 'string',
      teacherId: 'string',
      teacherName: 'string',
      teacherPresence: 'string',
      notes: 'string',
      academicYearId: 'string',
      semesterId: 'string',
      schoolId: 'string',
    },
  },

  startSession: {
    request: {
      sessionId: 'string',
      teacherPresence: 'string',
      schoolId: 'string',
      academicYearId: 'string',
      semesterId: 'string',
    },
    response: {
      success: 'boolean',
      session: 'TeachingSession',
    },
  },

  endSession: {
    request: {
      sessionId: 'string',
      notes: 'string',
      schoolId: 'string',
      academicYearId: 'string',
      semesterId: 'string',
    },
    response: {
      success: 'boolean',
      session: 'TeachingSession',
    },
  },

  cancelSession: {
    request: {
      sessionId: 'string',
      reason: 'string',
      schoolId: 'string',
      academicYearId: 'string',
      semesterId: 'string',
    },
    response: {
      success: 'boolean',
      session: 'TeachingSession',
    },
  },
}
