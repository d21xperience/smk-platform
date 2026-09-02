// FILE: src/adapters/api/HomeroomProgressApiAdapter.js
// STATUS: NEW
// STATUS IMPLEMENTASI: COMPLETE

import { api } from '@/boot/axios'

/**
 * Homeroom Progress API Adapter
 * Menggunakan Axios untuk mengambil data nilai siswa dari backend.
 * Hanya Adapter yang boleh menggunakan Axios.
 * Data mentah dikembalikan untuk diproses Domain Engine di Service layer.
 */
export class HomeroomProgressApiAdapter {
  /**
   * Mengambil data mentah nilai siswa per kelas berdasarkan wali kelas.
   * @param {Object} context - Operational context
   * @param {string} context.userId - User ID wali kelas
   * @param {string} context.schoolId - School ID
   * @param {string} context.academicYearId - Academic Year ID
   * @param {string} context.semesterId - Semester ID
   * @returns {Promise<Object>} Data mentah nilai siswa per kelas
   */
  async fetchProgressByHomeroom(context) {
    const { userId, schoolId, academicYearId, semesterId } = context

    const response = await api.get('/v1/homeroom/progress', {
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
      subjects: Array.isArray(data.subjects) ? data.subjects : [],
      students: Array.isArray(data.students) ? data.students : [],
    }
  }
}
