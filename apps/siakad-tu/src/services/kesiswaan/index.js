import { KesiswaanDashboardQueryService } from './KesiswaanDashboardQueryService.js'
import { kesiswaanDashboardApi } from '@/adapters/kesiswaanDashboardApi.js'

/**
 * Factory untuk membuat instance KesiswaanDashboardQueryService
 * dengan dependency injection default (menggunakan kesiswaanDashboardApi).
 *
 * @param {Object} dependencies - Dependencies override (untuk testing atau custom)
 * @param {Object} dependencies.adapter - Adapter override (default: kesiswaanDashboardApi)
 * @returns {KesiswaanDashboardQueryService}
 */
export function createKesiswaanDashboardQueryService(dependencies = {}) {
  const adapter = dependencies.adapter || kesiswaanDashboardApi

  return new KesiswaanDashboardQueryService({ adapter })
}

// Ekspor default instance untuk kemudahan penggunaan
export const kesiswaanDashboardQueryService = createKesiswaanDashboardQueryService()

// Ekspor class juga agar bisa di-instantiate manual jika diperlukan
export { KesiswaanDashboardQueryService }
