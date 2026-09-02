// apps/siakad-tu/src/events/DomainEvent.js

/**
 * Base class untuk semua Domain Event di SIakad TU.
 * Engine menghasilkan event ini, tidak boleh ada dependensi Vue/Pinia.
 */
export class DomainEvent {
  constructor(eventName, aggregateId, payload = {}) {
    this.eventName = eventName;
    this.aggregateId = aggregateId;
    this.payload = payload;
    this.occurredOn = new Date().toISOString();
  }
}