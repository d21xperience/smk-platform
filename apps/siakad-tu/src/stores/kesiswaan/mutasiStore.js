// apps/siakad-tu/src/stores/kesiswaan/mutasiStore.js

import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { useMutationCommandService, useMutationQueryService } from '@/services/mutation/serviceFactory'
import { useContextStore } from '@/stores/contextStore'

export const useMutasiStore = defineStore('mutasi', () => {
  const contextStore = useContextStore()
  const commandService = useMutationCommandService()
  const queryService = useMutationQueryService()

  const list = ref([])
  const isLoading = ref(false)
  const error = ref(null)

  const riwayatMasuk = computed(() => list.value.filter((m) => m.mutationType === 'MASUK'))
  const riwayatKeluar = computed(() => list.value.filter((m) => m.mutationType === 'KELUAR'))
  const pendingCount = computed(() => list.value.filter((m) => m.status === 'PENDING').length)

  async function fetchMutations(filters = {}) {
    isLoading.value = true
    error.value = null
    try {
      const context = contextStore.currentContext
      const result = await queryService.getMutations(context, filters)
      if (result.success) {
        list.value = result.data.items || []
      } else {
        error.value = result.error
      }
    } catch (err) {
      error.value = { code: 'UNEXPECTED_ERROR', message: err.message }
    } finally {
      isLoading.value = false
    }
  }

  async function createMutation(payload) {
    isLoading.value = true
    error.value = null
    try {
      const result = await commandService.createMutation(payload)
      if (result.success) {
        await fetchMutations()
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

  async function approveMutation(mutationId, catatan = '') {
    isLoading.value = true
    error.value = null
    try {
      const result = await commandService.approveMutation(mutationId, catatan)
      if (result.success) {
        await fetchMutations()
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

  async function rejectMutation(mutationId, catatan = '') {
    isLoading.value = true
    error.value = null
    try {
      const result = await commandService.rejectMutation(mutationId, catatan)
      if (result.success) {
        await fetchMutations()
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

  return {
    list,
    isLoading,
    error,
    riwayatMasuk,
    riwayatKeluar,
    pendingCount,
    fetchMutations,
    createMutation,
    approveMutation,
    rejectMutation,
  }
})
