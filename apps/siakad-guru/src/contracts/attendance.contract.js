/**
 * FRONTEND DRIVEN CONTRACT
 * Attendance API Contract
 */
export const AttendanceContract = {
  loadOrCreate: {
    request: {
      classId: 'string',
      subjectId: 'string',
      date: 'string (YYYY-MM-DD)',
      teacherId: 'string',
      academicYearId: 'string',
      semesterId: 'string',
    },
    response: {
      success: 'boolean',
      data: 'AttendanceSession',
      errors: 'array',
    },
  },
  upsertRecord: {
    request: {
      studentId: 'string',
      status: 'string (present|sick|permit|absent|late)',
      note: 'string (optional)',
    },
    response: {
      success: 'boolean',
      data: 'AttendanceRecord',
      errors: 'array',
    },
  },
  submit: {
    request: {},
    response: {
      success: 'boolean',
      data: 'AttendanceSession',
      errors: 'array',
    },
  },
}
