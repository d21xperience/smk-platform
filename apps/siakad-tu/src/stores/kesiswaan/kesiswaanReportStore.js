// apps/siakad-tu/src/stores/kesiswaan/kesiswaanReportStore.js

import { defineStore } from 'pinia'
import { ref } from 'vue'
import { kesiswaanReportRealAdapter } from '@/adapters/api/kesiswaanReportRealAdapter.js'

export const useKesiswaanReportStore = defineStore('kesiswaanReport', () => {
  const stats = ref([])
  const totalMasuk = ref(0)
  const totalKeluar = ref(0)
  const totalAkhir = ref(0)
  const loading = ref(false)
  const error = ref(null)

  async function fetchStatistik() {
    if (loading.value) return
    loading.value = true
    error.value = null
    try {
      const result = await kesiswaanReportRealAdapter.getStatistik()
      if (result.success) {
        stats.value = result.data.stats || []
        totalMasuk.value = result.data.totalMasuk || 0
        totalKeluar.value = result.data.totalKeluar || 0
        totalAkhir.value = result.data.totalAkhir || 0
      } else {
        error.value = result.error
        stats.value = []
      }
    } finally {
      loading.value = false
    }
  }

  async function downloadExcel() {
    return await kesiswaanReportRealAdapter.downloadExcel()
  }

  async function downloadPDF() {
    return await kesiswaanReportRealAdapter.downloadPDF()
  }

  return {
    stats,
    totalMasuk,
    totalKeluar,
    totalAkhir,
    loading,
    error,
    fetchStatistik,
    downloadExcel,
    downloadPDF,
  }
})
