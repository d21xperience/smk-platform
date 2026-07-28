// @/stores/assessment.store.js
import { defineStore } from 'pinia'
import { AssessmentEngine, DefaultPredicates } from '@/engine/AssessmentEngine'
import { getEventDispatcher } from '@/boot/events'
import { DomainEvent } from '@/events/DomainEvent'
import { DomainEventType } from '@/events/DomainEventType'

export const useAssessmentStore = defineStore('assessment', {
  state: () => ({
    assessmentService: null,
    session: null, // AssessmentSession
    components: [], // AssessmentComponent[]
    finalResult: null, // AssessmentFinalResult
    selectedComponentId: null, // komponen yang sedang diedit skornya

    loading: {
      init: false,
      saving: false,
      submitting: false,
    },
    error: null,
  }),

  getters: {
    // Komponen yang sedang dipilih
    selectedComponent: (state) =>
      state.components.find((c) => c.id === state.selectedComponentId) || null,

    // Semua komponen sudah lengkap skornya?
    isAllComponentsScored: (state) => {
      if (!state.session || state.components.length === 0) return false
      // Gunakan Engine untuk validasi submit (mengecek komponen wajib)
      const validation = AssessmentEngine.validateForSubmit({
        ...state.session,
        components: state.components,
      })
      return validation.valid
    },

    // Hasil kalkulasi per siswa (real-time preview)
    previewGrades: (state) => {
      if (!state.session || state.components.length === 0) return []
      return AssessmentEngine.generateGrades(
        { ...state.session, components: state.components },
        DefaultPredicates,
      )
    },

    // Statistik kelas (rata-rata, tertinggi, terendah)
    classStats: (state) => {
      const grades = AssessmentEngine.generateGrades(
        { ...state.session, components: state.components },
        DefaultPredicates,
      )
      return AssessmentEngine.calculateClassStats(grades)
    },

    // Status sesi saat ini
    canEdit: (state) => state.session?.status === 'draft',
    isSubmitted: (state) => state.session?.status === 'submitted',
  },

  actions: {
    setService(service) {
      this.assessmentService = service
    },

    // --- Inisialisasi ---
    async loadOrCreateSession(teachingSessionId, options) {
      if (!this.assessmentService) throw new Error('Service tidak tersedia')
      this.loading.init = true
      this.error = null
      try {
        const session = await this.assessmentService.fetchOrCreateSession(
          teachingSessionId,
          options,
        )
        this.session = session
        this.components = session.components || []
        if (this.components.length > 0 && !this.selectedComponentId) {
          this.selectedComponentId = this.components[0].id
        }
      } catch (err) {
        this.error = err.message
        throw err
      } finally {
        this.loading.init = false
      }
    },

    // --- Manajemen Komponen ---
    selectComponent(componentId) {
      this.selectedComponentId = componentId
    },

    addComponent(component) {
      this.components.push(component)
    },

    removeComponent(componentId) {
      this.components = this.components.filter((c) => c.id !== componentId)
      if (this.selectedComponentId === componentId) {
        this.selectedComponentId = this.components[0]?.id || null
      }
    },

    // --- Skor ---
    updateScore(componentId, studentId, score, notes = '') {
      const component = this.components.find((c) => c.id === componentId)
      if (!component) return
      let entry = component.scores.find((s) => s.studentId === studentId)
      if (entry) {
        entry.score = score
        entry.notes = notes
      } else {
        // Ambil nama dari entri lain atau placeholder
        const existing = component.scores[0]
        component.scores.push(
          AssessmentEngine.createScoreEntry({
            componentId,
            studentId,
            studentName: existing?.studentName || `Siswa ${studentId}`,
            score,
            notes,
          }),
        )
      }
    },

    // Batch update skor untuk satu komponen
    updateScores(componentId, scoresArray) {
      const component = this.components.find((c) => c.id === componentId)
      if (!component) return
      component.scores = scoresArray.map((s) => AssessmentEngine.createScoreEntry(s))
    },

    // --- Draft / Simpan ---
    async saveDraft() {
      if (!this.assessmentService || !this.session) return
      this.loading.saving = true
      this.error = null
      try {
        const updated = await this.assessmentService.saveDraft({
          ...this.session,
          components: this.components,
        })
        this.session = updated
        // Dispatch event
        this._dispatch(DomainEventType.ASSESSMENT_DRAFT_SAVED, {
          sessionId: updated.id,
          className: updated.className,
          subject: updated.subject,
          savedAt: new Date().toISOString(),
        })
      } catch (err) {
        this.error = err.message
        throw err
      } finally {
        this.loading.saving = false
      }
    },

    // --- Submit ---
    async submitAssessment() {
      if (!this.assessmentService || !this.session) throw new Error('Tidak ada sesi')
      // Validasi via Engine
      const validation = AssessmentEngine.validateForSubmit({
        ...this.session,
        components: this.components,
      })
      if (!validation.valid) {
        throw new Error(validation.errors.join(' '))
      }
      this.loading.submitting = true
      this.error = null
      try {
        const updated = await this.assessmentService.submit(this.session.id)
        this.session = updated
        // Generate final result
        this.finalResult = {
          sessionId: updated.id,
          className: updated.className,
          subject: updated.subject,
          grades: this.previewGrades,
          generatedAt: new Date().toISOString(),
        }
        // Dispatch event
        this._dispatch(DomainEventType.ASSESSMENT_SUBMITTED, {
          sessionId: updated.id,
          className: updated.className,
          subject: updated.subject,
          submittedAt: new Date().toISOString(),
          stats: this.classStats,
        })
      } catch (err) {
        this.error = err.message
        throw err
      } finally {
        this.loading.submitting = false
      }
    },

    // --- Helper ---
    _dispatch(type, payload) {
      const dispatcher = getEventDispatcher()
      dispatcher.dispatch(new DomainEvent(type, payload))
    },

    reset() {
      this.session = null
      this.components = []
      this.finalResult = null
      this.selectedComponentId = null
      this.loading = { init: false, saving: false, submitting: false }
      this.error = null
    },
  },
})
