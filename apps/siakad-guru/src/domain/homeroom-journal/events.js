// FILE: src/domain/homeroom-journal/events.js
// STATUS: NEW
// STATUS IMPLEMENTASI: COMPLETE

/**
 * Homeroom Journal Domain Events
 * Event types dan factory functions untuk publish melalui EventDispatcher.
 * Events ini dipublish oleh Service setelah mutasi berhasil.
 */

export const HOMEROOM_JOURNAL_EVENTS = {
  NOTE_CREATED: 'homeroom.journal.note.created',
  NOTE_SUBMITTED: 'homeroom.journal.note.submitted',
  NOTE_DELETED: 'homeroom.journal.note.deleted',
}

/**
 * Factory untuk event Note Created.
 * @param {Object} note - Catatan yang baru dibuat.
 * @param {string} userId - User ID yang membuat.
 * @returns {Object} Event payload.
 */
export function createNoteCreatedEvent(note, userId) {
  return {
    type: HOMEROOM_JOURNAL_EVENTS.NOTE_CREATED,
    payload: {
      noteId: note.noteId,
      studentId: note.studentId,
      studentName: note.studentName,
      category: note.category,
      priority: note.priority,
      title: note.title,
      status: note.status,
      noteDate: note.noteDate,
      createdBy: userId,
      timestamp: new Date().toISOString(),
    },
  }
}

/**
 * Factory untuk event Note Submitted.
 * @param {Object} note - Catatan yang di-submit.
 * @param {string} userId - User ID yang submit.
 * @returns {Object} Event payload.
 */
export function createNoteSubmittedEvent(note, userId) {
  return {
    type: HOMEROOM_JOURNAL_EVENTS.NOTE_SUBMITTED,
    payload: {
      noteId: note.noteId,
      studentId: note.studentId,
      studentName: note.studentName,
      category: note.category,
      title: note.title,
      status: note.status,
      submittedBy: userId,
      timestamp: new Date().toISOString(),
    },
  }
}

/**
 * Factory untuk event Note Deleted.
 * @param {string} noteId - ID catatan yang dihapus.
 * @param {string} userId - User ID yang menghapus.
 * @returns {Object} Event payload.
 */
export function createNoteDeletedEvent(noteId, userId) {
  return {
    type: HOMEROOM_JOURNAL_EVENTS.NOTE_DELETED,
    payload: {
      noteId: noteId,
      deletedBy: userId,
      timestamp: new Date().toISOString(),
    },
  }
}
