import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { useAuthStore } from '@/stores/authStore'
import { penutupanAdapter } from '@/adapters/api/keuangan/penutupanAdapter'

export const usePenutupanStore = defineStore('keuangan-penutupan', () => {
  const authStore = useAuthStore()
  const currentUser = computed(() => authStore.user || authStore.currentUser || {})
  const schoolID = computed(() => currentUser.value.school_id || '11111111-1111-1111-1111-111111111111')
  const periodID = computed(() => currentUser.value.academic_period_id || '20261')
  const userID = computed(() => currentUser.value.id || currentUser.value.user_id || '')

  function baseParams(extra = {}) {
    return { school_id: schoolID.value, academic_period_id: periodID.value, ...extra }
  }

  const penutupanList = ref([])
  const penutupanDetail = ref(null)
  const penutupanLoading = ref(false)
  const penutupanPagination = ref({ total: 0, limit: 20, offset: 0 })

  const previewSnapshot = ref(null)
  const previewLoading = ref(false)

  async function listPenutupan(params = {}) {
    penutupanLoading.value = true
    try {
      const res = await penutupanAdapter.listPenutupan(baseParams(params))
      penutupanList.value = res.data || []
      if (res.pagination) penutupanPagination.value = res.pagination
      return res
    } finally {
      penutupanLoading.value = false
    }
  }

  async function getPenutupan(id) {
    penutupanLoading.value = true
    try {
      const res = await penutupanAdapter.getPenutupan(id, baseParams())
      penutupanDetail.value = res.data
      return res
    } finally {
      penutupanLoading.value = false
    }
  }

  /**
   * Preview snapshot sebelum tutup buku
   * @param {number} bulan
   * @param {number} tahun
   */
  async function loadPreview(bulan, tahun) {
    previewLoading.value = true
    try {
      const res = await penutupanAdapter.previewSnapshot(baseParams({
        periode_bulan: bulan,
        periode_tahun: tahun,
      }))
      previewSnapshot.value = res.data
      return res
    } finally {
      previewLoading.value = false
    }
  }

  /**
   * Tutup buku
   * @param {Object} payload — { periode_bulan, periode_tahun, saldo_fisik_kas, catatan }
   */
  async function tutupBuku(payload) {
    penutupanLoading.value = true
    try {
      const body = { ...payload, closed_by: payload.closed_by || userID.value }
      const res = await penutupanAdapter.tutupBuku(body, baseParams())
      penutupanList.value = [res.data, ...penutupanList.value]
      penutupanDetail.value = res.data
      return res
    } finally {
      penutupanLoading.value = false
    }
  }

  /**
   * Sign penutupan (TTD Kepsek)
   * @param {string} id
   * @param {string} pin — 6 digit
   */
  async function signPenutupan(id, pin) {
    penutupanLoading.value = true
    try {
      const res = await penutupanAdapter.signPenutupan(id, { pin }, baseParams())
      const idx = penutupanList.value.findIndex(x => x.id === id)
      if (idx >= 0) penutupanList.value[idx] = res.data
      penutupanDetail.value = res.data
      return res
    } finally {
      penutupanLoading.value = false
    }
  }

  function reset() {
    penutupanList.value = []
    penutupanDetail.value = null
    previewSnapshot.value = null
  }

  return {
    schoolID, periodID, userID,
    penutupanList, penutupanDetail, penutupanLoading, penutupanPagination,
    listPenutupan, getPenutupan,
    previewSnapshot, previewLoading, loadPreview,
    tutupBuku, signPenutupan,
    reset,
  }
})
