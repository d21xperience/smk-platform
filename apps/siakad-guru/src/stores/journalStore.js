import { defineStore } from 'pinia'
import { journalService } from '../boot/services.js'

export const useJournalStore = defineStore('journal', {
  state: () => ({
    journal: null,
    loading: false,
    saving: false,
    error: null
  }),

  getters: {
    hasJournal: (state) => !!state.journal,
    isDraft: (state) => state.journal?.status === 'draft',
    isSubmitted: (state) => state.journal?.status === 'submitted'
  },

  actions: {
    async loadJournal({ teachingSessionId }) {
      this.loading = true
      this.error = null
      try {
        this.journal = await journalService.loadJournal({ teachingSessionId })
      } catch (err) {
        this.error = err.message
      } finally {
        this.loading = false
      }
    },

    async createJournal({
      teachingSessionId,
      classId,
      className,
      subjectId,
      subjectName,
      date,
      schoolId,
      academicYearId,
      semesterId
    }) {
      this.loading = true
      this.error = null
      try {
        this.journal = await journalService.createJournal({
          teachingSessionId,
          classId,
          className,
          subjectId,
          subjectName,
          date,
          schoolId,
          academicYearId,
          semesterId
        })
      } catch (err) {
        this.error = err.message
      } finally {
        this.loading = false
      }
    },

    async updateJournal({ journalId, material, activities, reflection }) {
      this.saving = true
      this.error = null
      try {
        this.journal = await journalService.updateJournal({
          journalId,
          material,
          activities,
          reflection
        })
      } catch (err) {
        this.error = err.message
      } finally {
        this.saving = false
      }
    },

    async submitJournal({ journalId }) {
      this.loading = true
      this.error = null
      try {
        this.journal = await journalService.submitJournal({ journalId })
      } catch (err) {
        this.error = err.message
      } finally {
        this.loading = false
      }
    }
  }
})
