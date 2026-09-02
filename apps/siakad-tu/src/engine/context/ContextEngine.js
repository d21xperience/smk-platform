// apps/siakad-tu/src/engine/context/ContextEngine.js

import { OperationalContext } from '@/domain/context/OperationalContext'
import {
  OperationalContextChanged,
  OperationalContextCleared,
} from '@/events/context/ContextEvents'

/**
 * ContextEngine — Pure JavaScript.
 *
 * TIDAK BOLEH import: vue, pinia, axios, router, localStorage.
 * Menerima data mentah, mengembalikan context + events.
 *
 * Tanggung jawab:
 * 1. Validasi bahwa context lengkap dan konsisten
 * 2. Derivasi periodId dari academicYear + semester
 * 3. Menghasilkan Domain Events saat context berubah
 */
export class ContextEngine {
  /**
   * Membangun OperationalContext dari data mentah (API response / localStorage)
   * @param {Object} rawData - { schoolId, schoolName, academicYear, semester }
   * @returns {{ context: OperationalContext, events: DomainEvent[] }}
   */
  establishContext(rawData) {
    // 1. Validasi input
    if (!rawData.schoolId) {
      throw new Error('ContextEngine: Tidak dapat membangun context tanpa schoolId.')
    }
    if (!rawData.academicYear || !rawData.semester) {
      throw new Error('ContextEngine: academicYear dan semester wajib dipilih.')
    }

    // 2. Bangun Aggregate (validasi internal terjadi di constructor)
    const context = new OperationalContext(rawData)

    if (!context.isValid()) {
      throw new Error('ContextEngine: Context yang dihasilkan tidak valid.')
    }

    // 3. Hasilkan event
    const events = [new OperationalContextChanged(context)]

    return { context, events }
  }

  /**
   * Validasi apakah context saat ini masih valid untuk melakukan transaksi
   * @param {OperationalContext} context
   * @returns {{ valid: boolean, reason: string }}
   */
  validateForTransaction(context) {
    if (!context) {
      return { valid: false, reason: 'Tidak ada Operational Context yang aktif.' }
    }
    if (!context.isValid()) {
      return { valid: false, reason: 'Operational Context tidak lengkap.' }
    }
    return { valid: true, reason: '' }
  }

  /**
   * Clear context (saat logout)
   * @returns {{ events: DomainEvent[] }}
   */
  clearContext() {
    return { events: [new OperationalContextCleared()] }
  }
}
