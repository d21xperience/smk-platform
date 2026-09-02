import { computed, onMounted } from 'vue'
import { useStudentAffairsStore } from '@/stores/kesiswaan/studentAffairsStore'

export function useStudentList() {
  const store = useStudentAffairsStore()

  onMounted(() => {
    store.loadStudents()
  })

  // Watch filters untuk debounce atau auto-fetch jika diperlukan
  // Saat ini sudah dihandle oleh store.setFilter yang langsung memanggil loadStudents

  const items = computed(() => store.students)
  const loading = computed(() => store.loading)
  const error = computed(() => store.error)
  const filters = computed(() => store.filters)

  const stats = computed(() => ({
    total: items.value.length,
    active: items.value.filter((s) => s.status === 'AKTIF').length,
  }))

  const pagination = computed(() => ({
    page: 1,
    totalPages: 1,
    total: items.value.length,
    limit: 10,
    hasMore: false,
  }))

  return {
    items,
    stats,
    pagination,
    filters,
    loading,
    error,
    setSearch: (val) => store.setFilter('search', val),
    setStatusFilter: (val) => store.setFilter('status', val),
    setJurusanFilter: (val) => store.setFilter('jurusan', val),
    setTingkatFilter: (val) => store.setFilter('tingkat', val),
    clearFilters: () => store.clearFilters(),
    clearError: () => {
      store.error = null
    },
    nextPage: () => {}, // Implementasi paginasi backend di Phase 2
    prevPage: () => {},
    goToPage: () => {},
  }
}
