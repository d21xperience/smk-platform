import { defineStore } from 'pinia'
import { createDispositionService } from '../services/DispositionService.js'
import { DispositionMockAdapter } from '../adapters/mock/DispositionMockAdapter.js'

const service = createDispositionService({ adapter: DispositionMockAdapter })

export const useDispositionStore = defineStore('disposition', {
  state: () => ({
    pendingIncomingLetters: [], // Kotak masuk disposisi untuk Kepsek
    myTasks: [], // Tugas disposisi untuk Staf/Waka/Guru
    stafOptions: [],
    isLoading: false,
    error: null,
  }),
  actions: {
    async loadPendingIncomingLetters(context) {
      this.isLoading = true
      this.error = null
      try {
        this.pendingIncomingLetters = await service.getPendingDispositions(context)
      } catch (e) {
        this.error = e.message || 'Gagal memuat kotak masuk'
      } finally {
        this.isLoading = false
      }
    },

    async loadMyTasks(context, userId) {
      this.isLoading = true
      this.error = null
      try {
        this.myTasks = await service.getMyTasks(context, userId)
      } catch (e) {
        this.error = e.message || 'Gagal memuat tugas'
      } finally {
        this.isLoading = false
      }
    },

    async createDisposition(context, payload) {
      this.isLoading = true
      this.error = null
      try {
        await service.createDisposition(context, payload)
        // Reload kotak masuk setelah disposisi dibuat
        await this.loadPendingIncomingLetters(context)
      } catch (e) {
        this.error = e.message || 'Gagal mengirim disposisi'
        throw e
      } finally {
        this.isLoading = false
      }
    },

    async submitFollowUp(context, dispositionId, payload) {
      this.isLoading = true
      this.error = null
      try {
        await service.submitFollowUp(context, dispositionId, payload)
        // Reload tugas setelah tindak lanjut disimpan
        await this.loadMyTasks(context, context.userId)
      } catch (e) {
        this.error = e.message || 'Gagal menyimpan tindak lanjut'
        throw e
      } finally {
        this.isLoading = false
      }
    },

    async loadStafOptions(context) {
      try {
        this.stafOptions = await service.getStafOptions(context)
      } catch (e) {
        console.log(e)
        this.stafOptions = []
      }
    },
  },
})
