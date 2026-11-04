// apps/siakad-tu/src/stores/student/studentListStore.js

import { defineStore } from 'pinia'
import { ref } from 'vue'
import { useStudentQueryService } from '@/services/student/serviceFactory.js'
import { useContextStore } from '@/stores/contextStore.js'
import { studentApi } from '@/adapters/api/studentApi.js' // Import langsung

export const useStudentListStore = defineStore('studentList', () => {
  const students = ref([])
  const isLoading = ref(false)
  const error = ref(null)
  const pagination = ref({ page: 1, limit: 20, total: 0, totalPages: 0 })
  
  const filters = ref({ status: null, search: '', tingkat: null, jurusan: null, kelas: null })
  const initialized = ref(false) 

  const queryService = useStudentQueryService()
  const contextStore = useContextStore()

  async function fetchStudents(customFilters = {}) {
    isLoading.value = true
    error.value = null

    try {
      const context = contextStore.currentContext
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
        initialized.value = true
      } else {
        error.value = result.error
      }
    } catch (err) {
      error.value = { code: 'UNEXPECTED_ERROR', message: err.message }
    } finally {
      isLoading.value = false
    }
  }

  function updateFilters(newFilters) {
    filters.value = { ...newFilters }
  }

  async function deleteStudent(studentId) {
    isLoading.value = true
    try {
      // Panggil langsung studentApi (bukan queryService.adapter)
      const result = await studentApi.deleteStudent(studentId)
      
      if (result.success) {
        removeItem(studentId) 
        return true
      }
      error.value = result.error
      return false
    } catch (err) {
      error.value = { code: 'UNEXPECTED_ERROR', message: err.message }
      return false
    } finally {
      isLoading.value = false
    }
  }

  function removeItem(studentId) {
    students.value = students.value.filter(s => s.id !== studentId)
    pagination.value.total = Math.max(0, pagination.value.total - 1)
  }

  function reset() {
    students.value = []
    error.value = null
    pagination.value = { page: 1, limit: 20, total: 0, totalPages: 0 }
    filters.value = { status: null, search: '', tingkat: null, jurusan: null, kelas: null }
    initialized.value = false
  }

  return {
    students, isLoading, error, pagination, filters, initialized,
    fetchStudents, updateFilters, deleteStudent, removeItem, reset,
  }
})
