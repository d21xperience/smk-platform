// src/events/EventDispatcher.js
export class EventDispatcher {
  constructor() {
    this.listeners = new Map()
    this.wildcardListeners = []
  }

  on(event, listener, options = {}) {
    const { priority = 0, once = false } = options
    if (typeof listener !== 'function') {
      throw new Error('Listener must be a function.')
    }

    if (event === '*') {
      this.wildcardListeners.push({ fn: listener, priority, once })
      this.wildcardListeners.sort((a, b) => b.priority - a.priority)
      return () => this.off('*', listener)
    }

    if (!this.listeners.has(event)) {
      this.listeners.set(event, [])
    }
    const listeners = this.listeners.get(event)
    listeners.push({ fn: listener, priority, once })
    listeners.sort((a, b) => b.priority - a.priority)
    return () => this.off(event, listener)
  }

  once(event, listener, options = {}) {
    return this.on(event, listener, { ...options, once: true })
  }

  off(event, listener) {
    if (event === '*') {
      this.wildcardListeners = this.wildcardListeners.filter((l) => l.fn !== listener)
      return
    }
    if (this.listeners.has(event)) {
      const listeners = this.listeners.get(event)
      const filtered = listeners.filter((l) => l.fn !== listener)
      if (filtered.length === 0) {
        this.listeners.delete(event)
      } else {
        this.listeners.set(event, filtered)
      }
    }
  }

  removeAllListeners(event) {
    if (event) {
      this.listeners.delete(event)
    } else {
      this.listeners.clear()
      this.wildcardListeners = []
    }
  }

  // Method emit (synchronous)
  emit(event, payload) {
    const listeners = this.listeners.get(event) || []
    // Execute wildcard listeners first
    for (const wildcard of this.wildcardListeners) {
      if (wildcard.once) {
        this.off('*', wildcard.fn)
      }
      wildcard.fn({ event, payload })
    }
    // Execute specific listeners
    for (const listener of listeners) {
      if (listener.once) {
        this.off(event, listener.fn)
      }
      listener.fn(payload)
    }
  }

  // Method emitAsync (asynchronous)
  async emitAsync(event, payload) {
    const listeners = this.listeners.get(event) || []
    for (const wildcard of this.wildcardListeners) {
      if (wildcard.once) {
        this.off('*', wildcard.fn)
      }
      await wildcard.fn({ event, payload })
    }
    for (const listener of listeners) {
      if (listener.once) {
        this.off(event, listener.fn)
      }
      await listener.fn(payload)
    }
  }

  listenerCount(event) {
    if (event === '*') return this.wildcardListeners.length
    return this.listeners.get(event)?.length || 0
  }
}
