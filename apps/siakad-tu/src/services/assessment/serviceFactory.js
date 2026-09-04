// apps/siakad-tu/src/services/assessment/serviceFactory.js

import { assessmentApi } from '@/adapters/api/assessmentApi.js'

export function createAssessmentCommandService(dependencies = {}) {
  const adapter = dependencies.adapter || assessmentApi
  const { AssessmentCommandService } = require('./AssessmentCommandService.js')
  return new AssessmentCommandService({ adapter })
}

export function createAssessmentQueryService(dependencies = {}) {
  const adapter = dependencies.adapter || assessmentApi
  const { AssessmentQueryService } = require('./AssessmentQueryService.js')
  return new AssessmentQueryService({ adapter })
}

export const useAssessmentCommandService = () => createAssessmentCommandService()
export const useAssessmentQueryService = () => createAssessmentQueryService()
