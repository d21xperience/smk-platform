import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

export const useAcademicStore = defineStore('academic', () => {
  // Simpan kode periode terpadu (Default awal: 20261)
  const activePeriodCode = ref(localStorage.getItem('academic_period_code') || '20261')

  // GETTER: Memecah kode periode menjadi teks Label Tahun Ajaran untuk dibaca Pegawai TU
  const activeYearLabel = computed(() => {
    // Ambil 4 digit pertama (misal: "2026")
    const startYear = parseInt(activePeriodCode.value.substring(0, 4))
    return `${startYear}/${startYear + 1}`
  })

  // GETTER: Memecah kode periode untuk mendapatkan teks Semester
  const activeSemesterLabel = computed(() => {
    // Ambil digit terakhir ("1" atau "2")
    const sem = activePeriodCode.value.slice(-1)
    return sem === '1' ? 'Ganjil' : 'Genap'
  })

  // ACTION: Menyimpan kode periode baru dari pilihan TU
  function setActivePeriod(code) {
    activePeriodCode.value = code.toString()
    localStorage.setItem('academic_period_code', code.toString())
  }

  return {
    activePeriodCode,
    activeYearLabel,
    activeSemesterLabel,
    setActivePeriod,
  }
})
