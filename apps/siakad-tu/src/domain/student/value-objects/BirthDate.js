// apps/siakad-tu/src/domain/student/value-objects/BirthDate.js

/**
 * Value Object: BirthDate
 *
 * Menyimpan tanggal lahir dan menghitung usia secara dinamis.
 * Validasi: tidak boleh di masa depan, tidak lebih dari 100 tahun.
 */
export class BirthDate {
  #date

  constructor(date) {
    const parsed = date instanceof Date ? date : new Date(date)

    if (isNaN(parsed.getTime())) {
      throw new Error(`Tanggal lahir tidak valid: "${date}".`)
    }

    const now = new Date()
    if (parsed > now) {
      throw new Error('Tanggal lahir tidak boleh di masa depan.')
    }

    const hundredYearsAgo = new Date()
    hundredYearsAgo.setFullYear(now.getFullYear() - 100)
    if (parsed < hundredYearsAgo) {
      throw new Error('Tanggal lahir tidak boleh lebih dari 100 tahun yang lalu.')
    }

    this.#date = parsed
    Object.freeze(this)
  }

  get date() {
    return this.#date
  }

  /**
   * Hitung usia dalam tahun
   */
  get ageInYears() {
    const now = new Date()
    let age = now.getFullYear() - this.#date.getFullYear()
    const monthDiff = now.getMonth() - this.#date.getMonth()

    if (monthDiff < 0 || (monthDiff === 0 && now.getDate() < this.#date.getDate())) {
      age--
    }

    return age
  }

  /**
   * Format untuk ditampilkan (DD/MM/YYYY)
   */
  get displayFormat() {
    const day = String(this.#date.getDate()).padStart(2, '0')
    const month = String(this.#date.getMonth() + 1).padStart(2, '0')
    const year = this.#date.getFullYear()
    return `${day}/${month}/${year}`
  }

  equals(other) {
    if (!(other instanceof BirthDate)) return false
    return this.#date.getTime() === other.#date.getTime()
  }

  toJSON() {
    return this.#date.toISOString()
  }

  static fromJSON(isoString) {
    return new BirthDate(isoString)
  }
}
