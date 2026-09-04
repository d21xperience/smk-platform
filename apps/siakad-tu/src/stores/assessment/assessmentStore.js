// apps/siakad-tu/src/stores/assessment/assessmentStore.js

import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { useAssessmentCommandService } from '@/services/assessment/serviceFactory'
import { useAssessmentQueryService } from '@/services/assessment/serviceFactory'
import { useContextStore } from '@/stores/contextStore'

export const useAssessmentStore = defineStore('assessment', () => {
  // === STATE ===
  const assessments = ref([])
  const currentAssessment = ref(null)
  const gradeSummary = ref(null)
  const reportCard = ref(null)
  const classGradeSummary = ref(null)
  const isLoading = ref(false)
  const error = ref(null)
  const success = ref(null)
  const pagination = ref({ page: 1, limit: 20, total: 0, totalPages: 0 })
  const filters = ref({ studentId: null, subjectId: null, type: null })

  // === SERVICES ===
  const commandService = useAssessmentCommandService()
  const queryService = useAssessmentQueryService()
  const contextStore = useContextStore()

  // === GETTERS ===
  const dailyAssessments = computed(() => assessments.value.filter((a) => a.type === 'DAILY'))

  const midtermAssessments = computed(() => assessments.value.filter((a) => a.type === 'MIDTERM'))

  const finalAssessments = computed(() => assessments.value.filter((a) => a.type === 'FINAL'))

  // === ACTIONS ===

  async function fetchAssessments(customFilters = {}) {
    isLoading.value = true
    error.value = null

    try {
      const context = contextStore.currentContext
      const mergedFilters = { ...filters.value, ...customFilters }

      const result = await queryService.getAssessments(context, mergedFilters)

      if (result.success) {
        assessments.value = result.data.items || []
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

  async function fetchAssessmentById(assessmentId) {
    isLoading.value = true
    error.value = null

    try {
      const context = contextStore.currentContext
      const result = await queryService.getAssessmentById(assessmentId, context)

      if (result.success) {
        currentAssessment.value = result.data
      } else {
        error.value = result.error
      }
    } catch (err) {
      error.value = { code: 'UNEXPECTED_ERROR', message: err.message }
    } finally {
      isLoading.value = false
    }
  }

  async function fetchGradeSummary(studentId, subjectId) {
    isLoading.value = true
    error.value = null

    try {
      const context = contextStore.currentContext
      const result = await queryService.getGradeSummary(studentId, subjectId, context)

      if (result.success) {
        gradeSummary.value = result.data
      } else {
        error.value = result.error
      }
    } catch (err) {
      error.value = { code: 'UNEXPECTED_ERROR', message: err.message }
    } finally {
      isLoading.value = false
    }
  }

  async function fetchStudentReportCard(studentId) {
    isLoading.value = true
    error.value = null

    try {
      const context = contextStore.currentContext
      const result = await queryService.getStudentReportCard(studentId, context)

      if (result.success) {
        reportCard.value = result.data
      } else {
        error.value = result.error
      }
    } catch (err) {
      error.value = { code: 'UNEXPECTED_ERROR', message: err.message }
    } finally {
      isLoading.value = false
    }
  }

  async function fetchClassGradeSummary(classId, subjectId) {
    isLoading.value = true
    error.value = null

    try {
      const context = contextStore.currentContext
      const result = await queryService.getClassGradeSummary(classId, subjectId, context)

      if (result.success) {
        classGradeSummary.value = result.data
      } else {
        error.value = result.error
      }
    } catch (err) {
      error.value = { code: 'UNEXPECTED_ERROR', message: err.message }
    } finally {
      isLoading.value = false
    }
  }

  async function inputScore(commandData) {
    isLoading.value = true
    error.value = null
    success.value = null

    try {
      const context = contextStore.currentContext
      const result = await commandService.inputScore(commandData, context)

      if (result.success) {
        success.value = { message: 'Nilai berhasil diinput' }
        currentAssessment.value = result.data
        return result.data
      } else {
        error.value = result.error
        return null
      }
    } catch (err) {
      error.value = { code: 'UNEXPECTED_ERROR', message: err.message }
      return null
    } finally {
      isLoading.value = false
    }
  }

  async function correctScore(commandData) {
    isLoading.value = true
    error.value = null
    success.value = null

    try {
      const context = contextStore.currentContext
      const result = await commandService.correctScore(commandData, context)

      if (result.success) {
        success.value = { message: 'Nilai berhasil dikoreksi' }
        currentAssessment.value = result.data

        const index = assessments.value.findIndex(
          (a) => a.assessmentId === result.data.assessmentId,
        )
        if (index !== -1) {
          assessments.value[index] = result.data
        }

        return result.data
      } else {
        error.value = result.error
        return null
      }
    } catch (err) {
      error.value = { code: 'UNEXPECTED_ERROR', message: err.message }
      return null
    } finally {
      isLoading.value = false
    }
  }

  async function deleteAssessment(assessmentId) {
    isLoading.value = true
    error.value = null
    success.value = null

    try {
      const context = contextStore.currentContext
      const result = await commandService.deleteAssessment(assessmentId, context)

      if (result.success) {
        success.value = { message: 'Nilai berhasil dihapus' }
        assessments.value = assessments.value.filter((a) => a.assessmentId !== assessmentId)
        return true
      } else {
        error.value = result.error
        return false
      }
    } catch (err) {
      error.value = { code: 'UNEXPECTED_ERROR', message: err.message }
      return false
    } finally {
      isLoading.value = false
    }
  }

  function reset() {
    assessments.value = []
    currentAssessment.value = null
    gradeSummary.value = null
    reportCard.value = null
    classGradeSummary.value = null
    error.value = null
    success.value = null
    pagination.value = { page: 1, limit: 20, total: 0, totalPages: 0 }
  }

  function updateFilters(newFilters) {
    filters.value = { ...filters.value, ...newFilters }
  }

  return {
    assessments,
    currentAssessment,
    gradeSummary,
    reportCard,
    classGradeSummary,
    isLoading,
    error,
    success,
    pagination,
    filters,
    dailyAssessments,
    midtermAssessments,
    finalAssessments,
    fetchAssessments,
    fetchAssessmentById,
    fetchGradeSummary,
    fetchStudentReportCard,
    fetchClassGradeSummary,
    inputScore,
    correctScore,
    deleteAssessment,
    reset,
    updateFilters,
  }
})
