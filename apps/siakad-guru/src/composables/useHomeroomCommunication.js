// FILE: src/composables/useHomeroomCommunication.js
// STATUS: NEW
// STATUS IMPLEMENTASI: COMPLETE

import { computed } from 'vue'
import { useHomeroomCommunicationStore } from '../stores/homeroomCommunicationStore.js'
import { useContext } from './useContext.js'
import { useAuth } from './useAuth.js'

/**
 * Composable Facade untuk UI Homeroom Communication.
 * Composable bertanggung jawab:
 * - Mengambil context dari useContext()
 * - Mengambil auth dari useAuth()
 * - Meneruskan parameter eksplisit ke Store
 * Page yang memutuskan kapan memanggil loadCommunications(), createCommunication(), dll.
 */
export function useHomeroomCommunication() {
  const store = useHomeroomCommunicationStore()
  const { currentContext } = useContext()
  const { currentUser } = useAuth()

  const communications = computed(() => store.communications)
  const summary = computed(() => store.summary)
  const loading = computed(() => store.loading)
  const creating = computed(() => store.creating)
  const updating = computed(() => store.updating)
  const error = computed(() => store.error)
  const hasCommunications = computed(() => store.hasCommunications)
  const hasSummary = computed(() => store.hasSummary)
  const pendingCommunications = computed(() => store.pendingCommunications)
  const followedUpCommunications = computed(() => store.followedUpCommunications)
  const closedCommunications = computed(() => store.closedCommunications)
  const hasPendingCommunications = computed(() => store.hasPendingCommunications)

  /**
   * Memuat daftar log komunikasi wali kelas.
   */
  async function loadCommunications() {
    if (!currentContext.value || !currentUser.value) return
    await store.loadCommunications({
      userId: currentUser.value.id,
      schoolId: currentContext.value.schoolId,
      academicYearId: currentContext.value.academicYearId,
      semesterId: currentContext.value.semesterId,
    })
  }

  /**
   * Membuat log komunikasi baru.
   * @param {Object} commData - Data log komunikasi dari form
   * @returns {Promise<Object>} Log komunikasi yang tersimpan
   */
  async function createCommunication(commData) {
    if (!currentContext.value || !currentUser.value) return null
    return store.createCommunication({
      userId: currentUser.value.id,
      schoolId: currentContext.value.schoolId,
      academicYearId: currentContext.value.academicYearId,
      semesterId: currentContext.value.semesterId,
      commData,
    })
  }

  /**
   * Mengupdate status tindak lanjut log komunikasi.
   * @param {Object} communication - Log komunikasi yang akan diupdate
   * @param {string} newStatus - Status baru (FOLLOWED_UP atau CLOSED)
   * @param {string} followUpNote - Catatan tindak lanjut (opsional)
   * @returns {Promise<Object>} Log komunikasi yang terupdate
   */
  async function updateStatus(communication, newStatus, followUpNote) {
    if (!currentContext.value || !currentUser.value) return null
    return store.updateStatus({
      userId: currentUser.value.id,
      schoolId: currentContext.value.schoolId,
      academicYearId: currentContext.value.academicYearId,
      semesterId: currentContext.value.semesterId,
      communication,
      newStatus,
      followUpNote,
    })
  }

  /**
   * Mengosongkan data store.
   */
  function clearData() {
    store.clearData()
  }

  return {
    communications,
    summary,
    loading,
    creating,
    updating,
    error,
    hasCommunications,
    hasSummary,
    pendingCommunications,
    followedUpCommunications,
    closedCommunications,
    hasPendingCommunications,
    loadCommunications,
    createCommunication,
    updateStatus,
    clearData,
  }
}
