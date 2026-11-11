// apps/siakad-tu/src/services/registration/serviceFactory.js

import { registrationApi } from '@/adapters/api/registrationApi.js'
import { RegistrationCommandService } from './RegistrationCommandService.js'
import { RegistrationQueryService } from './RegistrationQueryService.js'

export function createRegistrationCommandService(dependencies = {}) {
  const adapter = dependencies.adapter || registrationApi
  return new RegistrationCommandService({ adapter })
}

export function createRegistrationQueryService(dependencies = {}) {
  const adapter = dependencies.adapter || registrationApi
  return new RegistrationQueryService({ adapter })
}

export const useRegistrationCommandService = () => createRegistrationCommandService()
export const useRegistrationQueryService = () => createRegistrationQueryService()
