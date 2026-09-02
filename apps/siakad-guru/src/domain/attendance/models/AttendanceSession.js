import { SESSION_STATUS } from './AttendanceStatus.js'

export class AttendanceSession {
  constructor({
    id,
    teachingSessionId,
    classId,
    className,
    subjectId,
    subjectName,
    date,
    status,
    records,
    academicYearId,
    semesterId,
    schoolId
  }) {
    this.id = id
    this.teachingSessionId = teachingSessionId
    this.classId = classId
    this.className = className
    this.subjectId = subjectId
    this.subjectName = subjectName
    this.date = date
    this.status = status || SESSION_STATUS.DRAFT
    this.records = records || []
    this.academicYearId = academicYearId
    this.semesterId = semesterId
    this.schoolId = schoolId
  }

  isDraft() {
    return this.status === SESSION_STATUS.DRAFT
  }

  isSubmitted() {
    return this.status === SESSION_STATUS.SUBMITTED
  }

  getTotalStudents() {
    return this.records.length
  }

  getPresentCount() {
    return this.records.filter(r => r.status === 'present').length
  }

  getSickCount() {
    return this.records.filter(r => r.status === 'sick').length
  }

  getPermissionCount() {
    return this.records.filter(r => r.status === 'permission').length
  }

  getAbsentCount() {
    return this.records.filter(r => r.status === 'absent').length
  }

  isValid() {
    return !!(
      this.id &&
      this.teachingSessionId &&
      this.classId &&
      this.subjectId &&
      this.date &&
      this.status &&
      this.academicYearId &&
      this.semesterId &&
      this.schoolId
    )
  }
}
