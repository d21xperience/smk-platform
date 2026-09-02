// apps/siakad-tu/src/services/student/StudentQueryService.js

/**
 * StudentQueryService — Orkestrator untuk read operations.
 *
 * TANGGUNG JAWAB:
 * 1. Memanggil Adapter langsung untuk fetch data
 * 2. Return data dalam format yang konsisten
 *
 * CATATAN PENTING:
 * - Query TIDAK melewati Engine karena tidak butuh business logic
 * - Ini adalah Lightweight CQRS: read dan write dipisah
 * - Query bisa di-cache di masa depan tanpa mengubah Engine
 *
 * Di masa depan, QueryService bisa membaca dari Projection (Read Model)
 * alih-alih langsung ke Adapter, untuk performa yang lebih baik.
 */
export class StudentQueryService {
  /**
   * @param {Object} dependencies
   * @param {Object} dependencies.adapter - studentApi (mock atau real)
   */
  constructor({ adapter }) {
    if (!adapter) throw new Error('StudentQueryService: adapter wajib.')
    this.adapter = adapter
  }

  /**
   * Helper: format error response
   * @private
   */
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

  /**
   * Ambil data siswa berdasarkan ID
   * @param {string} studentId
   * @param {OperationalContext} context
   * @returns {Promise<{ success, data, error }>}
   */
  async getStudentById(studentId, context) {
    const result = await this.adapter.getStudentById(studentId, context)
    if (!result.success) {
      return this._formatError(result.error)
    }
    return { success: true, data: result.data, error: null }
  }

  /**
   * Ambil daftar siswa dengan filter
   * @param {OperationalContext} context
   * @param {Object} filters - { classId, status, search, page, limit }
   * @returns {Promise<{ success, data, error }>}
   */
  async getStudents(context, filters = {}) {
    const result = await this.adapter.getStudents(context, filters)
    if (!result.success) {
      return this._formatError(result.error)
    }
    return { success: true, data: result.data, error: null }
  }

  /**
   * Cek ketersediaan NISN
   * @param {string} nisn
   * @param {OperationalContext} context
   * @returns {Promise<{ success, data: { available, existingStudentId }, error }>}
   */
  async checkNisnAvailability(nisn, context) {
    const result = await this.adapter.checkNisnAvailability(nisn, context)
    if (!result.success) {
      return this._formatError(result.error)
    }
    return { success: true, data: result.data, error: null }
  }

  /**
   * Cek ketersediaan NIS
   * @param {string} nis
   * @param {OperationalContext} context
   * @returns {Promise<{ success, data: { available, existingStudentId }, error }>}
   */
  async checkNisAvailability(nis, context) {
    const result = await this.adapter.checkNisAvailability(nis, context)
    if (!result.success) {
      return this._formatError(result.error)
    }
    return { success: true, data: result.data, error: null }
  }
}
