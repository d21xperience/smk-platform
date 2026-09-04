// apps/siakad-tu/src/stores/student/studentListStore.js

import { defineStore } from 'pinia'
import { ref } from 'vue'
import { useStudentQueryService } from '@/services/student/serviceFactory.js'
import { useContextStore } from '@/stores/contextStore.js'

export const useStudentListStore = defineStore('studentList', () => {
  const students = ref([])
  const isLoading = ref(false)
  const error = ref(null)
  const pagination = ref({ page: 1, limit: 20, total: 0, totalPages: 0 })
  const filters = ref({ status: null, search: '' })

  const queryService = useStudentQueryService()
  const contextStore = useContextStore()

  async function fetchStudents(customFilters = {}) {
    console.log('[StudentListStore] fetchStudents called')

    isLoading.value = true
    error.value = null

    try {
      const context = contextStore.current

      // Jika context null, axios interceptor akan menangani fallback-nya
      // Tapi kita tetap perlu mengirim object kosong atau dummy ke service agar tidak error
      const finalContext = context || { schoolId: 'school-debug-001', periodId: 'period-2024-1' }

      const mergedFilters = { ...filters.value, ...customFilters }

      const result = await queryService.getStudents(finalContext, mergedFilters)

      if (result.success) {
        students.value = result.data.items || []
        pagination.value = {
          page: result.data.page || 1,
          limit: result.data.limit || 20,
          total: result.data.total || 0,
          totalPages: result.data.totalPages || 0,
        }
        console.log('[StudentListStore] Students loaded:', students.value.length)
      } else {
        console.error('[StudentListStore] Query failed:', result.error)
        error.value = result.error
      }
    } catch (err) {
      console.error('[StudentListStore] Unexpected error:', err)
      error.value = { code: 'UNEXPECTED_ERROR', message: err.message }
    } finally {
      isLoading.value = false
    }
  }

  function updateFilters(newFilters) {
    filters.value = { ...filters.value, ...newFilters }
  }

  function reset() {
    students.value = []
    error.value = null
    pagination.value = { page: 1, limit: 20, total: 0, totalPages: 0 }
  }

  return {
    students,
    isLoading,
    error,
    pagination,
    filters,
    fetchStudents,
    updateFilters,
    reset,
  }
})
