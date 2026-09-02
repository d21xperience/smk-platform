import { defineStore } from 'pinia'
import { StudentProgressService } from '@/services/StudentProgressService'
import { StudentProgress } from '@/domain/studentProgress/models/StudentProgress'

export const useStudentProgressStore = defineStore('studentProgress', {
  state: () => ({
    classProgress: [], // Array of StudentProgress (untuk homeroom)
    currentStudentProgress: null, // Detail per siswa
    isLoading: false,
    error: null,
  }),
  getters: {
    classSummary: (state) => {
      const total = state.classProgress.length
      const naik = state.classProgress.filter((p) => p.promotionStatus === 'NAIK').length
      const tidakNaik = state.classProgress.filter((p) => p.promotionStatus === 'TIDAK_NAIK').length
      const perhatian = state.classProgress.filter((p) => p.promotionStatus === 'PERHATIAN').length
      return { total, naik, tidakNaik, perhatian }
    },
  },
  actions: {
    async loadClassProgress(classId, semesterId) {
      this.isLoading = true
      this.error = null
      try {
        const data = await StudentProgressService.getClassProgress(classId, semesterId)
        this.classProgress = data.map((p) => new StudentProgress(p))
        return this.classProgress
      } catch (err) {
        this.error = err.message
        throw err
      } finally {
        this.isLoading = false
      }
    },

    async loadStudentProgress(studentId, semesterId) {
      this.isLoading = true
      this.error = null
      try {
        const data = await StudentProgressService.getStudentProgress(studentId, semesterId)
        this.currentStudentProgress = new StudentProgress(data)
        return this.currentStudentProgress
      } catch (err) {
        this.error = err.message
        throw err
      } finally {
        this.isLoading = false
      }
    },

    async addTeacherNote(studentId, content, type = 'GENERAL') {
      this.isLoading = true
      try {
        const data = await StudentProgressService.addTeacherNote(studentId, content, type)
        if (this.currentStudentProgress && this.currentStudentProgress.studentId === studentId) {
          this.currentStudentProgress.addTeacherNote(data)
        }
        // Update class list if exists
        const classIdx = this.classProgress.findIndex((p) => p.studentId === studentId)
        if (classIdx !== -1) {
          this.classProgress[classIdx].addTeacherNote(data)
        }
        return data
      } catch (err) {
        this.error = err.message
        throw err
      } finally {
        this.isLoading = false
      }
    },

    async addAchievement(studentId, achievement) {
      this.isLoading = true
      try {
        const data = await StudentProgressService.addAchievement(studentId, achievement)
        if (this.currentStudentProgress && this.currentStudentProgress.studentId === studentId) {
          this.currentStudentProgress.addAchievement(data)
        }
        const classIdx = this.classProgress.findIndex((p) => p.studentId === studentId)
        if (classIdx !== -1) {
          this.classProgress[classIdx].addAchievement(data)
        }
        return data
      } catch (err) {
        this.error = err.message
        throw err
      } finally {
        this.isLoading = false
      }
    },

    async addViolation(studentId, violation) {
      this.isLoading = true
      try {
        const data = await StudentProgressService.addViolation(studentId, violation)
        if (this.currentStudentProgress && this.currentStudentProgress.studentId === studentId) {
          this.currentStudentProgress.addViolation(data)
        }
        const classIdx = this.classProgress.findIndex((p) => p.studentId === studentId)
        if (classIdx !== -1) {
          this.classProgress[classIdx].addViolation(data)
        }
        return data
      } catch (err) {
        this.error = err.message
        throw err
      } finally {
        this.isLoading = false
      }
    },

    clear() {
      this.classProgress = []
      this.currentStudentProgress = null
      this.error = null
    },
  },
})
