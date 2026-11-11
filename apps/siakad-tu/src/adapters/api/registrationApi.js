// apps/siakad-tu/src/adapters/api/registrationApi.js

import { registrationMockAdapter } from '../mock/registrationMockAdapter.js'
import { registrationRealAdapter } from './registrationRealAdapter.js'

// PAKSA MENGGUNAKAN REAL API UNTUK DEVELOPMENT
const forceMock = false
const useApi = true

export const registrationApi = forceMock || !useApi ? registrationMockAdapter : registrationRealAdapter

export function __forceMockRegistrationAdapter() {
  return registrationMockAdapter
}

export function __forceRealRegistrationAdapter() {
  return registrationRealAdapter
}
