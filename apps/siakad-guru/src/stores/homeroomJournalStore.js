// FILE: src/stores/homeroomJournalStore.js
// STATUS: NEW
// STATUS IMPLEMENTASI: COMPLETE

import { defineStore } from 'pinia'
import { homeroomJournalService } from '../boot/services.js'

/**
 * Homeroom Journal Store
 * Mengelola state catatan wali kelas (CRUD).
 * Pola mengikuti homeroomBillingStore:
 * - Tidak import store lain
 * - Action menerima parameter eksplisit dari composable
 * - Store memanggil Application Service, bukan Axios
 */
export const useHomeroomJournalStore = defineStore('homeroomJournal', {
  state: () => ({
    notes: [],
    summary: null,
    loading: false,
    creating: false,
    submitting: false,
    deleting: false,
    error: null,
  }),

  getters: {
    hasNotes: (state) => state.notes.length > 0,
    hasSummary: (state) => !!state.summary,
    draftNotes: (state) => state.notes.filter(n => n.status === 'DRAFT'),
    submittedNotes: (state) => state.notes.filter(n => n.status === 'SUBMITTED'),
    hasDraftNotes: (state) => state.notes.some(n => n.status === 'DRAFT'),
  },

  actions: {
    /**
     * Memuat daftar catatan wali kelas.
     * @param {Object} params - Parameter dari composable
     */
    async loadNotes({ userId, schoolId, academicYearId, semesterId }) {
      this.loading = true
      this.error = null
      try {
        const result = await homeroomJournalService.getNotes({
          userId,
          schoolId,
          academicYearId,
          semesterId,
        })
        this.notes = result.notes
        this.summary = result.summary
      } catch (err) {
        this.error = err.message
        this.notes = []
        this.summary = null
      } finally {
        this.loading = false
      }
    },

    /**
     * Membuat catatan baru.
     * Setelah berhasil, note ditambahkan ke array dan summary di-rebuild.
     * @param {Object} params - Parameter dari composable
     * @param {Object} params.noteData - Data catatan dari form
     */
    async createNote({ userId, schoolId, academicYearId, semesterId, noteData }) {
      this.creating = true
      this.error = null
      try {
        const savedNote = await homeroomJournalService.createNote(
          { userId, schoolId, academicYearId, semesterId },
          noteData,
        )
        this.notes.unshift(savedNote)
        this.summary = homeroomJournalService.buildSummary(this.notes)
        return savedNote
      } catch (err) {
        this.error = err.message
        throw err
      } finally {
        this.creating = false
      }
    },

    /**
     * Meng-submit catatan (finalisasi DRAFT → SUBMITTED).
     * Setelah berhasil, note di array di-update dan summary di-rebuild.
     * @param {Object} params - Parameter dari composable
     * @param {Object} params.note - Catatan yang akan di-submit
     */
    async submitNote({ userId, schoolId, academicYearId, semesterId, note }) {
      this.submitting = true
      this.error = null
      try {
        const updatedNote = await homeroomJournalService.submitNote(
          { userId, schoolId, academicYearId, semesterId },
          note,
        )
        const noteIndex = this.notes.findIndex(n => n.noteId === updatedNote.noteId)
        if (noteIndex !== -1) {
          this.notes[noteIndex] = updatedNote
        }
        this.summary = homeroomJournalService.buildSummary(this.notes)
        return updatedNote
      } catch (err) {
        this.error = err.message
        throw err
      } finally {
        this.submitting = false
      }
    },

    /**
     * Menghapus catatan draft.
     * Setelah berhasil, note dihapus dari array dan summary di-rebuild.
     * @param {Object} params - Parameter dari composable
     * @param {Object} params.note - Catatan yang akan dihapus
     */
    async deleteNote({ userId, schoolId, academicYearId, semesterId, note }) {
      this.deleting = true
      this.error = null
      try {
        const result = await homeroomJournalService.deleteNote(
          { userId, schoolId, academicYearId, semesterId },
          note,
        )
        this.notes = this.notes.filter(n => n.noteId !== note.noteId)
        this.summary = homeroomJournalService.buildSummary(this.notes)
        return result
      } catch (err) {
        this.error = err.message
        throw err
      } finally {
        this.deleting = false
      }
    },

    /**
     * Mengosongkan state store.
     */
    clearData() {
      this.notes = []
      this.summary = null
      this.loading = false
      this.creating = false
      this.submitting = false
      this.deleting = false
      this.error = null
    },
  },
})
