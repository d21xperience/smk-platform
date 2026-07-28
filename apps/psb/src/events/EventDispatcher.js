/**
 * Event Dispatcher — Pure JavaScript
 * Tidak bergantung pada framework.
 * Bertanggung jawab:
 * - Menerima DomainEvent dari Engine/Store
 * - Meneruskan event ke Service (untuk persistence/log)
 * - Mendistribusikan event ke semua listener yang terdaftar
 */
export class EventDispatcher {
  constructor() {
    // Map<eventType, Set<callback>>
    this._listeners = new Map()
  }

  /**
   * Kirim event ke semua listener + persistence handler.
   * @param {import('./DomainEvent').DomainEvent} event
   */
  dispatch(event) {
    // 1. Kirim ke service handler (persistence/log/audit)
    if (this._serviceHandler) {
      this._serviceHandler(event).catch((err) => {
        console.error('[EventDispatcher] Service handler error:', err)
      })
    }

    // 2. Distribusikan ke listener
    const listeners = this._listeners.get(event.type)
    if (listeners) {
      for (const callback of listeners) {
        try {
          callback(event)
        } catch (err) {
          console.error(`[EventDispatcher] Listener error for ${event.type}:`, err)
        }
      }
    }

    // 3. Kirim juga ke wildcard listener "*" (mendengarkan semua event)
    const wildcardListeners = this._listeners.get('*')
    if (wildcardListeners) {
      for (const callback of wildcardListeners) {
        try {
          callback(event)
        } catch (err) {
          console.error('[EventDispatcher] Wildcard listener error:', err)
        }
      }
    }
  }

  /**
   * Daftarkan listener untuk event tertentu.
   * @param {string} eventType - dari DomainEventType
   * @param {Function} callback - function(event: DomainEvent)
   */
  on(eventType, callback) {
    if (!this._listeners.has(eventType)) {
      this._listeners.set(eventType, new Set())
    }
    this._listeners.get(eventType).add(callback)
  }

  /**
   * Hapus listener.
   * @param {string} eventType
   * @param {Function} callback
   */
  off(eventType, callback) {
    const listeners = this._listeners.get(eventType)
    if (listeners) {
      listeners.delete(callback)
    }
  }

  /**
   * Hapus semua listener (untuk reset/cleanup).
   */
  removeAllListeners() {
    this._listeners.clear()
  }

  /**
   * Pasang service handler — dipanggil saat event akan disimpan ke backend/audit.
   * Handler menerima DomainEvent dan mengembalikan Promise.
   * @param {Function} handler - async (event) => { ... }
   */
  setServiceHandler(handler) {
    this._serviceHandler = handler
  }
}

// export class EventDispatcher {
//   constructor() {
//     this._listeners = new Map()
//   }

//   dispatch(event) {
//     if (this._serviceHandler) {
//       this._serviceHandler(event).catch((err) =>
//         console.error('[EventDispatcher] Service handler error:', err),
//       )
//     }
//     const listeners = this._listeners.get(event.type)
//     if (listeners) {
//       for (const callback of listeners) {
//         try {
//           callback(event)
//         } catch (err) {
//           console.error(`[EventDispatcher] Listener error for ${event.type}:`, err)
//         }
//       }
//     }
//     const wildcardListeners = this._listeners.get('*')
//     if (wildcardListeners) {
//       for (const callback of wildcardListeners) {
//         try {
//           callback(event)
//         } catch (err) {
//           console.error('[EventDispatcher] Wildcard listener error:', err)
//         }
//       }
//     }
//   }

//   on(eventType, callback) {
//     if (!this._listeners.has(eventType)) this._listeners.set(eventType, new Set())
//     this._listeners.get(eventType).add(callback)
//   }

//   off(eventType, callback) {
//     const listeners = this._listeners.get(eventType)
//     if (listeners) listeners.delete(callback)
//   }

//   removeAllListeners() {
//     this._listeners.clear()
//   }

//   setServiceHandler(handler) {
//     this._serviceHandler = handler
//   }
// }
