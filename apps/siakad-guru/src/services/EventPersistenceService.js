/**
 * Handler yang dipanggil oleh EventDispatcher setiap ada event.
 * Bertugas menyimpan event ke backend/mock untuk persistence.
 * @param {import('@/events/DomainEvent').DomainEvent} event
 */
export async function eventPersistenceHandler(event) {
  // Di development: log ke console
  if (process.env.DEV) {
    console.log(`[EventPersistence] ${event.type}`, event.payload)
  }

  // Panggil adapter untuk simpan ke backend (bisa via HTTP atau mock)
  // Untuk saat ini, kita biarkan mock/adapter menangani.
  // Nanti bisa diganti dengan panggilan API:
  // await apiClient.post('/events', event.toJSON())

  // Simulasi async
  return Promise.resolve()
}
