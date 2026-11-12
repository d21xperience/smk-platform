// apps/siakad-tu/src/stores/kesiswaan/pklPlottingStore.js

import { defineStore } from 'pinia'
import { ref } from 'vue'
import { pklPlottingRealAdapter } from '@/adapters/api/pklPlottingRealAdapter.js'

export const usePklPlottingStore = defineStore('pklPlotting', () => {
  const list = ref([])
  const loading = ref(false)
  const error = ref(null)

  async function fetchList() {
    if (loading.value) return
    loading.value = true
    error.value = null
    try {
      const result = await pklPlottingRealAdapter.list()
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
    return await pklPlottingRealAdapter.create(payload)
  }

  async function update(id, payload) {
    return await pklPlottingRealAdapter.update(id, payload)
  }

  async function remove(id) {
    return await pklPlottingRealAdapter.delete(id)
  }

  return {
    list,
    loading,
    error,
    fetchList,
    create,
    update,
    remove,
  }
})
