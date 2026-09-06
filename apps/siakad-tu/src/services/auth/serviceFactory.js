import { authApi } from '@/adapters/api/authApi.js'
import { AuthService } from './AuthService.js'

/**
 * Service Factory — Dependency Injection container untuk Auth Domain.
 *
 * Membuat instance AuthService dengan dependencies yang sudah di-inject.
 *
 * Keuntungan DI:
 * 1. Mudah di-test (bisa inject mock dependencies)
 * 2. Mudah di-swap (mock → real adapter)
 * 3. Singleton per application lifecycle
 * 4. Dependencies eksplisit, tidak tersembunyi
 *
 * Penggunaan:
 *   import { authService } from './serviceFactory.js'
 *   await authService.login(credentials)
 *
 * Untuk testing:
 *   import { createAuthService } from './serviceFactory.js'
 *   const testService = createAuthService({ adapter: mockAdapter })
 */

/**
 * Buat instance AuthService dengan dependencies custom
 * (untuk testing atau environment khusus)
 *
 * @param {Object} dependencies
 * @param {Object} [dependencies.adapter] - AuthAdapter (default: authApi)
 * @returns {AuthService}
 */
export function createAuthService(dependencies = {}) {
  const adapter = dependencies.adapter || authApi

  return new AuthService({ adapter })
}

/**
 * Singleton instance untuk penggunaan normal di aplikasi.
 * Dibuat sekali, dipakai di seluruh aplikasi.
 */
export const authService = createAuthService()

/**
 * Helper: Get singleton instance
 * Untuk konsistensi dengan pattern useXxx() di composables
 */
export function useAuthService() {
  return authService
}
