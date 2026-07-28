import { computed } from 'vue'
import { useContextStore } from '@/stores/context.store'

export function useContext() {
  const store = useContextStore()

  return {
    // Operational Context
    operational: computed(() => store.operational),
    isOperationalContextReady: computed(() => store.isOperationalContextReady),

    // Query Context
    query: computed(() => store.query),

    // Data Master (dropdown)
    academicYears: computed(() => store.academicYears),
    semesters: computed(() => store.semesters),
    academicPeriods: computed(() => store.academicPeriods),

    // Status loading
    loading: computed(() => store.loading),

    // Actions untuk master data
    loadAcademicYears: store.loadAcademicYears,
    loadSemesters: store.loadSemesters,
    loadAcademicPeriods: store.loadAcademicPeriods,

    // Actions untuk operational context
    setOperationalContext: store.setOperationalContext,

    // Actions untuk query context
    setQueryContext: store.setQueryContext,
    clearQueryContext: store.clearQueryContext,
  }
}
