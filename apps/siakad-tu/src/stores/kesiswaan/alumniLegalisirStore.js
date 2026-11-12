// apps/siakad-tu/src/stores/kesiswaan/alumniLegalisirStore.js

import { defineStore } from 'pinia'
import { ref } from 'vue'
import { alumniLegalisirRealAdapter } from '@/adapters/api/alumniLegalisirRealAdapter.js'

export const useAlumniLegalisirStore = defineStore('alumniLegalisir', () => {
  const list = ref([])
  const loading = ref(false)
  const error = ref(null)

  async function fetchList() {
    if (loading.value) return
    loading.value = true
    error.value = null
    try {
      const result = await alumniLegalisirRealAdapter.list()
      if (result.success) {
        list.value = result.data
      } else {
        error.value = result.error
        list.value = []
      }
    } finally {
      loading.value = false
    }
  }

  async function create(payload) {
    return await alumniLegalisirRealAdapter.create(payload)
  }

  async function updateStatus(id, status) {
    return await alumniLegalisirRealAdapter.updateStatus(id, status)
  }

  async function remove(id) {
    return await alumniLegalisirRealAdapter.delete(id)
  }

  return {
    list,
    loading,
    error,
    fetchList,
    create,
    updateStatus,
    remove,
  }
})
