import { defineStore } from 'pinia'
import { ref } from 'vue'
import { api } from '@/boot/axios'

export const useSarprasDashboardStore = defineStore('sarpras-dashboard', () => {
  const stats = ref({
    total_assets: 0,
    total_value: 0,
    condition_baik: 0,
    condition_rr: 0,
    condition_rb: 0,
    active_borrowings: 0,
    pending_disposals: 0,
    low_stock_items: 0,
  })
  const isLoading = ref(false)

  const fetchStats = async () => {
    isLoading.value = true
    try {
      const res = await api.get('/sarpras/dashboard')
      if (res.data.data) {
        stats.value = res.data.data
      }
    } catch (e) {
      console.error('Gagal mengambil data dashboard:', e)
    } finally {
      isLoading.value = false
    }
  }

  const downloadComprehensiveReport = async (startDate, endDate) => {
    const res = await api.get(
      `/sarpras/reports/comprehensive?start_date=${startDate}&end_date=${endDate}`,
      { responseType: 'blob' },
    )
    const url = window.URL.createObjectURL(new Blob([res.data], { type: 'application/pdf' }))
    const link = document.createElement('a')
    link.href = url
    link.setAttribute('download', `Laporan_Komprehensif_${startDate}.pdf`)
    document.body.appendChild(link)
    link.click()
    link.remove()
  }

  return { stats, isLoading, fetchStats, downloadComprehensiveReport }
})
