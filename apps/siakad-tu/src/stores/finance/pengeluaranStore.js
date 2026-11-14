import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { useAuthStore } from '@/stores/authStore'
import { pengeluaranAdapter } from '@/adapters/api/keuangan/pengeluaranAdapter'

export const usePengeluaranStore = defineStore('keuangan-pengeluaran', () => {
  const authStore = useAuthStore()
  const currentUser = computed(() => authStore.user || authStore.currentUser || {})
  const schoolID = computed(() => currentUser.value.school_id || '11111111-1111-1111-1111-111111111111')
  const periodID = computed(() => currentUser.value.academic_period_id || '20261')
  const userID = computed(() => currentUser.value.id || currentUser.value.user_id || '')

  function baseParams(extra = {}) {
    return { school_id: schoolID.value, academic_period_id: periodID.value, ...extra }
  }

  // ============ PENGAJUAN ============
  const pengajuan = ref([])
  const pengajuanDetail = ref(null)
  const pengajuanLoading = ref(false)
  const pengajuanPagination = ref({ total: 0, limit: 50, offset: 0 })

  async function listPengajuan(params = {}) {
    pengajuanLoading.value = true
    try {
      const res = await pengeluaranAdapter.listPengajuan(baseParams(params))
      pengajuan.value = res.data || []
      if (res.pagination) pengajuanPagination.value = res.pagination
      return res
    } finally {
      pengajuanLoading.value = false
    }
  }

  async function getPengajuan(id) {
    pengajuanLoading.value = true
    try {
      const res = await pengeluaranAdapter.getPengajuan(id, baseParams())
      pengajuanDetail.value = res.data
      return res
    } finally {
      pengajuanLoading.value = false
    }
  }

  async function createPengajuan(payload) {
    pengajuanLoading.value = true
    try {
      const body = {
        ...payload,
        pemohon_user_id: payload.pemohon_user_id || userID.value,
      }
      const res = await pengeluaranAdapter.createPengajuan(body, baseParams())
      pengajuan.value = [res.data, ...pengajuan.value]
      return res
    } finally {
      pengajuanLoading.value = false
    }
  }

  // Verifikasi Bendahara: aksi = 'SETUJUI' | 'TOLAK' | 'REVISI'
  async function verifikasiBendahara(id, aksi, catatan = '') {
    pengajuanLoading.value = true
    try {
      const res = await pengeluaranAdapter.verifikasiBendahara(id, {
        aksi,
        bendahara_id: userID.value,
        catatan,
      }, baseParams())
      // Update in-list
      const idx = pengajuan.value.findIndex(x => x.id === id)
      if (idx >= 0 && res.data) pengajuan.value[idx] = res.data
      pengajuanDetail.value = res.data
      return res
    } finally {
      pengajuanLoading.value = false
    }
  }

  // Approval Kepala: aksi = 'SETUJUI' | 'TOLAK'
  // Saat SETUJUI → auto-generate BKK (no_bkk di pengajuan, pengeluaran row dibuat)
  async function approvalKepsek(id, aksi, catatan = '') {
    pengajuanLoading.value = true
    try {
      const res = await pengeluaranAdapter.approvalKepsek(id, {
        aksi,
        kepala_id: userID.value,
        catatan,
      }, baseParams())
      const idx = pengajuan.value.findIndex(x => x.id === id)
      if (idx >= 0 && res.data?.pengajuan) pengajuan.value[idx] = res.data.pengajuan
      pengajuanDetail.value = res.data?.pengajuan
      return res
    } finally {
      pengajuanLoading.value = false
    }
  }

  // ============ PENGELUARAN (auto BKK) ============
  const pengeluaran = ref([])
  const pengeluaranDetail = ref(null)
  const pengeluaranLoading = ref(false)
  const pengeluaranPagination = ref({ total: 0, limit: 50, offset: 0 })

  async function listPengeluaran(params = {}) {
    pengeluaranLoading.value = true
    try {
      const res = await pengeluaranAdapter.listPengeluaran(baseParams(params))
      pengeluaran.value = res.data || []
      if (res.pagination) pengeluaranPagination.value = res.pagination
      return res
    } finally {
      pengeluaranLoading.value = false
    }
  }

  async function getPengeluaran(noBKK) {
    pengeluaranLoading.value = true
    try {
      const res = await pengeluaranAdapter.getPengeluaran(noBKK, baseParams())
      pengeluaranDetail.value = res.data
      return res
    } finally {
      pengeluaranLoading.value = false
    }
  }

  function reset() {
    pengajuan.value = []
    pengajuanDetail.value = null
    pengeluaran.value = []
    pengeluaranDetail.value = null
  }

  return {
    schoolID, periodID, userID,
    // Pengajuan
    pengajuan, pengajuanDetail, pengajuanLoading, pengajuanPagination,
    listPengajuan, getPengajuan, createPengajuan, verifikasiBendahara, approvalKepsek,
    // Pengeluaran
    pengeluaran, pengeluaranDetail, pengeluaranLoading, pengeluaranPagination,
    listPengeluaran, getPengeluaran,
    // Global
    reset,
  }
})
