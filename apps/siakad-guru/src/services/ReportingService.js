import { TeacherDashboardProjection } from '../projection/TeacherDashboardProjection.js'

export class ReportingService {
  constructor({ reportingAdapter }) {
    this.reportingAdapter = reportingAdapter
  }

  async loadDashboard({ schoolId, academicYearId, semesterId, teacherId }) {
    const rawData = await this.reportingAdapter.fetchDashboardSummary({
      schoolId,
      academicYearId,
      semesterId,
      teacherId,
    })

    const summary = TeacherDashboardProjection.buildSummary({
      teachingData: {
        totalSessions: rawData.totalTeachingSessions,
        completedSessions: rawData.completedTeachingSessions,
      },
      attendanceData: {
        totalSessions: rawData.totalAttendanceSessions,
        submittedSessions: rawData.submittedAttendanceSessions,
      },
      assessmentData: {
        totalAssessments: rawData.totalAssessments,
        finalizedAssessments: rawData.finalizedAssessments,
      },
      journalData: {
        totalJournals: rawData.totalJournals,
        submittedJournals: rawData.submittedJournals,
      },
      progressData: {
        totalRecords: rawData.totalProgressRecords,
      },
      period: rawData.period,
    })

    return summary
  }
}
