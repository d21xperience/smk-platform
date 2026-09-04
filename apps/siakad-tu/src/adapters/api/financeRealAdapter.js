// apps/siakad-tu/src/adapters/api/financeRealAdapter.js

import { api } from '@/boot/axios.js'

/**
 * Helper: Map Backend Invoice DTO ke Frontend Contract
 */
const mapBackendToInvoice = (backendDto) => {
  if (!backendDto) return null

  return {
    invoiceId: backendDto.invoice_id || backendDto.id,
    studentId: backendDto.student_id,
    type: backendDto.type, // 'TUITION', 'BUILDING_FEE', 'EXAM_FEE', 'ACTIVITY_FEE', dll
    amount: backendDto.amount,
    paidAmount: backendDto.paid_amount || 0,
    outstandingAmount: backendDto.outstanding_amount || backendDto.amount,
    status: backendDto.status, // 'UNPAID', 'PARTIAL', 'PAID', 'OVERDUE'
    dueDate: backendDto.due_date || backendDto.dueDate,
    description: backendDto.description,
    periodId: backendDto.period_id,
    createdAt: backendDto.created_at || backendDto.createdAt,
    updatedAt: backendDto.updated_at || backendDto.updatedAt,
    payments: backendDto.payments ? backendDto.payments.map(mapBackendToPayment) : [],
  }
}

/**
 * Helper: Map Backend Payment DTO ke Frontend Contract
 */
const mapBackendToPayment = (backendDto) => {
  if (!backendDto) return null

  return {
    paymentId: backendDto.payment_id || backendDto.id,
    invoiceId: backendDto.invoice_id,
    amount: backendDto.amount,
    paymentDate: backendDto.payment_date || backendDto.paymentDate,
    paymentMethod: backendDto.payment_method || backendDto.paymentMethod, // 'CASH', 'TRANSFER', 'QRIS'
    referenceNumber: backendDto.reference_number || backendDto.referenceNumber,
    note: backendDto.note,
    createdAt: backendDto.created_at || backendDto.createdAt,
  }
}

/**
 * Helper: Map Backend Finance Summary DTO ke Frontend Contract
 */
const mapBackendToFinanceSummary = (backendDto) => {
  if (!backendDto) return null

  return {
    studentId: backendDto.student_id,
    totalInvoiced: backendDto.total_invoiced || 0,
    totalPaid: backendDto.total_paid || 0,
    totalOutstanding: backendDto.total_outstanding || 0,
    overdueAmount: backendDto.overdue_amount || 0,
    invoiceCount: backendDto.invoice_count || 0,
    paidInvoiceCount: backendDto.paid_invoice_count || 0,
    unpaidInvoiceCount: backendDto.unpaid_invoice_count || 0,
  }
}

/**
 * Helper: Konversi ApplicationError ke format Service yang konsisten
 */
const handleApiError = (error) => ({
  success: false,
  data: null,
  error: {
    code: error.code || 'API_ERROR',
    message: error.message || 'Terjadi kesalahan pada server.',
    details: error.details || null,
  },
})

const handleSuccess = (data) => ({
  success: true,
  data,
  error: null,
})

export const financeRealAdapter = {
  /**
   * Buat invoice baru untuk siswa
   */
  async createInvoice(command, context) {
    try {
      const payload = {
        studentId: command.studentId,
        type: command.type,
        amount: command.amount,
        dueDate: command.dueDate,
        description: command.description,
        periodId: context.periodId,
      }
      const response = await api.post('/invoices', payload)
      return handleSuccess(mapBackendToInvoice(response.data))
    } catch (error) {
      return handleApiError(error)
    }
  },

  /**
   * Catat pembayaran (partial atau full)
   */
  async recordPayment(invoiceId, command) {
    try {
      const payload = {
        amount: command.amount,
        paymentDate: command.paymentDate || new Date().toISOString(),
        paymentMethod: command.paymentMethod,
        referenceNumber: command.referenceNumber,
        note: command.note,
      }
      const response = await api.post(`/invoices/${invoiceId}/payments`, payload)
      return handleSuccess(mapBackendToInvoice(response.data)) // Backend mengembalikan invoice yang sudah diupdate
    } catch (error) {
      return handleApiError(error)
    }
  },

  /**
   * Ambil detail invoice by ID
   */
  async getInvoiceById(invoiceId) {
    try {
      const response = await api.get(`/invoices/${invoiceId}`)
      return handleSuccess(mapBackendToInvoice(response.data))
    } catch (error) {
      return handleApiError(error)
    }
  },

  /**
   * Ambil daftar invoice dengan filter
   */
  async getInvoices(context, filters = {}) {
    try {
      const params = { schoolId: context.schoolId, periodId: context.periodId, ...filters }
      const response = await api.get('/invoices', { params })

      const mappedData = response.data.items
        ? { ...response.data, items: response.data.items.map(mapBackendToInvoice) }
        : response.data

      return handleSuccess(mappedData)
    } catch (error) {
      return handleApiError(error)
    }
  },

  /**
   * Ambil ringkasan keuangan siswa
   */
  async getFinanceSummary(studentId, context) {
    try {
      const response = await api.get(`/students/${studentId}/finance-summary`, {
        params: { periodId: context.periodId },
      })
      return handleSuccess(mapBackendToFinanceSummary(response.data))
    } catch (error) {
      return handleApiError(error)
    }
  },

  /**
   * Ambil daftar invoice yang masih outstanding (belum lunas)
   */
  async getOutstandingInvoices(context, filters = {}) {
    try {
      const params = {
        schoolId: context.schoolId,
        periodId: context.periodId,
        status: 'OUTSTANDING', // Backend akan filter UNPAID + PARTIAL
        ...filters,
      }
      const response = await api.get('/invoices/outstanding', { params })

      const mappedData = response.data.items
        ? { ...response.data, items: response.data.items.map(mapBackendToInvoice) }
        : response.data

      return handleSuccess(mappedData)
    } catch (error) {
      return handleApiError(error)
    }
  },

  /**
   * Hapus invoice (hanya jika belum ada pembayaran)
   */
  async deleteInvoice(invoiceId) {
    try {
      await api.delete(`/invoices/${invoiceId}`)
      return handleSuccess({ deleted: true })
    } catch (error) {
      return handleApiError(error)
    }
  },
}
