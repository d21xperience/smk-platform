// src/composables/useReporting.js
import { computed } from 'vue'
import { useReportingStore } from 'src/stores/reporting.store'

export function useReporting() {
  const store = useReportingStore()

  return {
    reportData: computed(() => store.reportData),
    reportColumns: computed(() => store.reportColumns),
    reportMetadata: computed(() => store.reportMetadata),
    loading: computed(() => store.loading),
    error: computed(() => store.error),
    filterValues: computed(() => store.filterValues),

    loadReport: store.loadReport,
    applyFilters: store.applyFilters,
    exportReport: store.exportReport,
  }
}
