import { defineStore } from 'pinia'
import { AssessmentService } from '@/services/AssessmentService'
import { AssessmentSession } from '@/domain/assessment/models/AssessmentSession'

export const useAssessmentStore = defineStore('assessment', {
  state: () => ({
    currentSession: null,
    isLoading: false,
    error: null,
  }),
  getters: {
    isDraft: (state) => state.currentSession?.isDraft() ?? false,
    isFinalized: (state) => state.currentSession?.isFinalized() ?? false,
    components: (state) => state.currentSession?.components || [],
    results: (state) => state.currentSession?.results || [],
  },
  actions: {
    async loadOrCreateSession(params) {
      this.isLoading = true
      try {
        const data = await AssessmentService.loadOrCreateSession(params)
        this.currentSession = new AssessmentSession(data)
        return this.currentSession
      } catch (err) {
        this.error = err.message
        throw err
      } finally {
        this.isLoading = false
      }
    },
    async calculateGrade(studentScores) {
      if (!this.currentSession) throw new Error('Tidak ada sesi penilaian.')
      this.isLoading = true
      try {
        const data = await AssessmentService.calculateGrade(this.currentSession.id, studentScores)
        this.currentSession = new AssessmentSession(data)
        return this.currentSession
      } catch (err) {
        this.error = err.message
        throw err
      } finally {
        this.isLoading = false
      }
    },
    async finalize() {
      if (!this.currentSession) throw new Error('Tidak ada sesi penilaian.')
      this.isLoading = true
      try {
        const data = await AssessmentService.finalizeAssessment(this.currentSession.id)
        this.currentSession = new AssessmentSession(data)
        return this.currentSession
      } catch (err) {
        this.error = err.message
        throw err
      } finally {
        this.isLoading = false
      }
    },
    clear() {
      this.currentSession = null
      this.error = null
    },
  },
})
