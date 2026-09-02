/*
 * @sdp/context — Contract: contextContract
 *
 * Mendefinisikan interface/data contract untuk Operational Context.
 * Contract ini menjadi acuan untuk:
 * - Mock Adapter
 * - API Adapter
 * - Backend payload
 * - Storage serialization
 *
 * Frontend Driven API Contract:
 * Frontend mendefinisikan data yang dibutuhkan.
 * Backend mengikuti contract tersebut.
 */
/*
 * Shape untuk School data.
 * @typedef {Object} SchoolData
 * @property {string} id
 * @property {string} name
 * @property {string} npsn
 * @property {string} address
 * @property {string|null} logoUrl
 * @property {boolean} isActive
 */
/*
 * Shape untuk AcademicYear data.
 * @typedef {Object} AcademicYearData
 * @property {string} id
 * @property {string} name
 * @property {string} startDate - ISO date string
 * @property {string} endDate - ISO date string
 * @property {boolean} isActive
 */
/*
 * Shape untuk Semester data.
 * @typedef {Object} SemesterData
 * @property {string} id
 * @property {string} academicYearId
 * @property {string} name
 * @property {string} code
 * @property {string} startDate - ISO date string
 * @property {string} endDate - ISO date string
 * @property {boolean} isActive
 */
/*
 * Shape untuk UserData dalam context.
 * @typedef {Object} ContextUserData
 * @property {string} id
 * @property {string} role
 * @property {Array<string>} permissions
 */
/*
 * Shape untuk serialized OperationalContext.
 * Digunakan untuk persistence dan transport.
 * @typedef {Object} SerializedContext
 * @property {string} schoolId
 * @property {string} academicYearId
 * @property {string} semesterId
 * @property {string} userId
 * @property {string} role
 * @property {Array<string>} permissions
 */
export const ContextContract = {
  /*
   * Mengambil daftar sekolah yang dapat diakses oleh user.
   * @returns {Promise<Array<SchoolData>>}
   */
  async fetchSchools() {
    throw new Error('Method fetchSchools() must be implemented by adapter')
  },
  /*
   * Mengambil daftar tahun ajaran untuk sebuah sekolah.
   * @param {string} schoolId
   * @returns {Promise<Array<AcademicYearData>>}
   */
  async fetchAcademicYears(schoolId) {
    throw new Error('Method fetchAcademicYears() must be implemented by adapter')
  },
  /*
   * Mengambil daftar semester untuk sebuah tahun ajaran.
   * @param {string} academicYearId
   * @returns {Promise<Array<SemesterData>>}
   */
  async fetchSemesters(academicYearId) {
    throw new Error('Method fetchSemesters() must be implemented by adapter')
  },
  /*
   * Menyimpan context ke storage.
   * @param {SerializedContext} contextData
   * @returns {Promise<void>}
   */
  async persistContext(contextData) {
    throw new Error('Method persistContext() must be implemented by adapter')
  },
  /*
   * Memuat context dari storage.
   * @returns {Promise<SerializedContext|null>}
   */
  async loadPersistedContext() {
    throw new Error('Method loadPersistedContext() must be implemented by adapter')
  },
  /*
   * Menghapus context dari storage.
   * @returns {Promise<void>}
   */
  async clearPersistedContext() {
    throw new Error('Method clearPersistedContext() must be implemented by adapter')
  }
}
/*
 * Storage key constant untuk context persistence.
 * Digunakan oleh storage adapter.
 */
export const CONTEXT_STORAGE_KEY = 'sdp_operational_context'
/*
 * Membuat serialized context dari OperationalContext entity.
 * @param {import('../models/OperationalContext.js').OperationalContext} context
 * @returns {SerializedContext}
 */
export function serializeContext(context) {
  if (!context||!context.isValid()) {
    return null
  }
  return {
    schoolId: context.school.id,
    academicYearId: context.academicYear.id,
    semesterId: context.semester.id,
    userId: context.userId,
    role: context.role,
    permissions: [...context.permissions]
  }
}
/*
 * Validasi serialized context data.
 * @param {Object} data
 * @returns {boolean}
 */
export function isValidSerializedContext(data) {
  if (!data) return false
  return !!(
    data.schoolId &&
    data.academicYearId &&
    data.semesterId &&
    data.userId &&
    data.role
  )
}
