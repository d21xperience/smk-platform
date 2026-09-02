// apps/siakad-guru/src/services/BaseService.js
import { eventBus } from '@/events/EventBus'

export class BaseService {
  constructor(contextStore, dispatcher = null) {
    this.contextStore = contextStore
    this.dispatcher = dispatcher // bisa null
  }

  getContext() {
    return this.contextStore.current
  }

  // Gunakan eventBus langsung, atau jika ada dispatcher yang disediakan, gunakan itu
  async dispatchEvent(event) {
    try {
      if (this.dispatcher && typeof this.dispatcher.dispatch === 'function') {
        await this.dispatcher.dispatch(event)
      } else {
        // Fallback: gunakan eventBus jika ada
        if (eventBus && typeof eventBus.emit === 'function') {
          eventBus.emit(event.type, event.payload)
        } else {
          console.warn('[BaseService] No dispatcher or eventBus available for event:', event.type)
        }
      }
    } catch (error) {
      console.error('[BaseService] Failed to dispatch event:', error)
    }
  }
}
