// apps/siakad-tu/src/services/finance/serviceFactory.js

import { financeApi } from '@/adapters/api/financeApi.js'

export function createFinanceCommandService(dependencies = {}) {
  const adapter = dependencies.adapter || financeApi
  const { FinanceCommandService } = require('./FinanceCommandService.js')
  return new FinanceCommandService({ adapter })
}

export function createFinanceQueryService(dependencies = {}) {
  const adapter = dependencies.adapter || financeApi
  const { FinanceQueryService } = require('./FinanceQueryService.js')
  return new FinanceQueryService({ adapter })
}

export const useFinanceCommandService = () => createFinanceCommandService()
export const useFinanceQueryService = () => createFinanceQueryService()
