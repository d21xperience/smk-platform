export function createDashboardService({ adapter }) {
  return {
    async getStats(context) {
      return adapter.getStats(context)
    },
    async getTrend(context) {
      return adapter.getTrend(context)
    },
  }
}
