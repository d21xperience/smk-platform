export class TeacherPresence {
  constructor({ id, sessionId, teacherId, checkInTime, checkOutTime, status }) {
    this.id = id
    this.sessionId = sessionId
    this.teacherId = teacherId
    this.checkInTime = checkInTime // ISO timestamp
    this.checkOutTime = checkOutTime
    this.status = status // 'on_time', 'late', 'absent'
  }
}
