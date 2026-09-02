import { defineStore } from 'pinia'
import { FinanceService } from '@/services/FinanceService'
import { StudentFinancialSummary } from '@/domain/finance/models/value-objects/StudentFinancialSummary'

export const useFinanceStore = defineStore('finance', {
  state: () => ({
    classSummaries: [], // Array of StudentFinancialSummary untuk satu kelas
    currentStudentSummary: null, // Detail per siswa
    isLoading: false,
    error: null,
  }),
  getters: {
    overdueStudents: (state) => state.classSummaries.filter((s) => s.isOverdue()),
    paidStudents: (state) => state.classSummaries.filter((s) => s.isPaid()),
    totalBalance: (state) => state.classSummaries.reduce((sum, s) => sum + s.balance, 0),
  },
  actions: {
    async fetchClassFinancialSummary(classId) {
      this.isLoading = true
      this.error = null
      try {
        const data = await FinanceService.getClassFinancialSummary(classId)
        this.classSummaries = data.map((s) => new StudentFinancialSummary(s))
        return this.classSummaries
      } catch (err) {
        this.error = err.message
        throw err
      } finally {
        this.isLoading = false
      }
    },

    async fetchStudentFinancialDetail(studentId) {
      this.isLoading = true
      try {
        const data = await FinanceService.getStudentFinancialDetail(studentId)
        this.currentStudentSummary = new StudentFinancialSummary(data)
        return this.currentStudentSummary
      } catch (err) {
        this.error = err.message
        throw err
      } finally {
        this.isLoading = false
      }
    },

    clear() {
      this.classSummaries = []
      this.currentStudentSummary = null
      this.error = null
    },
  },
})
