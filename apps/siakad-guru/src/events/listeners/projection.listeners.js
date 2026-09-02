import { EVENTS } from '../EventCatalog'
import { projectionManager } from '@/projection/ProjectionManager'

export function registerProjectionListeners(eventBus) {
  // Daftarkan listener untuk semua event yang akan mengupdate projection
  // Mapping event -> projection name
  const eventProjectionMap = {
    [EVENTS.ATTENDANCE_SUBMITTED]: 'attendance',
    [EVENTS.ATTENDANCE_SUBMITTED]: 'statistics', // juga update statistics
    [EVENTS.SESSION_ENDED]: 'teaching', // nanti
    [EVENTS.JOURNAL_FINALIZED]: 'teaching',
  }

  // Karena satu event bisa update multiple projection, kita buat listener per event
  // yang akan update semua projection yang terdaftar.

  // Ambil semua event yang memiliki mapping
  const uniqueEvents = new Set(Object.keys(eventProjectionMap))
  for (const event of uniqueEvents) {
    eventBus.on(
      event,
      async (payload) => {
        // Cari projection yang perlu diupdate
        const projectionNames = []
        for (const [evt, proj] of Object.entries(eventProjectionMap)) {
          if (evt === event) {
            projectionNames.push(proj)
          }
        }
        // Untuk setiap projection, apply event
        for (const name of projectionNames) {
          try {
            const projection = projectionManager.getProjection(name)
            // Karena event payload mungkin berbeda, kita perlu bungkus dalam format event
            const wrappedEvent = { type: event, payload }
            projection.applyEvent(wrappedEvent)
          } catch (error) {
            console.error(`[ProjectionListener] Error updating projection ${name}:`, error)
          }
        }
      },
      { priority: 5 },
    )
  }

  // Tambahkan listener khusus untuk wildcard event jika ada
  // Untuk keperluan reset atau logging
  console.log('[ProjectionListener] Registered projection listeners.')
}
