// FILE: src/services/HomeroomProgressService.js
// STATUS: NEW
// STATUS IMPLEMENTASI: COMPLETE

import { HomeroomProgressEngine } from '../domain/homeroom-progress/HomeroomProgressEngine.js'
import { mapHomeroomProgress } from '../contracts/homeroomProgressContract.js'

/**
 * Homeroom Progress Application Service
 * Orchestration: Adapter → Domain Engine → Contract Mapping.
 * Jalur dengan Business Rule:
 * Service memanggil Adapter untuk raw data, lalu Domain Engine untuk processing.
 * Read-only: tidak ada mutasi data.
 */
export class HomeroomProgressService {
  constructor({ homeroomProgressAdapter }) {
    this.adapter = homeroomProgressAdapter
    this.engine = new HomeroomProgressEngine()
  }

  /**
   * Mengambil rekap nilai siswa per kelas berdasarkan wali kelas.
   * Flow: Adapter (raw) → Engine (process) → Contract (map to UI shape).
   * @param {Object} params - Parameter operasional
   * @param {string} params.userId - User ID wali kelas
   * @param {string} params.schoolId - School ID
   * @param {string} params.academicYearId - Academic Year ID
   * @param {string} params.semesterId - Semester ID
   * @returns {Promise<Object>} Rekap nilai yang telah diproses
   */
  async getProgressSummary({ userId, schoolId, academicYearId, semesterId }) {
    if (!userId) {
      throw new Error('userId is required to fetch homeroom progress.')
    }
    if (!academicYearId) {
      throw new Error('academicYearId is required to fetch homeroom progress.')
    }

    if (!this.adapter) {
      throw new Error('Adapter is not initialized in HomeroomProgressService.')
    }

    const rawData = await this.adapter.fetchProgressByHomeroom({
      userId,
      schoolId,
      academicYearId,
      semesterId,
    })

    const processedData = this.engine.processClassProgress(rawData)

    return mapHomeroomProgress(processedData)
  }
}
