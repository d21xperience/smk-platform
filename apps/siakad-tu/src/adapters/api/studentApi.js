// apps/siakad-tu/src/adapters/api/studentApi.js
import { studentMockAdapter } from '../mock/studentMockAdapter.js'
import { studentRealAdapter } from './studentRealAdapter.js'

// PAKSA MENGGUNAKAN REAL API
const forceMock = false
const useApi = true 

export const studentApi = forceMock || !useApi ? studentMockAdapter : studentRealAdapter

console.log('useApi:', useApi)
console.log('useForceMock:', forceMock)
