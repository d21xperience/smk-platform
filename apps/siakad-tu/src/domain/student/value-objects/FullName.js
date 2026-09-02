// apps/siakad-tu/src/domain/student/value-objects/FullName.js

/**
 * Value Object: FullName
 *
 * Terdiri dari firstName, middleName (optional), lastName (optional)
 * Aturan: firstName wajib, lainnya optional
 */
export class FullName {
  #firstName
  #middleName
  #lastName

  constructor(firstName, middleName = '', lastName = '') {
    if (!firstName || String(firstName).trim() === '') {
      throw new Error('Nama depan (firstName) tidak boleh kosong.')
    }

    this.#firstName = String(firstName).trim()
    this.#middleName = middleName ? String(middleName).trim() : ''
    this.#lastName = lastName ? String(lastName).trim() : ''

    Object.freeze(this)
  }

  get firstName() {
    return this.#firstName
  }
  get middleName() {
    return this.#middleName
  }
  get lastName() {
    return this.#lastName
  }

  /**
   * Nama lengkap untuk ditampilkan
   */
  get displayName() {
    return [this.#firstName, this.#middleName, this.#lastName].filter(Boolean).join(' ')
  }

  equals(other) {
    if (!(other instanceof FullName)) return false
    return (
      this.#firstName === other.#firstName &&
      this.#middleName === other.#middleName &&
      this.#lastName === other.#lastName
    )
  }

  toJSON() {
    return {
      firstName: this.#firstName,
      middleName: this.#middleName,
      lastName: this.#lastName,
    }
  }

  static fromJSON(data) {
    return new FullName(data.firstName, data.middleName, data.lastName)
  }
}
