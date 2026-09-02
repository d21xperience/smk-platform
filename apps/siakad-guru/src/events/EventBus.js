import { EventDispatcher } from './EventDispatcher'

// Singleton instance
export const eventBus = new EventDispatcher()

// Export untuk kemudahan penggunaan
export const emit = (event, payload) => eventBus.emit(event, payload)
export const emitAsync = (event, payload) => eventBus.emitAsync(event, payload)
export const on = (event, listener, options) => eventBus.on(event, listener, options)
export const once = (event, listener, options) => eventBus.once(event, listener, options)
export const off = (event, listener) => eventBus.off(event, listener)
