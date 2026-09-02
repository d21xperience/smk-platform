// apps/siakad-tu/src/domain/student/value-objects/ContactInfo.js

/**
 * Value Object: ContactInfo
 * Phone dan Email (keduanya optional, tapi minimal satu harus ada)
 */
export class ContactInfo {
  #phone
  #email

  constructor({ phone = '', email = '' } = {}) {
    this.#phone = phone ? String(phone).trim() : ''
    this.#email = email ? String(email).trim() : ''

    // Validasi format jika ada
    if (this.#phone && !/^[\d\s\-+()]{8,20}$/.test(this.#phone)) {
      throw new Error(`Format nomor telepon tidak valid: "${this.#phone}".`)
    }

    if (this.#email && !/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(this.#email)) {
      throw new Error(`Format email tidak valid: "${this.#email}".`)
    }

    // Minimal satu harus ada
    if (!this.#phone && !this.#email) {
      throw new Error('Minimal salah satu dari phone atau email harus diisi.')
    }

    Object.freeze(this)
  }

  get phone() {
    return this.#phone
  }
  get email() {
    return this.#email
  }

  equals(other) {
    if (!(other instanceof ContactInfo)) return false
    return this.#phone === other.#phone && this.#email === other.#email
  }

  toJSON() {
    return { phone: this.#phone, email: this.#email }
  }

  static fromJSON(data) {
    return new ContactInfo(data)
  }
}
