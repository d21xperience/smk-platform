// FILE: src/domain/homeroom-journal/HomeroomJournalEngine.js
// STATUS: NEW
// STATUS IMPLEMENTASI: COMPLETE

/**
 * Homeroom Journal Engine
 * Pure JavaScript Domain Engine untuk business rules catatan wali kelas.
 * TIDAK BOLEH import: Vue, Pinia, Axios, Quasar, HTTP, Database, Browser API.
 *
 * Business Rules:
 * - Kategori: POSITIVE, NEGATIVE, NEUTRAL, INCIDENT
 * - Prioritas: LOW, MEDIUM, HIGH, URGENT
 * - Status: DRAFT → SUBMITTED (one-way transition)
 * - Hanya DRAFT yang bisa di-submit dan di-delete
 * - Judul: min 3, max 200 karakter
 * - Deskripsi: min 10, max 1000 karakter
 * - Tanggal catatan tidak boleh di masa depan
 * - studentId wajib diisi
 */

export const NOTE_CATEGORY = {
  POSITIVE: 'POSITIVE',
  NEGATIVE: 'NEGATIVE',
  NEUTRAL: 'NEUTRAL',
  INCIDENT: 'INCIDENT',
}

export const NOTE_PRIORITY = {
  LOW: 'LOW',
  MEDIUM: 'MEDIUM',
  HIGH: 'HIGH',
  URGENT: 'URGENT',
}

export const NOTE_STATUS = {
  DRAFT: 'DRAFT',
  SUBMITTED: 'SUBMITTED',
}

const VALID_CATEGORIES = Object.values(NOTE_CATEGORY)
const VALID_PRIORITIES = Object.values(NOTE_PRIORITY)
// const VALID_STATUSES = Object.values(NOTE_STATUS)

const TITLE_MIN_LENGTH = 3
const TITLE_MAX_LENGTH = 200
const DESCRIPTION_MIN_LENGTH = 10
const DESCRIPTION_MAX_LENGTH = 1000

export class HomeroomJournalEngine {
  /**
   * Validasi input catatan baru.
   * Melempar Error jika ada validasi yang gagal.
   * @param {Object} noteData - Data catatan dari form
   * @param {string} noteData.studentId - ID siswa (wajib)
   * @param {string} noteData.category - Kategori catatan
   * @param {string} noteData.priority - Prioritas catatan
   * @param {string} noteData.title - Judul catatan
   * @param {string} noteData.description - Deskripsi catatan
   * @param {string} noteData.noteDate - Tanggal catatan (YYYY-MM-DD)
   */
  validateNoteInput(noteData) {
    if (!noteData) {
      throw new Error('Data catatan tidak boleh kosong.')
    }

    if (!noteData.studentId || String(noteData.studentId).trim() === '') {
      throw new Error('Siswa harus dipilih untuk catatan ini.')
    }

    if (!VALID_CATEGORIES.includes(noteData.category)) {
      throw new Error(`Kategori tidak valid. Harus salah satu dari: ${VALID_CATEGORIES.join(', ')}`)
    }

    if (!VALID_PRIORITIES.includes(noteData.priority)) {
      throw new Error(
        `Prioritas tidak valid. Harus salah satu dari: ${VALID_PRIORITIES.join(', ')}`,
      )
    }

    const title = String(noteData.title || '').trim()
    if (title.length < TITLE_MIN_LENGTH) {
      throw new Error(`Judul catatan minimal ${TITLE_MIN_LENGTH} karakter.`)
    }
    if (title.length > TITLE_MAX_LENGTH) {
      throw new Error(`Judul catatan maksimal ${TITLE_MAX_LENGTH} karakter.`)
    }

    const description = String(noteData.description || '').trim()
    if (description.length < DESCRIPTION_MIN_LENGTH) {
      throw new Error(`Deskripsi catatan minimal ${DESCRIPTION_MIN_LENGTH} karakter.`)
    }
    if (description.length > DESCRIPTION_MAX_LENGTH) {
      throw new Error(`Deskripsi catatan maksimal ${DESCRIPTION_MAX_LENGTH} karakter.`)
    }

    if (!noteData.noteDate || String(noteData.noteDate).trim() === '') {
      throw new Error('Tanggal catatan wajib diisi.')
    }

    this.validateNoteDate(noteData.noteDate)
  }

  /**
   * Validasi tanggal catatan tidak boleh di masa depan.
   * @param {string} dateStr - Tanggal dalam format YYYY-MM-DD.
   */
  validateNoteDate(dateStr) {
    const noteDate = new Date(dateStr)
    if (Number.isNaN(noteDate.getTime())) {
      throw new Error('Format tanggal tidak valid. Gunakan format YYYY-MM-DD.')
    }

    const today = new Date()
    today.setHours(23, 59, 59, 999)

    if (noteDate.getTime() > today.getTime()) {
      throw new Error('Tanggal catatan tidak boleh di masa depan.')
    }
  }

