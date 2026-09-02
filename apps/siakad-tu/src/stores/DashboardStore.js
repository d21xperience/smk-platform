import { defineStore } from 'pinia'
import { createDashboardService } from '../services/DashboardService.js'
import { DashboardMockAdapter } from '../adapters/mock/DashboardMockAdapter.js'

const service = createDashboardService({ adapter: DashboardMockAdapter })

export const useDashboardStore = defineStore('dashboard', {
  state: () => ({
    stats: { masuk: 0, keluar: 0, disposisiPending: 0 },
    trend: { categories: [], seriesMasuk: [], seriesKeluar: [] },
    isLoading: false,
  }),
  actions: {
    async loadStats(context) {
      this.isLoading = true
      try {
        this.stats = await service.getStats(context)
      } catch (e) {
        console.error('Failed to load stats', e)
      } finally {
        this.isLoading = false
      }
    },

    async loadTrend(context) {
      try {
        this.trend = await service.getTrend(context)
      } catch (e) {
        console.error('Failed to load trend', e)
      }
    },
  },
})
