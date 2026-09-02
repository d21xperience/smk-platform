// apps/siakad-tu/src/events/eventDispatcher.js

/**
 * Event Dispatcher sederhana untuk komunikasi antar modul.
 * Engine menghasilkan event, Dispatcher menyebarkan, Listener bereaksi.
 *
 * Ini adalah fondasi. Di masa depan bisa di-upgrade ke Redis Pub/Sub
 * untuk komunikasi lintas service.
 */
class EventDispatcher {
  constructor() {
    this._listeners = new Map()
  }

  /**
   * Daftarkan listener untuk event tertentu
   * @param {string} eventName
   * @param {Function} callback
   */
  on(eventName, callback) {
    if (!this._listeners.has(eventName)) {
      this._listeners.set(eventName, [])
    }
    this._listeners.get(eventName).push(callback)
  }

  /**
   * Hapus listener
   */
  off(eventName, callback) {
    if (!this._listeners.has(eventName)) return
    const callbacks = this._listeners.get(eventName)
    this._listeners.set(
      eventName,
      callbacks.filter((cb) => cb !== callback),
    )
  }

  /**
   * Dispatch event ke semua listener
   * @param {DomainEvent} event
   */
  dispatch(event) {
    const eventName = event.eventName || event.constructor.name
    if (!this._listeners.has(eventName)) return

    this._listeners.get(eventName).forEach((callback) => {
      try {
        callback(event)
      } catch (error) {
        console.error(`[EventDispatcher] Error di listener "${eventName}":`, error)
      }
    })
  }
}

// Singleton instance
export const eventDispatcher = new EventDispatcher()
