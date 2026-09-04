// apps/siakad-tu/src/adapters/utils/idGenerator.js

/**
 * ID Generator — Menghasilkan ID unik untuk setiap entitas.
 *
 * Format: {prefix}_{timestamp}_{random}
 * - prefix: Identifier domain (e.g., 'student', 'invoice')
 * - timestamp: Unix timestamp dalam milidetik
 * - random: 6 digit random number untuk menghindari collision
 *
 * Contoh: student_1704067200000_123456
 */

/**
 * Generate random number dengan crypto API
 * @returns {string} 6 digit random number
 */
function generateRandom() {
  const array = new Uint32Array(1)
  crypto.getRandomValues(array)
  return (array[0] % 1000000).toString().padStart(6, '0')
}

/**
 * Generate ID dengan format standar
 * @param {string} prefix - Prefix domain
 * @returns {string} ID unik
 */
function generateId(prefix) {
  const timestamp = Date.now()
  const random = generateRandom()
  return `${prefix}_${timestamp}_${random}`
}

export const idGenerator = {
  // === EXISTING IDs ===

  /**
   * Generate Student ID
   * @returns {string} e.g., 'student_1704067200000_123456'
   */
  studentId: () => generateId('student'),

  /**
   * Generate Enrollment ID
   * @returns {string} e.g., 'enrollment_1704067200000_123456'
   */
  enrollmentId: () => generateId('enrollment'),

  /**
   * Generate Class ID
   * @returns {string} e.g., 'class_1704067200000_123456'
   */
  classId: () => generateId('class'),

  /**
   * Generate Subject ID
   * @returns {string} e.g., 'subject_1704067200000_123456'
   */
  subjectId: () => generateId('subject'),

  /**
   * Generate Period ID
   * @returns {string} e.g., 'period_1704067200000_123456'
   */
  periodId: () => generateId('period'),

  /**
   * Generate School ID
   * @returns {string} e.g., 'school_1704067200000_123456'
   */
  schoolId: () => generateId('school'),

  /**
   * Generate User ID
   * @returns {string} e.g., 'user_1704067200000_123456'
   */
  userId: () => generateId('user'),

  // === NEW IDs untuk Domain Baru ===

  /**
   * Generate Registration ID (PPDB)
   * @returns {string} e.g., 'reg_1704067200000_123456'
   */
  registrationId: () => generateId('reg'),

  /**
   * Generate Mutation ID
   * @returns {string} e.g., 'mutation_1704067200000_123456'
   */
  mutationId: () => generateId('mutation'),

  /**
   * Generate Invoice ID (Finance)
   * @returns {string} e.g., 'inv_1704067200000_123456'
   */
  invoiceId: () => generateId('inv'),

  /**
   * Generate Payment ID (Finance)
   * @returns {string} e.g., 'payment_1704067200000_123456'
   */
  paymentId: () => generateId('payment'),

  /**
   * Generate Attendance Session ID
   * @returns {string} e.g., 'session_1704067200000_123456'
   */
  sessionId: () => generateId('session'),

  /**
   * Generate Attendance Record ID
   * @returns {string} e.g., 'record_1704067200000_123456'
   */
  recordId: () => generateId('record'),

  /**
   * Generate Assessment ID (Nilai)
   * @returns {string} e.g., 'assess_1704067200000_123456'
   */
  assessmentId: () => generateId('assess'),

  // === UTILITY METHODS ===

  /**
   * Generate UUID v4 (untuk kasus khusus yang butuh UUID standar)
   * @returns {string} UUID v4 format
   */
  uuid: () => {
    return 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, function (c) {
      const r = (Math.random() * 16) | 0
      const v = c === 'x' ? r : (r & 0x3) | 0x8
      return v.toString(16)
    })
  },

  /**
   * Generate Transaction ID (untuk payment gateway, dll)
   * @returns {string} e.g., 'trx_1704067200000_123456'
   */
  transactionId: () => generateId('trx'),

  /**
   * Generate Receipt ID (untuk bukti pembayaran)
   * @returns {string} e.g., 'receipt_1704067200000_123456'
   */
  receiptId: () => generateId('receipt'),

  /**
   * Generate Letter ID (untuk surat masuk/keluar)
   * @returns {string} e.g., 'letter_1704067200000_123456'
   */
  letterId: () => generateId('letter'),

  /**
   * Generate Disposition ID (untuk disposisi surat)
   * @returns {string} e.g., 'disp_1704067200000_123456'
   */
  dispositionId: () => generateId('disp'),

  /**
   * Generate Classification ID (untuk klasifikasi surat)
   * @returns {string} e.g., 'class_1704067200000_123456'
   */
  classificationId: () => generateId('class'),
}
