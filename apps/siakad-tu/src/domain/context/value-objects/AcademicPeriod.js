// apps/siakad-tu/src/domain/context/AcademicPeriod.js

/**
 * Value Object: AcademicPeriod
 * Mengkapsulasi logika derivasi periodId dari academicYear + semester.
 * Immutable setelah dibuat.
 *
 * Format periodId: YYYY + S
 * Contoh: "2026/2027" + semester 1 → "20261"
 *         "2026/2027" + semester 2 → "20262"
 */
export class AcademicPeriod {
  constructor(academicYear, semester) {
    this._validate(academicYear, semester)
    this._academicYear = academicYear
    this._semester = semester
    this._periodId = this._derivePeriodId(academicYear, semester)
  }

  _validate(academicYear, semester) {
    // Validasi format academicYear: "YYYY/YYYY"
    const yearPattern = /^\d{4}\/\d{4}$/
    if (!yearPattern.test(academicYear)) {
      throw new Error(
        `Format academicYear tidak valid: "${academicYear}". Harus "YYYY/YYYY", contoh "2026/2027".`,
      )
    }

    // Validasi konsistensi tahun
    const [startYear, endYear] = academicYear.split('/').map(Number)
    if (endYear !== startYear + 1) {
      throw new Error(
        `AcademicYear tidak konsisten: "${academicYear}". Tahun akhir harus tahun awal + 1.`,
      )
    }

    // Validasi semester
    if (semester !== 1 && semester !== 2) {
      throw new Error(`Semester harus 1 atau 2. Diterima: ${semester}`)
    }
  }

  _derivePeriodId(academicYear, semester) {
    const startYear = academicYear.split('/')[0]
    return `${startYear}${semester}`
  }

  get academicYear() {
    return this._academicYear
  }
  get semester() {
    return this._semester
  }
  get periodId() {
    return this._periodId
  }

  get semesterLabel() {
    return `${this._semester}`
  }

  get displayLabel() {
    return `${this._academicYear} - ${this.semesterLabel}`
  }

  /**
   * Cek apakah period ini sama dengan period lain
   */
  equals(other) {
    if (!(other instanceof AcademicPeriod)) return false
    return this._periodId === other._periodId
  }

  /**
   * Serialize untuk dikirim ke backend
   */
  toJSON() {
    return {
      academicYear: this._academicYear,
      semester: this._semester,
      periodId: this._periodId,
    }
  }
}
