import { kesiswaanDashboardMockAdapter } from './mock/kesiswaanDashboardMockAdapter.js'

/**
 * Factory untuk adapter dashboard kesiswaan.
 * Saat ini masih menggunakan mock adapter.
 * Nanti bisa diganti dengan real adapter saat VITE_USE_MOCK=false.
 */
const isDevelopment = import.meta.env?.DEV ?? true

// Untuk sementara, selalu gunakan mock adapter
// TODO: ganti ke real adapter setelah backend siap
export const kesiswaanDashboardApi = isDevelopment
  ? kesiswaanDashboardMockAdapter
  : kesiswaanDashboardMockAdapter

/**
 * Fungsi untuk memaksa menggunakan mock adapter (misal untuk testing)
 */
export function __forceMockAdapter() {
  return kesiswaanDashboardMockAdapter
}
