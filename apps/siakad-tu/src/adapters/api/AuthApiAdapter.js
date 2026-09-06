import { authMockAdapter } from '../mock/auth.mock.js'
/**
Auth API Factory.
Memilih implementasi adapter berdasarkan environment:
development/test → authMockAdapter (in-memory)
production → authRealAdapter (Axios → TU-Core)
Kontrak (method signature) HARUS identik antara mock dan real.
Service tidak perlu tahu implementasi mana yang dipakai.
Catatan: Real adapter akan dibuat saat backend auth sudah siap.
*/
// TODO: Di Langkah integrasi backend, ganti dengan:
// import { authRealAdapter } from './authRealAdapter.js'
const isDevelopment = import.meta.env?.QCLI_USE_API ?? true
/**
Adapter aktif untuk auth API.
Saat ini: Mock Adapter.
Nanti (production): Real Adapter.
*/
export const authApi = isDevelopment ? authMockAdapter : authMockAdapter // TODO: ganti ke authRealAdapter saat siap
/**
Helper: Force mock adapter (untuk testing)
*/
export function __forceMockAuthAdapter() {
  return authMockAdapter
}
