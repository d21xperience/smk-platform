import { defineStore } from 'pinia'
import { attendanceService } from '../boot/services.js'

export const useAttendanceStore = defineStore('attendance', {
  state: () => ({
    students: [],
    session: null,
    loading: false,
    saving: false,
    error: null,
  }),

  getters: {
    hasSession: (state) => !!state.session,
    isDraft: (state) => state.session?.status === 'draft',
    isSubmitted: (state) => state.session?.status === 'submitted',
    getSummary: (state) => {
      if (!state.session) return null
      return {
        total: state.session.getTotalStudents(),
        present: state.session.getPresentCount(),
        sick: state.session.getSickCount(),
        permission: state.session.getPermissionCount(),
        absent: state.session.getAbsentCount(),
      }
    },
  },

  actions: {
    async loadSession({ teachingSessionId }) {
      this.loading = true
      this.error = null
      try {
        this.session = await attendanceService.loadAttendanceSession({ teachingSessionId })
        this.students = this.session ? [...this.session.records] : []
      } catch (err) {
        this.error = err.message
      } finally {
        this.loading = false
      }
    },

    async createSession({
      teachingSessionId,
      classId,
      className,
      subjectId,
      subjectName,
      date,
      schoolId,
      academicYearId,
      semesterId,
    }) {
      this.loading = true
      this.error = null
      try {
        this.session = await attendanceService.createAttendanceSession({
          teachingSessionId,
          classId,
          className,
          subjectId,
          subjectName,
          date,
          schoolId,
          academicYearId,
          semesterId,
        })
        this.students = [...this.session.records]
      } catch (err) {
        this.error = err.message
      } finally {
        this.loading = false
      }
    },

    async saveDraft({ sessionId, records }) {
      this.saving = true
      this.error = null
      try {
        this.session = await attendanceService.saveDraft({ sessionId, records })
        this.students = [...this.session.records]
      } catch (err) {
        this.error = err.message
      } finally {
        this.saving = false
      }
    },

    async submitAttendance({ sessionId }) {
      this.loading = true
      this.error = null
      try {
        this.session = await attendanceService.submitAttendance({ sessionId })
        this.students = [...this.session.records]
      } catch (err) {
        this.error = err.message
      } finally {
        this.loading = false
      }
    },

    updateRecordStatus(studentId, status) {
      const record = this.students.find((r) => r.studentId === studentId)
      if (record) {
        record.status = status
      }
    },

    updateRecordNote(studentId, note) {
      const record = this.students.find((r) => r.studentId === studentId)
      if (record) {
        record.note = note
      }
    },
  },
})
