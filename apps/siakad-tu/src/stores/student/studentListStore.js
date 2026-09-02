// apps/siakad-tu/src/stores/student/studentListStore.js

import { defineStore } from 'pinia'

/**
 * studentListStore — Mengelola state daftar siswa.
 *
 * STATE ONLY. Tidak ada logika bisnis.
 *
 * Tanggung jawab:
 * - Menyimpan daftar siswa yang sedang ditampilkan
 * - Menyimpan pagination state
 * - Menyimpan filter state
 * - Menyimpan loading & error state
 *
 * TIDAK BOLEH:
 * - Validasi data
 * - Menghitung aturan bisnis
 * - Memanggil Engine langsung
 */
export const useStudentListStore = defineStore('studentList', {
  state: () => ({
    /** @type {Array<Object>} Daftar siswa */
    items: [],
    /** @type {number} Total siswa (untuk pagination) */
    total: 0,
    /** @type {number} Halaman saat ini */
    page: 1,
    /** @type {number} Jumlah item per halaman */
    limit: 20,
    /** @type {number} Total halaman */
    totalPages: 0,
    /** @type {Object} Filter aktif */
    filters: {
      status: '',
      classId: '',
      search: '',
    },
    /** @type {boolean} Loading state */
    loading: false,
    /** @type {Object|null} Error state */
    error: null,
    /** @type {boolean} Apakah data sudah pernah di-fetch */
    initialized: false,
  }),

  getters: {
    /**
     * Apakah masih ada halaman berikutnya
     */
    hasMore: (state) => state.page < state.totalPages,

    /**
     * Apakah list kosong
     */
    isEmpty: (state) => state.items.length === 0 && !state.loading,

    /**
     * Apakah ada filter aktif
     */
    hasActiveFilters: (state) => {
      return !!(state.filters.status || state.filters.classId || state.filters.search)
    },

    /**
     * Ringkasan statistik
     */
    stats: (state) => ({
      total: state.total,
      active: state.items.filter((i) => i.status === 'ACTIVE').length,
      currentPage: state.page,
      totalPages: state.totalPages,
    }),
  },

  actions: {
    /**
     * Set items dari hasil fetch
     * Dipanggil oleh Composable setelah Service return data
     */
    setItems({ items, total, page, totalPages }) {
      this.items = items
      this.total = total
      this.page = page
      this.totalPages = totalPages
      this.initialized = true
    },

    /**
     * Set loading state
     */
    setLoading(loading) {
      this.loading = loading
      if (loading) this.error = null
    },

    /**
     * Set error state
     */
    setError(error) {
      this.error = error
      this.loading = false
    },

    /**
     * Clear error
     */
    clearError() {
      this.error = null
    },

    /**
     * Update filter dan reset ke halaman 1
     */
    setFilters(filters) {
      this.filters = { ...this.filters, ...filters }
      this.page = 1
      this.initialized = false // force refetch
    },

    /**
     * Navigasi ke halaman tertentu
     */
    setPage(page) {
      if (page < 1 || page > this.totalPages) return
      this.page = page
      this.initialized = false // force refetch
    },

    /**
     * Tambah item ke list (setelah create berhasil)
     */
    addItem(item) {
      // Tambahkan di awal list
      this.items.unshift(item)
      this.total += 1
    },

    /**
     * Update item di list (setelah update berhasil)
     */
    updateItem(item) {
      const index = this.items.findIndex((i) => i.studentId === item.studentId)
      if (index !== -1) {
        this.items[index] = item
      }
    },

    /**
     * Hapus item dari list (setelah transfer/graduate)
     */
    removeItem(studentId) {
      this.items = this.items.filter((i) => i.studentId !== studentId)
      this.total -= 1
    },

    /**
     * Reset seluruh state (saat context berubah)
     */
    reset() {
      this.items = []
      this.total = 0
      this.page = 1
      this.totalPages = 0
      this.filters = { status: '', classId: '', search: '' }
      this.loading = false
      this.error = null
      this.initialized = false
    },
  },
})
