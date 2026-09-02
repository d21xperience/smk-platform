// apps/siakad-tu/src/domain/student/value-objects/NIS.js

/**
 * Value Object: NIS (Nomor Induk Siswa)
 *
 * Aturan:
 * - Unik dalam satu sekolah (validasi dilakukan di Repository/Service)
 * - Format bebas, tapi tidak boleh kosong
 * - Immutable
 */
export class NIS {
  #value

  constructor(value) {
    if (!value || String(value).trim() === '') {
      throw new Error('NIS tidak boleh kosong.')
    }
    this.#value = String(value).trim()
    Object.freeze(this)
  }

  get value() {
    return this.#value
  }

  equals(other) {
    if (!(other instanceof NIS)) return false
    return this.#value === other.#value
  }

  toString() {
    return this.#value
  }

  toJSON() {
    return this.#value
  }

  static fromJSON(value) {
    return new NIS(value)
  }
}
