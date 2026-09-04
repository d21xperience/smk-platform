// apps/siakad-tu/src/services/registration/serviceFactory.js

import { registrationApi } from '@/adapters/api/registrationApi.js'

/**
 * Service Factory — Dependency Injection container untuk Registration Domain.
 *
 * Membuat instance RegistrationCommandService dan RegistrationQueryService
 * dengan dependencies yang sudah di-inject.
 *
 * Keuntungan DI:
 * 1. Mudah di-test (bisa inject mock dependencies)
 * 2. Mudah di-swap (mock → real adapter)
 * 3. Singleton per application lifecycle
 * 4. Dependencies eksplisit, tidak tersembunyi
 */

/**
 * Buat instance RegistrationCommandService dengan dependencies custom
 */
export function createRegistrationCommandService(dependencies = {}) {
  const adapter = dependencies.adapter || registrationApi

  // Dynamic import untuk menghindari circular dependency
  const { RegistrationCommandService } = require('./RegistrationCommandService.js')

  return new RegistrationCommandService({ adapter })
}

/**
 * Buat instance RegistrationQueryService dengan dependencies custom
 */
export function createRegistrationQueryService(dependencies = {}) {
  const adapter = dependencies.adapter || registrationApi

  const { RegistrationQueryService } = require('./RegistrationQueryService.js')

  return new RegistrationQueryService({ adapter })
}

/**
 * Singleton instances untuk penggunaan normal di aplikasi.
 */
export const useRegistrationCommandService = () => createRegistrationCommandService()
export const useRegistrationQueryService = () => createRegistrationQueryService()
