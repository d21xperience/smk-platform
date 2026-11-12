// apps/siakad-tu/src/adapters/api/mutationApi.js

import { mutationMockAdapter } from '../mock/mutationMockAdapter.js'
import { mutationRealAdapter } from './mutationRealAdapter.js'

// PAKSA MENGGUNAKAN REAL API
const forceMock = false
const useApi = true

export const mutationApi = forceMock || !useApi ? mutationMockAdapter : mutationRealAdapter

export function __forceMockMutationAdapter() {
  return mutationMockAdapter
}

export function __forceRealMutationAdapter() {
  return mutationRealAdapter
}
