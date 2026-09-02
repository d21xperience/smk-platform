// FILE: src/contracts/homeroomCommunicationContract.js
// STATUS: NEW
// STATUS IMPLEMENTASI: COMPLETE

/**
 * Homeroom Communication Contract
 * Menentukan bentuk data yang dibutuhkan oleh UI Log Komunikasi Wali Kelas.
 * Frontend-First: hanya field yang diperlukan UI yang dipetakan.
 */

export const CONTRACT_COMM_CATEGORY = {
  PHONE_CALL: 'PHONE_CALL',
  WHATSAPP: 'WHATSAPP',
  MEETING: 'MEETING',
  HOME_VISIT: 'HOME_VISIT',
  OTHER: 'OTHER',
}

export const CONTRACT_COMM_DIRECTION = {
  INBOUND: 'INBOUND',
  OUTBOUND: 'OUTBOUND',
  BIDIRECTIONAL: 'BIDIRECTIONAL',
}

export const CONTRACT_COMM_STATUS = {
  PENDING_FOLLOW_UP: 'PENDING_FOLLOW_UP',
  FOLLOWED_UP: 'FOLLOWED_UP',
  CLOSED: 'CLOSED',
}

/**
 * Memetakan satu log komunikasi dari adapter/domain ke kontrak UI.
 * @param {Object} raw - Data log komunikasi mentah.
 * @returns {Object} Data log komunikasi sesuai kontrak UI.
 */
export function mapHomeroomCommunication(raw) {
  if (!raw) return null

  return {
    communicationId: raw.communicationId || null,
    studentId: raw.studentId || '',
    studentName: raw.studentName || '',
    category: raw.category || CONTRACT_COMM_CATEGORY.OTHER,
    direction: raw.direction || CONTRACT_COMM_DIRECTION.OUTBOUND,
    title: raw.title || '',
    summary: raw.summary || '',
    communicationDate: raw.communicationDate || '',
    followUpNote: raw.followUpNote || '',
    status: raw.status || CONTRACT_COMM_STATUS.PENDING_FOLLOW_UP,
    createdAt: raw.createdAt || '',
    createdBy: raw.createdBy || '',
  }
}

/**
 * Memetakan array log komunikasi ke kontrak UI.
 * @param {Array} rawList - Array data log komunikasi mentah.
 * @returns {Array} Array data log komunikasi sesuai kontrak UI.
 */
export function mapHomeroomCommunications(rawList) {
  if (!Array.isArray(rawList)) return []
  return rawList.map(mapHomeroomCommunication).filter((c) => c && c.communicationId)
}

/**
 * Memetakan summary dari engine ke kontrak UI.
 * @param {Object} summary - Summary dari engine.
 * @returns {Object} Summary sesuai kontrak UI.
 */
export function mapHomeroomCommunicationSummary(summary) {
  if (!summary) {
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

  return {
    totalCommunications: summary.totalCommunications || 0,
    pendingCount: summary.pendingCount || 0,
    followedUpCount: summary.followedUpCount || 0,
    closedCount: summary.closedCount || 0,
    inboundCount: summary.inboundCount || 0,
    outboundCount: summary.outboundCount || 0,
    bidirectionalCount: summary.bidirectionalCount || 0,
  }
}
