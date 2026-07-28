// src/engine/TeachingEngine.js

export const SessionStatus = {
  SCHEDULED: 'scheduled',
  STARTED: 'started',
  IN_PROGRESS: 'in_progress',
  COMPLETED: 'completed',
  LOCKED: 'locked',
}

const TRANSITIONS = {
  [SessionStatus.SCHEDULED]: [SessionStatus.STARTED],
  [SessionStatus.STARTED]: [SessionStatus.IN_PROGRESS],
  [SessionStatus.IN_PROGRESS]: [SessionStatus.COMPLETED],
  [SessionStatus.COMPLETED]: [SessionStatus.LOCKED],
  [SessionStatus.LOCKED]: [],
}

export const AuditAction = {
  SESSION_STARTED: 'SESSION_STARTED',
  TEACHER_CHECK_IN: 'TEACHER_CHECK_IN',
  TEACHER_CHECK_OUT: 'TEACHER_CHECK_OUT',
  SESSION_IN_PROGRESS: 'SESSION_IN_PROGRESS',
  ATTENDANCE_SAVED: 'ATTENDANCE_SAVED',
  ATTENDANCE_SUBMITTED: 'ATTENDANCE_SUBMITTED',
  JOURNAL_SAVED: 'JOURNAL_SAVED',
  SESSION_COMPLETED: 'SESSION_COMPLETED',
  SESSION_LOCKED: 'SESSION_LOCKED',
}

export class TeachingEngine {
  /**
   * Melakukan transisi status.
   * @param {string} currentStatus
   * @param {string} targetStatus
   * @returns {string} status baru
   * @throws {Error} jika transisi tidak valid
   */
  static transition(currentStatus, targetStatus) {
    const allowed = TRANSITIONS[currentStatus]
    if (!allowed || !allowed.includes(targetStatus)) {
      throw new Error(`Transisi tidak valid: ${currentStatus} → ${targetStatus}`)
    }
    return targetStatus
  }

  /**
   * Cek apakah transisi diizinkan.
   */
  static canTransition(currentStatus, targetStatus) {
    const allowed = TRANSITIONS[currentStatus]
    return allowed?.includes(targetStatus) ?? false
  }

  /**
   * Apakah status terminal (tidak bisa berubah lagi).
   */
  static isTerminalStatus(status) {
    return status === SessionStatus.LOCKED
  }

  /**
   * Apakah sesi masih bisa diedit (absensi, jurnal, presence).
   */
  static isEditable(status) {
    return [SessionStatus.SCHEDULED, SessionStatus.STARTED, SessionStatus.IN_PROGRESS].includes(
      status,
    )
  }

  /**
   * Hitung status kehadiran guru berdasarkan waktu check-in.
   * @param {string|Date} checkInTime ISO
   * @param {string|Date} scheduledStartTime ISO
   * @param {number} toleranceMinutes default 15
   * @returns {'on_time'|'late'}
   */
  static calculatePresenceStatus(checkInTime, scheduledStartTime, toleranceMinutes = 15) {
    const checkIn = new Date(checkInTime).getTime()
    const scheduled = new Date(scheduledStartTime).getTime()
    const diffMs = checkIn - scheduled
    const diffMinutes = diffMs / (1000 * 60)
    return diffMinutes <= toleranceMinutes ? 'on_time' : 'late'
  }

  /**
   * Buat objek audit trail standar.
   * @param {number|string} sessionId
   * @param {string} action (dari AuditAction)
   * @param {string} details
   * @returns {{ sessionId, timestamp: string, action: string, details: string }}
   */
  static createAuditEntry(sessionId, action, details = '') {
    return {
      sessionId,
      timestamp: new Date().toISOString(),
      action,
      details,
    }
  }

  /**
   * Validasi apakah sesi bisa diselesaikan (completed).
   * Syarat: attendance sudah disubmit, journal sudah diisi.
   * @param {boolean} isAttendanceSubmitted
   * @param {boolean} isJournalFilled
   * @returns {{ valid: boolean, errors: string[] }}
   */
  static validateCompletion(isAttendanceSubmitted, isJournalFilled) {
    const errors = []
    if (!isAttendanceSubmitted) errors.push('Absensi harus sudah disubmit.')
    if (!isJournalFilled) errors.push('Jurnal mengajar harus diisi.')
    return { valid: errors.length === 0, errors }
  }
}
