export class DashboardSummary {
  constructor({
    totalTeachingSessions,
    completedTeachingSessions,
    totalAttendanceSessions,
    submittedAttendanceSessions,
    totalAssessments,
    finalizedAssessments,
    totalJournals,
    submittedJournals,
    totalProgressRecords,
    period
  }) {
    this.totalTeachingSessions = totalTeachingSessions || 0
    this.completedTeachingSessions = completedTeachingSessions || 0
    this.totalAttendanceSessions = totalAttendanceSessions || 0
    this.submittedAttendanceSessions = submittedAttendanceSessions || 0
    this.totalAssessments = totalAssessments || 0
    this.finalizedAssessments = finalizedAssessments || 0
    this.totalJournals = totalJournals || 0
    this.submittedJournals = submittedJournals || 0
    this.totalProgressRecords = totalProgressRecords || 0
    this.period = period || null
  }

  isValid() {
    return !!(this.period && this.period.academicYearId && this.period.semesterId)
  }
}
