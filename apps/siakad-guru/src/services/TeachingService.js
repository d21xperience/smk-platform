export class TeachingService {
  constructor(adapter) {
    this.adapter = adapter
  }

  // --- Jadwal ---
  async fetchTeacherSchedule(academicYearId, semesterId) {
    return this.adapter.getTeacherSchedule(academicYearId, semesterId)
  }

  // --- Sesi (lifecycle) ---
  async fetchSessionsByDate(date) {
    return this.adapter.getTeachingSessions(date)
  }

  async startSession(sessionId) {
    return this.adapter.startSession(sessionId)
  }

  async beginProgress(sessionId) {
    return this.adapter.beginProgress(sessionId)
  }

  async completeSession(sessionId) {
    return this.adapter.completeSession(sessionId)
  }

  async lockSession(sessionId) {
    return this.adapter.lockSession(sessionId)
  }

  // --- Teacher Presence ---
  async checkInTeacher(sessionId, teacherId, checkInTime) {
    return this.adapter.checkInTeacher(sessionId, teacherId, checkInTime)
  }

  async checkOutTeacher(sessionId, checkOutTime) {
    return this.adapter.checkOutTeacher(sessionId, checkOutTime)
  }

  // --- Audit Trail ---
  async getAuditTrail(sessionId) {
    return this.adapter.getAuditTrail(sessionId)
  }

  async addAuditEntry(sessionId, action, details) {
    return this.adapter.addAuditEntry(sessionId, action, details)
  }

  // --- Jurnal ---
  async saveJournal(sessionId, journalData) {
    return this.adapter.saveJournal(sessionId, journalData)
  }

  async getJournalBySession(sessionId) {
    return this.adapter.getJournalBySession(sessionId)
  }
}
