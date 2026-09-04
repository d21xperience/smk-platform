// apps/siakad-tu/src/adapters/api/mutationApi.js

import { mutationMockAdapter } from '../mock/mutationMockAdapter.js'
import { mutationRealAdapter } from './mutationRealAdapter.js'

const useApi = import.meta.env?.QCLI_USE_API === 'true'
const forceMock = import.meta.env?.QCLI_MOCK_MODE === 'true'

export const mutationApi = forceMock || !useApi ? mutationMockAdapter : mutationRealAdapter

export function __forceMockMutationAdapter() {
  return mutationMockAdapter
}

export function __forceRealMutationAdapter() {
  return mutationRealAdapter
}
