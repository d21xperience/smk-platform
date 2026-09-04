// apps/siakad-tu/src/services/mutation/serviceFactory.js

import { mutationApi } from '@/adapters/api/mutationApi.js'

export function createMutationCommandService(dependencies = {}) {
  const adapter = dependencies.adapter || mutationApi
  const { MutationCommandService } = require('./MutationCommandService.js')
  return new MutationCommandService({ adapter })
}

export function createMutationQueryService(dependencies = {}) {
  const adapter = dependencies.adapter || mutationApi
  const { MutationQueryService } = require('./MutationQueryService.js')
  return new MutationQueryService({ adapter })
}

export const useMutationCommandService = () => createMutationCommandService()
export const useMutationQueryService = () => createMutationQueryService()
