/*
 * @sdp/context — Domain Model: OperationalContext
 *
 * Mengagregasi seluruh konteks yang wajib ada dalam setiap transaksi operasional.
 *
 * Context minimal:
 * - SchoolID
 * - AcademicYearID (ROOT)
 * - SemesterID
 * - UserID
 * - Role
 * - Permissions
 *
 * Academic Year adalah root operational context.
 * Tidak boleh ada transaksi operasional tanpa Academic Year dan Semester.
 *
 * Pure JavaScript — TIDAK bergantung pada Vue, Pinia, Quasar, Axios.
 */
export class OperationalContext {
  /*
   * @param {Object} params
   * @param {Object} params.school - School entity (with toJSON/fromJSON)
   * @param {Object} params.academicYear - AcademicYear entity
   * @param {Object} params.semester - Semester entity
   * @param {string} params.userId - User ID
   * @param {string} params.role - User role
   * @param {Array<string>} params.permissions - List of permission strings
   */
  constructor({ school, academicYear, semester, userId, role, permissions }) {
    this.school = school
    this.academicYear = academicYear
    this.semester = semester
    this.userId = userId
    this.role = role
    this.permissions = permissions||[]
  }
  /*
   * Memastikan seluruh komponen konteks wajib telah terisi.
   * @returns {boolean}
   */
  isValid() {
    return !!(
      this.school?.id &&
      this.academicYear?.id &&
      this.semester?.id &&
      this.userId &&
      this.role
    )
  }
  /*
   * Membentuk cache key yang context-aware.
   * Mencegah kebocoran data antar tahun ajaran/semester.
   *
   * Contoh: "attendance:school01:ay-2026:sem-20261:user123"
   *
   * @param {string} prefix - Prefix untuk domain tertentu
   * @returns {string|null}
   */
  getCacheKey(prefix) {
    if (!this.isValid()) return null
    return `${prefix}:${this.school.id}:${this.academicYear.id}:${this.semester.id}:${this.userId}`
  }
  /*
   * Membentuk storage key yang context-aware.
   * @param {string} prefix - Prefix untuk domain tertentu
   * @returns {string|null}
   */
  getStorageKey(prefix) {
    if (!this.isValid()) return null
    return `${prefix}:${this.school.id}:${this.academicYear.id}:${this.semester.id}`
  }
  /*
   * Memeriksa apakah user memiliki permission tertentu.
   * @param {string} permission - Contoh: "inventory.asset.read"
   * @returns {boolean}
   */
  hasPermission(permission) {
    return this.permissions.includes(permission)
  }
  /*
   * Memeriksa apakah user memiliki salah satu dari beberapa permission.
   * @param {Array<string>} permissions
   * @returns {boolean}
   */
  hasAnyPermission(permissions) {
    return permissions.some(p => this.permissions.includes(p))
  }
  /*
   * Memeriksa apakah user memiliki semua permission yang diberikan.
   * @param {Array<string>} permissions
   * @returns {boolean}
   */
  hasAllPermissions(permissions) {
    return permissions.every(p => this.permissions.includes(p))
  }
  /*
   * Memeriksa apakah user memiliki role tertentu.
   * @param {string} role
   * @returns {boolean}
   */
  hasRole(role) {
    return this.role === role
  }
  /*
   * Mendapatkan SchoolID.
   * @returns {string|null}
   */
  getSchoolId() {
    return this.school?.id||null
  }
  /*
   * Mendapatkan AcademicYearID.
   * @returns {string|null}
   */
  getAcademicYearId() {
    return this.academicYear?.id||null
  }
  /*
   * Mendapatkan SemesterID.
   * @returns {string|null}
   */
  getSemesterId() {
    return this.semester?.id||null
  }
  /*
   * Serialisasi ke plain object (untuk persistence/transport).
   * Tidak menggunakan toJSON() pada nested objects untuk menghindari
   * circular dependency. Aplikasi yang menggunakan package ini
   * bertanggung jawab untuk serialize nested objects.
   * @returns {Object}
   */
  toSerializable() {
    return {
      schoolId: this.school?.id||null,
      academicYearId: this.academicYear?.id||null,
      semesterId: this.semester?.id||null,
      userId: this.userId,
      role: this.role,
      permissions: [...this.permissions]
    }
  }
  /*
   * Membangun OperationalContext dari plain object IDs.
   * Digunakan untuk restore dari storage.
   * Aplikasi harus menyediakan resolver untuk mendapatkan
   * full entities dari IDs.
   *
   * @param {Object} data - Plain object dengan IDs
   * @param {Object} entities - Resolved entities
   * @param {Object} entities.school - School entity
   * @param {Object} entities.academicYear - AcademicYear entity
   * @param {Object} entities.semester - Semester entity
   * @returns {OperationalContext}
   */
  static fromSerializable(data, entities) {
    if (!data) return null
    return new OperationalContext({
      school: entities?.school||null,
      academicYear: entities?.academicYear||null,
      semester: entities?.semester||null,
      userId: data.userId,
      role: data.role,
      permissions: data.permissions||[]
    })
  }
}
