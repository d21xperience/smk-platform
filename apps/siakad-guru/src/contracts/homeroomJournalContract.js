// FILE: src/contracts/homeroomJournalContract.js
// STATUS: NEW
// STATUS IMPLEMENTASI: COMPLETE

/**
 * Homeroom Journal Contract
 * Menentukan bentuk data yang dibutuhkan oleh UI Catatan Wali Kelas.
 * Frontend-First: hanya field yang diperlukan UI yang dipetakan.
 */

export const CONTRACT_NOTE_CATEGORY = {
  POSITIVE: 'POSITIVE',
  NEGATIVE: 'NEGATIVE',
  NEUTRAL: 'NEUTRAL',
  INCIDENT: 'INCIDENT',
}

export const CONTRACT_NOTE_PRIORITY = {
  LOW: 'LOW',
  MEDIUM: 'MEDIUM',
  HIGH: 'HIGH',
  URGENT: 'URGENT',
}

export const CONTRACT_NOTE_STATUS = {
  DRAFT: 'DRAFT',
  SUBMITTED: 'SUBMITTED',
}

/**
 * Memetakan satu catatan dari adapter/domain ke kontrak UI.
 * @param {Object} raw - Data catatan mentah.
 * @returns {Object} Data catatan sesuai kontrak UI.
 */
export function mapHomeroomNote(raw) {
  if (!raw) return null

  return {
    noteId: raw.noteId || null,
    studentId: raw.studentId || '',
    studentName: raw.studentName || '',
    category: raw.category || CONTRACT_NOTE_CATEGORY.NEUTRAL,
    priority: raw.priority || CONTRACT_NOTE_PRIORITY.MEDIUM,
    title: raw.title || '',
    description: raw.description || '',
    noteDate: raw.noteDate || '',
    createdAt: raw.createdAt || '',
    createdBy: raw.createdBy || '',
    status: raw.status || CONTRACT_NOTE_STATUS.DRAFT,
  }
}

/**
 * Memetakan array catatan ke kontrak UI.
 * @param {Array} rawList - Array data catatan mentah.
 * @returns {Array} Array data catatan sesuai kontrak UI.
 */
export function mapHomeroomNotes(rawList) {
  if (!Array.isArray(rawList)) return []
  return rawList.map(mapHomeroomNote).filter((n) => n && n.noteId)
}

/**
 * Memetakan summary dari engine ke kontrak UI.
 * @param {Object} summary - Summary dari engine.
 * @returns {Object} Summary sesuai kontrak UI.
 */
export function mapHomeroomJournalSummary(summary) {
  if (!summary) {
    return {
      totalNotes: 0,
      draftCount: 0,
      submittedCount: 0,
      positiveCount: 0,
      negativeCount: 0,
      neutralCount: 0,
      incidentCount: 0,
      urgentCount: 0,
    }
  }

  return {
    totalNotes: summary.totalNotes || 0,
    draftCount: summary.draftCount || 0,
    submittedCount: summary.submittedCount || 0,
    positiveCount: summary.positiveCount || 0,
    negativeCount: summary.negativeCount || 0,
    neutralCount: summary.neutralCount || 0,
    incidentCount: summary.incidentCount || 0,
    urgentCount: summary.urgentCount || 0,
  }
}
