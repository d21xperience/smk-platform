const delay = (ms) => new Promise((resolve) => setTimeout(resolve, ms))

export const DashboardMockAdapter = {
  async getStats(context) {
    console.log(context)

    await delay(200)
    return {
      masuk: 24,
      keluar: 18,
      disposisiPending: 7,
    }
  },

  async getTrend(context) {
    console.log(context)
    await delay(300)
    return {
      categories: ['Jan', 'Feb', 'Mar', 'Apr', 'Mei', 'Jun'],
      seriesMasuk: [10, 15, 12, 20, 24, 18],
      seriesKeluar: [5, 8, 10, 15, 18, 12],
    }
  },
}
