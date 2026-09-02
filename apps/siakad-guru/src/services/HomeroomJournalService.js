// FILE: src/services/HomeroomJournalService.js
// STATUS: NEW
// STATUS IMPLEMENTASI: COMPLETE

import { HomeroomJournalEngine } from '../domain/homeroom-journal/HomeroomJournalEngine.js'
import {
  createNoteCreatedEvent,
  createNoteSubmittedEvent,
  createNoteDeletedEvent,
} from '../domain/homeroom-journal/events.js'
import {
  mapHomeroomNotes,
  mapHomeroomNote,
  mapHomeroomJournalSummary,
} from '../contracts/homeroomJournalContract.js'

/**
 * Homeroom Journal Application Service
 * Orchestration use case CRUD catatan wali kelas.
 * Jalur dengan Business Rule + Mutasi:
 * Service memanggil Domain Engine untuk validasi & state transition,
 * Adapter untuk persistence, dan EventDispatcher untuk publish events.
 */
export class HomeroomJournalService {
  constructor({ homeroomJournalAdapter, eventDispatcher }) {
    this.adapter = homeroomJournalAdapter
    this.eventDispatcher = eventDispatcher
    this.engine = new HomeroomJournalEngine()
  }

  /**
   * Mengambil daftar catatan wali kelas beserta summary.
   * @param {Object} context - Operational context
   * @param {string} context.userId - User ID wali kelas
   * @param {string} context.schoolId - School ID
   * @param {string} context.academicYearId - Academic Year ID
   * @param {string} context.semesterId - Semester ID
   * @returns {Promise<Object>} { notes, summary }
   */
  async getNotes(context) {
    if (!context.userId) {
      throw new Error('userId is required to fetch homeroom journal.')
    }
    if (!context.academicYearId) {
      throw new Error('academicYearId is required to fetch homeroom journal.')
    }
    if (!context.semesterId) {
      throw new Error('semesterId is required to fetch homeroom journal.')
    }

    const rawNotes = await this.adapter.fetchNotesByHomeroom(context)
    const notes = mapHomeroomNotes(rawNotes)
    const sortedNotes = this.engine.sortNotes(notes)
    const summary = mapHomeroomJournalSummary(this.engine.buildSummary(sortedNotes))

    return { notes: sortedNotes, summary }
  }

  /**
   * Membuat catatan baru.
   * Flow: Engine.validate → Engine.create → Adapter.save → Dispatch event.
   * @param {Object} context - Operational context
   * @param {Object} noteData - Data catatan dari form
   * @returns {Promise<Object>} Catatan yang tersimpan
   */
  async createNote(context, noteData) {
    if (!context.userId) {
      throw new Error('userId is required to create note.')
    }

    this.engine.validateNoteInput(noteData)

    const note = this.engine.createNote(noteData, context.userId)

    const savedNote = await this.adapter.saveNote(note, context)

    const mappedNote = mapHomeroomNote(savedNote)

    this.eventDispatcher.dispatch(createNoteCreatedEvent(mappedNote, context.userId))

    return mappedNote
  }

  /**
   * Meng-submit catatan (finalisasi DRAFT → SUBMITTED).
   * Flow: Engine.validateSubmit → Engine.applySubmit → Adapter.update → Dispatch event.
   * @param {Object} context - Operational context
   * @param {Object} note - Catatan yang akan di-submit
   * @returns {Promise<Object>} Catatan yang terupdate
   */
  async submitNote(context, note) {
    if (!context.userId) {
      throw new Error('userId is required to submit note.')
    }

    this.engine.validateSubmit(note)

    const submittedNote = this.engine.applySubmit(note)

    const updatedNote = await this.adapter.updateNote(submittedNote, context)

    const mappedNote = mapHomeroomNote(updatedNote)

    this.eventDispatcher.dispatch(createNoteSubmittedEvent(mappedNote, context.userId))

    return mappedNote
  }

  /**
   * Menghapus catatan draft.
   * Flow: Engine.validateDelete → Adapter.delete → Dispatch event.
   * @param {Object} context - Operational context
   * @param {Object} note - Catatan yang akan dihapus
   * @returns {Promise<Object>} Hasil penghapusan
   */
  async deleteNote(context, note) {
    if (!context.userId) {
      throw new Error('userId is required to delete note.')
    }

    this.engine.validateDelete(note)

    const result = await this.adapter.deleteNote(note.noteId, context)

    this.eventDispatcher.dispatch(createNoteDeletedEvent(note.noteId, context.userId))

    return result
  }

  /**
   * Bangun summary dari array catatan.
   * Digunakan oleh Store untuk rebuild summary setelah mutasi lokal.
   * @param {Array<Object>} notes - Array catatan.
   * @returns {Object} Summary yang telah dipetakan ke kontrak UI.
   */
  buildSummary(notes) {
    return mapHomeroomJournalSummary(this.engine.buildSummary(notes))
  }
}
