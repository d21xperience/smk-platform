import { defineStore } from 'pinia'
import { progressService } from '../boot/services.js'

export const useProgressStore = defineStore('progress', {
  state: () => ({
    students: [],
    progressRecords: [],
    loading: false,
    error: null,
  }),

  getters: {
    hasStudents: (state) => state.students.length > 0,
    hasProgress: (state) => state.progressRecords.length > 0,
  },

  actions: {
    async loadStudents({ classId }) {
      this.loading = true
      this.error = null
      try {
        this.students = await progressService.loadStudentsByClass({ classId })
      } catch (err) {
        this.error = err.message
      } finally {
        this.loading = false
      }
    },

    async loadProgress({ studentId, academicYearId, semesterId }) {
      this.loading = true
      this.error = null
      try {
        this.progressRecords = await progressService.loadProgressByStudent({
          studentId,
          academicYearId,
          semesterId,
        })
      } catch (err) {
        this.error = err.message
      } finally {
        this.loading = false
      }
    },

    async createProgress({
      studentId,
      studentName,
      classId,
      className,
      type,
      title,
      description,
      date,
      teacherId,
      teacherName,
      schoolId,
      academicYearId,
      semesterId,
    }) {
      this.loading = true
      this.error = null
      try {
        const progress = await progressService.createProgress({
          studentId,
          studentName,
          classId,
          className,
          type,
          title,
          description,
          date,
          teacherId,
          teacherName,
          schoolId,
          academicYearId,
          semesterId,
        })
        this.progressRecords.push(progress)
        return progress
      } catch (err) {
        this.error = err.message
        throw err
      } finally {
        this.loading = false
      }
    },
  },
})
