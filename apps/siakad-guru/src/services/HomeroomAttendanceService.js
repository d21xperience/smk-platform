// FILE: src/services/HomeroomAttendanceService.js
// STATUS: NEW
// STATUS IMPLEMENTASI: COMPLETE

import { HomeroomAttendanceEngine } from '../domain/homeroom-attendance/HomeroomAttendanceEngine.js'
import { mapHomeroomAttendance } from '../contracts/homeroomAttendanceContract.js'

/**
 * Homeroom Attendance Application Service
 * Orchestration: Adapter → Domain Engine → Contract Mapping.
 * Jalur dengan Business Rule:
 * Service memanggil Adapter untuk raw data, lalu Domain Engine untuk processing.
 * Read-only: tidak ada mutasi data.
 */
export class HomeroomAttendanceService {
  constructor({ homeroomAttendanceAdapter }) {
    this.adapter = homeroomAttendanceAdapter
    this.engine = new HomeroomAttendanceEngine()
  }

  /**
   * Mengambil rekap kehadiran siswa per kelas berdasarkan wali kelas.
   * Flow: Adapter (raw) → Engine (process) → Contract (map to UI shape).
   * @param {Object} params - Parameter operasional
   * @param {string} params.userId - User ID wali kelas
   * @param {string} params.schoolId - School ID
   * @param {string} params.academicYearId - Academic Year ID
   * @param {string} params.semesterId - Semester ID
   * @returns {Promise<Object>} Rekap kehadiran yang telah diproses
   */
  async getAttendanceSummary({ userId, schoolId, academicYearId, semesterId }) {
    if (!userId) {
      throw new Error('userId is required to fetch homeroom attendance.')
    }
    if (!academicYearId) {
      throw new Error('academicYearId is required to fetch homeroom attendance.')
    }
    if (!semesterId) {
      throw new Error('semesterId is required to fetch homeroom attendance.')
    }

    if (!this.adapter) {
      throw new Error('Adapter is not initialized in HomeroomAttendanceService.')
    }

    const rawData = await this.adapter.fetchAttendanceByHomeroom({
      userId,
      schoolId,
      academicYearId,
      semesterId,
    })

    const processedData = this.engine.processClassAttendance(rawData)

    return mapHomeroomAttendance(processedData)
  }
}
