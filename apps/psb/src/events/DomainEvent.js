/**
 * Domain Event — kontrak standar untuk semua event.
 * Tidak bergantung pada framework apa pun.
 */
export class DomainEvent {
  /**
   * @param {string} type - Nama event (dari DomainEventType)
   * @param {Object} payload - Data yang dibawa event
   * @param {string} [timestamp] - ISO timestamp, default sekarang
   */
  constructor(type, payload, timestamp = new Date().toISOString()) {
    this.type = type
    this.payload = payload
    this.timestamp = timestamp
  }
}
