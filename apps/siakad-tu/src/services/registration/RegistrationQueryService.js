// apps/siakad-tu/src/services/registration/RegistrationQueryService.js

/**
 * RegistrationQueryService — Orkestrator untuk read operations pada Registration Domain.
 *
 * TANGGUNG JAWAB:
 * 1. Memanggil Adapter langsung untuk fetch data
 * 2. Return data dalam format yang konsisten
 *
 * CATATAN:
 * - Query TIDAK melewati Engine karena tidak butuh business logic
 * - Ini adalah Lightweight CQRS: read dan write dipisah
 */
export class RegistrationQueryService {
  constructor({ adapter }) {
    if (!adapter) throw new Error('RegistrationQueryService: adapter wajib.')
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

  async getRegistrationById(registrationId, context) {
    const result = await this.adapter.getRegistrationById(registrationId, context)
    if (!result.success) return this._formatError(result.error)
    return { success: true, data: result.data, error: null }
  }

  async getRegistrations(context, filters = {}) {
    const result = await this.adapter.getRegistrations(context, filters)
    if (!result.success) return this._formatError(result.error)
    return { success: true, data: result.data, error: null }
  }
}
