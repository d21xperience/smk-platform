// apps/siakad-tu/src/stores/kesiswaan/alumniTrackingStore.js

import { defineStore } from 'pinia'
import { ref } from 'vue'
import { alumniTrackingRealAdapter } from '@/adapters/api/alumniTrackingRealAdapter.js'

export const useAlumniTrackingStore = defineStore('alumniTracking', () => {
  const list = ref([])
  const loading = ref(false)
  const error = ref(null)

  async function fetchList() {
    if (loading.value) return
    loading.value = true
    error.value = null
    try {
      const result = await alumniTrackingRealAdapter.list()
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

  async function updateAlumni(id, payload) {
    return await alumniTrackingRealAdapter.update(id, payload)
  }

  return {
    list,
    loading,
    error,
    fetchList,
    updateAlumni,
  }
})
