import { defineStore } from 'pinia'
import { createIncomingLetterService } from '../services/IncomingLetterService.js'
import { IncomingLetterMockAdapter } from '../adapters/mock/IncomingLetterMockAdapter.js'

// Initialize service with Mock Adapter
const service = createIncomingLetterService({ adapter: IncomingLetterMockAdapter })

export const useIncomingLetterStore = defineStore('incomingLetter', {
  state: () => ({
    list: [],
    currentItem: null,
    isLoading: false,
    error: null,
  }),
  actions: {
    async loadList(context, query) {
      this.isLoading = true
      this.error = null
      try {
        this.list = await service.listIncomingLetters(context, query)
      } catch (e) {
        this.error = e.message || 'Gagal memuat data'
        this.list = []
      } finally {
        this.isLoading = false
      }
    },

    async loadDetail(context, id) {
      this.isLoading = true
      this.error = null
      try {
        this.currentItem = await service.getIncomingLetterDetail(context, id)
      } catch (e) {
        this.error = e.message || 'Gagal memuat detail'
        this.currentItem = null
      } finally {
        this.isLoading = false
      }
    },

    async createItem(context, payload) {
      this.isLoading = true
      this.error = null
      try {
        const newItem = await service.createIncomingLetter(context, payload)
        this.list.unshift(newItem) // Add to top of list
        return newItem
      } catch (e) {
        this.error = e.message || 'Gagal menyimpan data'
        throw e
      } finally {
        this.isLoading = false
      }
    },

    async updateStatus(context, id, newStatus) {
      this.isLoading = true
      this.error = null
      try {
        const updated = await service.updateIncomingLetterStatus(context, id, newStatus)
        const index = this.list.findIndex((l) => l.id === id)
        if (index !== -1) this.list[index] = updated
        if (this.currentItem?.id === id) this.currentItem = updated
        return updated
      } catch (e) {
        this.error = e.message || 'Gagal mengubah status'
        throw e
      } finally {
        this.isLoading = false
      }
    },
  },
})
