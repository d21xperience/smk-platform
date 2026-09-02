export class ReportingMockAdapter {
  // eslint-disable-next-line no-unused-vars
  async fetchDashboardSummary({ schoolId, academicYearId, semesterId, teacherId }) {
    await new Promise(resolve => setTimeout(resolve, 600))

    return {
      totalTeachingSessions: 24,
      completedTeachingSessions: 18,
      totalAttendanceSessions: 18,
      submittedAttendanceSessions: 16,
      totalAssessments: 6,
      finalizedAssessments: 4,
      totalJournals: 18,
      submittedJournals: 15,
      totalProgressRecords: 12,
      period: {
        academicYearId,
        academicYearName: '2026/2027',
        semesterId,
        semesterName: semesterId.endsWith('1') ? 'Semester 1' : 'Semester 2'
      }
    }
  }
}
