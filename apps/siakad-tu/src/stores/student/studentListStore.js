// apps/siakad-tu/src/stores/student/studentListStore.js

import { defineStore } from 'pinia'
import { ref } from 'vue'
import { useStudentQueryService } from '@/services/student/serviceFactory'
import { useContextStore } from '@/stores/contextStore'

export const useStudentListStore = defineStore('studentList', () => {
  // === STATE ===
  const students = ref([])
  const isLoading = ref(false)
  const error = ref(null)
  const pagination = ref({ page: 1, limit: 20, total: 0, totalPages: 0 })
  const filters = ref({ status: null, search: '' })

  // === SERVICES ===
  const queryService = useStudentQueryService()
  const contextStore = useContextStore()

  // === ACTIONS ===

  async function fetchStudents(customFilters = {}) {
    console.log('fethcStudent->store', customFilters)
    isLoading.value = true
    error.value = null

    try {
      const context = contextStore.currentContext
      console.log('contecxt', context)
      console.log('fethcStudent->context', context)

      const mergedFilters = { ...filters.value, ...customFilters }

      const result = await queryService.getStudents(context, mergedFilters)

      if (result.success) {
        students.value = result.data.items || []
        pagination.value = {
          page: result.data.page || 1,
          limit: result.data.limit || 20,
          total: result.data.total || 0,
          totalPages: result.data.totalPages || 0,
        }
      } else {
        error.value = result.error
      }
    } catch (err) {
      error.value = { code: 'UNEXPECTED_ERROR', message: err.message }
    } finally {
      isLoading.value = false
    }
  }

  async function deleteStudent(studentId) {
    isLoading.value = true
    error.value = null

    try {
      // const context = contextStore.currentContext
      // TODO: Implement delete di service layer
      // const result = await commandService.deleteStudent(studentId, context)

      // Simulasi untuk sekarang
      students.value = students.value.filter((s) => s.studentId !== studentId)

      return true
    } catch (err) {
      error.value = { code: 'UNEXPECTED_ERROR', message: err.message }
      return false
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
    deleteStudent,
    updateFilters,
    reset,
  }
})
