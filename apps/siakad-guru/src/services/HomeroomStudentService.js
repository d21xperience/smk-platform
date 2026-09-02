// FILE: src/services/HomeroomStudentService.js
// STATUS: MODIFY (koreksi dari Batch 1)
// STATUS IMPLEMENTASI: COMPLETE

/**
 * Homeroom Student Application Service
 * Mengorkestrasi use case pengambilan data siswa wali kelas.
 * Read-only: tidak ada business rule kompleks.
 */
export class HomeroomStudentService {
  constructor({ homeroomStudentAdapter }) {
    this.adapter = homeroomStudentAdapter
  }

  /**
   * Mengambil daftar siswa wali kelas menggunakan operational context.
   * @param {Object} context - Operational context (schoolId, academicYearId, userId, dll)
   * @returns {Promise<Object>} Data kelas dan daftar siswa
   */
  async getHomeroomStudents(context) {
    if (!context) {
      throw new Error('Operational context is required to fetch homeroom students.')
    }

    if (!this.adapter) {
      throw new Error('Adapter is not initialized in HomeroomStudentService.')
    }

    return this.adapter.fetchStudentsByHomeroom(context)
  }
}
