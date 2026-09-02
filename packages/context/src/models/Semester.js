/*
 * @sdp/context — Domain Model: Semester
 *
 * Merepresentasikan Semester dalam sebuah Tahun Ajaran.
 *
 * Contoh untuk Academic Year 2026/2027:
 *   - Semester 1: code "20261"
 *   - Semester 2: code "20262"
 *
 * Pure JavaScript — TIDAK bergantung pada Vue, Pinia, Quasar, Axios.
 */
export class Semester {
  /*
   * @param {Object} params
   * @param {string} params.id - Unique identifier
   * @param {string} params.academicYearId - Reference ke AcademicYear
   * @param {string} params.name - Contoh: "Semester 1 (Ganjil)"
   * @param {string} params.code - Contoh: "20261", "20262"
   * @param {string|Date} params.startDate - Tanggal mulai
   * @param {string|Date} params.endDate - Tanggal selesai
   * @param {boolean} params.isActive - Status aktif
   */
  constructor({
    id,
    academicYearId,
    name,
    code,
    startDate,
    endDate,
    isActive,
  }) {
    if (!id) throw new Error("Semester ID is required");
    if (!academicYearId)
      throw new Error("Academic Year ID is required for Semester");

    this.id = id;
    this.academicYearId = academicYearId;
    this.name = name || "";
    this.code = code || "";
    this.startDate = startDate ? new Date(startDate) : null;
    this.endDate = endDate ? new Date(endDate) : null;
    this.isActive = Boolean(isActive);
  }
  /*
   * Memeriksa apakah semester sedang aktif berdasarkan tanggal.
   * @param {Date|null} checkDate - Tanggal untuk dicek (default: sekarang)
   * @returns {boolean}
   */
  isCurrentlyActive(checkDate) {
    if (!this.startDate || !this.endDate) return this.isActive;
    const date = checkDate || new Date();
    return date >= this.startDate && date <= this.endDate;
  }
  /*
   * Mendapatkan label display.
   * @returns {string}
   */
  getDisplayLabel() {
    return this.name || this.code;
  }
  /*
   * Serialisasi ke plain object.
   * @returns {Object}
   */
  toJSON() {
    return {
      id: this.id,
      academicYearId: this.academicYearId,
      name: this.name,
      code: this.code,
      startDate: this.startDate ? this.startDate.toISOString() : null,
      endDate: this.endDate ? this.endDate.toISOString() : null,
      isActive: this.isActive,
    };
  }
  /*
   * Deserialisasi dari plain object.
   * @param {Object} data
   * @returns {Semester}
   */
  static fromJSON(data) {
    if (!data) return null;
    return new Semester({
      id: data.id,
      academicYearId: data.academicYearId,
      name: data.name,
      code: data.code,
      startDate: data.startDate,
      endDate: data.endDate,
      isActive: data.isActive,
    });
  }
}
