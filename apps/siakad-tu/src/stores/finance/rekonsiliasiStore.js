import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { useAuthStore } from '@/stores/authStore'
import { rekonsiliasiAdapter } from '@/adapters/api/keuangan/rekonsiliasiAdapter'

export const useRekonsiliasiStore = defineStore('keuangan-rekonsiliasi', () => {
  const authStore = useAuthStore()
  const currentUser = computed(() => authStore.user || authStore.currentUser || {})
  const schoolID = computed(() => currentUser.value.school_id || '11111111-1111-1111-1111-111111111111')
  const periodID = computed(() => currentUser.value.academic_period_id || '20261')
  const userID = computed(() => currentUser.value.id || currentUser.value.user_id || '')

  function baseParams(extra = {}) {
    return { school_id: schoolID.value, academic_period_id: periodID.value, ...extra }
  }

  const rekonList = ref([])
  const rekonDetail = ref(null)
  const rekonLoading = ref(false)
  const rekonPagination = ref({ total: 0, limit: 20, offset: 0 })

  async function listRekonsiliasi(params = {}) {
    rekonLoading.value = true
    try {
      const res = await rekonsiliasiAdapter.listRekonsiliasi(baseParams(params))
      rekonList.value = res.data || []
      if (res.pagination) rekonPagination.value = res.pagination
      return res
    } finally {
      rekonLoading.value = false
    }
  }

  async function getRekonsiliasi(id) {
    rekonLoading.value = true
    try {
      const res = await rekonsiliasiAdapter.getRekonsiliasi(id, baseParams())
      rekonDetail.value = res.data
      return res
    } finally {
      rekonLoading.value = false
    }
  }

  /**
   * Jalankan rekonsiliasi baru
   * @param {Object} payload — { rekening_bank_id, periode_bulan, periode_tahun, mutasi_items: [...], catatan }
   */
  async function jalankanRekonsiliasi(payload) {
    rekonLoading.value = true
    try {
      const res = await rekonsiliasiAdapter.jalankanRekonsiliasi(payload, baseParams())
      rekonList.value = [res.data, ...rekonList.value]
      rekonDetail.value = res.data
      return res
    } finally {
      rekonLoading.value = false
    }
  }

  async function deleteRekonsiliasi(id) {
    rekonLoading.value = true
    try {
      await rekonsiliasiAdapter.deleteRekonsiliasi(id, baseParams())
      rekonList.value = rekonList.value.filter(x => x.id !== id)
      if (rekonDetail.value?.id === id) rekonDetail.value = null
    } finally {
      rekonLoading.value = false
    }
  }

  function reset() {
    rekonList.value = []
    rekonDetail.value = null
  }

  return {
    schoolID, periodID, userID,
    rekonList, rekonDetail, rekonLoading, rekonPagination,
    listRekonsiliasi, getRekonsiliasi, jalankanRekonsiliasi, deleteRekonsiliasi,
    reset,
  }
})
