// apps/siakad-tu/src/domain/student/value-objects/Gender.js

/**
 * Value Object: Gender
 * Enum: MALE, FEMALE
 */
export class Gender {
  static MALE = new Gender('MALE', 'Laki-laki')
  static FEMALE = new Gender('FEMALE', 'Perempuan')

  static VALUES = [Gender.MALE, Gender.FEMALE]

  #code
  #label

  constructor(code, label) {
    this.#code = code
    this.#label = label
    Object.freeze(this)
  }

  get code() {
    return this.#code
  }
  get label() {
    return this.#label
  }

  equals(other) {
    if (!(other instanceof Gender)) return false
    return this.#code === other.#code
  }

  toString() {
    return this.#code
  }

  toJSON() {
    return this.#code
  }

  static fromJSON(code) {
    const found = Gender.VALUES.find((g) => g.#code === code)
    if (!found) {
      throw new Error(`Gender tidak valid: "${code}". Harus MALE atau FEMALE.`)
    }
    return found
  }
}
