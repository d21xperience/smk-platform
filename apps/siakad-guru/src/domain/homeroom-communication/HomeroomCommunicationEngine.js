// FILE: src/domain/homeroom-communication/HomeroomCommunicationEngine.js
// STATUS: NEW
// STATUS IMPLEMENTASI: COMPLETE

/**
 * Homeroom Communication Engine
 * Pure JavaScript Domain Engine untuk business rules log komunikasi wali kelas dengan orang tua.
 * TIDAK BOLEH import: Vue, Pinia, Axios, Quasar, HTTP, Database, Browser API.
 *
 * Business Rules:
 * - Kategori: PHONE_CALL, WHATSAPP, MEETING, HOME_VISIT, OTHER
 * - Arah: INBOUND (dari ortu), OUTBOUND (ke ortu), BIDIRECTIONAL
 * - Status: PENDING_FOLLOW_UP → FOLLOWED_UP → CLOSED (one-way, tidak bisa mundur)
 * - PENDING_FOLLOW_UP bisa langsung ke CLOSED (tanpa follow-up)
 * - Judul: min 3, max 200 karakter
 * - Ringkasan: min 10, max 2000 karakter
 * - Follow-up note: opsional, max 1000 karakter
 * - Tanggal komunikasi tidak boleh di masa depan
 * - studentId wajib diisi
 */

export const COMM_CATEGORY = {
  PHONE_CALL: 'PHONE_CALL',
  WHATSAPP: 'WHATSAPP',
  MEETING: 'MEETING',
  HOME_VISIT: 'HOME_VISIT',
  OTHER: 'OTHER',
}

export const COMM_DIRECTION = {
  INBOUND: 'INBOUND',
  OUTBOUND: 'OUTBOUND',
  BIDIRECTIONAL: 'BIDIRECTIONAL',
}

export const COMM_STATUS = {
  PENDING_FOLLOW_UP: 'PENDING_FOLLOW_UP',
  FOLLOWED_UP: 'FOLLOWED_UP',
  CLOSED: 'CLOSED',
}

const VALID_CATEGORIES = Object.values(COMM_CATEGORY)
const VALID_DIRECTIONS = Object.values(COMM_DIRECTION)
const VALID_STATUSES = Object.values(COMM_STATUS)

const TITLE_MIN_LENGTH = 3
const TITLE_MAX_LENGTH = 200
const SUMMARY_MIN_LENGTH = 10
const SUMMARY_MAX_LENGTH = 2000
const FOLLOW_UP_NOTE_MAX_LENGTH = 1000

const VALID_TRANSITIONS = {
  [COMM_STATUS.PENDING_FOLLOW_UP]: [COMM_STATUS.FOLLOWED_UP, COMM_STATUS.CLOSED],
  [COMM_STATUS.FOLLOWED_UP]: [COMM_STATUS.CLOSED],
  [COMM_STATUS.CLOSED]: [],
}

export class HomeroomCommunicationEngine {
  /**
   * Validasi input log komunikasi baru.
   * Melempar Error jika ada validasi yang gagal.
   * @param {Object} commData - Data log komunikasi dari form
   */
  validateCommunicationInput(commData) {
    if (!commData) {
      throw new Error('Data komunikasi tidak boleh kosong.')
    }

    if (!commData.studentId || String(commData.studentId).trim() === '') {
      throw new Error('Siswa harus dipilih untuk log komunikasi ini.')
    }

    if (!VALID_CATEGORIES.includes(commData.category)) {
      throw new Error(`Kategori tidak valid. Harus salah satu dari: ${VALID_CATEGORIES.join(', ')}`)
    }

    if (!VALID_DIRECTIONS.includes(commData.direction)) {
      throw new Error(
        `Arah komunikasi tidak valid. Harus salah satu dari: ${VALID_DIRECTIONS.join(', ')}`,
      )
    }

    const title = String(commData.title || '').trim()
    if (title.length < TITLE_MIN_LENGTH) {
      throw new Error(`Judul minimal ${TITLE_MIN_LENGTH} karakter.`)
    }
    if (title.length > TITLE_MAX_LENGTH) {
      throw new Error(`Judul maksimal ${TITLE_MAX_LENGTH} karakter.`)
    }

    const summary = String(commData.summary || '').trim()
    if (summary.length < SUMMARY_MIN_LENGTH) {
      throw new Error(`Ringkasan minimal ${SUMMARY_MIN_LENGTH} karakter.`)
    }
    if (summary.length > SUMMARY_MAX_LENGTH) {
      throw new Error(`Ringkasan maksimal ${SUMMARY_MAX_LENGTH} karakter.`)
    }

    if (commData.followUpNote !== undefined && commData.followUpNote !== null) {
      const followUpNote = String(commData.followUpNote).trim()
      if (followUpNote.length > FOLLOW_UP_NOTE_MAX_LENGTH) {
        throw new Error(`Catatan tindak lanjut maksimal ${FOLLOW_UP_NOTE_MAX_LENGTH} karakter.`)
      }
    }

    if (!commData.communicationDate || String(commData.communicationDate).trim() === '') {
      throw new Error('Tanggal komunikasi wajib diisi.')
    }

    this.validateCommunicationDate(commData.communicationDate)
  }

