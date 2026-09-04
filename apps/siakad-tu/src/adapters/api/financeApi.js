// apps/siakad-tu/src/adapters/api/financeApi.js

import { financeMockAdapter } from '../mock/financeMockAdapter.js'
import { financeRealAdapter } from './financeRealAdapter.js'

const useApi = import.meta.env?.QCLI_USE_API === 'true'
const forceMock = import.meta.env?.QCLI_MOCK_MODE === 'true'

export const financeApi = forceMock || !useApi ? financeMockAdapter : financeRealAdapter

export function __forceMockFinanceAdapter() {
  return financeMockAdapter
}

export function __forceRealFinanceAdapter() {
  return financeRealAdapter
}
