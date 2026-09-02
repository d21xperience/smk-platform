import { EventCatalog } from './Catalog.js'

export class EventDispatcher {
  constructor() {
    this.listeners = new Map()
  }

  subscribe(eventType, callback) {
    if (!EventCatalog.isValidType(eventType)) {
      console.warn(`[EventDispatcher] Attempted to subscribe to unknown event type: ${eventType}`)
      return
    }

    if (!this.listeners.has(eventType)) {
      this.listeners.set(eventType, new Set())
    }

    this.listeners.get(eventType).add(callback)
  }

  unsubscribe(eventType, callback) {
    if (this.listeners.has(eventType)) {
      this.listeners.get(eventType).delete(callback)
      if (this.listeners.get(eventType).size === 0) {
        this.listeners.delete(eventType)
      }
    }
  }

  dispatch(event) {
    if (!event || !event.type) {
      console.error('[EventDispatcher] Invalid event object dispatched', event)
      return
    }

    if (!EventCatalog.isValidType(event.type)) {
      console.warn(`[EventDispatcher] Dispatched unknown event type: ${event.type}`)
    }

    const callbacks = this.listeners.get(event.type)
    if (callbacks) {
      callbacks.forEach((callback) => {
        try {
          callback(event.payload)
        } catch (error) {
          console.error(`[EventDispatcher] Error in event listener for ${event.type}:`, error)
        }
      })
    }
  }

  clear() {
    this.listeners.clear()
  }
}
