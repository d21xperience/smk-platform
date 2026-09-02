// apps/siakad-tu/src/adapters/api/studentApi.contract.js

/**
 * ============================================================
 * STUDENT API CONTRACT — FRONTEND DRIVEN
 * ============================================================
 *
 * Kontrak ini adalah SUMBER KEBENARAN untuk API kesiswaan.
 * Backend HARUS mengikuti kontrak ini, tidak boleh mengubah payload.
 *
 * Semua method menerima:
 * - command: Object berisi data operasi
 * - context: OperationalContext { schoolId, periodId, ... }
 *
 * Semua method mengembalikan Promise dengan format:
 * Success: { success: true, data: {...} }
 * Failure: { success: false, error: { code, message, details? } }
 *
 * Error Codes:
 * - NISN_DUPLICATE: NISN sudah terdaftar di sekolah ini
 * - NIS_DUPLICATE: NIS sudah digunakan di sekolah ini
 * - STUDENT_NOT_FOUND: Siswa tidak ditemukan
 * - SCHOOL_MISMATCH: Siswa terdaftar di sekolah lain
 * - INVALID_CONTEXT: OperationalContext tidak valid
 * - VALIDATION_ERROR: Data tidak valid (lihat details)
 * - NETWORK_ERROR: Gagal koneksi ke backend
 * ============================================================
 */

export const StudentApiContract = {
  /**
   * Mendaftarkan siswa baru
   * POST /api/students/register
   *
   * @param {Object} command
   * @param {string} command.studentId - ID siswa (dari Engine)
   * @param {string} command.nisn - NISN (10 digit)
   * @param {string} command.nis - NIS sekolah
   * @param {Object} command.fullName - { firstName, middleName, lastName }
   * @param {string} command.birthDate - ISO date string
   * @param {string} command.gender - 'MALE' | 'FEMALE'
   * @param {Object} command.address - { street, rtRw, village, district, city, postalCode }
   * @param {Object} command.contactInfo - { phone, email }
   * @param {Object} command.guardianInfo - { name, relation, phone, occupation }
   * @param {OperationalContext} context
   * @returns {Promise<{ success: boolean, data?: Object, error?: Object }>}
   */
  registerStudent: async (command, context) => {
    /* implement */
    console.log(command, context)
  },

  /**
   * Mendaftarkan siswa ke kelas
   * POST /api/students/:studentId/enroll
   *
   * @param {Object} command
   * @param {string} command.studentId
   * @param {string} command.enrollmentId
   * @param {string} command.classId
   * @param {string} [command.enrollmentDate] - ISO date string
   * @param {OperationalContext} context
   * @returns {Promise<{ success: boolean, data?: Object, error?: Object }>}
   */
  enrollStudent: async (command, context) => {
    console.log(command, context)
    /* implement */
  },

  /**
   * Update profil siswa
   * PATCH /api/students/:studentId/profile
   *
   * @param {Object} command
   * @param {string} command.studentId
   * @param {Object} [command.address]
   * @param {Object} [command.contactInfo]
   * @param {Object} [command.guardianInfo]
   * @param {OperationalContext} context
   * @returns {Promise<{ success: boolean, data?: Object, error?: Object }>}
   */
  updateStudentProfile: async (command, context) => {
    /* implement */
    console.log(command, context)
  },

  /**
   * Luluskan siswa
   * POST /api/students/:studentId/graduate
   *
   * @param {Object} command
   * @param {string} command.studentId
   * @param {string} [command.graduationDate]
   * @param {OperationalContext} context
   * @returns {Promise<{ success: boolean, data?: Object, error?: Object }>}
   */
  graduateStudent: async (command, context) => {
    /* implement */
    console.log(command, context)
  },

  /**
   * Pindahkan siswa
   * POST /api/students/:studentId/transfer
   *
   * @param {Object} command
   * @param {string} command.studentId
   * @param {string} command.targetSchool
   * @param {string} [command.reason]
   * @param {OperationalContext} context
   * @returns {Promise<{ success: boolean, data?: Object, error?: Object }>}
   */
  transferStudent: async (command, context) => {
    /* implement */
    console.log(command, context)
  },

  /**
   * Ambil data siswa berdasarkan ID
   * GET /api/students/:studentId
   *
   * @param {string} studentId
   * @param {OperationalContext} context
   * @returns {Promise<{ success: boolean, data?: Object, error?: Object }>}
   */
  getStudentById: async (studentId, context) => {
    /* implement */
    console.log(studentId, context)
  },

  /**
   * Ambil daftar siswa dengan filter
   * GET /api/students?schoolId=X&periodId=Y&classId=Z&status=ACTIVE
   *
   * @param {OperationalContext} context
   * @param {Object} filters
   * @param {string} [filters.classId]
   * @param {string} [filters.status]
   * @param {string} [filters.search] - search by name/NISN/NIS
   * @param {number} [filters.page=1]
   * @param {number} [filters.limit=20]
   * @returns {Promise<{ success: boolean, data?: { items: [], total: number, page: number } }>}
   */
  getStudents: async (context, filters = {}) => {
    /* implement */
    console.log(filters, context)
  },

  /**
   * Cek ketersediaan NISN
   * GET /api/students/check-nisn?nisn=X
   *
   * @param {string} nisn
   * @param {OperationalContext} context
   * @returns {Promise<{ success: boolean, data?: { available: boolean, existingStudentId?: string } }>}
   */
  checkNisnAvailability: async (nisn, context) => {
    /* implement */
    console.log(nisn, context)
  },

  /**
   * Cek ketersediaan NIS
   * GET /api/students/check-nis?nis=X
   *
   * @param {string} nis
   * @param {OperationalContext} context
   * @returns {Promise<{ success: boolean, data?: { available: boolean, existingStudentId?: string } }>}
   */
  checkNisAvailability: async (nis, context) => {
    /* implement */
    console.log(nis, context)
  },

  /**
   * Reset data mock (khusus mock adapter)
   * Hanya tersedia di mock adapter, tidak di real adapter.
   */
  __resetMockData: async () => {
    /* implement */
  },
}
