// apps/siakad-tu/src/stores/finance/financeStore.js

import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { useFinanceCommandService } from '@/services/finance/serviceFactory'
import { useFinanceQueryService } from '@/services/finance/serviceFactory'
import { useContextStore } from '@/stores/contextStore'

export const useFinanceStore = defineStore('finance', () => {
  // === STATE ===
  const invoices = ref([])
  const currentInvoice = ref(null)
  const financeSummary = ref(null)
  const outstandingInvoices = ref([])
  const isLoading = ref(false)
  const error = ref(null)
  const success = ref(null)
  const pagination = ref({ page: 1, limit: 20, total: 0, totalPages: 0 })
  const filters = ref({ status: null, type: null, studentId: null, search: '' })

  // === SERVICES ===
  const commandService = useFinanceCommandService()
  const queryService = useFinanceQueryService()
  const contextStore = useContextStore()

  // === GETTERS ===
  const unpaidInvoices = computed(() => invoices.value.filter((inv) => inv.status === 'UNPAID'))

  const partialInvoices = computed(() => invoices.value.filter((inv) => inv.status === 'PARTIAL'))

  const paidInvoices = computed(() => invoices.value.filter((inv) => inv.status === 'PAID'))

  const overdueInvoices = computed(() => invoices.value.filter((inv) => inv.status === 'OVERDUE'))

  const totalOutstanding = computed(() =>
    invoices.value.reduce((sum, inv) => sum + (inv.outstandingAmount || 0), 0),
  )

  // === ACTIONS ===

  async function fetchInvoices(customFilters = {}) {
    isLoading.value = true
    error.value = null

    try {
      const context = contextStore.currentContext
      const mergedFilters = { ...filters.value, ...customFilters }

      const result = await queryService.getInvoices(context, mergedFilters)

      if (result.success) {
        invoices.value = result.data.items || []
        pagination.value = {
          page: result.data.page || 1,
          limit: result.data.limit || 20,
          total: result.data.total || 0,
          totalPages: result.data.totalPages || 0,
        }
      } else {
        error.value = result.error
      }
    } catch (err) {
      error.value = { code: 'UNEXPECTED_ERROR', message: err.message }
    } finally {
      isLoading.value = false
    }
  }

  async function fetchInvoiceById(invoiceId) {
    isLoading.value = true
    error.value = null

    try {
      const context = contextStore.currentContext
      const result = await queryService.getInvoiceById(invoiceId, context)

      if (result.success) {
        currentInvoice.value = result.data
      } else {
        error.value = result.error
      }
    } catch (err) {
      error.value = { code: 'UNEXPECTED_ERROR', message: err.message }
    } finally {
      isLoading.value = false
    }
  }

  async function fetchFinanceSummary(studentId) {
    isLoading.value = true
    error.value = null

    try {
      const context = contextStore.currentContext
      const result = await queryService.getFinanceSummary(studentId, context)

      if (result.success) {
        financeSummary.value = result.data
      } else {
        error.value = result.error
      }
    } catch (err) {
      error.value = { code: 'UNEXPECTED_ERROR', message: err.message }
    } finally {
      isLoading.value = false
    }
  }

  async function fetchOutstandingInvoices(customFilters = {}) {
    isLoading.value = true
    error.value = null

    try {
      const context = contextStore.currentContext
      const result = await queryService.getOutstandingInvoices(context, customFilters)

      if (result.success) {
        outstandingInvoices.value = result.data.items || []
      } else {
        error.value = result.error
      }
    } catch (err) {
      error.value = { code: 'UNEXPECTED_ERROR', message: err.message }
    } finally {
      isLoading.value = false
    }
  }

  async function createInvoice(commandData) {
    isLoading.value = true
    error.value = null
    success.value = null

    try {
      const context = contextStore.currentContext
      const result = await commandService.createInvoice(commandData, context)

      if (result.success) {
        success.value = { message: 'Invoice berhasil dibuat' }
        currentInvoice.value = result.data
        return result.data
      } else {
        error.value = result.error
        return null
      }
    } catch (err) {
      error.value = { code: 'UNEXPECTED_ERROR', message: err.message }
      return null
    } finally {
      isLoading.value = false
    }
  }

  async function processPayment(commandData) {
    isLoading.value = true
    error.value = null
    success.value = null

    try {
      const context = contextStore.currentContext
      const result = await commandService.processPayment(commandData, context)

      if (result.success) {
        success.value = { message: 'Pembayaran berhasil dicatat' }
        currentInvoice.value = result.data

        // Update list
        const index = invoices.value.findIndex((inv) => inv.invoiceId === result.data.invoiceId)
        if (index !== -1) {
          invoices.value[index] = result.data
        }

        return result.data
      } else {
        error.value = result.error
        return null
      }
    } catch (err) {
      error.value = { code: 'UNEXPECTED_ERROR', message: err.message }
      return null
    } finally {
      isLoading.value = false
    }
  }

  async function deleteInvoice(invoiceId) {
    isLoading.value = true
    error.value = null
    success.value = null

    try {
      const context = contextStore.currentContext
      const result = await commandService.deleteInvoice(invoiceId, context)

      if (result.success) {
        success.value = { message: 'Invoice berhasil dihapus' }

        // Remove from list
        invoices.value = invoices.value.filter((inv) => inv.invoiceId !== invoiceId)

        return true
      } else {
        error.value = result.error
        return false
      }
    } catch (err) {
      error.value = { code: 'UNEXPECTED_ERROR', message: err.message }
      return false
    } finally {
      isLoading.value = false
    }
  }

  function reset() {
    invoices.value = []
    currentInvoice.value = null
    financeSummary.value = null
    outstandingInvoices.value = []
    error.value = null
    success.value = null
    pagination.value = { page: 1, limit: 20, total: 0, totalPages: 0 }
  }

  function updateFilters(newFilters) {
    filters.value = { ...filters.value, ...newFilters }
  }

  return {
    invoices,
    currentInvoice,
    financeSummary,
    outstandingInvoices,
    isLoading,
    error,
    success,
    pagination,
    filters,
    unpaidInvoices,
    partialInvoices,
    paidInvoices,
    overdueInvoices,
    totalOutstanding,
    fetchInvoices,
    fetchInvoiceById,
    fetchFinanceSummary,
    fetchOutstandingInvoices,
    createInvoice,
    processPayment,
    deleteInvoice,
    reset,
    updateFilters,
  }
})
