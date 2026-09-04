// apps/siakad-tu/src/adapters/api/registrationApi.js

import { registrationMockAdapter } from '../mock/registrationMockAdapter.js'
import { registrationRealAdapter } from './registrationRealAdapter.js'

const useApi = import.meta.env?.QCLI_USE_API === 'true'
const forceMock = import.meta.env?.QCLI_MOCK_MODE === 'true'

export const registrationApi =
  forceMock || !useApi ? registrationMockAdapter : registrationRealAdapter

export function __forceMockRegistrationAdapter() {
  return registrationMockAdapter
}

export function __forceRealRegistrationAdapter() {
  return registrationRealAdapter
}
