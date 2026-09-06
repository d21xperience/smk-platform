// apps/siakad-tu/src/domain/context/OperationalContext.js

import { AcademicPeriod } from './AcademicPeriod.js'

/**
 * Aggregate Root: OperationalContext
 *
 * Setiap transaksi dalam sistem HARUS berjalan di dalam OperationalContext.
 * Context ini dipilih satu kali setelah login dan digunakan oleh seluruh modul.
 *
 * Immutable setelah dibuat. Untuk mengganti context, buat instance baru
 * dan dispatch event OperationalContextChanged.
 */
export class OperationalContext {
  constructor({ schoolId, schoolName, academicYear, semester }) {
    if (!schoolId) throw new Error('OperationalContext: schoolId wajib diisi.')
    if (!schoolName) throw new Error('OperationalContext: schoolName wajib diisi.')

    this._schoolId = schoolId
    this._schoolName = schoolName
    this._period = new AcademicPeriod(academicYear, semester)
  }

  get schoolId() {
    return this._schoolId
  }
  get schoolName() {
    return this._schoolName
  }
  get academicYear() {
    return this._period.academicYear
  }
  get semester() {
    return this._period.semester
  }
  get periodId() {
    return this._period.periodId
  }
  get period() {
    return this._period
  }

  /**
   * Label lengkap untuk ditampilkan di UI (header, breadcrumb, dll)
   * Contoh: "SMK Pasundan Jatinangor | 2026/2027 - Semester 1"
   */
  get displayLabel() {
    return `${this._period.displayLabel}`
  }
  get schoolLabel() {
    return `${this._schoolName}`
  }

  /**
   * Cek apakah context ini lengkap dan valid untuk melakukan transaksi
   */
  isValid() {
    return !!(this._schoolId && this._period.periodId)
  }

  /**
   * Cek apakah context ini sama dengan context lain
   */
  equals(other) {
    if (!(other instanceof OperationalContext)) return false
    return this._schoolId === other._schoolId && this._period.equals(other._period)
  }

  /**
   * Serialize untuk header API
   */
  toHeaders() {
    return {
      'X-School-Id': this._schoolId,
      'X-Academic-Period': this._period.periodId,
    }
  }

  /**
   * Serialize untuk disimpan di localStorage / dikirim ke backend
   */
  toJSON() {
    return {
      schoolId: this._schoolId,
      schoolName: this._schoolName,
      ...this._period.toJSON(),
    }
  }

  /**
   * Rehydrate dari JSON (localStorage / API response)
   */
  static fromJSON(data) {
    return new OperationalContext({
      schoolId: data.schoolId,
      schoolName: data.schoolName,
      academicYear: data.academicYear,
      semester: data.semester,
    })
  }
}
