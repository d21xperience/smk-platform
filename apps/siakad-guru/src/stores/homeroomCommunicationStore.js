// FILE: src/stores/homeroomCommunicationStore.js
// STATUS: NEW
// STATUS IMPLEMENTASI: COMPLETE

import { defineStore } from 'pinia'
import { homeroomCommunicationService } from '../boot/services.js'

/**
 * Homeroom Communication Store
 * Mengelola state log komunikasi wali kelas dengan orang tua.
 * Pola mengikuti homeroomBillingStore:
 * - Tidak import store lain
 * - Action menerima parameter eksplisit dari composable
 * - Store memanggil Application Service, bukan Axios
 * Tidak ada operasi delete — log komunikasi tidak dihapus.
 */
export const useHomeroomCommunicationStore = defineStore('homeroomCommunication', {
  state: () => ({
    communications: [],
    summary: null,
    loading: false,
    creating: false,
    updating: false,
    error: null,
  }),

  getters: {
    hasCommunications: (state) => state.communications.length > 0,
    hasSummary: (state) => !!state.summary,
    pendingCommunications: (state) =>
      state.communications.filter((c) => c.status === 'PENDING_FOLLOW_UP'),
    followedUpCommunications: (state) =>
      state.communications.filter((c) => c.status === 'FOLLOWED_UP'),
    closedCommunications: (state) => state.communications.filter((c) => c.status === 'CLOSED'),
    hasPendingCommunications: (state) =>
      state.communications.some((c) => c.status === 'PENDING_FOLLOW_UP'),
  },

  actions: {
    /**
     * Memuat daftar log komunikasi wali kelas.
     * @param {Object} params - Parameter dari composable
     */
    async loadCommunications({ userId, schoolId, academicYearId, semesterId }) {
      this.loading = true
      this.error = null
      try {
        const result = await homeroomCommunicationService.getCommunications({
          userId,
          schoolId,
          academicYearId,
          semesterId,
        })
        this.communications = result.communications
        this.summary = result.summary
      } catch (err) {
        this.error = err.message
        this.communications = []
        this.summary = null
      } finally {
        this.loading = false
      }
    },

    /**
     * Membuat log komunikasi baru.
     * Setelah berhasil, communication ditambahkan ke array dan summary di-rebuild.
     * @param {Object} params - Parameter dari composable
     * @param {Object} params.commData - Data log komunikasi dari form
     */
    async createCommunication({ userId, schoolId, academicYearId, semesterId, commData }) {
      this.creating = true
      this.error = null
      try {
        const savedCommunication = await homeroomCommunicationService.createCommunication(
          { userId, schoolId, academicYearId, semesterId },
          commData,
        )
        this.communications.unshift(savedCommunication)
        this.summary = homeroomCommunicationService.buildSummary(this.communications)
        return savedCommunication
      } catch (err) {
        this.error = err.message
        throw err
      } finally {
        this.creating = false
      }
    },

    /**
     * Mengupdate status tindak lanjut log komunikasi.
     * Setelah berhasil, communication di array di-update dan summary di-rebuild.
     * @param {Object} params - Parameter dari composable
     * @param {Object} params.communication - Log komunikasi yang akan diupdate
     * @param {string} params.newStatus - Status baru
     * @param {string} params.followUpNote - Catatan tindak lanjut (opsional)
     */
    async updateStatus({
      userId,
      schoolId,
      academicYearId,
      semesterId,
      communication,
      newStatus,
      followUpNote,
    }) {
      this.updating = true
      this.error = null
      try {
        const updatedCommunication = await homeroomCommunicationService.updateStatus(
          { userId, schoolId, academicYearId, semesterId },
          communication,
          newStatus,
          followUpNote,
        )
        const commIndex = this.communications.findIndex(
          (c) => c.communicationId === updatedCommunication.communicationId,
        )
        if (commIndex !== -1) {
          this.communications[commIndex] = updatedCommunication
        }
        this.summary = homeroomCommunicationService.buildSummary(this.communications)
        return updatedCommunication
      } catch (err) {
        this.error = err.message
        throw err
      } finally {
        this.updating = false
      }
    },

    /**
     * Mengosongkan state store.
     */
    clearData() {
      this.communications = []
      this.summary = null
      this.loading = false
      this.creating = false
      this.updating = false
      this.error = null
    },
  },
})
