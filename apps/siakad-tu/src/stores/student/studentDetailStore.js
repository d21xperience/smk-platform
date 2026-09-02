// apps/siakad-tu/src/stores/student/studentDetailStore.js

import { defineStore } from 'pinia'

/**
 * studentDetailStore — Mengelola state detail siswa per ID.
 *
 * Menggunakan cache Map<studentId, data> untuk menghindari fetch berulang.
 *
 * STATE ONLY. Tidak ada logika bisnis.
 */
export const useStudentDetailStore = defineStore('studentDetail', {
  state: () => ({
    /** @type {Map<string, Object>} Cache data siswa by studentId */
    cache: new Map(),
    /** @type {Set<string>} IDs yang sedang di-fetch */
    loadingIds: new Set(),
    /** @type {Map<string, Object>} Error per studentId */
    errorMap: new Map(),
  }),

  getters: {
    /**
     * Ambil data siswa by ID (reactive via Map)
     */
    getById: (state) => (studentId) => {
      return state.cache.get(studentId) || null
    },

    /**
     * Cek apakah siswa sedang di-fetch
     */
    isLoading: (state) => (studentId) => {
      return state.loadingIds.has(studentId)
    },

    /**
     * Ambil error untuk siswa tertentu
     */
    getError: (state) => (studentId) => {
      return state.errorMap.get(studentId) || null
    },

    /**
     * Jumlah siswa di-cache
     */
    cacheSize: (state) => state.cache.size,
  },

  actions: {
    /**
     * Set data siswa ke cache
     */
    setItem(studentId, data) {
      // Pinia perlu plain object, bukan Map instance
      // Kita gunakan object biasa untuk reactivity
      this.cache = new Map(this.cache)
      this.cache.set(studentId, data)
      this.loadingIds = new Set(this.loadingIds)
      this.loadingIds.delete(studentId)
      this.errorMap = new Map(this.errorMap)
      this.errorMap.delete(studentId)
    },

    /**
     * Set loading state untuk studentId tertentu
     */
    setLoading(studentId, loading) {
      this.loadingIds = new Set(this.loadingIds)
      if (loading) {
        this.loadingIds.add(studentId)
      } else {
        this.loadingIds.delete(studentId)
      }
    },

    /**
     * Set error untuk studentId tertentu
     */
    setError(studentId, error) {
      this.errorMap = new Map(this.errorMap)
      this.errorMap.set(studentId, error)
      this.loadingIds = new Set(this.loadingIds)
      this.loadingIds.delete(studentId)
    },

    /**
     * Invalidate cache untuk studentId tertentu
     * (akan force refetch saat diakses lagi)
     */
    invalidate(studentId) {
      this.cache = new Map(this.cache)
      this.cache.delete(studentId)
    },

    /**
     * Clear seluruh cache (saat context berubah)
     */
    reset() {
      this.cache = new Map()
      this.loadingIds = new Set()
      this.errorMap = new Map()
    },
  },
})
