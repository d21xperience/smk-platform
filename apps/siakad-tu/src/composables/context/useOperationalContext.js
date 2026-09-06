// apps/siakad-tu/src/composables/context/useOperationalContext.js

import { computed, onMounted, ref } from 'vue'
import { useContextStore } from '@/stores/contextStore'
import { OperationalContext } from '@/domain/context/OperationalContext'

export function useOperationalContext() {
  const contextStore = useContextStore()

  // === STATE ===
  const isLoading = ref(false)
  const error = ref(null)

  // === STATE BARU: Untuk Context Selector ===
  const selectedYear = ref(null)
  const selectedSemester = ref(null)

  // Daftar tahun ajaran yang tersedia
  // TODO: Ke depannya load dari API backend
  const availableYears = ref([
    { value: '2024/2025', label: '2024/2025', isActive: true },
    { value: '2023/2024', label: '2023/2024', isActive: false },
    { value: '2022/2023', label: '2022/2023', isActive: false },
  ])

  const availableSemesters = ref([
    { value: 1, label: 'Semester 1 (Ganjil)' },
    { value: 2, label: 'Semester 2 (Genap)' },
  ])

  // === COMPUTED (EXISTING) ===
  const current = computed(() => contextStore.current)
  const isReady = computed(() => contextStore.isContextReady)
  const schoolId = computed(() => contextStore.schoolId)
  const schoolName = computed(() => contextStore.schoolName)
  const academicYear = computed(() => contextStore.academicYear)
  const semester = computed(() => contextStore.semester)
  const periodId = computed(() => contextStore.periodId)
  const displayLabel = computed(() => contextStore.displayLabel)
  const schoolLabel = computed(() => contextStore.schoolLabel)

  // === COMPUTED BARU: Untuk q-select options ===
  const yearOptions = computed(() =>
    availableYears.value.map((y) => ({ label: y.label, value: y.value })),
  )

  const semesterOptions = computed(() =>
    availableSemesters.value.map((s) => ({ label: s.label, value: s.value })),
  )

  // === ACTIONS (EXISTING) ===
  function setContext(contextData) {
    const ctx = new OperationalContext(contextData)
    if (ctx.isValid()) {
      contextStore.setContext(ctx)

      // Simpan ke localStorage
      localStorage.setItem('operational_context', JSON.stringify(contextData))
    } else {
      error.value = 'Context tidak valid'
    }
  }

  function clearContext() {
    contextStore.clearContext()
    localStorage.removeItem('operational_context')
    selectedYear.value = null
    selectedSemester.value = null
  }

  // === ACTIONS BARU: Untuk Context Selector ===

  /**
   * Handle perubahan tahun ajaran dari dropdown
   */
  function onYearChange(newYear) {
    selectedYear.value = newYear
    applyContext()
  }

  /**
   * Handle perubahan semester dari dropdown
   */
  function onSemesterChange(newSemester) {
    selectedSemester.value = newSemester
    applyContext()
  }

  /**
   * Terapkan context yang dipilih ke contextStore
   */
  function applyContext() {
    if (!selectedYear.value || !selectedSemester.value) return

    setContext({
      schoolId: schoolId.value || 'school-001',
      schoolName: schoolName.value || 'SMK Negeri 1',
      periodId: `${selectedYear.value}-sem-${selectedSemester.value}`,
      academicYear: selectedYear.value,
      semester: selectedSemester.value,
    })
  }

  /**
   * Set default context ke tahun ajaran aktif
   */
  function setDefaultContext() {
    const activeYear = availableYears.value.find((y) => y.isActive)
    if (activeYear) {
      selectedYear.value = activeYear.value
      selectedSemester.value = 1
      applyContext()
    }
  }

  // === LIFECYCLE (UPDATED) ===
  onMounted(() => {
    if (!isReady.value) {
      loadDefaultContext()
    }
  })

  function loadDefaultContext() {
    const saved = localStorage.getItem('operational_context')
    if (saved) {
      try {
        const parsed = JSON.parse(saved)
        selectedYear.value = parsed.academicYear
        selectedSemester.value = parsed.semester
        setContext(parsed)
      } catch {
        setDefaultContext()
      }
    } else {
      setDefaultContext()
    }
  }

  // === RETURN (UPDATED) ===
  return {
    // Existing
    current,
    isReady,
    schoolId,
    schoolName,
    academicYear,
    semester,
    periodId,
    displayLabel,
    isLoading,
    error,
    schoolLabel,
    setContext,
    clearContext,

    // Baru: Context Selector
    selectedYear,
    selectedSemester,
    availableYears,
    availableSemesters,
    yearOptions,
    semesterOptions,
    onYearChange,
    onSemesterChange,
    setDefaultContext,
  }
}
