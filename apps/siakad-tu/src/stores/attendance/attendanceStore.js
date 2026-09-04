// apps/siakad-tu/src/stores/attendance/attendanceStore.js

import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { useAttendanceCommandService } from '@/services/attendance/serviceFactory'
import { useAttendanceQueryService } from '@/services/attendance/serviceFactory'
import { useContextStore } from '@/stores/contextStore'

export const useAttendanceStore = defineStore('attendance', () => {
  // === STATE ===
  const sessions = ref([])
  const currentSession = ref(null)
  const attendanceRecords = ref([])
  const attendanceSummary = ref(null)
  const classSummary = ref([])
  const isLoading = ref(false)
  const error = ref(null)
  const success = ref(null)
  const pagination = ref({ page: 1, limit: 20, total: 0, totalPages: 0 })
  const filters = ref({ classId: null, subjectId: null, status: null, date: null })

  // === SERVICES ===
  const commandService = useAttendanceCommandService()
  const queryService = useAttendanceQueryService()
  const contextStore = useContextStore()

  // === GETTERS ===
  const openSessions = computed(() => sessions.value.filter((s) => s.status === 'OPEN'))

  const closedSessions = computed(() => sessions.value.filter((s) => s.status === 'CLOSED'))

  // === ACTIONS ===

  async function fetchSessions(customFilters = {}) {
    isLoading.value = true
    error.value = null

    try {
      const context = contextStore.currentContext
      const mergedFilters = { ...filters.value, ...customFilters }

      const result = await queryService.getSessions(context, mergedFilters)

      if (result.success) {
        sessions.value = result.data.items || []
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

  async function fetchSessionById(sessionId) {
    isLoading.value = true
    error.value = null

    try {
      const context = contextStore.currentContext
      const result = await queryService.getSessionById(sessionId, context)

      if (result.success) {
        currentSession.value = result.data
      } else {
        error.value = result.error
      }
    } catch (err) {
      error.value = { code: 'UNEXPECTED_ERROR', message: err.message }
    } finally {
      isLoading.value = false
    }
  }

  async function fetchAttendanceRecords(sessionId) {
    isLoading.value = true
    error.value = null

    try {
      const context = contextStore.currentContext
      const result = await queryService.getAttendanceRecords(sessionId, context)

      if (result.success) {
        attendanceRecords.value = result.data || []
      } else {
        error.value = result.error
      }
    } catch (err) {
      error.value = { code: 'UNEXPECTED_ERROR', message: err.message }
    } finally {
      isLoading.value = false
    }
  }

  async function fetchAttendanceSummary(studentId) {
    isLoading.value = true
    error.value = null

    try {
      const context = contextStore.currentContext
      const result = await queryService.getAttendanceSummary(studentId, context)

      if (result.success) {
        attendanceSummary.value = result.data
      } else {
        error.value = result.error
      }
    } catch (err) {
      error.value = { code: 'UNEXPECTED_ERROR', message: err.message }
    } finally {
      isLoading.value = false
    }
  }

  async function fetchClassAttendanceSummary(classId) {
    isLoading.value = true
    error.value = null

    try {
      const context = contextStore.currentContext
      const result = await queryService.getClassAttendanceSummary(classId, context)

      if (result.success) {
        classSummary.value = result.data || []
      } else {
        error.value = result.error
      }
    } catch (err) {
      error.value = { code: 'UNEXPECTED_ERROR', message: err.message }
    } finally {
      isLoading.value = false
    }
  }

  async function openSession(commandData) {
    isLoading.value = true
    error.value = null
    success.value = null

    try {
      const context = contextStore.currentContext
      const result = await commandService.openSession(commandData, context)

      if (result.success) {
        success.value = { message: 'Sesi absensi berhasil dibuka' }
        currentSession.value = result.data
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

  async function closeSession(sessionId) {
    isLoading.value = true
    error.value = null
    success.value = null

    try {
      const context = contextStore.currentContext
      const result = await commandService.closeSession(sessionId, context)

      if (result.success) {
        success.value = { message: 'Sesi absensi berhasil ditutup' }
        currentSession.value = result.data

        const index = sessions.value.findIndex((s) => s.sessionId === sessionId)
        if (index !== -1) {
          sessions.value[index] = result.data
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

  async function submitBulkAttendance(commandData) {
    isLoading.value = true
    error.value = null
    success.value = null

    try {
      const context = contextStore.currentContext
      const result = await commandService.submitBulkAttendance(commandData, context)

      if (result.success) {
        success.value = { message: 'Absensi berhasil dicatat' }
        attendanceRecords.value = result.data || []
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

  async function updateAttendanceRecord(commandData) {
    isLoading.value = true
    error.value = null
    success.value = null

    try {
      const context = contextStore.currentContext
      const result = await commandService.updateAttendanceRecord(commandData, context)

      if (result.success) {
        success.value = { message: 'Record absensi berhasil diupdate' }

        const index = attendanceRecords.value.findIndex((r) => r.recordId === result.data.recordId)
        if (index !== -1) {
          attendanceRecords.value[index] = result.data
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

  function reset() {
    sessions.value = []
    currentSession.value = null
    attendanceRecords.value = []
    attendanceSummary.value = null
    classSummary.value = []
    error.value = null
    success.value = null
    pagination.value = { page: 1, limit: 20, total: 0, totalPages: 0 }
  }

  function updateFilters(newFilters) {
    filters.value = { ...filters.value, ...newFilters }
  }

  return {
    sessions,
    currentSession,
    attendanceRecords,
    attendanceSummary,
    classSummary,
    isLoading,
    error,
    success,
    pagination,
    filters,
    openSessions,
    closedSessions,
    fetchSessions,
    fetchSessionById,
    fetchAttendanceRecords,
    fetchAttendanceSummary,
    fetchClassAttendanceSummary,
    openSession,
    closeSession,
    submitBulkAttendance,
    updateAttendanceRecord,
    reset,
    updateFilters,
  }
})
