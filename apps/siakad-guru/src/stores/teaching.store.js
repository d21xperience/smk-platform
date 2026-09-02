import { defineStore } from 'pinia'
import { teachingService } from '@/services/TeachingService'
import { TeachingSession } from '@/domain/teaching/models/TeachingSession'

export const useTeachingStore = defineStore('teaching', {
  state: () => ({
    todaySchedule: [], // Array dari TeachingSession
    activeSession: null, // TeachingSession yang sedang aktif
    isLoading: false,
    error: null,
    journalLoading: false,
    journalError: null,
    currentJournal: null, // TeachingJournal
  }),

  getters: {
    hasActiveSession: (state) => state.activeSession !== null,
    currentClassId: (state) => state.activeSession?.classId || null,
    currentSubjectId: (state) => state.activeSession?.subjectId || null,
    isJournalFinalized: (state) => state.currentJournal?.isFinalized() ?? false,
    isJournalDraft: (state) => state.currentJournal?.isDraft() ?? false,
  },

  actions: {
    async fetchTodaySchedule(teacherId, date) {
      this.isLoading = true
      this.error = null
      try {
        const schedule = await teachingService.getTodaySchedule(teacherId, date)
        this.todaySchedule = schedule.map((s) => new TeachingSession(s))

        // Cari sesi yang statusnya 'active' secara otomatis
        const active = this.todaySchedule.find((s) => s.isActive())
        if (active) {
          this.activeSession = active
        }
        return this.todaySchedule
      } catch (err) {
        this.error = err.message
        throw err
      } finally {
        this.isLoading = false
      }
    },

    async startSession(sessionId) {
      this.isLoading = true
      this.error = null
      try {
        const session = await teachingService.startSession(sessionId)
        this.activeSession = new TeachingSession(session)
        // Update di todaySchedule
        const index = this.todaySchedule.findIndex((s) => s.id === session.id)
        if (index !== -1) {
          this.todaySchedule[index] = new TeachingSession(session)
        }
        return this.activeSession
      } catch (err) {
        this.error = err.message
        throw err
      } finally {
        this.isLoading = false
      }
    },

    async endSession(sessionId) {

      this.isLoading = true
      this.error = null
      try {
        const session = await teachingService.endSession(sessionId)
        this.activeSession = null
        const index = this.todaySchedule.findIndex((s) => s.id === session.id)
        if (index !== -1) {
          this.todaySchedule[index] = new TeachingSession(session)
        }
        console.log('return',session)
        return session
      } catch (err) {
        this.error = err.message
        throw err
      } finally {
        this.isLoading = false
      }
    },

    async loadOrCreateJournal(sessionId) {
      this.journalLoading = true
      this.journalError = null
      try {
        const data = await teachingService.loadOrCreateJournal(sessionId)
        this.currentJournal = data
        return this.currentJournal
      } catch (err) {
        this.journalError = err.message
        throw err
      } finally {
        this.journalLoading = false
      }
    },

    async saveJournalDraft(payload) {
      if (!this.currentJournal) {
        throw new Error('Tidak ada jurnal aktif.')
      }
      this.journalLoading = true
      this.journalError = null
      try {
        const data = await teachingService.saveJournalDraft(this.currentJournal.id, payload)
        this.currentJournal = data
        return this.currentJournal
      } catch (err) {
        this.journalError = err.message
        throw err
      } finally {
        this.journalLoading = false
      }
    },

    async finalizeJournal() {
      if (!this.currentJournal) {
        throw new Error('Tidak ada jurnal aktif.')
      }
      this.journalLoading = true
      this.journalError = null
      try {
        const data = await teachingService.finalizeJournal(this.currentJournal.id)
        this.currentJournal = data
        return this.currentJournal
      } catch (err) {
        this.journalError = err.message
        throw err
      } finally {
        this.journalLoading = false
      }
    },

    async markAttendanceSubmitted(sessionId) {
      this.journalLoading = true
      try {
        const data = await teachingService.markAttendanceSubmitted(sessionId)
        if (this.currentJournal && this.currentJournal.sessionId === sessionId) {
          this.currentJournal.attendanceSubmitted = true
        }
        return data
      } catch (err) {
        this.journalError = err.message
        throw err
      } finally {
        this.journalLoading = false
      }
    },

    clearActiveSession() {
      this.activeSession = null
    },

    clearJournal() {
      this.currentJournal = null
      this.journalError = null
    },

    clear() {
      this.todaySchedule = []
      this.activeSession = null
      this.currentJournal = null
      this.isLoading = false
      this.error = null
      this.journalLoading = false
      this.journalError = null
    },
  },
})
