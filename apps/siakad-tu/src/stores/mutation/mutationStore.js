// apps/siakad-tu/src/stores/mutation/mutationStore.js

import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { useMutationCommandService } from '@/services/mutation/serviceFactory'
import { useMutationQueryService } from '@/services/mutation/serviceFactory'
import { useContextStore } from '@/stores/contextStore'

export const useMutationStore = defineStore('mutation', () => {
  // === STATE ===
  const mutations = ref([])
  const currentMutation = ref(null)
  const isLoading = ref(false)
  const error = ref(null)
  const success = ref(null)
  const pagination = ref({ page: 1, limit: 20, total: 0, totalPages: 0 })
  const filters = ref({ type: null, status: null, search: '' })

  // === SERVICES ===
  const commandService = useMutationCommandService()
  const queryService = useMutationQueryService()
  const contextStore = useContextStore()

  // === GETTERS ===
  const pendingMutations = computed(() =>
    mutations.value.filter(m => m.status === 'PENDING')
  )

  const approvedMutations = computed(() =>
    mutations.value.filter(m => m.status === 'APPROVED')
  )

  const mutationIn = computed(() =>
    mutations.value.filter(m => m.type === 'IN')
  )

  const mutationOut = computed(() =>
    mutations.value.filter(m => m.type === 'OUT')
  )

  // === ACTIONS ===

  async function fetchMutations(customFilters = {}) {
    isLoading.value = true
    error.value = null

    try {
      const context = contextStore.currentContext
      const mergedFilters = { ...filters.value, ...customFilters }

      const result = await queryService.getMutations(context, mergedFilters)

      if (result.success) {
        mutations.value = result.data.items || []
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

  async function fetchMutationById(mutationId) {
    isLoading.value = true
    error.value = null

    try {
      const context = contextStore.currentContext
      const result = await queryService.getMutationById(mutationId, context)

      if (result.success) {
        currentMutation.value = result.data
      } else {
        error.value = result.error
      }
    } catch (err) {
      error.value = { code: 'UNEXPECTED_ERROR', message: err.message }
    } finally {
      isLoading.value = false
    }
  }

  async function requestMutationIn(commandData) {
    isLoading.value = true
    error.value = null
    success.value = null

    try {
      const context = contextStore.currentContext
      const result = await commandService.requestMutationIn(commandData, context)

      if (result.success) {
        success.value = { message: 'Permohonan mutasi masuk berhasil diajukan' }
        currentMutation.value = result.data
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

  async function requestMutationOut(commandData) {
    isLoading.value = true
    error.value = null
    success.value = null

    try {
      const context = contextStore.currentContext
      const result = await commandService.requestMutationOut(commandData, context)

      if (result.success) {
        success.value = { message: 'Permohonan mutasi keluar berhasil diajukan' }
        currentMutation.value = result.data
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

  async function approveMutation(mutationId, commandData) {
    isLoading.value = true
    error.value = null
    success.value = null

    try {
      const context = contextStore.currentContext
      const result = await commandService.approveMutation(mutationId, commandData, context)

      if (result.success) {
        success.value = { message: 'Mutasi berhasil di-approve' }
        currentMutation.value = result.data

        const index = mutations.value.findIndex(m => m.mutationId === mutationId)
        if (index !== -1) {
          mutations.value[index] = result.data
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

  async function rejectMutation(mutationId, rejectionReason) {
    isLoading.value = true
    error.value = null
    success.value = null

    try {
      const context = contextStore.currentContext
      const result = await commandService.rejectMutation(mutationId, { rejectionReason }, context)

      if (result.success) {
        success.value = { message: 'Mutasi berhasil di-reject' }
        currentMutation.value = result.data

        const index = mutations.value.findIndex(m => m.mutationId === mutationId)
        if (index !== -1) {
          mutations.value[index] = result.data
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
    mutations.value = []
    currentMutation.value = null
    error.value = null
    success.value = null
    pagination.value = { page: 1, limit: 20, total: 0, totalPages: 0 }
  }

  function updateFilters(newFilters) {
    filters.value = { ...filters.value, ...newFilters }
  }

  return {
    mutations,
    currentMutation,
    isLoading,
    error,
    success,
    pagination,
    filters,
    pendingMutations,
    approvedMutations,
    mutationIn,
    mutationOut,
    fetchMutations,
    fetchMutationById,
    requestMutationIn,
    requestMutationOut,
    approveMutation,
    rejectMutation,
    reset,
    updateFilters,
  }
})
