// FILE: src/stores/homeroomAttendanceStore.js
// STATUS: NEW
// STATUS IMPLEMENTASI: COMPLETE

import { defineStore } from 'pinia'
import { homeroomAttendanceService } from '../boot/services.js'

/**
 * Homeroom Attendance Store
 * Mengelola state rekap kehadiran wali kelas.
 * Pola mengikuti homeroomBillingStore:
 * - Tidak import store lain
 * - Action menerima parameter eksplisit dari composable
 * - Store memanggil Application Service, bukan Axios
 */
export const useHomeroomAttendanceStore = defineStore('homeroomAttendance', {
  state: () => ({
    classId: null,
    className: null,
    months: [],
    summary: null,
    students: [],
    loading: false,
    error: null,
  }),

  getters: {
    hasData: (state) => state.students.length > 0,
    hasSummary: (state) => !!state.summary,
    hasMonths: (state) => state.months.length > 0,
  },

  actions: {
    /**
     * Memuat rekap kehadiran siswa berdasarkan parameter eksplisit.
     * @param {Object} params - Parameter dari composable
     * @param {string} params.userId - User ID wali kelas
     * @param {string} params.schoolId - School ID
     * @param {string} params.academicYearId - Academic Year ID
     * @param {string} params.semesterId - Semester ID
     */
    async loadAttendanceSummary({ userId, schoolId, academicYearId, semesterId }) {
      this.loading = true
      this.error = null
      try {
        const result = await homeroomAttendanceService.getAttendanceSummary({
          userId,
          schoolId,
          academicYearId,
          semesterId,
        })
        this.classId = result.classId
        this.className = result.className
        this.months = result.months
        this.summary = result.summary
        this.students = result.students
      } catch (err) {
        this.error = err.message
        this.classId = null
        this.className = null
        this.months = []
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
      this.months = []
      this.summary = null
      this.students = []
      this.error = null
    },
  },
})