  /**
   * Membuat objek catatan baru dengan status DRAFT.
   * Asumsi: noteData sudah divalidasi sebelumnya.
   * @param {Object} noteData - Data catatan yang sudah valid.
   * @param {string} noteData.studentId - ID siswa.
   * @param {string} noteData.studentName - Nama siswa.
   * @param {string} noteData.category - Kategori.
   * @param {string} noteData.priority - Prioritas.
   * @param {string} noteData.title - Judul.
   * @param {string} noteData.description - Deskripsi.
   * @param {string} noteData.noteDate - Tanggal catatan.
   * @param {string} createdBy - User ID pembuat.
   * @returns {Object} Objek catatan baru dengan status DRAFT.
   */
  createNote(noteData, createdBy) {
    return {
      noteId: null,
      studentId: noteData.studentId,
      studentName: noteData.studentName || '',
      category: noteData.category,
      priority: noteData.priority,
      title: String(noteData.title).trim(),
      description: String(noteData.description).trim(),
      noteDate: noteData.noteDate,
      createdAt: new Date().toISOString(),
      createdBy: createdBy,
      status: NOTE_STATUS.DRAFT,
    }
  }

  /**
   * Validasi bahwa catatan bisa di-submit.
   * Hanya DRAFT yang bisa di-submit.
   * @param {Object} note - Catatan yang akan di-submit.
   */
  validateSubmit(note) {
    if (!note) {
      throw new Error('Catatan tidak ditemukan.')
    }

    if (note.status !== NOTE_STATUS.DRAFT) {
      throw new Error(
        `Catatan dengan status ${note.status} tidak bisa di-submit. Hanya DRAFT yang bisa di-submit.`,
      )
    }
  }

  /**
   * Menerapkan submit pada catatan (ubah status DRAFT → SUBMITTED).
   * @param {Object} note - Catatan yang akan di-submit.
   * @returns {Object} Catatan dengan status SUBMITTED.
   */
  applySubmit(note) {
    return {
      ...note,
      status: NOTE_STATUS.SUBMITTED,
    }
  }

  /**
   * Validasi bahwa catatan bisa di-delete.
   * Hanya DRAFT yang bisa di-delete.
   * @param {Object} note - Catatan yang akan di-delete.
   */
  validateDelete(note) {
    if (!note) {
      throw new Error('Catatan tidak ditemukan.')
    }

    if (note.status !== NOTE_STATUS.DRAFT) {
      throw new Error(
        `Catatan dengan status ${note.status} tidak bisa dihapus. Hanya DRAFT yang bisa dihapus.`,
      )
    }
  }

  /**
   * Bangun summary dari array catatan.
   * @param {Array<Object>} notes - Array catatan.
   * @returns {Object} Summary catatan.
   */
  buildSummary(notes) {
    if (!Array.isArray(notes) || notes.length === 0) {
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

    let draftCount = 0
    let submittedCount = 0
    let positiveCount = 0
    let negativeCount = 0
    let neutralCount = 0
    let incidentCount = 0
    let urgentCount = 0

    for (const note of notes) {
      if (note.status === NOTE_STATUS.DRAFT) draftCount++
      if (note.status === NOTE_STATUS.SUBMITTED) submittedCount++
      if (note.category === NOTE_CATEGORY.POSITIVE) positiveCount++
      if (note.category === NOTE_CATEGORY.NEGATIVE) negativeCount++
      if (note.category === NOTE_CATEGORY.NEUTRAL) neutralCount++
      if (note.category === NOTE_CATEGORY.INCIDENT) incidentCount++
      if (note.priority === NOTE_PRIORITY.URGENT) urgentCount++
    }

    return {
      totalNotes: notes.length,
      draftCount: draftCount,
      submittedCount: submittedCount,
      positiveCount: positiveCount,
      negativeCount: negativeCount,
      neutralCount: neutralCount,
      incidentCount: incidentCount,
      urgentCount: urgentCount,
    }
  }

  /**
   * Urutkan catatan: URGENT dulu, lalu berdasarkan tanggal terbaru.
   * @param {Array<Object>} notes - Array catatan.
   * @returns {Array<Object>} Array catatan terurut.
   */
  sortNotes(notes) {
    if (!Array.isArray(notes)) return []

    const priorityOrder = { URGENT: 0, HIGH: 1, MEDIUM: 2, LOW: 3 }

    return [...notes].sort((a, b) => {
      const priorityDiff = (priorityOrder[a.priority] || 99) - (priorityOrder[b.priority] || 99)
      if (priorityDiff !== 0) return priorityDiff

      return new Date(b.noteDate) - new Date(a.noteDate)
    })
  }
}
