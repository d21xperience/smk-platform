import { DispositionStatus } from '../models/DispositionStatus.js'

export class DispositionEngine {
  /**
   * Memvalidasi apakah transisi status disposisi diperbolehkan.
   * @param {string} currentStatus - Status saat ini
   * @param {string} newStatus - Status baru yang diinginkan
   * @returns {boolean} True jika transisi valid
   */
  static canTransition(currentStatus, newStatus) {
    const validTransitions = {
      [DispositionStatus.PENDING]: [DispositionStatus.DIBACA, DispositionStatus.SELESAI],
      [DispositionStatus.DIBACA]: [DispositionStatus.SELESAI],
      [DispositionStatus.SELESAI]: [],
    }
    return validTransitions[currentStatus]?.includes(newStatus) || false
  }
}
