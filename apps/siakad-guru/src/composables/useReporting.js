import { computed } from 'vue'
import { useReportingStore } from '../stores/reportingStore.js'
import { useContextStore } from '../stores/contextStore.js'

export function useReporting() {
  const reportingStore = useReportingStore()
  const contextStore = useContextStore()

  const currentContext = computed(() => contextStore.currentContext)
  const dashboard = computed(() => reportingStore.dashboard)
  const isLoading = computed(() => reportingStore.loading)
  const error = computed(() => reportingStore.error)

  const loadDashboard = async () => {
    if (!currentContext.value) return
    await reportingStore.loadDashboard({
      schoolId: currentContext.value.schoolId,
      academicYearId: currentContext.value.academicYearId,
      semesterId: currentContext.value.semesterId,
      teacherId: currentContext.value.userId,
    })
  }

  return {
    dashboard,
    isLoading,
    error,
    loadDashboard,
  }
}
