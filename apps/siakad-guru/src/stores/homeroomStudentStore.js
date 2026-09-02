// FILE: src/stores/homeroomStudentStore.js
// STATUS: MODIFY (koreksi pola arsitektur)
// STATUS IMPLEMENTASI: COMPLETE

import { defineStore } from 'pinia'
import { homeroomStudentService } from '../boot/services.js'

export const useHomeroomStudentStore = defineStore('homeroomStudent', {
  state: () => ({
    students: [],
    classId: null,
    className: null,
    loading: false,
    error: null,
  }),

  getters: {
    activeStudents: (state) => state.students.filter((s) => s.status === 'ACTIVE'),
    totalActiveStudents: (state) => state.students.filter((s) => s.status === 'ACTIVE').length,
    totalStudents: (state) => state.students.length,
  },

  actions: {
    async loadStudents({ userId, schoolId, academicYearId }) {
      this.loading = true
      this.error = null
      try {
        const result = await homeroomStudentService.getHomeroomStudents({
          userId,
          schoolId,
          academicYearId,
        })
        this.classId = result.classId
        this.className = result.className
        this.students = result.students
      } catch (err) {
        this.error = err.message
        this.students = []
        this.classId = null
        this.className = null
      } finally {
        this.loading = false
      }
    },

    clearData() {
      this.students = []
      this.classId = null
      this.className = null
      this.error = null
    },
  },
})
