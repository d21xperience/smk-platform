import { computed } from 'vue'
import { useDashboardStore } from '../stores/DashboardStore.js'
import { useDispositionStore } from '../stores/DispositionStore.js'
import { useOperationalContext } from './useOperationalContext.js'

export function useDashboardStats() {
  const dashboardStore = useDashboardStore()
  const dispositionStore = useDispositionStore()
  const context = useOperationalContext()

  // Computed States
  const stats = computed(() => dashboardStore.stats)
  const trend = computed(() => dashboardStore.trend)
  const pendingDisposisi = computed(() => dispositionStore.pendingIncomingLetters)
  const tugasDisposisi = computed(() => dispositionStore.myTasks)
  const isLoading = computed(() => dashboardStore.isLoading)

  // Chart options for ApexCharts (UI Specific mapping)
  const chartOptions = computed(() => ({
    chart: { id: 'tren-surat', toolbar: { show: false } },
    xaxis: { categories: trend.value.categories || [] },
    stroke: { curve: 'smooth' },
  }))

  const chartSeries = computed(() => [
    { name: 'Surat Masuk', data: trend.value.seriesMasuk || [] },
    { name: 'Surat Keluar', data: trend.value.seriesKeluar || [] },
  ])

  // Actions
  async function loadDashboardData() {
    await Promise.all([
      dashboardStore.loadStats(context),
      dashboardStore.loadTrend(context),
      dispositionStore.loadPendingIncomingLetters(context),
      dispositionStore.loadMyTasks(context, context.userId),
    ])
  }

  return {
    stats,
    chartOptions,
    chartSeries,
    pendingDisposisi,
    tugasDisposisi,
    isLoading,
    loadDashboardData,
  }
}
