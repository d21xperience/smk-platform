// src/models/StudentProgress.js

export class StudentProgress {
  constructor({
    studentId,
    studentName,
    nis,
    className,
    academicYearId,
    semesterId,
    attendanceSummary = null, // AttendanceSummary
    assessmentSummaries = [], // AssessmentSummary[]
    teacherNotes = [], // TeacherNote[]
    homeroomNotes = [], // HomeroomNote[]
    achievements = [], // Achievement[]
    violations = [], // Violation[]
    counselingRecords = [], // CounselingRecord[] (opsional)
    lastUpdated = null,
  }) {
    this.studentId = studentId
    this.studentName = studentName
    this.nis = nis
    this.className = className
    this.academicYearId = academicYearId
    this.semesterId = semesterId
    this.attendanceSummary = attendanceSummary
    this.assessmentSummaries = assessmentSummaries
    this.teacherNotes = teacherNotes
    this.homeroomNotes = homeroomNotes
    this.achievements = achievements
    this.violations = violations
    this.counselingRecords = counselingRecords
    this.lastUpdated = lastUpdated || new Date().toISOString()
  }
}
