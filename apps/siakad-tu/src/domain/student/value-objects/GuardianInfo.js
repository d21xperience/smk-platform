// apps/siakad-tu/src/domain/student/value-objects/GuardianInfo.js

/**
 * Value Object: GuardianInfo
 * Data orang tua/wali siswa
 */
export class GuardianInfo {
  #name
  #relation // father, mother, guardian
  #phone
  #occupation

  constructor({ name, relation, phone = '', occupation = '' }) {
    if (!name || String(name).trim() === '') {
      throw new Error('Nama wali/orang tua tidak boleh kosong.')
    }
    if (!relation || !['father', 'mother', 'guardian'].includes(relation)) {
      throw new Error(`Relasi tidak valid: "${relation}". Harus father, mother, atau guardian.`)
    }

    this.#name = String(name).trim()
    this.#relation = relation
    this.#phone = phone ? String(phone).trim() : ''
    this.#occupation = occupation ? String(occupation).trim() : ''

    Object.freeze(this)
  }

  get name() {
    return this.#name
  }
  get relation() {
    return this.#relation
  }
  get phone() {
    return this.#phone
  }
  get occupation() {
    return this.#occupation
  }

  get relationLabel() {
    const labels = {
      father: 'Ayah',
      mother: 'Ibu',
      guardian: 'Wali',
    }
    return labels[this.#relation]
  }

  equals(other) {
    if (!(other instanceof GuardianInfo)) return false
    return (
      this.#name === other.#name &&
      this.#relation === other.#relation &&
      this.#phone === other.#phone
    )
  }

  toJSON() {
    return {
      name: this.#name,
      relation: this.#relation,
      phone: this.#phone,
      occupation: this.#occupation,
    }
  }

  static fromJSON(data) {
    return new GuardianInfo(data)
  }
}
