import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { kesiswaanDashboardQueryService } from '@/services/kesiswaan'
import { useContextStore } from './contextStore'

export const useKesiswaanStore = defineStore('kesiswaan', () => {
  // === STATE ===
  const dashboardData = ref(null)
  const isLoading = ref(false)
  const error = ref(null)

  // === GETTERS (computed) ===
  const hasData = computed(() => dashboardData.value !== null)
  const hasError = computed(() => error.value !== null)

  // === ACTIONS ===
  /**
   * Memuat data dashboard kesiswaan
   * @param {Object} context - Optional, jika tidak diberikan akan diambil dari contextStore
   * @returns {Promise<{ success: boolean, data?: object, error?: object }>}
   */
  async function loadDashboard(context) {
    // Jika context tidak diberikan, ambil dari contextStore
    let operationalContext = context
    if (!operationalContext) {
      const contextStore = useContextStore()
      operationalContext = contextStore.current
    }

    // Validasi context
    if (!operationalContext || !operationalContext.schoolId) {
      error.value = 'Context tidak valid atau schoolId tidak ditemukan.'
      isLoading.value = false
      return {
        success: false,
        data: null,
        error: { code: 'INVALID_CONTEXT', message: error.value },
      }
    }

    // Set loading
    isLoading.value = true
    error.value = null

    try {
      const result = await kesiswaanDashboardQueryService.loadDashboard(operationalContext)

      if (result.success) {
        dashboardData.value = result.data
        isLoading.value = false
        return {
          success: true,
          data: result.data,
          error: null,
        }
      } else {
        // Error dari service
        error.value = result.error?.message || 'Gagal memuat data dashboard.'
        dashboardData.value = null
        isLoading.value = false
        return {
          success: false,
          data: null,
          error: result.error,
        }
      }
    } catch (err) {
      // Unexpected error
      error.value = err.message || 'Terjadi kesalahan tak terduga.'
      dashboardData.value = null
      isLoading.value = false
      return {
        success: false,
        data: null,
        error: {
          code: 'UNEXPECTED_ERROR',
          message: error.value,
          details: err,
        },
      }
    }
  }

  /**
   * Reset store ke keadaan awal
   */
  function reset() {
    dashboardData.value = null
    isLoading.value = false
    error.value = null
  }

  // === RETURN ===
  return {
    // State
    dashboardData,
    isLoading,
    error,
    // Getters
    hasData,
    hasError,
    // Actions
    loadDashboard,
    reset,
  }
})
