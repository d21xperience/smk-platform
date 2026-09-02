// apps/siakad-tu/src/adapters/api/studentApi.js

import { studentMockAdapter } from '../mock/studentMockAdapter.js'

/**
 * Student API Factory.
 *
 * Memilih implementasi adapter berdasarkan environment:
 * - development/test → studentMockAdapter (localStorage)
 * - production → studentRealAdapter (Axios → edge-bff)
 *
 * Kontrak (method signature) HARUS identik antara mock dan real.
 * Service tidak perlu tahu implementasi mana yang dipakai.
 *
 * Catatan: Real adapter akan dibuat di Langkah 10 (Backend Implementation).
 */

// TODO: Di Langkah 10, ganti dengan:
// import { studentRealAdapter } from '../real/studentRealAdapter.js';

const isDevelopment = import.meta.env?.DEV ?? true

/**
 * Adapter aktif untuk student API.
 * Saat ini: Mock Adapter.
 * Nanti (production): Real Adapter.
 */
export const studentApi = isDevelopment ? studentMockAdapter : studentMockAdapter // TODO: ganti ke studentRealAdapter saat siap

/**
 * Helper: Force mock adapter (untuk testing)
 */
export function __forceMockAdapter() {
  return studentMockAdapter
}
