export class ContextService {
  /**
   * @param {Object} adapter - Objek dengan method getAcademicYears(), getSemesters(), getAcademicPeriods()
   */
  constructor(adapter) {
    this.adapter = adapter
  }

  /**
   * Ambil daftar semua tahun ajaran
   * @returns {Promise<import('src/models/AcademicYear').AcademicYear[]>}
   */
  async fetchAcademicYears() {
    // const tes = await this.adapter.getAcademicYears()
    return this.adapter.getAcademicYears()
  }

  /**
   * Ambil daftar semester berdasarkan tahun ajaran
   * @param {number} academicYearId
   * @returns {Promise<import('src/models/Semester').Semester[]>}
   */
  async fetchSemesters(academicYearId) {
    return this.adapter.getSemesters(academicYearId)
  }

  /**
   * Ambil daftar periode akademik berdasarkan semester
   * @param {number} semesterId
   * @returns {Promise<import('src/models/AcademicPeriod').AcademicPeriod[]>}
   */
  async fetchAcademicPeriods(semesterId) {
    return this.adapter.getAcademicPeriods(semesterId)
  }

  /**
   * Membentuk objek Operational Context dari data yang dipilih
   * @param {import('src/models/AcademicYear').AcademicYear} academicYear
   * @param {import('src/models/Semester').Semester} semester
   * @param {import('src/models/AcademicPeriod').AcademicPeriod} [academicPeriod=null]
   * @returns {{ academicYearId: number, academicYearName: string, semesterId: number, semesterName: string, academicPeriodId: number|null, academicPeriodName: string|null }}
   */
  buildOperationalContext(academicYear, semester, academicPeriod = null) {
    if (!academicYear || !semester) {
      throw new Error('AcademicYear dan Semester wajib diisi')
    }

    return {
      academicYearId: academicYear.id,
      academicYearName: academicYear.name,
      semesterId: semester.id,
      semesterName: semester.name,
      academicPeriodId: academicPeriod?.id || null,
      academicPeriodName: academicPeriod?.name || null,
    }
  }
}
