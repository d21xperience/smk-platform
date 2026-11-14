import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { useAuthStore } from '@/stores/authStore'
import { transaksiAdapter } from '@/adapters/api/keuangan/transaksiAdapter'

export const useTransaksiStore = defineStore('keuangan-transaksi', () => {
  const authStore = useAuthStore()
  const currentUser = computed(() => authStore.user || authStore.currentUser || {})
  const schoolID = computed(() => currentUser.value.school_id || '11111111-1111-1111-1111-111111111111')
  const periodID = computed(() => currentUser.value.academic_period_id || '20261')

  function baseParams(extra = {}) {
    return { school_id: schoolID.value, academic_period_id: periodID.value, ...extra }
  }

  // ============ TAGIHAN ============
  const tagihan = ref([])
  const tagihanDetail = ref(null)
  const tagihanLoading = ref(false)
  const tagihanPagination = ref({ total: 0, limit: 50, offset: 0 })

  async function listTagihan(params = {}) {
    tagihanLoading.value = true
    try {
      const res = await transaksiAdapter.listTagihan(baseParams(params))
      tagihan.value = res.data || []
      if (res.pagination) tagihanPagination.value = res.pagination
      return res
    } finally {
      tagihanLoading.value = false
    }
  }

  async function getTagihan(id) {
    tagihanLoading.value = true
    try {
      const res = await transaksiAdapter.getTagihan(id, baseParams())
      tagihanDetail.value = res.data
      return res
    } finally {
      tagihanLoading.value = false
    }
  }

  // ============ PENERIMAAN IURAN (auto BKM + auto BKU) ============
  const penerimaan = ref([])
  const penerimaanDetail = ref(null)
  const penerimaanLoading = ref(false)
  const penerimaanPagination = ref({ total: 0, limit: 50, offset: 0 })

  async function listPenerimaan(params = {}) {
    penerimaanLoading.value = true
    try {
      const res = await transaksiAdapter.listPenerimaan(baseParams(params))
      penerimaan.value = res.data || []
      if (res.pagination) penerimaanPagination.value = res.pagination
      return res
    } finally {
      penerimaanLoading.value = false
    }
  }

  async function createPenerimaan(payload) {
    penerimaanLoading.value = true
    try {
      const res = await transaksiAdapter.createPenerimaan(payload, baseParams())
      // Backend return { penerimaan, tagihan, bukti_kas } — refresh list
      await listPenerimaan()
      return res
    } finally {
      penerimaanLoading.value = false
    }
  }

  async function getPenerimaan(noBKM) {
    penerimaanLoading.value = true
    try {
      const res = await transaksiAdapter.getPenerimaan(noBKM, baseParams())
      penerimaanDetail.value = res.data
      return res
    } finally {
      penerimaanLoading.value = false
    }
  }

  // ============ BUKU KAS UMUM ============
  const bku = ref([])
  const bkuLoading = ref(false)
  const bkuPagination = ref({ total: 0, limit: 100, offset: 0 })
  const saldoKas = ref(0)

  async function listBKU(params = {}) {
    bkuLoading.value = true
    try {
      const res = await transaksiAdapter.listBKU(baseParams(params))
      bku.value = res.data || []
      if (res.pagination) bkuPagination.value = res.pagination
      return res
    } finally {
      bkuLoading.value = false
    }
  }

  async function getSaldo(params = {}) {
    const res = await transaksiAdapter.getBKUSaldo(baseParams(params))
    saldoKas.value = res.data?.saldo ?? 0
    return res
  }

  function reset() {
    tagihan.value = []
    tagihanDetail.value = null
    penerimaan.value = []
    penerimaanDetail.value = null
    bku.value = []
    saldoKas.value = 0
  }

  return {
    schoolID, periodID,
    // Tagihan
    tagihan, tagihanDetail, tagihanLoading, tagihanPagination, listTagihan, getTagihan,
    // Penerimaan
    penerimaan, penerimaanDetail, penerimaanLoading, penerimaanPagination, listPenerimaan, createPenerimaan, getPenerimaan,
    // BKU
    bku, bkuLoading, bkuPagination, saldoKas, listBKU, getSaldo,
    // Global
    reset,
  }
})
