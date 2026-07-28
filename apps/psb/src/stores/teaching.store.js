import { defineStore } from 'pinia'
import { TeachingEngine, SessionStatus, AuditAction } from '@/engine/TeachingEngine'
import { useContextStore } from '@/stores/context.store'

export const useTeachingStore = defineStore('teaching', {
  state: () => ({
    teachingService: null,

    // Jadwal guru
    schedule: [],

    // Sesi pada tanggal aktif
    sessions: [],

    // Sesi yang sedang dipilih/dikelola (aggregate root view)
    activeSessionId: null,

    // Teacher presence untuk sesi aktif
    teacherPresence: null,

    // Audit trail sesi aktif
    auditTrail: [],

    // Jurnal sesi aktif
    currentJournal: null,

    loading: {
      schedule: false,
      sessions: false,
      action: false, // start/complete/lock
      checkIn: false,
      journal: false,
    },
    error: null,
  }),

  getters: {
    activeSession: (state) => state.sessions.find((s) => s.id === state.activeSessionId) || null,

    isActiveSessionEditable: (state) => {
      const session = state.sessions.find((s) => s.id === state.activeSessionId)
      return session ? TeachingEngine.isEditable(session.status) : false
    },

    activeSessionStatus: (state) => {
      const session = state.sessions.find((s) => s.id === state.activeSessionId)
      return session?.status || null
    },

    // Cek apakah sesi bisa di-complete (perlu attendance submitted + journal filled)
    // eslint-disable-next-line no-unused-vars
    canCompleteSession: (state) => (attendanceSubmitted, journalFilled) => {
      return TeachingEngine.validateCompletion(attendanceSubmitted, journalFilled).valid
    },

    scheduleByDay: (state) => (day) => state.schedule.filter((s) => s.dayOfWeek === day),
  },

  actions: {
    setService(service) {
      this.teachingService = service
    },

    // --- Jadwal ---
    async loadSchedule() {
      const ctx = useContextStore()
      if (!ctx.isOperationalContextReady) throw new Error('Operational context belum dipilih')
      if (!this.teachingService) throw new Error('TeachingService belum diinisialisasi')

      this.loading.schedule = true
      this.error = null
      try {
        this.schedule = await this.teachingService.fetchTeacherSchedule(
          ctx.operational.academicYearId,
          ctx.operational.semesterId,
        )
      } catch (err) {
        this.error = err.message
        throw err
      } finally {
        this.loading.schedule = false
      }
    },

    // --- Sesi & Lifecycle ---
    async loadSessions(date) {
      if (!this.teachingService) throw new Error('TeachingService belum diinisialisasi')
      this.loading.sessions = true
      this.error = null
      try {
        this.sessions = await this.teachingService.fetchSessionsByDate(date)
        // Jika ada sesi yang sedang aktif (started/in_progress), jadikan aktif
        const ongoing = this.sessions.find((s) =>
          [SessionStatus.STARTED, SessionStatus.IN_PROGRESS].includes(s.status),
        )
        if (ongoing) {
          await this.setActiveSession(ongoing.id)
        } else if (this.sessions.length > 0) {
          await this.setActiveSession(this.sessions[0].id) // default ke pertama
        }
      } catch (err) {
        this.error = err.message
        throw err
      } finally {
        this.loading.sessions = false
      }
    },

    async setActiveSession(sessionId) {
      this.activeSessionId = sessionId
      // Muat data turunan: presence, audit, journal
      await Promise.all([this.loadTeacherPresence(), this.loadAuditTrail(), this.loadJournal()])
    },

    // --- State Machine Actions ---
    async startSession(sessionId) {
      if (!this.teachingService) throw new Error('Service tidak tersedia')
      this.loading.action = true
      try {
        const updated = await this.teachingService.startSession(sessionId)
        this._updateSessionInList(updated)
        this.activeSessionId = sessionId
        // Setelah start, seharusnya check-in teacher dilakukan bersamaan
        await this.loadTeacherPresence()
      } catch (err) {
        this.error = err.message
        throw err
      } finally {
        this.loading.action = false
      }
    },

    async beginProgress(sessionId) {
      if (!this.teachingService) throw new Error('Service tidak tersedia')
      this.loading.action = true
      try {
        const updated = await this.teachingService.beginProgress(sessionId)
        this._updateSessionInList(updated)
      } catch (err) {
        this.error = err.message
        throw err
      } finally {
        this.loading.action = false
      }
    },

    async completeSession(sessionId) {
      if (!this.teachingService) throw new Error('Service tidak tersedia')
      this.loading.action = true
      try {
        const updated = await this.teachingService.completeSession(sessionId)
        this._updateSessionInList(updated)
        await this.loadAuditTrail()
      } catch (err) {
        this.error = err.message
        throw err
      } finally {
        this.loading.action = false
      }
    },

    async lockSession(sessionId) {
      if (!this.teachingService) throw new Error('Service tidak tersedia')
      this.loading.action = true
      try {
        const updated = await this.teachingService.lockSession(sessionId)
        this._updateSessionInList(updated)
        await this.loadAuditTrail()
      } catch (err) {
        this.error = err.message
        throw err
      } finally {
        this.loading.action = false
      }
    },

    // --- Teacher Presence ---
    async checkIn(teacherId) {
      if (!this.activeSessionId || !this.teachingService) return
      this.loading.checkIn = true
      try {
        const presence = await this.teachingService.checkInTeacher(
          this.activeSessionId,
          teacherId,
          new Date().toISOString(),
        )
        this.teacherPresence = presence
        await this.loadAuditTrail()
      } catch (err) {
        this.error = err.message
        throw err
      } finally {
        this.loading.checkIn = false
      }
    },

    async checkOut() {
      if (!this.activeSessionId || !this.teachingService) return
      this.loading.checkIn = true
      try {
        const presence = await this.teachingService.checkOutTeacher(
          this.activeSessionId,
          new Date().toISOString(),
        )
        this.teacherPresence = presence
        await this.loadAuditTrail()
      } catch (err) {
        this.error = err.message
        throw err
      } finally {
        this.loading.checkIn = false
      }
    },

    async loadTeacherPresence() {
      if (!this.activeSessionId || !this.teachingService) return
      // Presence sekarang bagian dari session, tapi kita bisa ambil dari session yang sudah dimuat
      const session = this.sessions.find((s) => s.id === this.activeSessionId)
      this.teacherPresence = session?.teacherPresence || null
    },

    // --- Audit Trail ---
    async loadAuditTrail() {
      if (!this.activeSessionId || !this.teachingService) return
      try {
        this.auditTrail = await this.teachingService.getAuditTrail(this.activeSessionId)
      } catch {
        this.auditTrail = []
      }
    },

    async addAuditEntry(action, details) {
      if (!this.activeSessionId || !this.teachingService) return
      const entry = await this.teachingService.addAuditEntry(this.activeSessionId, action, details)
      this.auditTrail.push(entry)
    },

    // --- Jurnal ---
    async loadJournal() {
      if (!this.activeSessionId || !this.teachingService) return
      try {
        this.currentJournal = await this.teachingService.getJournalBySession(this.activeSessionId)
      } catch {
        this.currentJournal = null
      }
    },

    async saveJournal(data) {
      if (!this.activeSessionId || !this.teachingService) throw new Error('Tidak ada sesi aktif')
      this.loading.journal = true
      try {
        const saved = await this.teachingService.saveJournal(this.activeSessionId, data)
        this.currentJournal = saved
        await this.addAuditEntry(AuditAction.JOURNAL_SAVED, 'Jurnal disimpan')
        return saved
      } catch (err) {
        this.error = err.message
        throw err
      } finally {
        this.loading.journal = false
      }
    },

    // --- Helper ---
    _updateSessionInList(updated) {
      const idx = this.sessions.findIndex((s) => s.id === updated.id)
      if (idx !== -1) this.sessions[idx] = updated
    },

    reset() {
      this.schedule = []
      this.sessions = []
      this.activeSessionId = null
      this.teacherPresence = null
      this.auditTrail = []
      this.currentJournal = null
      this.error = null
    },
  },
})
