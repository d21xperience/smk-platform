// apps/siakad-tu/src/adapters/utils/idGenerator.js

/**
 * ID Generator — Pure JavaScript, tanpa dependensi eksternal.
 *
 * Menghasilkan UUID v4-like string.
 * Digunakan untuk generate studentId, enrollmentId, dll.
 *
 * Catatan: Untuk production, sebaiknya gunakan crypto.randomUUID()
 * jika tersedia di browser. Implementasi ini untuk fallback.
 */
export const idGenerator = {
  /**
   * Generate UUID v4-like string
   * @returns {string} UUID-like, contoh: "f47ac10b-58cc-4372-a567-0e02b2c3d479"
   */
  uuid() {
    // Gunakan crypto.randomUUID() jika tersedia (browser modern)
    if (typeof crypto !== 'undefined' && typeof crypto.randomUUID === 'function') {
      return crypto.randomUUID()
    }

    // Fallback: manual generation
    return 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, (c) => {
      const r = (Math.random() * 16) | 0
      const v = c === 'x' ? r : (r & 0x3) | 0x8
      return v.toString(16)
    })
  },

  /**
   * Generate short ID (untuk enrollmentId, dll)
   * Format: "enr_" + timestamp + random
   * @returns {string}
   */
  shortId(prefix = 'id') {
    const timestamp = Date.now().toString(36)
    const random = Math.random().toString(36).substring(2, 8)
    return `${prefix}_${timestamp}_${random}`
  },

  /**
   * Generate enrollment ID
   * @returns {string}
   */
  enrollmentId() {
    return this.shortId('enr')
  },

  /**
   * Generate student ID
   * @returns {string}
   */
  studentId() {
    return this.uuid()
  },
}
