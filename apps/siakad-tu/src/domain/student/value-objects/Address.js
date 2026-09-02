// apps/siakad-tu/src/domain/student/value-objects/Address.js

/**
 * Value Object: Address
 * Structured address sesuai standar Indonesia
 */
export class Address {
  #street
  #rtRw
  #village
  #district
  #city
  #postalCode

  constructor({ street, rtRw = '', village, district, city, postalCode }) {
    if (!street || String(street).trim() === '') {
      throw new Error('Alamat (street) tidak boleh kosong.')
    }
    if (!village || String(village).trim() === '') {
      throw new Error('Kelurahan (village) tidak boleh kosong.')
    }
    if (!city || String(city).trim() === '') {
      throw new Error('Kota (city) tidak boleh kosong.')
    }

    this.#street = String(street).trim()
    this.#rtRw = rtRw ? String(rtRw).trim() : ''
    this.#village = String(village).trim()
    this.#district = district ? String(district).trim() : ''
    this.#city = String(city).trim()
    this.#postalCode = postalCode ? String(postalCode).trim() : ''

    Object.freeze(this)
  }

  get street() {
    return this.#street
  }
  get rtRw() {
    return this.#rtRw
  }
  get village() {
    return this.#village
  }
  get district() {
    return this.#district
  }
  get city() {
    return this.#city
  }
  get postalCode() {
    return this.#postalCode
  }

  /**
   * Format lengkap untuk ditampilkan
   */
  get displayFormat() {
    const parts = [this.#street]
    if (this.#rtRw) parts.push(`RT/RW ${this.#rtRw}`)
    parts.push(this.#village)
    if (this.#district) parts.push(this.#district)
    parts.push(this.#city)
    if (this.#postalCode) parts.push(this.#postalCode)
    return parts.join(', ')
  }

  equals(other) {
    if (!(other instanceof Address)) return false
    return (
      this.#street === other.#street &&
      this.#rtRw === other.#rtRw &&
      this.#village === other.#village &&
      this.#district === other.#district &&
      this.#city === other.#city &&
      this.#postalCode === other.#postalCode
    )
  }

  toJSON() {
    return {
      street: this.#street,
      rtRw: this.#rtRw,
      village: this.#village,
      district: this.#district,
      city: this.#city,
      postalCode: this.#postalCode,
    }
  }

  static fromJSON(data) {
    return new Address(data)
  }
}
