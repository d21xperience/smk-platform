// FILE: src/adapters/api/HomeroomAttendanceApiAdapter.js
// STATUS: NEW
// STATUS IMPLEMENTASI: COMPLETE

import { api } from '@/boot/axios'

/**
 * Homeroom Attendance API Adapter
 * Menggunakan Axios untuk mengambil data kehadiran siswa dari backend.
 * Hanya Adapter yang boleh menggunakan Axios.
 * Data mentah dikembalikan untuk diproses Domain Engine di Service layer.
 */
export class HomeroomAttendanceApiAdapter {
  /**
   * Mengambil data mentah kehadiran siswa per kelas berdasarkan wali kelas.
   * @param {Object} context - Operational context
   * @param {string} context.userId - User ID wali kelas
   * @param {string} context.schoolId - School ID
   * @param {string} context.academicYearId - Academic Year ID
   * @param {string} context.semesterId - Semester ID
   * @returns {Promise<Object>} Data mentah kehadiran siswa per kelas
   */
  async fetchAttendanceByHomeroom(context) {
    const { userId, schoolId, academicYearId, semesterId } = context

    const response = await api.get('/v1/homeroom/attendance', {
      headers: {
        'X-School-Id': schoolId,
        'X-Academic-Year-Id': academicYearId,
        'X-Semester-Id': semesterId,
        'X-User-Id': userId,
      },
    })

    const data = response.data || {}

    return {
      classId: data.classId || null,
      className: data.className || null,
      months: Array.isArray(data.months) ? data.months : [],
      students: Array.isArray(data.students) ? data.students : [],
    }
  }
}
