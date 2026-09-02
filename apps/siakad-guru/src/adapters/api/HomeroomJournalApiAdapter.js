// FILE: src/adapters/api/HomeroomJournalApiAdapter.js
// STATUS: NEW
// STATUS IMPLEMENTASI: COMPLETE

import { api } from '@/boot/axios'
import { mapHomeroomNotes, mapHomeroomNote } from '@/contracts/homeroomJournalContract'

/**
 * Homeroom Journal API Adapter
 * Menggunakan Axios untuk operasi CRUD catatan wali kelas.
 * Hanya Adapter yang boleh menggunakan Axios.
 */
export class HomeroomJournalApiAdapter {
  /**
   * Mengambil daftar catatan wali kelas berdasarkan context.
   * @param {Object} context - Operational context
   * @param {string} context.userId - User ID wali kelas
   * @param {string} context.schoolId - School ID
   * @param {string} context.academicYearId - Academic Year ID
   * @param {string} context.semesterId - Semester ID
   * @returns {Promise<Array>} Array catatan
   */
  async fetchNotesByHomeroom(context) {
    const { userId, schoolId, academicYearId, semesterId } = context

    const response = await api.get('/v1/homeroom/journal', {
      headers: {
        'X-School-Id': schoolId,
        'X-Academic-Year-Id': academicYearId,
        'X-Semester-Id': semesterId,
        'X-User-Id': userId,
      },
    })

    const data = response.data || {}
    return mapHomeroomNotes(data.notes || [])
  }

  /**
   * Menyimpan catatan baru ke backend.
   * @param {Object} note - Objek catatan baru.
   * @param {Object} context - Operational context.
   * @returns {Promise<Object>} Catatan yang tersimpan dengan noteId.
   */
  async saveNote(note, context) {
    const { userId, schoolId, academicYearId, semesterId } = context

    const response = await api.post('/v1/homeroom/journal', note, {
      headers: {
        'X-School-Id': schoolId,
        'X-Academic-Year-Id': academicYearId,
        'X-Semester-Id': semesterId,
        'X-User-Id': userId,
      },
    })

    return mapHomeroomNote(response.data)
  }

  /**
   * Mengupdate catatan di backend.
   * @param {Object} note - Objek catatan yang sudah diupdate.
   * @param {Object} context - Operational context.
   * @returns {Promise<Object>} Catatan yang terupdate.
   */
  async updateNote(note, context) {
    const { userId, schoolId, academicYearId, semesterId } = context

    const response = await api.put(`/v1/homeroom/journal/${note.noteId}`, note, {
      headers: {
        'X-School-Id': schoolId,
        'X-Academic-Year-Id': academicYearId,
        'X-Semester-Id': semesterId,
        'X-User-Id': userId,
      },
    })

    return mapHomeroomNote(response.data)
  }

  /**
   * Menghapus catatan dari backend.
   * @param {string} noteId - ID catatan yang akan dihapus.
   * @param {Object} context - Operational context.
   * @returns {Promise<Object>} Hasil penghapusan.
   */
  async deleteNote(noteId, context) {
    const { userId, schoolId, academicYearId, semesterId } = context

    await api.delete(`/v1/homeroom/journal/${noteId}`, {
      headers: {
        'X-School-Id': schoolId,
        'X-Academic-Year-Id': academicYearId,
        'X-Semester-Id': semesterId,
        'X-User-Id': userId,
      },
    })

    return { success: true, deletedNoteId: noteId }
  }
}
