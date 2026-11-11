// apps/siakad-tu/src/stores/registration/registrationStore.js

import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { useRegistrationCommandService, useRegistrationQueryService } from '@/services/registration/serviceFactory'
import { useContextStore } from '@/stores/contextStore'

export const useRegistrationStore = defineStore('registration', () => {
  const registrations = ref([])
  const isLoading = ref(false)
  const error = ref(null)
  const filters = ref({ status: null, search: '' })

  const commandService = useRegistrationCommandService()
  const queryService = useRegistrationQueryService()
  const contextStore = useContextStore()

  const pendingVerifications = computed(() =>
    registrations.value.filter(r => r.status === 'REGISTERED' || r.status === 'SUBMITTED')
  )

  async function fetchRegistrations(customFilters = {}) {
    isLoading.value = true
    error.value = null
    try {
      const context = contextStore.currentContext
      const mergedFilters = { ...filters.value, ...customFilters }
      const result = await queryService.getRegistrations(context, mergedFilters)
      if (result.success) {
        registrations.value = result.data.items || []
      } else {
        error.value = result.error
      }
    } catch (err) {
      error.value = { code: 'UNEXPECTED_ERROR', message: err.message }
    } finally {
      isLoading.value = false
    }
  }

  async function createRegistration(commandData) {
    isLoading.value = true
    error.value = null
    try {
      const result = await commandService.createRegistration(commandData)
      if (result.success) {
        await fetchRegistrations() // Refresh list
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

  async function generateNisMassal(prefix) {
    isLoading.value = true
    error.value = null
    try {
      const result = await commandService.generateNisMassal(prefix)
      if (result.success) {
        await fetchRegistrations() // Refresh list
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

  async function activateRegistration(regId) {
    isLoading.value = true
    error.value = null
    try {
      const result = await commandService.activateRegistration(regId)
      if (result.success) {
        await fetchRegistrations() // Refresh list
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

  function updateFilters(newFilters) {
    filters.value = { ...filters.value, ...newFilters }
  }

  return {
    registrations,
    isLoading,
    error,
    filters,
    pendingVerifications,
    fetchRegistrations,
    createRegistration,
    generateNisMassal,
    activateRegistration,
    updateFilters,
  }
})
