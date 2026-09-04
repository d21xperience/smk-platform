// apps/siakad-tu/src/adapters/mock/financeMockAdapter.js

import { mockStorage } from './mockStorage.js'
import { idGenerator } from '../utils/idGenerator.js'

const _simulateLatency = (ms = 50) => new Promise((resolve) => setTimeout(resolve, ms))
const _error = (code, message, details = null) => ({
  success: false,
  data: null,
  error: { code, message, details },
})
const _success = (data) => ({ success: true, data, error: null })

/**
 * Helper: Hitung status invoice berdasarkan paidAmount
 */
const _calculateInvoiceStatus = (invoice) => {
  if (invoice.paidAmount >= invoice.amount) return 'PAID'
  if (invoice.paidAmount > 0) return 'PARTIAL'

  // Cek apakah overdue
  const now = new Date()
  const dueDate = new Date(invoice.dueDate)
  if (now > dueDate && invoice.paidAmount < invoice.amount) return 'OVERDUE'

  return 'UNPAID'
}

/**
 * Helper: Update outstanding amount
 */
// const _updateOutstanding = (invoice) => {
//   invoice.outstandingAmount = invoice.amount - invoice.paidAmount
//   invoice.status = _calculateInvoiceStatus(invoice)
//   return invoice
// }

