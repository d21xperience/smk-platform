// FILE: src/composables/useHomeroomJournal.js
// STATUS: NEW
// STATUS IMPLEMENTASI: COMPLETE

import { computed } from 'vue'
import { useHomeroomJournalStore } from '../stores/homeroomJournalStore.js'
import { useContext } from './useContext.js'
import { useAuth } from './useAuth.js'

/**
 * Composable Facade untuk UI Homeroom Journal.
 * Composable bertanggung jawab:
 * - Mengambil context dari useContext()
 * - Mengambil auth dari useAuth()
 * - Meneruskan parameter eksplisit ke Store
 * Page yang memutuskan kapan memanggil loadNotes(), createNote(), dll.
 */
export function useHomeroomJournal() {
  const store = useHomeroomJournalStore()
  const { currentContext } = useContext()
  const { currentUser } = useAuth()

  const notes = computed(() => store.notes)
  const summary = computed(() => store.summary)
  const loading = computed(() => store.loading)
  const creating = computed(() => store.creating)
  const submitting = computed(() => store.submitting)
  const deleting = computed(() => store.deleting)
  const error = computed(() => store.error)
  const hasNotes = computed(() => store.hasNotes)
  const hasSummary = computed(() => store.hasSummary)
  const draftNotes = computed(() => store.draftNotes)
  const submittedNotes = computed(() => store.submittedNotes)
  const hasDraftNotes = computed(() => store.hasDraftNotes)

  /**
   * Memuat daftar catatan wali kelas.
   */
  async function loadNotes() {
    if (!currentContext.value || !currentUser.value) return
    await store.loadNotes({
      userId: currentUser.value.id,
      schoolId: currentContext.value.schoolId,
      academicYearId: currentContext.value.academicYearId,
      semesterId: currentContext.value.semesterId,
    })
  }

  /**
   * Membuat catatan baru.
   * @param {Object} noteData - Data catatan dari form
   * @returns {Promise<Object>} Catatan yang tersimpan
   */
  async function createNote(noteData) {
    if (!currentContext.value || !currentUser.value) return null
    return store.createNote({
      userId: currentUser.value.id,
      schoolId: currentContext.value.schoolId,
      academicYearId: currentContext.value.academicYearId,
      semesterId: currentContext.value.semesterId,
      noteData,
    })
  }

  /**
   * Meng-submit catatan (finalisasi).
   * @param {Object} note - Catatan yang akan di-submit
   * @returns {Promise<Object>} Catatan yang terupdate
   */
  async function submitNote(note) {
    if (!currentContext.value || !currentUser.value) return null
    return store.submitNote({
      userId: currentUser.value.id,
      schoolId: currentContext.value.schoolId,
      academicYearId: currentContext.value.academicYearId,
      semesterId: currentContext.value.semesterId,
      note,
    })
  }

  /**
   * Menghapus catatan draft.
   * @param {Object} note - Catatan yang akan dihapus
   * @returns {Promise<Object>} Hasil penghapusan
   */
  async function deleteNote(note) {
    if (!currentContext.value || !currentUser.value) return null
    return store.deleteNote({
      userId: currentUser.value.id,
      schoolId: currentContext.value.schoolId,
      academicYearId: currentContext.value.academicYearId,
      semesterId: currentContext.value.semesterId,
      note,
    })
  }

  /**
   * Mengosongkan data store.
   */
  function clearData() {
    store.clearData()
  }

  return {
    notes,
    summary,
    loading,
    creating,
    submitting,
    deleting,
    error,
    hasNotes,
    hasSummary,
    draftNotes,
    submittedNotes,
    hasDraftNotes,
    loadNotes,
    createNote,
    submitNote,
    deleteNote,
    clearData,
  }
}
