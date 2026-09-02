import { computed, onMounted } from 'vue'
import { useStudentAffairsStore } from '@/stores/kesiswaan/studentAffairsStore'

export function useDashboard() {
  const store = useStudentAffairsStore()

  onMounted(() => {
    store.loadDashboard()
  })

  const metrics = computed(() => store.dashboardData?.metrics || {})
  const chartSeries = computed(() => store.dashboardData?.chartSeries || [])
  const reminders = computed(() => store.dashboardData?.reminders || [])
  const suratQueue = computed(() => store.dashboardData?.suratQueue || [])
  const activityLogs = computed(() => store.dashboardData?.activityLogs || [])
  const loading = computed(() => store.loading)

  return {
    metrics,
    chartSeries,
    reminders,
    suratQueue,
    activityLogs,
    loading,
  }
}
