import { ATTENDANCE_STATUS } from './AttendanceStatus.js'

export class AttendanceRecord {
  constructor({
    studentId,
    studentName,
    status,
    note
  }) {
    this.studentId = studentId
    this.studentName = studentName
    this.status = status || ATTENDANCE_STATUS.PRESENT
    this.note = note || ''
  }

  isPresent() {
    return this.status === ATTENDANCE_STATUS.PRESENT
  }

  isAbsent() {
    return this.status === ATTENDANCE_STATUS.ABSENT
  }

  isValid() {
    return !!(this.studentId && this.studentName && this.status)
  }
}
