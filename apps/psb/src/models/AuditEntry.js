export class AuditEntry {
  constructor({ id, sessionId, timestamp, action, details }) {
    this.id = id
    this.sessionId = sessionId
    this.timestamp = timestamp
    this.action = action // 'SESSION_STARTED', 'ATTENDANCE_SAVED', 'JOURNAL_SAVED', 'SESSION_COMPLETED', dll.
    this.details = details // string atau objek kecil
  }
}
