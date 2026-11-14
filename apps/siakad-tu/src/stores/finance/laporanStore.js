import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { useAuthStore } from '@/stores/authStore'
import { laporanAdapter } from '@/adapters/api/keuangan/laporanAdapter'

export const useLaporanStore = defineStore('keuangan-laporan', () => {
  const authStore = useAuthStore()
  const currentUser = computed(() => authStore.user || authStore.currentUser || {})
  const schoolID = computed(() => currentUser.value.school_id || '11111111-1111-1111-1111-111111111111')
  const periodID = computed(() => currentUser.value.academic_period_id || '20261')

  function baseParams(extra = {}) {
    return { school_id: schoolID.value, academic_period_id: periodID.value, ...extra }
  }

  // ============ NERACA ============
  const neraca = ref(null)
  const neracaLoading = ref(false)

  async function getNeraca(params = {}) {
    neracaLoading.value = true
    try {
      const res = await laporanAdapter.getNeraca(baseParams(params))
      neraca.value = res.data
      return res
    } finally {
      neracaLoading.value = false
    }
  }

  async function getNeracaSnapshot(penutupanID) {
    return getNeraca({ mode: 'SNAPSHOT', penutupan_id: penutupanID })
  }

  // ============ POSISI KEUANGAN (strict) ============
  const posisiKeuangan = ref(null)
  const posisiLoading = ref(false)

  async function getPosisiKeuangan(params = {}) {
    posisiLoading.value = true
    try {
      const res = await laporanAdapter.getPosisiKeuangan(baseParams(params))
      posisiKeuangan.value = res.data
      return res
    } finally {
      posisiLoading.value = false
    }
  }

  // ============ LRA ============
  const lra = ref(null)
  const lraLoading = ref(false)

  async function getLRA(tahun, sampaiBulan = null) {
    lraLoading.value = true
    try {
      const params = { tahun, sampai_bulan: sampaiBulan || new Date().getMonth() + 1 }
      const res = await laporanAdapter.getLRA(baseParams(params))
      lra.value = res.data
      return res
    } finally {
      lraLoading.value = false
    }
  }

  // ============ DASHBOARD ============
  const dashboard = ref(null)
  const dashboardLoading = ref(false)

  async function getDashboard() {
    dashboardLoading.value = true
    try {
      const res = await laporanAdapter.getDashboard(baseParams())
      dashboard.value = res.data
      return res
    } finally {
      dashboardLoading.value = false
    }
  }

  // ============ RINGKASAN BULAN ============
  const ringkasanBulan = ref(null)
  const ringkasanLoading = ref(false)

  async function getRingkasanBulan(bulan, tahun) {
    ringkasanLoading.value = true
    try {
      const res = await laporanAdapter.getRingkasanBulan(baseParams({ bulan, tahun }))
      ringkasanBulan.value = res.data
      return res
    } finally {
      ringkasanLoading.value = false
    }
  }

  function reset() {
    neraca.value = null
    posisiKeuangan.value = null
    lra.value = null
    dashboard.value = null
    ringkasanBulan.value = null
  }

  return {
    schoolID, periodID,
    // Neraca
    neraca, neracaLoading, getNeraca, getNeracaSnapshot,
    // Posisi
    posisiKeuangan, posisiLoading, getPosisiKeuangan,
    // LRA
    lra, lraLoading, getLRA,
    // Dashboard
    dashboard, dashboardLoading, getDashboard,
    // Ringkasan
    ringkasanBulan, ringkasanLoading, getRingkasanBulan,
    // Global
    reset,
  }
})
