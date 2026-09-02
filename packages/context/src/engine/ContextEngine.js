import { OperationalContext } from "../models/OperationalContext.js";
/*
 * @sdp/context — Engine: ContextEngine
 *
 * Berisi business rules untuk validasi dan pembentukan Operational Context.
 *
 * Pure JavaScript — TIDAK boleh mengetahui:
 * Vue, Pinia, Axios, HTTP, Database, Message Broker, localStorage, browser API
 *
 * Tanggung jawab:
 * - Validasi kelengkapan konteks
 * - Validasi relasi semester ↔ academic year
 * - Validasi rentang tanggal
 * - Pembentukan OperationalContext yang valid
 */
export class ContextEngine {
  /*
   * Memvalidasi bahwa semester yang dipilih benar-benar milik tahun ajaran yang dipilih.
   * @param {Object} semester - Semester entity
   * @param {Object} academicYear - AcademicYear entity
   * @throws {Error} jika semester bukan milik academic year
   * @returns {boolean} true jika valid
   */
  static validateSemesterBelongsToAcademicYear(semester, academicYear) {
    if (!semester || !academicYear) {
      throw new Error("Semester and Academic Year are required for validation");
    }
    if (semester.academicYearId !== academicYear.id) {
      throw new Error("Semester does not belong to the selected Academic Year");
    }
    return true;
  }
  /*
   * Memvalidasi rentang tanggal semester.
   * @param {Object} semester - Semester entity dengan startDate dan endDate
   * @throws {Error} jika tanggal tidak valid
   * @returns {boolean} true jika valid
   */
  static validateSemesterDates(semester) {
    if (!semester.startDate || !semester.endDate) {
      throw new Error("Semester start and end dates are required");
    }
    if (semester.startDate > semester.endDate) {
      throw new Error(
        "Invalid semester date range: Start date is after end date",
      );
    }
    return true;
  }
  /*
   * Memvalidasi rentang tanggal academic year.
   * @param {Object} academicYear - AcademicYear entity
   * @throws {Error} jika tanggal tidak valid
   * @returns {boolean} true jika valid
   */
  static validateAcademicYearDates(academicYear) {
    if (!academicYear.startDate || !academicYear.endDate) {
      throw new Error("Academic Year start and end dates are required");
    }
    if (academicYear.startDate > academicYear.endDate) {
      throw new Error(
        "Invalid academic year date range: Start date is after end date",
      );
    }
    return true;
  }
  /*
   * Memvalidasi bahwa semester berada dalam rentang academic year.
   * @param {Object} semester
   * @param {Object} academicYear
   * @throws {Error} jika semester di luar rentang academic year
   * @returns {boolean}
   */
  static validateSemesterWithinAcademicYear(semester, academicYear) {
    if (
      !semester.startDate ||
      !semester.endDate ||
      !academicYear.startDate ||
      !academicYear.endDate
    ) {
      return true; // Skip validation if dates are not available
    }
    if (
      semester.startDate < academicYear.startDate ||
      semester.endDate > academicYear.endDate
    ) {
      throw new Error("Semester dates must be within Academic Year date range");
    }
    return true;
  }
  /*
   * Memvalidasi user context.
   * @param {Object} user - Object minimal { id, role }
   * @throws {Error} jika user tidak valid
   * @returns {boolean}
   */
  static validateUser(user) {
    if (!user || !user.id || !user.role) {
      throw new Error(
        "Valid User with ID and Role is required for Operational Context",
      );
    }
    return true;
  }
  /*
   * Membangun dan memvalidasi Operational Context secara utuh.
   *
   * @param {Object} params
   * @param {Object} params.school - School entity
   * @param {Object} params.academicYear - AcademicYear entity
   * @param {Object} params.semester - Semester entity
   * @param {Object} params.user - Object minimal { id, role, permissions }
   * @throws {Error} jika validasi gagal
   * @returns {OperationalContext}
   */
  static buildContext({ school, academicYear, semester, user }) {
    // Validate school
    if (!school || !school.id) {
      throw new Error("Valid School is required for Operational Context");
    }
    // Validate academic year
    if (!academicYear || !academicYear.id) {
      throw new Error(
        "Valid Academic Year is required for Operational Context",
      );
    }
    // Validate semester
    if (!semester || !semester.id) {
      throw new Error("Valid Semester is required for Operational Context");
    }
    // Validate user
    this.validateUser(user);
    // Validate relationships
    this.validateSemesterBelongsToAcademicYear(semester, academicYear);
    this.validateSemesterDates(semester);
    // Build context
    const context = new OperationalContext({
      school,
      academicYear,
      semester,
      userId: user.id,
      role: user.role,
      permissions: user.permissions || [],
    });
    if (!context.isValid()) {
      throw new Error("Failed to build a valid Operational Context");
    }
    return context;
  }
  /*
   * Memeriksa apakah konteks masih valid (belum expired).
   * Berguna untuk check saat restore dari storage.
   *
   * @param {OperationalContext} context
   * @param {Date|null} checkDate - Tanggal untuk dicek (default: sekarang)
   * @returns {boolean}
   */
  static isContextStillValid(context, checkDate) {
    if (!context || !context.isValid()) {
      return false;
    }
    const date = checkDate || new Date();
    // Check if semester is still within date range
    if (context.semester?.startDate && context.semester?.endDate) {
      if (
        date < context.semester.startDate ||
        date > context.semester.endDate
      ) {
        return false;
      }
    }
    return true;
  }
}
