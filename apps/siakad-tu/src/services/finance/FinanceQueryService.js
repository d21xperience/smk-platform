// apps/siakad-tu/src/services/finance/FinanceQueryService.js

export class FinanceQueryService {
  constructor({ adapter }) {
    if (!adapter) throw new Error('FinanceQueryService: adapter wajib.')
    this.adapter = adapter
  }

  _formatError(error) {
    return {
      success: false,
      data: null,
      error: {
        code: error.code || 'QUERY_ERROR',
        message: error.message || 'Gagal mengambil data.',
        source: 'adapter',
        details: error.details || null,
      },
    }
  }

  async getInvoiceById(invoiceId, context) {
    const result = await this.adapter.getInvoiceById(invoiceId, context)
    if (!result.success) return this._formatError(result.error)
    return { success: true, data: result.data, error: null }
  }

  async getInvoices(context, filters = {}) {
    const result = await this.adapter.getInvoices(context, filters)
    if (!result.success) return this._formatError(result.error)
    return { success: true, data: result.data, error: null }
  }

  async getFinanceSummary(studentId, context) {
    const result = await this.adapter.getFinanceSummary(studentId, context)
    if (!result.success) return this._formatError(result.error)
    return { success: true, data: result.data, error: null }
  }

  async getOutstandingInvoices(context, filters = {}) {
    const result = await this.adapter.getOutstandingInvoices(context, filters)
    if (!result.success) return this._formatError(result.error)
    return { success: true, data: result.data, error: null }
  }
}
