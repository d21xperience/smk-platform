import { defineStore } from 'pinia'
import { ref } from 'vue'
import { StudentAffairsService } from '@/services/StudentAffairsService'

export const useStudentAffairsStore = defineStore('studentAffairs', () => {
  // State
  const students = ref([])
  const mutations = ref([])
  const dashboardData = ref(null)
  const loading = ref(false)
  const error = ref(null)

  // Filters (untuk Student List)
  const filters = ref({
    search: '',
    status: '',
    jurusan: 'Semua',
    tingkat: 'Semua',
  })

  // Actions: Dashboard
  async function loadDashboard() {
    loading.value = true
    error.value = null
    try {
      dashboardData.value = await StudentAffairsService.getDashboardMetrics()
    } catch (err) {
      error.value = err.message
    } finally {
      loading.value = false
    }
  }

  // Actions: Students
  async function loadStudents() {
    loading.value = true
    error.value = null
    try {
      const response = await StudentAffairsService.getStudents(filters.value)
      students.value = response.data
      return response
    } catch (err) {
      error.value = err.message
    } finally {
      loading.value = false
    }
  }

  function setFilter(key, value) {
    filters.value[key] = value
    loadStudents() // Auto-fetch saat filter berubah
  }

  function clearFilters() {
    filters.value = { search: '', status: '', jurusan: 'Semua', tingkat: 'Semua' }
    loadStudents()
  }

  // Actions: Mutations
  async function loadMutations(statusFilter = 'semua') {
    loading.value = true
    error.value = null
    try {
      mutations.value = await StudentAffairsService.getMutations(statusFilter)
    } catch (err) {
      error.value = err.message
    } finally {
      loading.value = false
    }
  }

  async function approveMutationAction(id, catatan) {
    loading.value = true
    try {
      await StudentAffairsService.approveMutation(id, catatan)
      // Reload data untuk refleksi perubahan
      await loadMutations()
      await loadStudents() // Karena status siswa mungkin berubah
    } catch (err) {
      error.value = err.message
      throw err
    } finally {
      loading.value = false
    }
  }

  async function rejectMutationAction(id, catatan) {
    loading.value = true
    try {
      await StudentAffairsService.rejectMutation(id, catatan)
      await loadMutations()
    } catch (err) {
      error.value = err.message
      throw err
    } finally {
      loading.value = false
    }
  }

  return {
    students,
    mutations,
    dashboardData,
    loading,
    error,
    filters,
    loadDashboard,
    loadStudents,
    loadMutations,
    setFilter,
    clearFilters,
    approveMutationAction,
    rejectMutationAction,
  }
})
