export class AttendanceService {
  /**
   * @param {Object} adapter - objek dengan kontrak attendance (mock/api)
   */
  constructor(adapter) {
    this.adapter = adapter
  }

  /**
   * Mendapatkan daftar siswa untuk sesi tertentu.
   * @param {number} sessionId
   * @param {string} className
   * @returns {Promise<Array<{ studentId: number|string, studentName: string, nis: string, status: string, note: string }>>}
   */
  async fetchStudents(sessionId, className) {
    return this.adapter.fetchStudentsForSession(sessionId, className)
  }

  /**
   * Mendapatkan attendance record yang sudah ada (draft/submitted) untuk sesi.
   * @param {number} sessionId
   * @returns {Promise<import('src/models/Attendance').AttendanceRecord|null>}
   */
  async fetchRecord(sessionId) {
    return this.adapter.fetchAttendanceRecord(sessionId)
  }

  /**
   * Menyimpan draft absensi (create atau update).
   * @param {Object} session - data sesi dari TeachingSession
   * @param {Array} students - array student attendance items
   * @returns {Promise<import('src/models/Attendance').AttendanceRecord>}
   */
  async saveDraft(session, students) {
    return this.adapter.saveDraft(session, students)
  }

  /**
   * Submit final absensi (status menjadi 'submitted').
   * @param {number} sessionId
   * @returns {Promise<import('src/models/Attendance').AttendanceRecord>}
   */
  async submit(sessionId) {
    return this.adapter.submitAttendance(sessionId)
  }

  /**
   * Mendapatkan ringkasan absensi (data mentah, summary dihitung Engine).
   * @param {number} sessionId
   * @returns {Promise<Object|null>}
   */
  async getSummary(sessionId) {
    return this.adapter.getAttendanceSummary(sessionId)
  }
}
