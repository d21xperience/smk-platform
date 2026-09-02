// apps/siakad-tu/src/engine/Result.js

/**
 * Result helper — konsisten return type untuk seluruh Engine.
 *
 * Semua Engine method HARUS mengembalikan Result, tidak boleh throw error
 * (kecuali untuk programming error seperti parameter salah type).
 *
 * Format:
 * Success: { success: true, data: {...}, events: [...], errors: [] }
 * Failure: { success: false, data: null, events: [], errors: [{code, message, field?}] }
 */
export class Result {
  /**
   * Buat Result sukses
   * @param {Object} data - Data hasil (aggregate, projection, dll)
   * @param {DomainEvent[]} events - Domain events yang dihasilkan
   */
  static success(data, events = []) {
    return {
      success: true,
      data,
      events: Array.isArray(events) ? events : [events],
      errors: [],
    }
  }

  /**
   * Buat Result gagal
   * @param {Array<{code: string, message: string, field?: string}>} errors
   */
  static failure(errors) {
    const errorArray = Array.isArray(errors) ? errors : [errors]
    return {
      success: false,
      data: null,
      events: [],
      errors: errorArray.map((e) => ({
        code: e.code || 'VALIDATION_ERROR',
        message: e.message || String(e),
        field: e.field || null,
      })),
    }
  }

  /**
   * Gabungkan beberapa Result. Jika semua sukses, return sukses dengan data array.
   * Jika ada yang gagal, return failure dengan gabungan errors.
   */
  static combine(results) {
    const failures = results.filter((r) => !r.success)
    if (failures.length > 0) {
      const allErrors = failures.flatMap((r) => r.errors)
      return Result.failure(allErrors)
    }

    return Result.success(
      results.map((r) => r.data),
      results.flatMap((r) => r.events),
    )
  }
}
