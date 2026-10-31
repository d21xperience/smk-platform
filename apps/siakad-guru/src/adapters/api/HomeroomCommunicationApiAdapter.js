// FILE: src/adapters/api/HomeroomCommunicationApiAdapter.js
// STATUS: NEW
// STATUS IMPLEMENTASI: COMPLETE

import { api } from '@/boot/axios'
import {
  mapHomeroomCommunications,
  mapHomeroomCommunication,
} from '@/contracts/homeroomCommunicationContract'

/**
 * Homeroom Communication API Adapter
 * Menggunakan Axios untuk operasi log komunikasi wali kelas.
 * Hanya Adapter yang boleh menggunakan Axios.
 * Operasi: fetch, save, update (tidak ada delete).
 */
export class HomeroomCommunicationApiAdapter {
  /**
   * Mengambil daftar log komunikasi wali kelas berdasarkan context.
   * @param {Object} context - Operational context
   * @param {string} context.userId - User ID wali kelas
   * @param {string} context.schoolId - School ID
   * @param {string} context.academicYearId - Academic Year ID
   * @param {string} context.semesterId - Semester ID
   * @returns {Promise<Array>} Array log komunikasi
   */
  async fetchCommunicationsByHomeroom(context) {
    const { userId, schoolId, academicYearId, semesterId } = context

    const response = await api.get('/v1/homeroom/communications', {
      headers: {
        'x-school-id': schoolId,
        'x-academic-year-id': academicYearId,
        'x-semester-id': semesterId,
        'x-user-id': userId,
      },
    })

    const data = response.data || {}
    return mapHomeroomCommunications(data.communications || [])
  }

  /**
   * Menyimpan log komunikasi baru ke backend.
   * @param {Object} communication - Objek log komunikasi baru.
   * @param {Object} context - Operational context.
   * @returns {Promise<Object>} Log komunikasi yang tersimpan dengan communicationId.
   */
  async saveCommunication(communication, context) {
    const { userId, schoolId, academicYearId, semesterId } = context

    const response = await api.post('/v1/homeroom/communications', communication, {
      headers: {
        'x-school-id': schoolId,
        'x-academic-year-id': academicYearId,
        'x-semester-id': semesterId,
        'x-user-id': userId,
      },
    })

    return mapHomeroomCommunication(response.data)
  }

  /**
   * Mengupdate log komunikasi di backend.
   * Digunakan untuk update status follow-up.
   * @param {Object} communication - Objek log komunikasi yang sudah diupdate.
   * @param {Object} context - Operational context.
   * @returns {Promise<Object>} Log komunikasi yang terupdate.
   */
  async updateCommunication(communication, context) {
    const { userId, schoolId, academicYearId, semesterId } = context

    const response = await api.put(
      `/v1/homeroom/communications/${communication.communicationId}`,
      communication,
      {
        headers: {
          'x-school-id': schoolId,
          'x-academic-year-id': academicYearId,
          'x-semester-id': semesterId,
          'x-user-id': userId,
        },
      },
    )

    return mapHomeroomCommunication(response.data)
  }
}
