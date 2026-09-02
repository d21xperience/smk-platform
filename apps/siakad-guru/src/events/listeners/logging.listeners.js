// import { EVENTS } from '../EventCatalog'

export function registerLoggingListeners(eventBus) {
  // Log semua event (untuk debugging di development)
  if (process.env.NODE_ENV === 'development') {
    eventBus.on(
      '*',
      (eventData) => {
        console.log(`[EVENT BUS] ${eventData.event}`, eventData.payload)
      },
      { priority: -100 }, // Priority rendah, agar tidak mengganggu logika bisnis
    )
  }
}
