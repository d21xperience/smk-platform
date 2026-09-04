// apps/siakad-tu/src/adapters/api/studentApi.js

import { studentMockAdapter } from '../mock/studentMockAdapter.js'
import { studentRealAdapter } from './studentRealAdapter.js'

/**
 * Student API Factory.
 *
 * Logika Switching:
 * 1. Jika QCLI_MOCK_MODE=true → Paksa pakai Mock (untuk testing UI)
 * 2. Jika QCLI_USE_API=true → Pakai Real API (connect ke port 8081)
 * 3. Default (DEV tanpa env khusus) → Pakai Mock (aman untuk development)
 */

const useApi = import.meta.env?.QCLI_USE_API === true
// console.log('useApi', import.meta.env?.QCLI_USE_API)
console.log('useApi', useApi)
const forceMock = import.meta.env?.QCLI_MOCK_MODE === true
console.log('useforceMock', forceMock)

export const studentApi = forceMock || !useApi ? studentMockAdapter : studentRealAdapter

/**
 * Helper: Force mock adapter (khusus untuk unit testing)
 */
export function __forceMockAdapter() {
  return studentMockAdapter
}

/**
 * Helper: Force real adapter (khusus untuk integration testing)
 */
export function __forceRealAdapter() {
  return studentRealAdapter
}