  /**
   * Validasi tanggal komunikasi tidak boleh di masa depan.
   * @param {string} dateStr - Tanggal dalam format YYYY-MM-DD.
   */
  validateCommunicationDate(dateStr) {
    const commDate = new Date(dateStr)
    if (Number.isNaN(commDate.getTime())) {
      throw new Error('Format tanggal tidak valid. Gunakan format YYYY-MM-DD.')
    }

    const today = new Date()
    today.setHours(23, 59, 59, 999)

    if (commDate.getTime() > today.getTime()) {
      throw new Error('Tanggal komunikasi tidak boleh di masa depan.')
    }
  }

  /**
   * Membuat objek log komunikasi baru dengan status PENDING_FOLLOW_UP.
   * Asumsi: commData sudah divalidasi sebelumnya.
   * @param {Object} commData - Data log komunikasi yang sudah valid.
   * @param {string} createdBy - User ID pembuat.
   * @returns {Object} Objek log komunikasi baru.
   */
  createCommunication(commData, createdBy) {
    return {
      communicationId: null,
      studentId: commData.studentId,
      studentName: commData.studentName || '',
      category: commData.category,
      direction: commData.direction,
      title: String(commData.title).trim(),
      summary: String(commData.summary).trim(),
      communicationDate: commData.communicationDate,
      followUpNote: String(commData.followUpNote || '').trim(),
      status: COMM_STATUS.PENDING_FOLLOW_UP,
      createdAt: new Date().toISOString(),
      createdBy: createdBy,
    }
  }

  /**
   * Validasi transisi status.
   * Hanya transisi maju yang diizinkan:
   * PENDING_FOLLOW_UP → FOLLOWED_UP atau CLOSED
   * FOLLOWED_UP → CLOSED
   * CLOSED → tidak bisa ke mana-mana
   * @param {Object} communication - Log komunikasi saat ini.
   * @param {string} newStatus - Status baru yang diminta.
   */
  validateStatusTransition(communication, newStatus) {
    if (!communication) {
      throw new Error('Log komunikasi tidak ditemukan.')
    }

    if (!VALID_STATUSES.includes(newStatus)) {
      throw new Error(`Status tidak valid. Harus salah satu dari: ${VALID_STATUSES.join(', ')}`)
    }

    const currentStatus = communication.status
    const allowedTransitions = VALID_TRANSITIONS[currentStatus] || []

    if (!allowedTransitions.includes(newStatus)) {
      throw new Error(
        `Tidak bisa mengubah status dari ${currentStatus} ke ${newStatus}. ` +
          `Transisi yang diizinkan: ${allowedTransitions.length > 0 ? allowedTransitions.join(', ') : 'tidak ada'}.`,
      )
    }
  }

  /**
   * Menerapkan perubahan status pada log komunikasi.
   * @param {Object} communication - Log komunikasi yang akan diupdate.
   * @param {string} newStatus - Status baru.
   * @param {string} followUpNote - Catatan tindak lanjut (opsional).
   * @returns {Object} Log komunikasi dengan status baru.
   */
  applyStatusChange(communication, newStatus, followUpNote) {
    const updatedFollowUpNote =
      followUpNote !== undefined && followUpNote !== null
        ? String(followUpNote).trim()
        : communication.followUpNote

    return {
      ...communication,
      status: newStatus,
      followUpNote: updatedFollowUpNote,
    }
  }

  /**
   * Bangun summary dari array log komunikasi.
   * @param {Array<Object>} communications - Array log komunikasi.
   * @returns {Object} Summary log komunikasi.
   */
  buildSummary(communications) {
    if (!Array.isArray(communications) || communications.length === 0) {
      return {
        totalCommunications: 0,
        pendingCount: 0,
        followedUpCount: 0,
        closedCount: 0,
        inboundCount: 0,
        outboundCount: 0,
        bidirectionalCount: 0,
      }
    }

    let pendingCount = 0
    let followedUpCount = 0
    let closedCount = 0
    let inboundCount = 0
    let outboundCount = 0
    let bidirectionalCount = 0

    for (const comm of communications) {
      if (comm.status === COMM_STATUS.PENDING_FOLLOW_UP) pendingCount++
      if (comm.status === COMM_STATUS.FOLLOWED_UP) followedUpCount++
      if (comm.status === COMM_STATUS.CLOSED) closedCount++
      if (comm.direction === COMM_DIRECTION.INBOUND) inboundCount++
      if (comm.direction === COMM_DIRECTION.OUTBOUND) outboundCount++
      if (comm.direction === COMM_DIRECTION.BIDIRECTIONAL) bidirectionalCount++
    }

    return {
      totalCommunications: communications.length,
      pendingCount: pendingCount,
      followedUpCount: followedUpCount,
      closedCount: closedCount,
      inboundCount: inboundCount,
      outboundCount: outboundCount,
      bidirectionalCount: bidirectionalCount,
    }
  }

  /**
   * Urutkan log komunikasi: PENDING dulu, lalu FOLLOWED_UP, lalu CLOSED.
   * Dalam setiap grup, urutkan berdasarkan tanggal terbaru.
   * @param {Array<Object>} communications - Array log komunikasi.
   * @returns {Array<Object>} Array log komunikasi terurut.
   */
  sortCommunications(communications) {
    if (!Array.isArray(communications)) return []

    const statusOrder = {
      [COMM_STATUS.PENDING_FOLLOW_UP]: 0,
      [COMM_STATUS.FOLLOWED_UP]: 1,
      [COMM_STATUS.CLOSED]: 2,
    }

    return [...communications].sort((a, b) => {
      const statusDiff = (statusOrder[a.status] || 99) - (statusOrder[b.status] || 99)
      if (statusDiff !== 0) return statusDiff

      return new Date(b.communicationDate) - new Date(a.communicationDate)
    })
  }
}