export const financeMockAdapter = {
  async createInvoice(command, context) {
    await _simulateLatency(100)

    const invoiceData = {
      invoiceId: idGenerator.invoiceId(),
      studentId: command.studentId,
      type: command.type,
      amount: command.amount,
      paidAmount: 0,
      outstandingAmount: command.amount,
      status: 'UNPAID',
      dueDate: command.dueDate,
      description: command.description,
      periodId: context.periodId,
      schoolId: context.schoolId,
      payments: [],
      createdAt: new Date().toISOString(),
      updatedAt: new Date().toISOString(),
    }

    const invoices = mockStorage.read(context.schoolId, 'invoices')
    invoices.push(invoiceData)
    mockStorage.write(context.schoolId, 'invoices', invoices)

    return _success(invoiceData)
  },

  async recordPayment(invoiceId, command, context) {
    await _simulateLatency(150)

    const invoices = mockStorage.read(context.schoolId, 'invoices')
    const index = invoices.findIndex((inv) => inv.invoiceId === invoiceId)

    if (index === -1) return _error('NOT_FOUND', 'Invoice tidak ditemukan.')

    const invoice = invoices[index]

    // Validasi: Pembayaran tidak boleh melebihi outstanding
    if (command.amount > invoice.outstandingAmount) {
      return _error(
        'INVALID_AMOUNT',
        `Jumlah pembayaran (${command.amount}) melebihi sisa tagihan (${invoice.outstandingAmount}).`,
      )
    }

    // Validasi: Invoice tidak boleh sudah PAID
    if (invoice.status === 'PAID') {
      return _error('ALREADY_PAID', 'Invoice sudah lunas.')
    }

    // Catat pembayaran
    const payment = {
      paymentId: idGenerator.paymentId(),
      invoiceId: invoiceId,
      amount: command.amount,
      paymentDate: command.paymentDate || new Date().toISOString(),
      paymentMethod: command.paymentMethod || 'CASH',
      referenceNumber: command.referenceNumber || null,
      note: command.note || null,
      createdAt: new Date().toISOString(),
    }

    invoice.payments.push(payment)
    invoice.paidAmount += command.amount
    invoice.outstandingAmount = invoice.amount - invoice.paidAmount
    invoice.status = _calculateInvoiceStatus(invoice)
    invoice.updatedAt = new Date().toISOString()

    invoices[index] = invoice
    mockStorage.write(context.schoolId, 'invoices', invoices)

    return _success(invoice)
  },

  async getInvoiceById(invoiceId, context) {
    await _simulateLatency(30)
    const invoices = mockStorage.read(context.schoolId, 'invoices')
    const data = invoices.find((inv) => inv.invoiceId === invoiceId)

    if (!data) return _error('NOT_FOUND', 'Invoice tidak ditemukan.')
    return _success(data)
  },

  async getInvoices(context, filters = {}) {
    await _simulateLatency(100)
    let invoices = mockStorage.read(context.schoolId, 'invoices')

    // Filter by studentId
    if (filters.studentId) {
      invoices = invoices.filter((inv) => inv.studentId === filters.studentId)
    }

    // Filter by type
    if (filters.type) {
      invoices = invoices.filter((inv) => inv.type === filters.type)
    }

    // Filter by status
    if (filters.status) {
      invoices = invoices.filter((inv) => inv.status === filters.status)
    }

    // Filter by periodId
    if (filters.periodId) {
      invoices = invoices.filter((inv) => inv.periodId === filters.periodId)
    }

    // Search by description
    if (filters.search) {
      const term = filters.search.toLowerCase()
      invoices = invoices.filter((inv) => inv.description?.toLowerCase().includes(term))
    }

    const page = filters.page || 1
    const limit = filters.limit || 20
    const total = invoices.length
    const items = invoices.slice((page - 1) * limit, page * limit)

    return _success({ items, total, page, limit, totalPages: Math.ceil(total / limit) })
  },

  async getFinanceSummary(studentId, context) {
    await _simulateLatency(50)
    const invoices = mockStorage.read(context.schoolId, 'invoices')

    const studentInvoices = invoices.filter(
      (inv) => inv.studentId === studentId && inv.periodId === context.periodId,
    )

    const summary = {
      studentId,
      totalInvoiced: studentInvoices.reduce((sum, inv) => sum + inv.amount, 0),
      totalPaid: studentInvoices.reduce((sum, inv) => sum + inv.paidAmount, 0),
      totalOutstanding: studentInvoices.reduce((sum, inv) => sum + inv.outstandingAmount, 0),
      overdueAmount: studentInvoices
        .filter((inv) => inv.status === 'OVERDUE')
        .reduce((sum, inv) => sum + inv.outstandingAmount, 0),
      invoiceCount: studentInvoices.length,
      paidInvoiceCount: studentInvoices.filter((inv) => inv.status === 'PAID').length,
      unpaidInvoiceCount: studentInvoices.filter((inv) => inv.status !== 'PAID').length,
    }

    return _success(summary)
  },

  async getOutstandingInvoices(context, filters = {}) {
    await _simulateLatency(100)
    let invoices = mockStorage.read(context.schoolId, 'invoices')

    // Filter hanya yang outstanding (UNPAID, PARTIAL, OVERDUE)
    invoices = invoices.filter(
      (inv) => inv.status === 'UNPAID' || inv.status === 'PARTIAL' || inv.status === 'OVERDUE',
    )

    // Filter by studentId
    if (filters.studentId) {
      invoices = invoices.filter((inv) => inv.studentId === filters.studentId)
    }

    // Sort by dueDate (yang paling dekat duluan)
    invoices.sort((a, b) => new Date(a.dueDate) - new Date(b.dueDate))

    const page = filters.page || 1
    const limit = filters.limit || 20
    const total = invoices.length
    const items = invoices.slice((page - 1) * limit, page * limit)

    return _success({ items, total, page, limit, totalPages: Math.ceil(total / limit) })
  },

  async deleteInvoice(invoiceId, context) {
    await _simulateLatency(100)
    const invoices = mockStorage.read(context.schoolId, 'invoices')
    const index = invoices.findIndex((inv) => inv.invoiceId === invoiceId)

    if (index === -1) return _error('NOT_FOUND', 'Invoice tidak ditemukan.')

    // Validasi: Hanya bisa hapus jika belum ada pembayaran
    if (invoices[index].paidAmount > 0) {
      return _error('HAS_PAYMENTS', 'Tidak dapat menghapus invoice yang sudah ada pembayarannya.')
    }

    invoices.splice(index, 1)
    mockStorage.write(context.schoolId, 'invoices', invoices)

    return _success({ deleted: true })
  },
}
