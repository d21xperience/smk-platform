// apps/siakad-tu/src/adapters/api/assessmentApi.js

import { assessmentMockAdapter } from '../mock/assessmentMockAdapter.js'
import { assessmentRealAdapter } from './assessmentRealAdapter.js'

const useApi = import.meta.env?.QCLI_USE_API === 'true'
const forceMock = import.meta.env?.QCLI_MOCK_MODE === 'true'

export const assessmentApi = forceMock || !useApi ? assessmentMockAdapter : assessmentRealAdapter

export function __forceMockAssessmentAdapter() {
  return assessmentMockAdapter
}

export function __forceRealAssessmentAdapter() {
  return assessmentRealAdapter
}
