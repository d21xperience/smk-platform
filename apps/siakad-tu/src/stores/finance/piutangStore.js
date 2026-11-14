import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { useAuthStore } from '@/stores/authStore'
import { piutangAdapter } from '@/adapters/api/keuangan/piutangAdapter'

export const usePiutangStore = defineStore('keuangan-piutang', () => {
  const authStore = useAuthStore()
  const currentUser = computed(() => authStore.user || authStore.currentUser || {})
  const schoolID = computed(() => currentUser.value.school_id || '11111111-1111-1111-1111-111111111111')
  const periodID = computed(() => currentUser.value.academic_period_id || '20261')

  function baseParams(extra = {}) {
    return { school_id: schoolID.value, academic_period_id: periodID.value, ...extra }
  }

  // ============ LIST FLAT (per tagihan) ============
  const piutang = ref([])
  const piutangLoading = ref(false)
  const piutangPagination = ref({ total: 0, limit: 50, offset: 0 })

  async function listPiutang(params = {}) {
    piutangLoading.value = true
    try {
      const res = await piutangAdapter.listPiutang(baseParams(params))
      piutang.value = res.data || []
      if (res.pagination) piutangPagination.value = res.pagination
      return res
    } finally {
      piutangLoading.value = false
    }
  }

  // ============ LIST PER-SISWA (grouped) ============
  const perSiswa = ref([])
  const perSiswaLoading = ref(false)
  const perSiswaPagination = ref({ total: 0, limit: 50, offset: 0 })

  async function listPerSiswa(params = {}) {
    perSiswaLoading.value = true
    try {
      const res = await piutangAdapter.listPiutangPerSiswa(baseParams(params))
      perSiswa.value = res.data || []
      if (res.pagination) perSiswaPagination.value = res.pagination
      return res
    } finally {
      perSiswaLoading.value = false
    }
  }

  // ============ DETAIL PER NIS ============
  const detailSiswa = ref(null)
  const detailLoading = ref(false)

  async function getPiutangByNIS(nis) {
    detailLoading.value = true
    try {
      const res = await piutangAdapter.getPiutangByNIS(nis, baseParams())
      detailSiswa.value = res.data
      return res
    } finally {
      detailLoading.value = false
    }
  }

  // ============ AGING REPORT ============
  const agingReport = ref(null)
  const agingLoading = ref(false)

  async function getAgingReport() {
    agingLoading.value = true
    try {
      const res = await piutangAdapter.getAgingReport(baseParams())
      agingReport.value = res.data
      return res
    } finally {
      agingLoading.value = false
    }
  }

  function reset() {
    piutang.value = []
    perSiswa.value = []
    detailSiswa.value = null
    agingReport.value = null
  }

  return {
    schoolID, periodID,
    // Flat list
    piutang, piutangLoading, piutangPagination, listPiutang,
    // Per-siswa
    perSiswa, perSiswaLoading, perSiswaPagination, listPerSiswa,
    // Detail
    detailSiswa, detailLoading, getPiutangByNIS,
    // Aging
    agingReport, agingLoading, getAgingReport,
    reset,
  }
})
