// apps/siakad-tu/src/stores/kesiswaan/studentBukuIndukStore.js

import { defineStore } from 'pinia'
import { ref } from 'vue'
import { studentBukuIndukRealAdapter } from '@/adapters/api/studentBukuIndukRealAdapter.js'

export const useStudentBukuIndukStore = defineStore('studentBukuInduk', () => {
  const students = ref([])
  const activeStudent = ref(null)
  const loading = ref(false)
  const loadingDetail = ref(false)
  const error = ref(null)

  async function fetchList() {
    if (loading.value) return
    loading.value = true
    error.value = null
    try {
      const result = await studentBukuIndukRealAdapter.list()
      if (result.success) {
        students.value = result.data
      } else {
        error.value = result.error
        students.value = []
      }
    } finally {
      loading.value = false
    }
  }

  async function fetchDetail(studentId) {
    if (!studentId) return
    loadingDetail.value = true
    error.value = null
    try {
      const result = await studentBukuIndukRealAdapter.getById(studentId)
      if (result.success) {
        activeStudent.value = result.data
      } else {
        error.value = result.error
        activeStudent.value = null
      }
    } finally {
      loadingDetail.value = false
    }
  }

  async function downloadPDF(studentId) {
    if (!studentId) return { success: false, error: 'ID required' }
    return await studentBukuIndukRealAdapter.downloadPDF(studentId)
  }

  function clearActive() {
    activeStudent.value = null
  }

  function reset() {
    students.value = []
    activeStudent.value = null
    error.value = null
  }

  return {
    students,
    activeStudent,
    loading,
    loadingDetail,
    error,
    fetchList,
    fetchDetail,
    downloadPDF,
    clearActive,
    reset,
  }
})
