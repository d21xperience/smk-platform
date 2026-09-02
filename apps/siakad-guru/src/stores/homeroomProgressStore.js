// FILE: src/stores/homeroomProgressStore.js
// STATUS: NEW
// STATUS IMPLEMENTASI: COMPLETE

import { defineStore } from 'pinia'
import { homeroomProgressService } from '../boot/services.js'

/**
 * Homeroom Progress Store
 * Mengelola state rekap nilai wali kelas.
 * Pola mengikuti homeroomBillingStore:
 * - Tidak import store lain
 * - Action menerima parameter eksplisit dari composable
 * - Store memanggil Application Service, bukan Axios
 */
export const useHomeroomProgressStore = defineStore('homeroomProgress', {
  state: () => ({
    classId: null,
    className: null,
    summary: null,
    students: [],
    loading: false,
    error: null,
  }),

  getters: {
    hasData: (state) => state.students.length > 0,
    hasSummary: (state) => !!state.summary,
  },

  actions: {
    /**
     * Memuat rekap nilai siswa berdasarkan parameter eksplisit.
     * @param {Object} params - Parameter dari composable
     * @param {string} params.userId - User ID wali kelas
     * @param {string} params.schoolId - School ID
     * @param {string} params.academicYearId - Academic Year ID
     * @param {string} params.semesterId - Semester ID
     */
    async loadProgressSummary({ userId, schoolId, academicYearId, semesterId }) {
      this.loading = true
      this.error = null
      try {
        const result = await homeroomProgressService.getProgressSummary({
          userId,
          schoolId,
          academicYearId,
          semesterId,
        })
        this.classId = result.classId
        this.className = result.className
        this.summary = result.summary
        this.students = result.students
      } catch (err) {
        this.error = err.message
        this.classId = null
        this.className = null
        this.summary = null
        this.students = []
      } finally {
        this.loading = false
      }
    },

    /**
     * Mengosongkan state store.
     */
    clearData() {
      this.classId = null
      this.className = null
      this.summary = null
      this.students = []
      this.error = null
    },
  },
})
