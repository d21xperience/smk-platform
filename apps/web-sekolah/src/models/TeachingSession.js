export class TeachingSession {
  constructor({
    id,
    scheduleId,
    date,
    startTime,
    endTime,
    subject,
    className,
    status, // 'scheduled' | 'started' | 'in_progress' | 'completed' | 'locked'
    teacherPresence, // instance TeacherPresence atau null
    auditTrail = [], // array of AuditEntry
  }) {
    this.id = id
    this.scheduleId = scheduleId
    this.date = date
    this.startTime = startTime
    this.endTime = endTime
    this.subject = subject
    this.className = className
    this.status = status
    this.teacherPresence = teacherPresence
    this.auditTrail = auditTrail
  }
}
