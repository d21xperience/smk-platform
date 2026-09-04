// apps/siakad-tu/src/stores/registration/registrationStore.js

import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { useRegistrationCommandService } from '@/services/registration/serviceFactory'
import { useRegistrationQueryService } from '@/services/registration/serviceFactory'
import { useContextStore } from '@/stores/contextStore'

export const useRegistrationStore = defineStore('registration', () => {
  // === STATE ===
  const registrations = ref([])
  const currentRegistration = ref(null)
  const isLoading = ref(false)
  const error = ref(null)
  const success = ref(null)
  const pagination = ref({ page: 1, limit: 20, total: 0, totalPages: 0 })
  const filters = ref({ status: null, search: '' })

  // === SERVICES ===
  const commandService = useRegistrationCommandService()
  const queryService = useRegistrationQueryService()
  const contextStore = useContextStore()

  // === GETTERS ===
  const draftRegistrations = computed(() =>
    registrations.value.filter(r => r.status === 'DRAFT')
  )

  const pendingVerifications = computed(() =>
    registrations.value.filter(r => r.status === 'SUBMITTED' || r.status === 'VERIFIED')
  )

  const approvedRegistrations = computed(() =>
    registrations.value.filter(r => r.status === 'APPROVED')
  )

  // === ACTIONS ===

  /**
   * Fetch daftar registrasi dengan filter
   */
  async function fetchRegistrations(customFilters = {}) {
    isLoading.value = true
    error.value = null
    success.value = null

    try {
      const context = contextStore.currentContext
      const mergedFilters = { ...filters.value, ...customFilters }

      const result = await queryService.getRegistrations(context, mergedFilters)

      if (result.success) {
        registrations.value = result.data.items || []
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

  /**
   * Fetch detail registrasi by ID
   */
  async function fetchRegistrationById(registrationId) {
    isLoading.value = true
    error.value = null

    try {
      const context = contextStore.currentContext
      const result = await queryService.getRegistrationById(registrationId, context)

      if (result.success) {
        currentRegistration.value = result.data
      } else {
        error.value = result.error
      }
    } catch (err) {
      error.value = { code: 'UNEXPECTED_ERROR', message: err.message }
    } finally {
      isLoading.value = false
    }
  }

  /**
   * Buat draft registrasi baru
   */
  async function createDraft(commandData) {
    isLoading.value = true
    error.value = null
    success.value = null

    try {
      const context = contextStore.currentContext
      const result = await commandService.createDraft(commandData, context)

      if (result.success) {
        success.value = { message: 'Draft registrasi berhasil dibuat' }
        currentRegistration.value = result.data
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

  /**
   * Submit draft untuk verifikasi
   */
  async function submitRegistration(registrationId) {
    isLoading.value = true
    error.value = null
    success.value = null

    try {
      const context = contextStore.currentContext
      const result = await commandService.submitRegistration(registrationId, context)

      if (result.success) {
        success.value = { message: 'Registrasi berhasil di-submit untuk verifikasi' }
        currentRegistration.value = result.data

        // Update list jika ada
        const index = registrations.value.findIndex(r => r.registrationId === registrationId)
        if (index !== -1) {
          registrations.value[index] = result.data
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

  /**
   * Approve registrasi (trigger auto-create student)
   */
  async function approveRegistration(registrationId, commandData) {
    isLoading.value = true
    error.value = null
    success.value = null

    try {
      const context = contextStore.currentContext
      const result = await commandService.approveRegistration(registrationId, commandData, context)

      if (result.success) {
        success.value = { message: 'Registrasi berhasil di-approve. Siswa otomatis terdaftar.' }
        currentRegistration.value = result.data

        // Update list
        const index = registrations.value.findIndex(r => r.registrationId === registrationId)
        if (index !== -1) {
          registrations.value[index] = result.data
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

  /**
   * Reject registrasi
   */
  async function rejectRegistration(registrationId, reason) {
    isLoading.value = true
    error.value = null
    success.value = null

    try {
      const context = contextStore.currentContext
      const result = await commandService.rejectRegistration(registrationId, { reason }, context)

      if (result.success) {
        success.value = { message: 'Registrasi berhasil di-reject' }
        currentRegistration.value = result.data

        // Update list
        const index = registrations.value.findIndex(r => r.registrationId === registrationId)
        if (index !== -1) {
          registrations.value[index] = result.data
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

  /**
   * Reset state
   */
  function reset() {
    registrations.value = []
    currentRegistration.value = null
    error.value = null
    success.value = null
    pagination.value = { page: 1, limit: 20, total: 0, totalPages: 0 }
  }

  /**
   * Update filters
   */
  function updateFilters(newFilters) {
    filters.value = { ...filters.value, ...newFilters }
  }

  return {
    // State
    registrations,
    currentRegistration,
    isLoading,
    error,
    success,
    pagination,
    filters,

    // Getters
    draftRegistrations,
    pendingVerifications,
    approvedRegistrations,

    // Actions
    fetchRegistrations,
    fetchRegistrationById,
    createDraft,
    submitRegistration,
    approveRegistration,
    rejectRegistration,
    reset,
    updateFilters,
  }
})
