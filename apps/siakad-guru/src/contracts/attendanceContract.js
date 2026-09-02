export const AttendanceContract = {
  fetchStudentsByClass: {
    request: {
      schoolId: 'string',
      academicYearId: 'string',
      semesterId: 'string',
      classId: 'string',
    },
    response: [
      {
        studentId: 'string',
        studentName: 'string',
        classId: 'string',
        className: 'string',
      },
    ],
  },

  loadAttendanceSession: {
    request: {
      schoolId: 'string',
      academicYearId: 'string',
      semesterId: 'string',
      teachingSessionId: 'string',
    },
    response: {
      id: 'string',
      teachingSessionId: 'string',
      classId: 'string',
      className: 'string',
      subjectId: 'string',
      subjectName: 'string',
      date: 'string',
      status: 'string', // draft, submitted
      records: [
        {
          studentId: 'string',
          studentName: 'string',
          status: 'string', // present, sick, permission, absent
          note: 'string',
        },
      ],
      academicYearId: 'string',
      semesterId: 'string',
      schoolId: 'string',
    },
  },

  saveAttendanceDraft: {
    request: {
      sessionId: 'string',
      records: [
        {
          studentId: 'string',
          status: 'string',
          note: 'string',
        },
      ],
      schoolId: 'string',
      academicYearId: 'string',
      semesterId: 'string',
    },
    response: {
      success: 'boolean',
      session: 'AttendanceSession',
    },
  },

  submitAttendance: {
    request: {
      sessionId: 'string',
      schoolId: 'string',
      academicYearId: 'string',
      semesterId: 'string',
    },
    response: {
      success: 'boolean',
      session: 'AttendanceSession',
    },
  },
}
