export class StudentAttendanceItem {
  constructor({ studentId, studentName, nis, status = 'unmarked', note = '' }) {
    this.studentId = studentId
    this.studentName = studentName
    this.nis = nis
    this.status = status // 'hadir', 'izin', 'sakit', 'alpha', 'unmarked'
    this.note = note
  }
}

export class AttendanceRecord {
  constructor({ id, sessionId, date, className, subject, students, status = 'draft' }) {
    this.id = id
    this.sessionId = sessionId // dari TeachingSession
    this.date = date
    this.className = className
    this.subject = subject
    this.students = students // array of StudentAttendanceItem
    this.status = status // 'draft', 'submitted'
    this.lastSavedAt = null
  }
}
