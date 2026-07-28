import { defineStore } from 'pinia'

export const useContextStore = defineStore('context', {
  state: () => ({
    // Operational Context: dipilih sekali, digunakan untuk transaksi (absensi, nilai, jurnal, rapor)
    operational: {
      academicYearId: null,
      academicYearName: '',
      semesterId: null,
      semesterName: '',
      academicPeriodId: null,
      academicPeriodName: '',
    },

    // Data master untuk dropdown (diisi dari API/mock)
    academicYears: [],
    semesters: [],
    academicPeriods: [],

    // Query Context: filter sementara untuk pencarian/laporan historis
    query: {
      academicYearId: null,
      semesterId: null,
      academicPeriodId: null,
    },

    // Status loading
    loading: {
      academicYears: false,
      semesters: false,
      academicPeriods: false,
    },
  }),

  getters: {
    // Cek apakah operational context sudah di-set
    isOperationalContextReady: (state) =>
      !!state.operational.academicYearId && !!state.operational.semesterId,

    activeAcademicYearId: (state) => state.operational.academicYearId,
    activeSemesterId: (state) => state.operational.semesterId,
    activeAcademicPeriodId: (state) => state.operational.academicPeriodId,
  },

  actions: {
    /**
     * Injeksi ContextService dari luar.
     * @param {import('src/services/ContextService').ContextService} service
     */
    setService(service) {
      this.contextService = service
    },

    // --- Data Master ---

    async loadAcademicYears() {
      if (!this.contextService) throw new Error('ContextService belum diinisialisasi')
      this.loading.academicYears = true
      try {
        this.academicYears = await this.contextService.fetchAcademicYears()
      } finally {
        this.loading.academicYears = false
      }
    },

    async loadSemesters(academicYearId) {
      if (!this.contextService) throw new Error('ContextService belum diinisialisasi')
      this.loading.semesters = true
      try {
        this.semesters = await this.contextService.fetchSemesters(academicYearId)
      } finally {
        this.loading.semesters = false
      }
    },

    async loadAcademicPeriods(semesterId) {
      if (!this.contextService) throw new Error('ContextService belum diinisialisasi')
      this.loading.academicPeriods = true
      try {
        this.academicPeriods = await this.contextService.fetchAcademicPeriods(semesterId)
      } finally {
        this.loading.academicPeriods = false
      }
    },

    // --- Operational Context ---

    /**
     * Set operational context dari model yang dipilih.
     * @param {import('src/models/AcademicYear').AcademicYear} academicYear
     * @param {import('src/models/Semester').Semester} semester
     * @param {import('src/models/AcademicPeriod').AcademicPeriod} [academicPeriod]
     */
    setOperationalContext(academicYear, semester, academicPeriod = null) {
      if (!this.contextService) throw new Error('ContextService belum diinisialisasi')
      this.operational = this.contextService.buildOperationalContext(
        academicYear,
        semester,
        academicPeriod,
      )
    },

    // --- Query Context ---

    /**
     * Set query context (filter historis). Tidak memengaruhi operational.
     * @param {{ academicYearId?: number, semesterId?: number, academicPeriodId?: number }} filter
     */
    setQueryContext(filter) {
      this.query = { ...filter }
    },

    /**
     * Reset query context ke default.
     */
    clearQueryContext() {
      this.query = {
        academicYearId: null,
        semesterId: null,
        academicPeriodId: null,
      }
    },
  },
})
