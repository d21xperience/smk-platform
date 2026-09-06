export class KesiswaanDashboardQueryService {
  constructor({ adapter } = {}) {
    if (!adapter) {
      throw new Error('KesiswaanDashboardQueryService: adapter wajib.')
    }
    this.adapter = adapter
  }

  /**
   * Helper untuk memformat error
   * @param {string} source - Sumber error (misal: 'adapter')
   * @param {Object} error - Error object dari adapter
   * @returns {Object} Response error dengan format { success, data, error }
   */
  _formatError(source, error) {
    return {
      success: false,
      data: null,
      error: {
        code: error.code || `${source.toUpperCase()}_ERROR`,
        message: error.message || 'Terjadi kesalahan.',
        source,
        details: error.details || null,
      },
    }
  }

  /**
   * Memuat data dashboard kesiswaan
   * @param {OperationalContext} context - Konteks operasional (schoolId, periodId, dll)
   * @returns {Promise<{ success: boolean, data?: object, error?: object }>}
   */
  async loadDashboard(context) {
    if (!context || !context.schoolId) {
      return this._formatError('service', {
        code: 'INVALID_CONTEXT',
        message: 'Context tidak valid atau schoolId tidak ditemukan.',
        details: { context },
      })
    }

    try {
      const result = await this.adapter.loadDashboard(context)

      if (!result.success) {
        return this._formatError('adapter', result.error)
      }

      return {
        success: true,
        data: result.data,
        error: null,
      }
    } catch (error) {
      return this._formatError('service', {
        code: 'UNEXPECTED_ERROR',
        message: error.message || 'Terjadi kesalahan tak terduga.',
        details: error,
      })
    }
  }
}
