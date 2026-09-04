// apps/siakad-tu/src/services/finance/FinanceCommandService.js

export class FinanceCommandService {
  constructor({ adapter }) {
    if (!adapter) throw new Error('FinanceCommandService: adapter wajib.')
    this.adapter = adapter
  }

  _formatError(source, error) {
    return {
      success: false,
      data: null,
      error: {
        code: error.code || `${source}_ERROR`,
        message: error.message || 'Terjadi kesalahan.',
        source,
        details: error.details || null,
      },
    }
  }

  async createInvoice(commandData, context) {
    if (!commandData.studentId) {
      return this._formatError('validation', {
        code: 'STUDENT_REQUIRED',
        message: 'Siswa wajib dipilih.',
      })
    }

    if (!commandData.amount || commandData.amount <= 0) {
      return this._formatError('validation', {
        code: 'INVALID_AMOUNT',
        message: 'Jumlah tagihan harus lebih dari 0.',
      })
    }

    if (!commandData.type) {
      return this._formatError('validation', {
        code: 'TYPE_REQUIRED',
        message: 'Jenis tagihan wajib dipilih.',
      })
    }

    const result = await this.adapter.createInvoice(commandData, context)
    if (!result.success) return this._formatError('adapter', result.error)
    return { success: true, data: result.data, error: null }
  }

  async processPayment(commandData, context) {
    if (!commandData.invoiceId) {
      return this._formatError('validation', {
        code: 'INVOICE_REQUIRED',
        message: 'Invoice wajib dipilih.',
      })
    }

    if (!commandData.amount || commandData.amount <= 0) {
      return this._formatError('validation', {
        code: 'INVALID_AMOUNT',
        message: 'Jumlah pembayaran harus lebih dari 0.',
      })
    }

    // Ambil invoice untuk validasi outstanding
    const invoiceResult = await this.adapter.getInvoiceById(commandData.invoiceId, context)
    if (!invoiceResult.success) {
      return this._formatError('adapter', invoiceResult.error)
    }

    const invoice = invoiceResult.data
    if (commandData.amount > invoice.outstandingAmount) {
      return this._formatError('business_rule', {
        code: 'PAYMENT_EXCEEDS_OUTSTANDING',
        message: `Jumlah pembayaran melebihi sisa tagihan (Rp ${invoice.outstandingAmount.toLocaleString()}).`,
      })
    }

    const result = await this.adapter.recordPayment(commandData.invoiceId, commandData, context)
    if (!result.success) return this._formatError('adapter', result.error)
    return { success: true, data: result.data, error: null }
  }

  async deleteInvoice(invoiceId, context) {
    const result = await this.adapter.deleteInvoice(invoiceId, context)
    if (!result.success) return this._formatError('adapter', result.error)
    return { success: true, data: result.data, error: null }
  }
}
