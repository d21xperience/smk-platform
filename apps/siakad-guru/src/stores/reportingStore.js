import { defineStore } from 'pinia'
import { reportingService } from '../boot/services.js'

export const useReportingStore = defineStore('reporting', {
  state: () => ({
    dashboard: null,
    loading: false,
    error: null,
  }),

  getters: {
    hasDashboard: (state) => !!state.dashboard,
  },

  actions: {
    async loadDashboard({ schoolId, academicYearId, semesterId, teacherId }) {
      this.loading = true
      this.error = null
      try {
        this.dashboard = await reportingService.loadDashboard({
          schoolId,
          academicYearId,
          semesterId,
          teacherId,
        })
      } catch (err) {
        this.error = err.message
      } finally {
        this.loading = false
      }
    },
  },
})
