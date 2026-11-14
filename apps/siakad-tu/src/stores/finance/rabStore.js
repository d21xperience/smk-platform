import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { useAuthStore } from '@/stores/authStore'
import { rabAdapter } from '@/adapters/api/keuangan/rabAdapter'

export const useRabStore = defineStore('keuangan-rab', () => {
  const authStore = useAuthStore()
  const currentUser = computed(() => authStore.user || authStore.currentUser || {})
  const schoolID = computed(() => currentUser.value.school_id || '11111111-1111-1111-1111-111111111111')
  const periodID = computed(() => currentUser.value.academic_period_id || '20261')
  const userID = computed(() => currentUser.value.id || currentUser.value.user_id || '')

  function baseParams(extra = {}) {
    return { school_id: schoolID.value, academic_period_id: periodID.value, ...extra }
  }

  // ============ RAB ============
  const rabList = ref([])
  const rabDetail = ref(null)
  const rabLoading = ref(false)
  const rabPagination = ref({ total: 0, limit: 20, offset: 0 })

  async function listRab(params = {}) {
    rabLoading.value = true
    try {
      const res = await rabAdapter.listRab(baseParams(params))
      rabList.value = res.data || []
      if (res.pagination) rabPagination.value = res.pagination
      return res
    } finally {
      rabLoading.value = false
    }
  }

  async function getRab(id) {
    rabLoading.value = true
    try {
      const res = await rabAdapter.getRab(id, baseParams())
      rabDetail.value = res.data
      return res
    } finally {
      rabLoading.value = false
    }
  }

  async function createRab(payload) {
    rabLoading.value = true
    try {
      const body = { ...payload, disusun_oleh: payload.disusun_oleh || userID.value }
      const res = await rabAdapter.createRab(body, baseParams())
      rabList.value = [res.data, ...rabList.value]
      return res
    } finally {
      rabLoading.value = false
    }
  }

  async function updateRab(id, payload) {
    rabLoading.value = true
    try {
      const res = await rabAdapter.updateRab(id, payload, baseParams())
      const idx = rabList.value.findIndex(x => x.id === id)
      if (idx >= 0) rabList.value[idx] = res.data
      rabDetail.value = res.data
      return res
    } finally {
      rabLoading.value = false
    }
  }

  async function sahkanRab(id, catatan = '') {
    rabLoading.value = true
    try {
      const res = await rabAdapter.sahkanRab(id, {
        disahkan_oleh: userID.value,
        catatan,
      }, baseParams())
      const idx = rabList.value.findIndex(x => x.id === id)
      if (idx >= 0) rabList.value[idx] = res.data
      rabDetail.value = res.data
      return res
    } finally {
      rabLoading.value = false
    }
  }

  async function deleteRab(id) {
    rabLoading.value = true
    try {
      await rabAdapter.deleteRab(id, baseParams())
      rabList.value = rabList.value.filter(x => x.id !== id)
    } finally {
      rabLoading.value = false
    }
  }

  // ============ REALISASI ============
  const realisasi = ref(null)
  const realisasiLoading = ref(false)

  async function getRealisasi(tahun, sampaiBulan = null) {
    realisasiLoading.value = true
    try {
      const params = { sampai_bulan: sampaiBulan || new Date().getMonth() + 1 }
      const res = await rabAdapter.getRealisasi(tahun, baseParams(params))
      realisasi.value = res.data
      return res
    } finally {
      realisasiLoading.value = false
    }
  }

  function reset() {
    rabList.value = []
    rabDetail.value = null
    realisasi.value = null
  }

  return {
    schoolID, periodID,
    rabList, rabDetail, rabLoading, rabPagination,
    listRab, getRab, createRab, updateRab, sahkanRab, deleteRab,
    realisasi, realisasiLoading, getRealisasi,
    reset,
  }
})
