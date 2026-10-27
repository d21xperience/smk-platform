import { computed, ref } from 'vue'
import { useRouter } from 'vue-router'
import { useContextStore } from '../../stores/contextStore.js'
import { useAuthStore } from '../../stores/authStore.js'
import { ContextEngine } from '../../domain/context/engine/ContextEngine.js'
import { eventDispatcher } from '../../events/eventDispatcher.js'
import {
  mockContextData,
  mockYearOptions,
  mockSemesterOptions,
} from '../../adapters/mock/data/context.mock.js'

export function useOperationalContext() {
  const store = useContextStore()
  const authStore = useAuthStore()
  const engine = new ContextEngine()
  const router = useRouter()

  const current = computed(() => store.current)
  const isReady = computed(() => store.isContextReady)
  const displayLabel = computed(() => store.displayLabel)

  const selectedYear = computed(() => store.current?.academicYear || mockContextData.academicYear)
  const selectedSemester = computed(() => store.current?.semester || mockContextData.semester)

  const yearOptions = ref(mockYearOptions)
  const semesterOptions = ref(mockSemesterOptions)

  // Otorisasi Role
  const hasYearChangeAuthority = computed(() => authStore.hasPermission('academic_year:write'))

  const isLoading = ref(false)
  const error = ref(null)

  async function selectContext(contextData) {
    isLoading.value = true
    error.value = null
    try {
      const { context, events } = engine.establishContext(contextData)
      store.setContext(context)
      localStorage.setItem('operational_context', JSON.stringify(context.toJSON()))
      events.forEach((event) => eventDispatcher.dispatch(event))
      return { success: true }
    } catch (err) {
      error.value = err.message
      return { success: false, error: err.message }
    } finally {
      isLoading.value = false
    }
  }

  async function onYearChange(year) {
    if (!hasYearChangeAuthority.value) return
    await selectContext({
      schoolId: store.current?.schoolId || mockContextData.schoolId,
      schoolName: store.current?.schoolName || mockContextData.schoolName,
      academicYear: year,
      semester: selectedSemester.value,
    })
  }

  async function onSemesterChange(semester) {
    await selectContext({
      schoolId: store.current?.schoolId || mockContextData.schoolId,
      schoolName: store.current?.schoolName || mockContextData.schoolName,
      academicYear: selectedYear.value,
      semester: semester,
    })
  }

  function restoreContext() {
    try {
      const saved = localStorage.getItem('operational_context')
      if (!saved) return false
      const data = JSON.parse(saved)
      const { context } = engine.establishContext(data)
      store.setContext(context)
      return true
    } catch {
      store.clearContext()
      localStorage.removeItem('operational_context')
      return false
    }
  }

  function clearContext() {
    const { events } = engine.clearContext()
    store.clearContext()
    localStorage.removeItem('operational_context')
    events.forEach((event) => eventDispatcher.dispatch(event))
  }

  function requireContext() {
    const validation = engine.validateForTransaction(store.current)
    if (!validation.valid) {
      router.push({ name: 'context-selection' })
      throw new Error(validation.reason)
    }
    return store.current
  }

  return {
    current,
    isReady,
    displayLabel,
    selectedYear,
    selectedSemester,
    yearOptions,
    semesterOptions,
    hasYearChangeAuthority,
    isLoading,
    error,
    selectContext,
    onYearChange,
    onSemesterChange,
    restoreContext,
    clearContext,
    requireContext,
  }
}
