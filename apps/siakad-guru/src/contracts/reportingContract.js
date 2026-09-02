export const ReportingContract = {
  fetchDashboardSummary: {
    request: {
      schoolId: 'string',
      academicYearId: 'string',
      semesterId: 'string',
      teacherId: 'string',
    },
    response: {
      totalTeachingSessions: 'number',
      completedTeachingSessions: 'number',
      totalAttendanceSessions: 'number',
      submittedAttendanceSessions: 'number',
      totalAssessments: 'number',
      finalizedAssessments: 'number',
      totalJournals: 'number',
      submittedJournals: 'number',
      totalProgressRecords: 'number',
      period: {
        academicYearId: 'string',
        academicYearName: 'string',
        semesterId: 'string',
        semesterName: 'string',
      },
    },
  },
}
