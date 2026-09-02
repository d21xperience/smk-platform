import { DashboardSummary } from '../domain/reporting/models/DashboardSummary.js'

export class TeacherDashboardProjection {
  static buildSummary({
    teachingData,
    attendanceData,
    assessmentData,
    journalData,
    progressData,
    period,
  }) {
    return new DashboardSummary({
      totalTeachingSessions: teachingData?.totalSessions || 0,
      completedTeachingSessions: teachingData?.completedSessions || 0,
      totalAttendanceSessions: attendanceData?.totalSessions || 0,
      submittedAttendanceSessions: attendanceData?.submittedSessions || 0,
      totalAssessments: assessmentData?.totalAssessments || 0,
      finalizedAssessments: assessmentData?.finalizedAssessments || 0,
      totalJournals: journalData?.totalJournals || 0,
      submittedJournals: journalData?.submittedJournals || 0,
      totalProgressRecords: progressData?.totalRecords || 0,
      period,
    })
  }
}
