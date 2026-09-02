import { defineStore } from 'pinia'
import { assessmentService } from '../boot/services.js'

export const useAssessmentStore = defineStore('assessment', {
  state: () => ({
    assessment: null,
    result: null,
    loading: false,
    saving: false,
    error: null,
  }),

  getters: {
    hasAssessment: (state) => !!state.assessment,
    isDraft: (state) => state.assessment?.status === 'draft',
    isFinalized: (state) => state.assessment?.status === 'finalized',
  },

  actions: {
    async loadAssessment({ classId, subjectId, academicYearId, semesterId }) {
      console.log('😥cek...')
      this.loading = true
      this.error = null
      try {
        this.assessment = await assessmentService.loadAssessment({
          classId,
          subjectId,
          academicYearId,
          semesterId,
        })
        this.result = null
      } catch (err) {
        this.error = err.message
      } finally {
        this.loading = false
      }
    },

    async loadAssessmentById({ assessmentId }) {
      this.loading = true
      this.error = null
      try {
        this.assessment = await assessmentService.loadAssessmentById({ assessmentId })
        this.result = null
      } catch (err) {
        this.error = err.message
      } finally {
        this.loading = false
      }
    },

    async createAssessment({
      classId,
      className,
      subjectId,
      subjectName,
      date,
      components,
      schoolId,
      academicYearId,
      semesterId,
    }) {
      this.loading = true
      this.error = null
      try {
        this.assessment = await assessmentService.createAssessment({
          classId,
          className,
          subjectId,
          subjectName,
          date,
          components,
          schoolId,
          academicYearId,
          semesterId,
        })
        this.result = null
      } catch (err) {
        this.error = err.message
      } finally {
        this.loading = false
      }
    },

    async updateScore({ assessmentId, componentId, score }) {
      this.saving = true
      this.error = null
      try {
        this.assessment = await assessmentService.updateScore({
          assessmentId,
          componentId,
          score,
        })
      } catch (err) {
        this.error = err.message
      } finally {
        this.saving = false
      }
    },

    async finalizeAssessment({ assessmentId }) {
      this.loading = true
      this.error = null
      try {
        const { session, result } = await assessmentService.finalizeAssessment({ assessmentId })
        this.assessment = session
        this.result = result
      } catch (err) {
        this.error = err.message
      } finally {
        this.loading = false
      }
    },

    getPreviewResult() {
      if (!this.assessment) return null
      try {
        return assessmentService.calculatePreview(this.assessment)
      } catch (err) {
        console.log('😭', err)
        return null
      }
    },
  },
})
