import { defineStore } from 'pinia'
import { homeroomBillingService } from '../boot/services.js'

export const useHomeroomBillingStore = defineStore('homeroomBilling', {
  state: () => ({
    assignment: null,
    bills: [],
    summary: null,
    loading: false,
    error: null,
  }),

  getters: {
    hasAssignment: (state) => !!state.assignment,
    hasData: (state) => state.bills.length > 0,
    hasSummary: (state) => !!state.summary,
    outstandingBills: (state) => state.bills.filter((b) => b.isOutstanding()),
    settledBills: (state) => state.bills.filter((b) => b.isSettled()),
  },

  actions: {
    async loadHomeroomAssignment({ teacherId }) {
      this.loading = true
      this.error = null
      try {
        this.assignment = await homeroomBillingService.getHomeroomAssignment({ teacherId })
        return this.assignment
      } catch (err) {
        this.error = err.message
        return null
      } finally {
        this.loading = false
      }
    },

    async loadBillingData({ classId, academicYearId, semesterId, teacherId }) {
      this.loading = true
      this.error = null
      try {
        const readModel = await homeroomBillingService.getBillingData({
          classId,
          academicYearId,
          semesterId,
          teacherId,
        })
        this.bills = readModel.bills
        this.summary = readModel.summary
      } catch (err) {
        this.error = err.message
        this.bills = []
        this.summary = null
      } finally {
        this.loading = false
      }
    },

    clearData() {
      this.assignment = null
      this.bills = []
      this.summary = null
      this.error = null
    },
  },
})
