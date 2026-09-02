import { defineStore } from 'pinia'
import { teachingService } from '../boot/services.js'

export const useTeachingStore = defineStore('teaching', {
  state: () => ({
    schedules: [],
    sessions: [],
    activeSession: null,
    loading: false,
    error: null,
  }),

  getters: {
    hasSchedules: (state) => state.schedules.length > 0,
    hasSessions: (state) => state.sessions.length > 0,
    getScheduledSessions: (state) => state.sessions.filter((s) => s.status === 'scheduled'),
    getActiveSessions: (state) => state.sessions.filter((s) => s.status === 'active'),
    getCompletedSessions: (state) => state.sessions.filter((s) => s.status === 'completed'),
  },

  actions: {
    async loadSchedules({ schoolId, academicYearId, semesterId, teacherId }) {
      this.loading = true
      this.error = null
      try {
        this.schedules = await teachingService.loadSchedulesByContext({
          schoolId,
          academicYearId,
          semesterId,
          teacherId,
        })
      console.log('😥',this.schedules)

      } catch (err) {
        this.error = err.message
      } finally {
        this.loading = false
      }
    },

    async loadSessions({ schoolId, academicYearId, semesterId, teacherId, date }) {
      this.loading = true
      this.error = null
      try {
        this.sessions = await teachingService.loadSessionsByDate({
          schoolId,
          academicYearId,
          semesterId,
          teacherId,
          date,
        })
      } catch (err) {
        this.error = err.message
      } finally {
        this.loading = false
      }
    },

    async createSession({ scheduleId, date, schoolId, academicYearId, semesterId, teacherId }) {
      this.loading = true
      this.error = null
      try {
        const session = await teachingService.createSession({
          scheduleId,
          date,
          schoolId,
          academicYearId,
          semesterId,
          teacherId,
        })
        this.sessions.push(session)
        return session
      } catch (err) {
        this.error = err.message
        throw err
      } finally {
        this.loading = false
      }
    },

    async startSession({ sessionId, teacherPresence, schoolId, academicYearId, semesterId }) {
      this.loading = true
      this.error = null
      try {
        const updatedSession = await teachingService.startSession({
          sessionId,
          teacherPresence,
          schoolId,
          academicYearId,
          semesterId,
        })
        this.updateSessionInList(updatedSession)
        this.activeSession = updatedSession
        return updatedSession
      } catch (err) {
        this.error = err.message
        throw err
      } finally {
        this.loading = false
      }
    },

    async endSession({ sessionId, notes, schoolId, academicYearId, semesterId }) {
      this.loading = true
      this.error = null
      try {
        const updatedSession = await teachingService.endSession({
          sessionId,
          notes,
          schoolId,
          academicYearId,
          semesterId,
        })
        this.updateSessionInList(updatedSession)
        this.activeSession = null
        return updatedSession
      } catch (err) {
        this.error = err.message
        throw err
      } finally {
        this.loading = false
      }
    },

    async cancelSession({ sessionId, reason, schoolId, academicYearId, semesterId }) {
      this.loading = true
      this.error = null
      try {
        const updatedSession = await teachingService.cancelSession({
          sessionId,
          reason,
          schoolId,
          academicYearId,
          semesterId,
        })
        this.updateSessionInList(updatedSession)
        if (this.activeSession && this.activeSession.id === sessionId) {
          this.activeSession = null
        }
        return updatedSession
      } catch (err) {
        this.error = err.message
        throw err
      } finally {
        this.loading = false
      }
    },

    updateSessionInList(updatedSession) {
      const index = this.sessions.findIndex((s) => s.id === updatedSession.id)
      if (index !== -1) {
        this.sessions[index] = updatedSession
      }
    },

    setActiveSession(session) {
      this.activeSession = session
    },
  },
})
