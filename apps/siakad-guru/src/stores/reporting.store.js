import { defineStore } from 'pinia'
import { ReportingService } from '@/services/ReportingService'
import { DashboardSummary } from '@/domain/reporting/models/DashboardSummary'

export const useReportingStore = defineStore('reporting', {
  state: () => ({
    dashboardStats: null, // instance of DashboardSummary
    currentReport: null, // ReportResult
    isLoading: false,
    error: null,
  }),
  getters: {
    hasDashboard: (state) => state.dashboardStats !== null,
  },
  actions: {
    async loadDashboardStats() {
      this.isLoading = true
      this.error = null
      try {
        const data = await ReportingService.getDashboardStats()
        this.dashboardStats = new DashboardSummary(data)
        return this.dashboardStats
      } catch (err) {
        this.error = err.message
        throw err
      } finally {
        this.isLoading = false
      }
    },

    async generateClassReport(classId) {
      this.isLoading = true
      this.error = null
      try {
        const data = await ReportingService.generateClassReport(classId)
        this.currentReport = data
        return data
      } catch (err) {
        this.error = err.message
        throw err
      } finally {
        this.isLoading = false
      }
    },

    async exportReport(format, filename) {
      if (!this.currentReport) {
        throw new Error('Tidak ada laporan yang siap diekspor.')
      }
      await ReportingService.exportReport(this.currentReport, format, filename)
    },

    clear() {
      this.dashboardStats = null
      this.currentReport = null
      this.error = null
    },
  },
})
